package taskruns

import (
	"fmt"
	"strings"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/validator"
	v1 "k8s.io/api/admission/v1"
)

func validateTaskrunCreate() controller.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*controller.Result, error) {
		taskrun, err := parseTaskrun(r.Object.Raw)
		if err != nil {
			return &controller.Result{Message: err.Error()}, nil
		}

		//handle the case, when taskruns are instantiated from
		//webhook -> triggers -> pipelineruns
		//currently, we only verify for manually created taskruns
		annotations := taskrun.GetAnnotations()
		messageAnnotationFound := false
		signatureAnnotationFound := false

		var sig, msg string
		for key, val := range annotations {
			if strings.EqualFold(key, controller.SignatureAnnotation) {
				if val != "" {
					signatureAnnotationFound = true
					sig = val
				}
			}
			if strings.EqualFold(key, controller.MessageAnnotation) {
				if val != "" {
					messageAnnotationFound = true
					msg = val
				}
			}
		}

		if !signatureAnnotationFound && !messageAnnotationFound {
			return &controller.Result{Message: "signature or message annotation not found"}, nil
		}

		isVerified, err := validator.ValidateYAMLObject(msg, sig)
		if err != nil || !isVerified {
			return &controller.Result{Message: "Signature validation failed"}, nil
		}

		//Validate all step images
		failedImgs := []string{}
		for _, step := range taskrun.Spec.TaskSpec.Steps {
			stepImg := step.Image
			isVerified, err := validator.ValidateImage(stepImg, "")
			if err != nil || !isVerified {
				failedImgs = append(failedImgs, stepImg)
			}
		}
		if len(failedImgs) != 0 {
			failedMsg := fmt.Sprintf("Signature validation failed for images: [%s]", strings.Join(failedImgs, ","))
			return &controller.Result{Message: failedMsg}, nil
		}

		return &controller.Result{Allowed: true}, nil
	}
}
