package canary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"leet/models"

	"github.com/rs/zerolog/log"
)

func RunCanary(serverDomain string) {

	// create
	request, err := http.NewRequest("POST",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(`{"Title": "test", "Completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	var task models.Task
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		panic(err)
	}
	if "test" != task.Title {
		panic(err)
	}

	// verify get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	var tasks []models.Task
	err = json.Unmarshal(buf.Bytes(), &tasks)
	if err != nil {
		panic(err)
	}
	foundTask := false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	if true != foundTask {
		panic(err)
	}

	// update
	// id := task.ID
	request, err = http.NewRequest("PUT",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(`{"title": "changedit", "completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		panic(err)
	}
	if "changedit" != task.Title {
		panic(err)
	}

	// delete
	request, err = http.NewRequest("DELETE",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}

	// verify no get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	response, err = c.Do(request)
	if err != nil {
		panic(err)
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if http.StatusOK != response.StatusCode {
		panic(err)
	}
	err = json.Unmarshal(buf.Bytes(), &tasks)
	if err != nil {
		panic(err)
	}
	foundTask = false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	if false != foundTask {
		panic(err)
	}
}
func main() {
	log.Info().Msg("Initializing canary..")

	var serverDomain = "https://tasks.dev.leetcyber.com"

	RunCanary(serverDomain)
	log.Info().Msg("Canary finished successfully..")

}
