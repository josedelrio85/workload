package workload

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	voalarm "github.com/bysidecar/voalarm"
)

// Response blablab
type Response struct {
	Writer  http.ResponseWriter
	Code    int
	Message string    `json:"message,omitempty"`
	Project []Project `json:"project"`
	Work    []Work    `json:"work"`
	Status  []Status  `json:"status"`
	Person  []Person  `json:"person"`
	Error   error
}

// Result blablab
type Result struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Project []Project `json:"project,omitempty"`
	Work    []Work    `json:"work,omitempty"`
	Status  []Status  `json:"status,omitempty"`
	Person  []Person  `json:"person,omitempty"`
}

// Ko generates log, alarm and response when an error occurs and returns an http code and message
func (r *Response) Ko() {
	if r.Error != nil {
		e := &errorLogger{r.Message, r.Code, r.Error, logError(r.Error)}
		e.sendAlarm()
	}

	result := Result{
		Success: false,
		Message: r.Message,
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(r.Code)
	json.NewEncoder(r.Writer).Encode(result)

}

// Ok calls response function with proper data to generate an OK response
func (r *Response) Ok(tr int64) {

	result := Result{
		Success: true,
		Message: r.Message,
	}

	switch tr {
	case 1:
		result.Project = r.Project
		break
	case 2:
		result.Work = r.Work
		break
	case 3:
		result.Status = r.Status
		break
	case 4:
		result.Person = r.Person
		break
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(r.Code)
	json.NewEncoder(r.Writer).Encode(result)
}

// errorLogger is a struct to handle error properties
type errorLogger struct {
	msg    string
	status int
	err    error
	log    string
}

// sendAlarm to VictorOps plattform and format the error for more info
func (e *errorLogger) sendAlarm() {
	e.msg = fmt.Sprintf("Workload -> %s", e.msg)
	log.Println(e.log)

	mstype := voalarm.Acknowledgement
	switch e.status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err := alarm.SendAlarm(e.msg, mstype, e.err)
	if err != nil {
		log.Fatalf(e.msg)
	}
}

// logError obtains a trace of the line and file where the error happens
func logError(err error) string {
	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
}
