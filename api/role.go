package api

import (
	"github.com/gorilla/mux"
	"main.go/models"
	"main.go/util"
	"net/http"
)

// GetRoles -> get lists of roles
func GetRoles(w http.ResponseWriter, r *http.Request) {
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

	var roles []*models.Role

	if order != "" {
		err = DB.Model(&models.Role{}).Order(order).Offset(offset).Limit(pageSize).Find(&roles).Error
	} else {
		err = DB.Model(&models.Role{}).Offset(offset).Limit(pageSize).Find(&roles).Error
	}

	if err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONData(w, roles, http.StatusOK)
}

// AddRole -> add 1 role
func AddRole(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}

	if err := readJSON(r, role); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}
	if err := DB.Save(role).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}
	
	sendJSONMessage(w, http.StatusCreated, "Great!", "Role " + role.RoleName + " successfully added.")
}

// GetRole by ID
func GetRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	role := &models.Role{}
	if DB.First(role, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	sendJSONData(w, role, http.StatusOK)
}

// UpdateRole by ID
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	role := &models.Role{}
	if DB.First(role, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}

	updated := &models.Role{}
	if err := readJSON(r, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := util.Copy(role, updated); err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	if err := DB.Save(role).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "Role updated to " + updated.RoleName)
}

// DeleteRole by ID
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	role := &models.Role{}

	if DB.First(role, id).Error != nil {
		sendJSONMessage(w, http.StatusNotFound, "Oops", "Data that you're looking for is not found.")
		return
	}
	if err := DB.Delete(role).Error; err != nil {
		sendJSONMessage(w, http.StatusInternalServerError, "Ooops!", err.Error())
		return
	}

	sendJSONMessage(w, http.StatusOK, "Great!", "Role " + role.RoleName + " successfully deleted.")
}