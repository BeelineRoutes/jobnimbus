/** ****************************************************************************************************************** **
	Calls related to jobs

    
** ****************************************************************************************************************** **/

package jobnimbus 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
    "time"
    "encoding/json"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type JobStatus int 
const (
    JobStatus_lead          JobStatus = 355
    JobStatus_scheduled     JobStatus = 356

)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Job struct {
    Id string `json:"jnid,omitempty"`
    Name string `json:"name"`
    DateStart int64 `json:"date_start"`
    DateEnd int64 `json:"date_end"`
    Active bool `json:"is_active,omitempty"`
    Archived bool `json:"is_archived,omitempty"`
    Type string `json:"type,omitempty"`
    Status JobStatus `json:"status,omitempty"`
    StatusName string `json:"status_name,omitempty"`
    Tags []string `json:"tags,omitempty"`
    SalesRep string `json:"sales_rep"`
    SalesRepName string `json:"sales_rep_name"`
    RecordType string `json:"record_type_name"`
    MobilePhone string `json:"parent_mobile_phone"`
    Address1 string `json:"address_line1"`
    Address2 string `json:"address_line2"`
    City string `json:"city"`
    State string `json:"state_text"`
    Zip string `json:"zip"`
    Primary struct {
        Id string `json:"id"`
        Email string `json:"email,omitempty"`
        Name string `json:"name,omitempty"`
    } `json:"primary"`
    Geo struct {
        Lat float64 `json:"lat"`
        Lon float64 `json:"lon"`
    } `json:"geo"`
}

type rangeFilter struct {
    Range struct {
        DateStart struct {
            Gte int64 `json:"gte"`
            Lte int64 `json:"lte"`
        } `json:"date_start"`
    } `json:"range"`
}
type jobFilter struct {
    Must []rangeFilter `json:"must"`
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// returns all jobs that match our conditions
// this defaults to 1000 jobs returned, which i'm hoping is enough, so i'm not going to loop through them
func (this *Jobnimbus) ListJobs (ctx context.Context, token string, start, end time.Time) ([]*Job, error) {
    params := url.Values{}

    // create our filter for the time range to look for jobs
    filter := &jobFilter{}
    rFilter := rangeFilter{}
    if start.IsZero() == false {
        rFilter.Range.DateStart.Gte = start.Unix()
    }
    if end.IsZero() == false {
        rFilter.Range.DateStart.Lte = end.Unix()
    }
    filter.Must = append (filter.Must, rFilter) // still append it, this is how we search for "unscheduled" jobs

    filterJson, err := json.Marshal(filter)
    if err != nil { return nil, errors.WithStack(err) } // bail

    params.Set("filter", string(filterJson))

    params.Set("fields", "date_end,date_start,is_active,is_archived,type,status,tags,sales_rep,record_type_name,primary,geo,name,parent_mobile_phone,address_line1,address_line2,city,state_text,zip") // these are really the items we care about

    var resp struct {
        Results []*Job
    }
    
    err = this.send (ctx, 0, http.MethodGet, token, fmt.Sprintf("jobs?%s", params.Encode()), nil, &resp)
    if err != nil { return nil, err } // bail
    
    // we're here, we're good
    return resp.Results, nil 
}

// updates the start/end time for a job
// the time needs to be set to whatever timezone the user's account is in
func (this *Jobnimbus) UpdateJobSchedule (ctx context.Context, token, jobId, salesRep string, startTime time.Time, duration time.Duration) error {
    var data struct {
        DateEnd int64 `json:"date_end"`
        DateStart int64 `json:"date_start"`
        SalesRep string `json:"sales_rep"`
    }
    
    data.DateEnd = startTime.Add(duration).Unix()
    data.DateStart = startTime.Unix()
    data.SalesRep = salesRep

    err := this.send (ctx, 0, http.MethodPut, token, fmt.Sprintf("jobs/%s", jobId), data, nil)
    if err != nil { return err } // bail
    
    // we're here, we're good
    return nil
}

// creates a new job in the system
func (this *Jobnimbus) CreateJob (ctx context.Context, token string, job *Job) (string, error) {
    resp := &Job{}
    
    err := this.send (ctx, 0, http.MethodPost, token, "jobs", job, resp)
    if err != nil { return "", err } // bail
    
    // we're here, we're good
    return resp.Id, nil
}

