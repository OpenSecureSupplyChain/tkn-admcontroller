package pipelines

import (
	"strings"

	"github.com/tkn-admcontroller/pkg/controller"

	v1 "k8s.io/api/admission/v1"
)

const signKeyAnnotation = "ossc.sigstore.tapestry.dev/transparency"

func validateCreate() controller.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*controller.Result, error) {
		pipeline, err := parsePipeline(r.Object.Raw)
		if err != nil {
			return &controller.Result{Message: err.Error()}, nil
		}

		annotations := pipeline.GetAnnotations()
		// fmt.Printf("annotations: %v\n", annotations)
		// for _, t := range pipeline.Spec.Tasks {
		// 	if strings.HasSuffix(t.Name, "app") {
		// 		return &controller.Result{Message: "You cannot use the tag `app` in a task name."}, nil
		// 	}
		// }

		annotationFound := false
		for key, val := range annotations {
			if strings.EqualFold(key, signKeyAnnotation) {
				if val != "" {
					annotationFound = true
				}
			}
		}

		if !annotationFound {
			return &controller.Result{Message: "sigstore sign annotation not found"}, nil
		}
		return &controller.Result{Allowed: true}, nil
	}
}
