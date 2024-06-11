
package jobnimbus 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestJobs (t *testing.T) {
	w := &Jobnimbus{}
	cfg := getRealConfig(t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	jobs, err := w.ListJobs (ctx, cfg.Token, time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 1, 1))
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.NotEqual (t, "", jobs[0].Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Primary.Id, "not filled in")
	assert.Equal (t, true, jobs[0].Geo.Lat > 1, "not filled in")
	
	/*
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	*/
}

