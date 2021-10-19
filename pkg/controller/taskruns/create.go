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

		fmt.Printf("need to verify taskref: %s\n", taskrun.Spec.TaskRef.Name)
		//TODO: retrieve referenced task definition and verify if it signed
		// also, verify all step images if they are signed
		return &controller.Result{Allowed: true}, nil
	}
}
