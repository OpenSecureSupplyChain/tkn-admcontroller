package pipelineruns

import (
	"fmt"
	"strings"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/validator"
	v1 "k8s.io/api/admission/v1"
)

func validatePipelinerunCreate() controller.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*controller.Result, error) {
		pipelinerun, err := parsePipelinerun(r.Object.Raw)
		if err != nil {
			return &controller.Result{Message: err.Error()}, nil
		}

		//handle the case, when pipelineruns are instantiated from
		//webhook -> triggers -> pipelineruns
		//currently, we only verify for manually created pipelineruns
		prAnnotations := pipelinerun.GetAnnotations()
		messageAnnotationFound := false
		signatureAnnotationFound := false

		var sig, msg string
		for key, val := range prAnnotations {
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
			return &controller.Result{Message: "signature or message annotation not found for pipelinerun"}, nil
		}

		isVerified, err := validator.ValidateYAMLObject(msg, sig)
		if err != nil || !isVerified {
			return &controller.Result{Message: "Signature validation failed for pipelinerun"}, nil
		}
		fmt.Printf("need to verify pipelineref: %s\n", pipelinerun.Spec.PipelineRef.Name)
		//TODO: retrieve referenced pipeline definition and verify if it signed
		// also, verify all tasks and all step images if they are signed

		// //validate all tasks and task->step images
		// failedTasks := []string{}
		// failedTaskImgs := []string{}
		// for _, pt := range pipelinerun.Spec.PipelineSpec.Tasks {
		// 	tAnnotations := pt.TaskSpecMetadata().Annotations
		// 	for key, val := range tAnnotations {
		// 		if strings.EqualFold(key, controller.SignatureAnnotation) {
		// 			if val != "" {
		// 				signatureAnnotationFound = true
		// 				sig = val
		// 			}
		// 		}
		// 		if strings.EqualFold(key, controller.MessageAnnotation) {
		// 			if val != "" {
		// 				messageAnnotationFound = true
		// 				msg = val
		// 			}
		// 		}
		// 	}

		// 	if !signatureAnnotationFound && !messageAnnotationFound {
		// 		failedTasks = append(failedTasks, pt.Name)
		// 	}

		// 	isVerified, err := validator.ValidateYAMLObject(msg, sig)
		// 	if err != nil || !isVerified {
		// 		failedTasks = append(failedTasks, pt.Name)
		// 	}

		// 	//Validate all step images
		// 	for _, step := range pt.TaskSpec.Steps {
		// 		stepImg := step.Image
		// 		isVerified, err := validator.ValidateImage(stepImg, "")
		// 		if err != nil || !isVerified {
		// 			failedTaskImgs = append(failedTaskImgs, stepImg)
		// 		}
		// 	}
		// }
		// if len(failedTasks) != 0 {
		// 	failedMsg := fmt.Sprintf("Signature validation failed for task specs: [%s]", strings.Join(failedTasks, ","))
		// 	return &controller.Result{Message: failedMsg}, nil
		// }
		// if len(failedTaskImgs) != 0 {
		// 	failedMsg := fmt.Sprintf("Signature validation failed for images: [%s]", strings.Join(failedTaskImgs, ","))
		// 	return &controller.Result{Message: failedMsg}, nil
		// }
		return &controller.Result{Allowed: true}, nil
	}
}
