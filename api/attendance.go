package api

import (
	"../models"
	"../util"
	"github.com/gorilla/mux"
	"net/http"
)

// GetAttendances -> get lists of attendances
func GetAttendances(w http.ResponseWriter, r *http.Request) {
	page, err := readInt(r, "page", 1)
	if err != nil || page < 1 {
		sendJSONMessage(w, http.StatusBadRequest, "Ooops!", err.Error())
	}
	pageSize, err := readInt(r, "page_size", 20)
	if err != nil || pageSize <= 0 {
		sendJSONMessage(w, http.StatusBadRequest, "Ooops!", err.Error())
	}
	offset := (page - 1) * pageSize

	order := r.FormValue("order")

	var attendances []*models.Attendance

	if order != "" {
		err = DB.Preload("Student").Preload("Schedule").Preload("Schedule.Lecturer").Order(order).Offset(offset).Limit(pageSize).Find(&attendances).Error
	} else {
		err = DB.Preload("Student").Preload("Schedule").Preload("Schedule.Lecturer").Offset(offset).Limit(pageSize).Find(&attendances).Error
	}

	if err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONData(w, attendances, http.StatusOK)
}

// AddAttendance -> add 1 schedule
func AddAttendance(w http.ResponseWriter, r *http.Request) {
	attendance := &models.Attendance{}

	if err := readJSON(r, attendance); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if validate := validateAttendance(w, attendance); validate == true {
		if err := DB.Save(attendance).Error; err != nil {
			sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
			return
		}

		sendJSONMessage(w, http.StatusCreated, "Great!", "attendance successfully added.")
	}
}

// GetAttendance by ID
func GetAttendance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	attendance := &models.Attendance{}
	if DB.Preload("Student").Preload("Schedule").Preload("Schedule.Lecturer").First(attendance, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	sendJSONData(w, attendance, http.StatusOK)
}

// UpdateAttendance by ID
func UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	attendance := &models.Attendance{}
	if DB.First(attendance, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}

	updated := &models.Attendance{}
	if err := readJSON(r, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := util.Copy(attendance, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := DB.Save(attendance).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "attendance successfully updated")
}

// DeleteAttendance by ID
func DeleteAttendance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	attendance := &models.Attendance{}

	if DB.First(attendance, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	if err := DB.Delete(attendance).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "attendance successfully deleted.")
}

func validateAttendance(w http.ResponseWriter, attendance *models.Attendance) bool {
	if DB.Where("student_id = ?", attendance.StudentID).Where("schedule_id = ?", attendance.ScheduleID).First(&attendance).Error == nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're trying to input is already stored. We can't save the same data!")
		return false
	}
	return true
}