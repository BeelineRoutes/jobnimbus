/** ****************************************************************************************************************** **
	The actual sending and receiving stuff
	Reused for most of the calls to Jobnimbus
	
** ****************************************************************************************************************** **/

package jobnimbus 

import (
    "github.com/pkg/errors"

    "fmt"
    "net/http"
    "context"
    "encoding/json"
    "io/ioutil"
    "bytes"
	"time"
	"math"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// handles making the request and reading the results from it 
// if there's an error the Error object will be set, otherwise it will be nil
func (this *Jobnimbus) finish (req *http.Request, out interface{}) error {
	resp, err := http.DefaultClient.Do (req)
	
	if err != nil { return errors.WithStack (err) }
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll (resp.Body)

    if resp.StatusCode > 399 { 
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			// special error 
			return errors.Wrapf (ErrAuthExpired, "Unauthorized : %d : %s", resp.StatusCode, string(body))

		case http.StatusTooManyRequests:
			return errors.Wrapf (ErrQuota, "Quota : %d : %s", resp.StatusCode, string(body))
		}
		// just a default
		err = errors.Wrapf (ErrUnexpected, "Jobnimbus Error : %d : %s", resp.StatusCode, string(body))

		return err
    }
	
	if out != nil { 
		err = errors.WithStack (json.Unmarshal (body, out))
		if err != nil {
			err = errors.Wrap (err, string(body)) // if it didn't unmarshal, include the body so we know what it did look like
		}
	}
	
	return err // we're good
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// this recurses
// retries itself on a 429 - ErrQuota
func (this *Jobnimbus) send (ctx context.Context, retries int, requestType, token, link string, in, out interface{}) error {
	if ctx.Err() != nil { return ctx.Err() } // bail on a context timeout

	var jstr []byte 
	var err error 

	header := make(map[string]string)
	header["Authorization"] = "bearer " + token

	if in != nil {
		jstr, err = json.Marshal (in)
		if err != nil { return errors.WithStack (err) }

		header["Content-Type"] = "application/json"
	}
	
	req, err := http.NewRequestWithContext (ctx, requestType, fmt.Sprintf ("%s/%s", apiURL, link), bytes.NewBuffer(jstr))
	if err != nil { return errors.Wrap (err, link) }

	for key, val := range header { req.Header.Set (key, val) }
	err = this.finish (req, out)

	switch errors.Cause(err) {
	case ErrQuota:
		if retries < 6 { // 6 gives 1 + 3 + 7 + 15 + 31 + 63 seconds wait
			time.Sleep (time.Second * time.Duration(math.Pow(2, float64(retries)))) // exp timeout for sleeping
			return this.send (ctx, retries +1, requestType, token, link, in, out)
		}
	}
	
	return errors.Wrapf (err, " %s : %s", link, string(jstr))
}
