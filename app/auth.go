package app

import (
	"github.com/PierreBougon/Bym-BackEnd/models"
	u "github.com/PierreBougon/Bym-BackEnd/utils"

	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sendErrorJson = func(w http.ResponseWriter, errMsg string, httpStatus int) {
			response := make(map[string]interface{})

			response = u.Message(false, errMsg)
			w.WriteHeader(httpStatus)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
		}

		notAuth := []string{"", "/api", "/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                           //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		if tokenHeader == "" {                       //Token is missing, returns with error code 403 Unauthorized
			sendErrorJson(w, "Missing auth token", http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			sendErrorJson(w, "Invalid/Malformed auth token", http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			sendErrorJson(w, "Malformed authentication token", http.StatusForbidden)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			sendErrorJson(w, "Token is not valid.", http.StatusForbidden)
			return
		}

		account := models.GetUser(tk.UserId)
		if account == nil {
			sendErrorJson(w, "Account does not exist", http.StatusForbidden)
			return
		}
		if tk.TokenVersion != account.TokenVersion {
			sendErrorJson(w, "Token is not valid anymore.", http.StatusForbidden)
			return
		}
		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		//fmt.Sprintf("User %", tk.Username) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
