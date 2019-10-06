package api

import (
	"../models"
	"../util"
	"github.com/gorilla/mux"
	"net/http"
)

// GetSchedules -> get lists of schedules
func GetSchedules(w http.ResponseWriter, r *http.Request) {
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

	var schedules []*models.Schedule

	if order != "" {
		err = DB.Preload("Lecturer").Order(order).Offset(offset).Limit(pageSize).Find(&schedules).Error
	} else {
		err = DB.Preload("Lecturer").Offset(offset).Limit(pageSize).Find(&schedules).Error
	}

	if err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONData(w, schedules, http.StatusOK)
}

// AddSchedule -> add 1 schedule
func AddSchedule(w http.ResponseWriter, r *http.Request) {
	schedule := &models.Schedule{}

	if err := readJSON(r, schedule); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if validate := validateSchedule(w, schedule); validate == true {
		if err := DB.Save(schedule).Error; err != nil {
			sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
			return
		}

		sendJSONMessage(w, http.StatusCreated, "Great!", "schedule " + schedule.CourseName + " successfully added.")
	}
}

// GetSchedule by ID
func GetSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	schedule := &models.Schedule{}
	if DB.Preload("Lecturer").First(schedule, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	sendJSONData(w, schedule, http.StatusOK)
}

// UpdateSchedule by ID
func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	schedule := &models.Schedule{}
	if DB.First(schedule, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}

	updated := &models.Schedule{}
	if err := readJSON(r, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := util.Copy(schedule, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := DB.Save(schedule).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "schedule updated to " + updated.CourseName)
}

// DeleteSchedule by ID
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	schedule := &models.Schedule{}

	if DB.First(schedule, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	if err := DB.Delete(schedule).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "schedule " + schedule.CourseName + " successfully deleted.")
}

func validateSchedule(w http.ResponseWriter, schedule *models.Schedule) bool {
	if DB.Where("time_start = ?", schedule.TimeStart).Where("time_end = ?", schedule.TimeEnd).Where("day = ?", schedule.Day).Where("class_name = ?", schedule.ClassName).First(&schedule).Error == nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're trying to input is already stored. We can't save the same data!")
		return false
	}
	return true
}