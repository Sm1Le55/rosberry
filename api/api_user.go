/*
 * Rosberry Mobile APP API (Test work)
 *
 * Test work for rosberry
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rosberry/database"
	"rosberry/model"
	"strconv"

	"github.com/gorilla/mux"
)

func RegistrationUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)
	var data model.UserRegData
	err := decoder.Decode(&data)
	fmt.Printf("request: %v\n", data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}", err)))
		return
	}

	err = database.Registration(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)
	var data model.UserLoginData
	err := decoder.Decode(&data)
	fmt.Printf("request: %v\n", data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}", err)))
		return
	}

	info, err := database.Login(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	//clear access key in db
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	email := r.FormValue("email")

	err := database.Logout(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetDisplaySettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	userIDString := mux.Vars(r)["userId"]
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings, err := database.DisplaySettingsQuery(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(settings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func SetDisplaySettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)
	var data model.DisplaySettings
	err := decoder.Decode(&data)
	fmt.Printf("request: %v\n", data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}", err)))
		return
	}

	err = database.SaveDisplaySettings(&data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
