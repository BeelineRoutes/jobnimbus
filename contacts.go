/** ****************************************************************************************************************** **
	Calls related to contacts (customers)

    
** ****************************************************************************************************************** **/

package jobnimbus 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
    "encoding/json"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//


  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Contact struct {
    Id string `json:"jnid,omitempty"`
    Company string `json:"company"`
    Country string `json:"country_name"`
    DisplayName string `json:"display_name"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Email string `json:"email"`
    Activbe bool `json:"is_active,omitempty"`
    Archived bool `json:"is_archived,omitempty"`
    Tags []string `json:"tags,omitempty"`
    Address1 string `json:"address_line1"`
    Address2 string `json:"address_line2"`
    City string `json:"city"`
    State string `json:"state_text"`
    Zip string `json:"zip"`
    HomePhone string `json:"home_phone"`
    MobilePhone string `json:"mobile_phone"`
    Geo struct {
        Lat float64 `json:"lat"`
        Lon float64 `json:"lon"`
    } `json:"geo"`
}

type shouldFilter struct {
    Prefix map[string]string
}

type boolFilter struct {
    Bool struct {
        Should []shouldFilter `json:"should"`
    } `json:"bool"`
}

type contactFilter struct {
    Must []boolFilter `json:"must"`
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// returns all jobs that match our conditions
// this defaults to 1000 jobs returned, which i'm hoping is enough, so i'm not going to loop through them
func (this *Jobnimbus) ListContacts (ctx context.Context, token string, searchStreet, searchFirst, searchLast string) ([]*Contact, error) {
    params := url.Values{}

    // create our filter for the time range to look for jobs
    filter := &contactFilter{}
    bFilter := boolFilter{}
    
    firstName := make(map[string]string)
    firstName["first_name"] = searchFirst 

    lastName := make(map[string]string)
    lastName["last_name"] = searchLast 

    street := make(map[string]string)
    street["address_line1"] = searchStreet 

    bFilter.Bool.Should = append(bFilter.Bool.Should, shouldFilter{firstName}, shouldFilter{lastName}, shouldFilter{street})

    filter.Must = append(filter.Must, bFilter)

    filterJson, err := json.Marshal(filter)
    if err != nil { return nil, errors.WithStack(err) } // bail

    params.Set("filter", string(filterJson))

    params.Set("fields", "address_line1,address_line2,city,company,country_name,jnid,display_name,first_name,last_name,email,geo,home_phone,is_active,is_archived,mobile_phone,state_text,tags,zip") // these are really the items we care about

    var resp struct {
        Results []*Contact
    }
    
    err = this.send (ctx, 0, http.MethodGet, token, fmt.Sprintf("contacts?%s", params.Encode()), nil, &resp)
    if err != nil { return nil, err } // bail
    
    // we're here, we're good
    return resp.Results, nil 
}

// creates a new contact in the system
func (this *Jobnimbus) CreateContact (ctx context.Context, token string, contact *Contact) (string, error) {
    resp := &Contact{}
    
    err := this.send (ctx, 0, http.MethodPost, token, "contacts", contact, resp)
    if err != nil { return "", err } // bail
    
    // we're here, we're good
    return resp.Id, nil
}

