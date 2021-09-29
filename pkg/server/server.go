package server

import (
	"fmt"
	"net/http"

	"github.com/tkn-admcontroller/pkg/controller/pipelines"
	"github.com/tkn-admcontroller/pkg/handlers"
)

//NewServer :
func NewServer(port string) *http.Server {
	// Instances hooks

	pipelineValidation := pipelines.NewValidationHook()

	// Routers
	ah := handlers.NewAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", handlers.Healthz())
	mux.Handle("/validate/pipelines", ah.Serve(pipelineValidation))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
