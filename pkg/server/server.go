package server

import (
	"fmt"
	"net/http"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller/pipelineruns"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller/pipelines"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller/taskruns"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller/tasks"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/handlers"
)

//NewServer :
func NewServer(port string) *http.Server {
	// Instances hooks

	pipelineValidation := pipelines.NewValidationHook()
	taskValidation := tasks.NewValidationHook()
	taskrunsValidation := taskruns.NewValidationHook()
	pipelinerunsValidation := pipelineruns.NewValidationHook()

	// Routers
	ah := handlers.NewAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", handlers.Healthz())
	mux.Handle("/validate/pipelines", ah.Serve(pipelineValidation))
	mux.Handle("/validate/pipelineruns", ah.Serve(pipelinerunsValidation))
	mux.Handle("/validate/tasks", ah.Serve(taskValidation))
	mux.Handle("/validate/taskruns", ah.Serve(taskrunsValidation))

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
