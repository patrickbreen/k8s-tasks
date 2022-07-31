package tasks

import (
	"encoding/json"
	"fmt"
	"leet/models"
	"leet/util"
	"net/http"
	"time"

	"github.com/golang/gddo/httputil/header"
	"github.com/rs/zerolog/log"
)

func myDeserialize(w http.ResponseWriter, req *http.Request) (*json.Decoder, int) {
	if req.Header.Get("Content-Type") == "" {
		return nil, http.StatusUnsupportedMediaType
	}
	if value, _ := header.ParseValueAndParams(req.Header, "Content-Type"); value != "application/json" {
		return nil, http.StatusUnsupportedMediaType
	}
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec, 0
}

func TasksHandler(w http.ResponseWriter, req *http.Request) {
	claims, err := util.ValidateAndGetClaims(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	log.Info().Msg(fmt.Sprintf("user claims: %v\n", claims))
	if req.Method == "GET" {
		GetTasks(w, req)
	} else if req.Method == "POST" {
		CreateTask(w, req)
	} else if req.Method == "PUT" {
		UpdateTask(w, req)
	} else if req.Method == "DELETE" {
		DeleteTask(w, req)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetTasks(w http.ResponseWriter, req *http.Request) {
	var tasks []models.Task
	db := util.GetDB()
	db.Find(&tasks)
	b, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func CreateTask(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	var db = util.GetDB()

	// check deserialize new thing
	dec, status_code := myDeserialize(w, req)
	if status_code != 0 {
		w.WriteHeader(status_code)
		return
	}
	if err := dec.Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// do thing
	db.Create(&task)
	b, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func UpdateTask(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	db := util.GetDB()

	// check thing to update
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := ids[0]
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check deserialize new thing
	dec, status_code := myDeserialize(w, req)
	if status_code != 0 {
		w.WriteHeader(status_code)
		return
	}
	if err := dec.Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// do thing
	task.UpdatedAt = time.Now()
	db.Save(&task)
	// write back new object as json 200
	b, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func DeleteTask(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	db := util.GetDB()

	// check thing to update
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := ids[0]
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// do thing
	db.Delete(&task)
}
