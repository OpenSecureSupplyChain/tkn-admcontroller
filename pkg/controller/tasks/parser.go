package tasks

import (
	"encoding/json"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// NewValidationHook creates a new instance of task validation hook
func NewValidationHook() controller.Hook {
	return controller.Hook{
		Create: validateTaskCreate(),
	}
}

func parseTask(object []byte) (*v1beta1.Task, error) {
	var task v1beta1.Task
	if err := json.Unmarshal(object, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
