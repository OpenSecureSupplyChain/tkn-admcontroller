package pipelines

import (
	"strings"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"

	v1 "k8s.io/api/admission/v1"
)

const messageAnnotation = "cosign.sigstore.dev/message"
const signatureAnnotation = "cosign.sigstore.dev/signature"

func validateCreate() controller.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*controller.Result, error) {
		pipeline, err := parsePipeline(r.Object.Raw)
		if err != nil {
			return &controller.Result{Message: err.Error()}, nil
		}

		annotations := pipeline.GetAnnotations()

		messageAnnotationFound := false
		signatureAnnotationFound := false
		for key, val := range annotations {
			if strings.EqualFold(key, signatureAnnotation) {
				if val == "" {
					signatureAnnotationFound = true
				}
			}
			if strings.EqualFold(key, messageAnnotation) {
				if val != "" {
					messageAnnotationFound = true
				}
			}
		}

		if !signatureAnnotationFound && !messageAnnotationFound{
			return &controller.Result{Message: "signature or message annotation not found"}, nil
		}

		// TODO: verify signature logic

		return &controller.Result{Allowed: true}, nil
	}
}


