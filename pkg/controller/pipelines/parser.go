package pipelines

import (
	"encoding/json"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tkn-admcontroller/pkg/controller"
)

// NewValidationHook creates a new instance of pods validation hook
func NewValidationHook() controller.Hook {
	return controller.Hook{
		Create: validateCreate(),
	}
}

func parsePipeline(object []byte) (*v1beta1.Pipeline, error) {
	var pipeline v1beta1.Pipeline
	if err := json.Unmarshal(object, &pipeline); err != nil {
		return nil, err
	}
	return &pipeline, nil
}
