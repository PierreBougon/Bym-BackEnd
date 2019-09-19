package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

var UpdateAccount = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	var account = &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	account.ID = user
	resp := account.UpdateAccount()
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var UpdatePassword = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}
	resp := models.UpdatePassword(user, account.Password)
	u.Respond(w, resp)
}
