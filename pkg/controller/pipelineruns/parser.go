package pipelineruns

import (
	"encoding/json"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// NewValidationHook creates a new instance of pipelinerun validation hook
func NewValidationHook() controller.Hook {
	return controller.Hook{
		Create: validatePipelinerunCreate(),
	}
}

func parsePipelinerun(object []byte) (*v1beta1.PipelineRun, error) {
	var pipelinerun v1beta1.PipelineRun
	if err := json.Unmarshal(object, &pipelinerun); err != nil {
		return nil, err
	}
	return &pipelinerun, nil
}
