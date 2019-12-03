package controllers

import (
	"github.com/PierreBougon/Bym-BackEnd/app/models"
	u "github.com/PierreBougon/Bym-BackEnd/app/utils"

	"encoding/json"
	"net/http"
)

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondBadRequest(w)
		return
	}

	resp := models.Login(account.Email, account.Password)
	if resp["status"] == false {
		u.RespondUnauthorized(w)
	} else {
		u.Respond(w, resp)
	}
}
