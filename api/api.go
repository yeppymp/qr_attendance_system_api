package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

// DB instance
var DB *gorm.DB

// Message JSON struct
type Message struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

// Init the router
func Init(router *mux.Router) {

	// ROLES
	router.HandleFunc("/roles", GetRoles).Methods("GET")
	router.HandleFunc("/roles", AddRole).Methods("POST")
	router.HandleFunc("/role/{id}", GetRole).Methods("GET")
	router.HandleFunc("/role/{id}", UpdateRole).Methods("PUT")
	router.HandleFunc("/role/{id}", DeleteRole).Methods("DELETE")

	// USERS
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", AddUser).Methods("POST")
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")

	// SCHEDULES
	router.HandleFunc("/schedules", GetSchedules).Methods("GET")
	router.HandleFunc("/schedules", AddSchedule).Methods("POST")
	router.HandleFunc("/schedule/{id}", GetSchedule).Methods("GET")
	router.HandleFunc("/schedule/{id}", UpdateSchedule).Methods("PUT")
	router.HandleFunc("/schedule/{id}", DeleteSchedule).Methods("DELETE")

	// ATTENDANCES
	router.HandleFunc("/attendances", GetAttendances).Methods("GET")
	router.HandleFunc("/attendances", AddAttendance).Methods("POST")
	router.HandleFunc("/attendance/{id}", GetAttendance).Methods("GET")
	router.HandleFunc("/attendance/{id}", UpdateAttendance).Methods("PUT")
	router.HandleFunc("/attendance/{id}", DeleteAttendance).Methods("DELETE")

	router.NotFoundHandler = http.HandlerFunc(notFound)
}

func sendJSONMessage(w http.ResponseWriter, status int, title string, msg string) {
	data := Message{title, msg}
	ms, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(ms)
}

func sendJSONData(w http.ResponseWriter, d interface{}, status int) {
	data, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(data)
}

func readInt(r *http.Request, param string, v int64) (int64, error) {
	p := r.FormValue(param)
	if p == "" {
		return v, nil
	}
	return strconv.ParseInt(p, 10, 64)
}

func readJSON(r *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, v)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	message := Message{"I'm sorry", "URI or pages that you're trying to access is not found."}
	ms, _ := json.Marshal(message)
	w.WriteHeader(http.StatusNotFound)
	w.Write(ms)
}