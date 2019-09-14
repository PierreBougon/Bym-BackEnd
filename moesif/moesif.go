package moesif

import (
	"github.com/moesif/moesifmiddleware-go"
	"net/http"
	"os"
)


var options = fetchMoesifOptions()

func fetchMoesifOptions() map[string]interface{}{
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
	if options == nil {
		return h
	}
	return moesifmiddleware.MoesifMiddleware(h, options)
}

func CatchOutgoingCalls() {
	if options != nil {
		moesifmiddleware.StartCaptureOutgoing(options)
	}
}