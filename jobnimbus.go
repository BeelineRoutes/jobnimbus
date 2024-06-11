/** ****************************************************************************************************************** **
    Jobnimbus API wrapper
    written for GoLang
    Created 2024-06-05 by Nathan Thomas 
    Courtesy of BeelineRoutes.com

    current docs in v1
    https://documenter.getpostman.com/view/3919598/S11PpG4x#getting-started

** ****************************************************************************************************************** **/

package jobnimbus 

import (
    "github.com/pkg/errors"

    // "fmt"
    "net/http"
    "encoding/json"
    "os"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

const apiURL = "https://app.jobnimbus.com/api1"

var (
    ErrUnexpected       = errors.New("idk...")
	ErrNotFound 		= errors.New("Item was not found")
	ErrTooManyRecords	= errors.New("Too many records returned")
    ErrAuthExpired      = errors.New("Auth Expired")
    ErrQuota            = errors.New("Too many requests - quota limit")
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Config struct {
    Token string 
}

func (this Config) Valid () bool {
    if len(this.Token) < 10 { return false } // i'm making these 10 so the example_config comes back as false
    
    return true 
}


//----- ERRORS ---------------------------------------------------------------------------------------------------------//
type Error struct {
	Msg string
	StatusCode int
}

func (this *Error) Err () error {
	if this == nil { return nil } // no error
	switch this.StatusCode {
	case http.StatusUnauthorized:
        return errors.Wrapf (ErrAuthExpired, "Unauthorized : %d : %s", this.StatusCode, this.Msg)
	}
	// just a default
	return errors.Wrapf (ErrUnexpected, "Jobnimbus Error : %d : %s", this.StatusCode, this.Msg)
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CLASS -----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Jobnimbus struct {
	
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// this is just used for local testing
// so you don't have to keep your actual tokens in the repo
func parseConfig (jsonFile string) (*Config, error) {
	config, err := os.Open(jsonFile)
	if err != nil { return nil, errors.WithStack (err) }

	jsonParser := json.NewDecoder (config)

    ret := &Config{}
	err = jsonParser.Decode (ret)
    return ret, errors.WithStack(err)
}
