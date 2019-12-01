package moesif

import (
	"fmt"
	"github.com/moesif/moesifmiddleware-go"
	"net/http"
	"os"
)

func MiddlewareWrapper(h http.Handler) http.Handler {
	options := fetchMoesifOptions()
	if options == nil {
		return h
	}

	return moesifmiddleware.MoesifMiddleware(h, options)
}

func fetchMoesifOptions() map[string]interface{} {
	appId, exist := os.LookupEnv("moesif_app_id")

	if !exist || appId == "" {
		fmt.Println("Moesif middleware is not used. appid does not exist or is empty.")
		return nil
	}

	return map[string]interface{}{
		"Application_Id":           appId,
		"Log_Body":                 true,
		"Capture_Outoing_Requests": false,
		"Identify_User": func(r *http.Request, recorder moesifmiddleware.MoesifResponseRecorder) string {
			if user := r.Context().Value("user"); user != nil {
				return fmt.Sprint(user)
			}
			return "unauthenticated"
		},
	}
}
