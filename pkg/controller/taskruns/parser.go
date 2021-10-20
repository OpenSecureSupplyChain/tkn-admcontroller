package taskruns

import (
	"encoding/json"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// NewValidationHook creates a new instance of taskrun validation hook
func NewValidationHook() controller.Hook {
	return controller.Hook{
		Create: validateTaskrunCreate(),
	}
}

func parseTaskrun(object []byte) (*v1beta1.TaskRun, error) {
	var taskRun v1beta1.TaskRun
	if err := json.Unmarshal(object, &taskRun); err != nil {
		return nil, err
	}
	return &taskRun, nil
}
