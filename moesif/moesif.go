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

func fetchMoesifOptions() map[string]interface{}{
	appId := os.Getenv("moesif_app_id")

	if appId == "" {
		return nil
	}

	return map[string]interface{} {
		"Application_Id": appId,
		"Log_Body": true,
		"Capture_Outoing_Requests": true,
		"Identify_User": func (r *http.Request, recorder moesifmiddleware.MoesifResponseRecorder) string {
			return fmt.Sprint(r.Context().Value("user"))
		},
	}
}
