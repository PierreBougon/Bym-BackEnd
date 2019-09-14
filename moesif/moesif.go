package moesif

import (
	"github.com/moesif/moesifmiddleware-go"
	"net/http"
	"os"
)

func moesifOptions() map[string]interface{}{
	appId := os.Getenv("moesif_app_id")

	if appId == "" {
		return nil
	}
	return map[string]interface{} {
		"Application_Id": appId,
		"Log_Body": true,
	}
}

func MiddlewareWrapper(h http.Handler) http.Handler {
	options := moesifOptions()
	if options == nil {
		return h
	}
	return moesifmiddleware.MoesifMiddleware(h, options)
}