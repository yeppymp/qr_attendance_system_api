package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"main.go/models"
	"main.go/util"
	"net/http"
	"strings"
)

// UserInput class
type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetUsers -> get lists of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
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

	var users []*models.User

	if order != "" {
		err = DB.Preload("Role").Order(order).Offset(offset).Limit(pageSize).Find(&users).Error
	} else {
		err = DB.Preload("Role").Offset(offset).Limit(pageSize).Find(&users).Error
	}

	if err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONData(w, users, http.StatusOK)
}

// AddUser function
func AddUser(w http.ResponseWriter, r *http.Request) {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}
	
	user := &models.User{}
	json.Unmarshal(body, user)
	validateUser(w, user, input.Password)

	if err := user.SetPassword(input.Password); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
	}

	if err := DB.Save(user).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusCreated, "Great!", "User " + user.Name + " successfully added.")
}

// GetUser by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	user := &models.User{}
	if DB.First(user, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	sendJSONData(w, user, http.StatusOK)
}

// UpdateUser by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	user := &models.User{}
	if DB.First(user, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}

	updated := &models.User{}
	if err := readJSON(r, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := util.Copy(user, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := DB.Save(user).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "User " + updated.Name + " updated.")
}

// DeleteUser by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	user := &models.User{}

	if DB.First(user, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	if err := DB.Delete(user).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "User " + user.Name + " successfully deleted.")
}

func validateUser(w http.ResponseWriter, user *models.User, password string) {
	if !strings.Contains(user.Email, "@") {
		sendJSONMessage(w, http.StatusBadRequest, "Oops", "Please enter a valid email!")
	}
	
	if password == "" {
		sendJSONMessage(w, http.StatusBadRequest, "Oops", "Please enter a password!")
	}
}