/** ****************************************************************************************************************** **
	Calls related to tasks

    
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

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type owner struct {
    Id string `json:"id"`
}

type Task struct {
    Id string `json:"jnid,omitempty"`
    DateStart int64 `json:"date_start"`
    DateEnd int64 `json:"date_end"`
    Description string `json:"description,omitempty"`
    Active bool `json:"is_active,omitempty"`
    Archived bool `json:"is_archived,omitempty"`
    Completed bool `json:"is_completed,omitempty"`
    RecordType string `json:"record_type_name"`
    Type string `json:"type,omitempty"`
    Title string `json:"title"`
    Tags []string `json:"tags,omitempty"`
    Owners []owner `json:"owners"`
    Related []struct {
        Id string `json:"id"`
        Type string `json:"type"`
    } `json:"related,omitempty"`
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// returns all tasks that match our conditions
// this defaults to 1000 tasks returned, which i'm hoping is enough, so i'm not going to loop through them
func (this *Jobnimbus) ListTasks (ctx context.Context, token string, start, end time.Time) ([]*Task, error) {
    params := url.Values{}

    // create our filter for the time range to look for tasks
    filter := &jobFilter{}
    rFilter := rangeFilter{}
    if start.IsZero() == false {
        rFilter.Range.DateStart.Gte = start.Unix()
    }
    if end.IsZero() == false {
        rFilter.Range.DateStart.Lte = end.Unix()
    }
    filter.Must = append (filter.Must, rFilter) // still append it, this is how we search for "unscheduled" tasks

    filterJson, err := json.Marshal(filter)
    if err != nil { return nil, errors.WithStack(err) } // bail

    params.Set("filter", string(filterJson))

    params.Set("fields", "date_end,date_start,is_active,is_archived,is_completed,type,tags,description,owners,record_type_name,title") // these are really the items we care about

    var resp struct {
        Results []*Task
    }
    
    err = this.send (ctx, 0, http.MethodGet, token, fmt.Sprintf("tasks?%s", params.Encode()), nil, &resp)
    if err != nil { return nil, err } // bail
    
    // we're here, we're good
    return resp.Results, nil 
}

// updates the start/end time for a task
// the time needs to be set to whatever timezone the user's account is in
func (this *Jobnimbus) UpdateTaskSchedule (ctx context.Context, token, taskId, salesRep string, startTime time.Time, duration time.Duration) error {
    var data struct {
        DateEnd int64 `json:"date_end"`
        DateStart int64 `json:"date_start"`
        Owners []owner `json:"owners"`
    }
    
    data.DateEnd = startTime.Add(duration).Unix()
    data.DateStart = startTime.Unix()
    data.Owners = append(data.Owners, owner { salesRep })

    err := this.send (ctx, 0, http.MethodPut, token, fmt.Sprintf("tasks/%s", taskId), data, nil)
    if err != nil { return err } // bail
    
    // we're here, we're good
    return nil
}

// creates a new task in the system
func (this *Jobnimbus) CreateTask (ctx context.Context, token string, task *Task) (string, error) {
    return "", nil // can't figure this out yet
    /*
    resp := &Task{}
    
    err := this.send (ctx, 0, http.MethodPost, token, "tasks", task, resp)
    if err != nil { return "", err } // bail
    
    // we're here, we're good
    return resp.Id, nil
    */
}

