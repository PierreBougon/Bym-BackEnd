package controllers

import (
	"encoding/json"
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"
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

var DeleteAccount = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	resp := (&models.Account{}).DeleteAccount(user)
	if resp["status"] == false {
		w.WriteHeader(http.StatusBadRequest)
	}
	u.Respond(w, resp)
}

var GetAccount = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	acc := models.GetUser(user) //Create account
	if acc == nil {
		u.RespondBadRequestWithMessage(w, "Error cannot find account")
		return
	}
	resp := u.Message(true, "success")
	resp["account"] = acc
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
