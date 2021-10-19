package pipelines

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/ecdsa"
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/controller"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"github.com/sigstore/sigstore/pkg/signature"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/klog"
)

const messageAnnotation = "cosign.sigstore.dev/message"
const signatureAnnotation = "cosign.sigstore.dev/signature"

func validateCreate() controller.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*controller.Result, error) {

		// Retrieve the cosign pub key, this must be loaded as a secret
		// kubectl create secret generic cosign-pub --from-file=./cosign.pub
		cosignPub, err := ioutil.ReadFile("/etc/cosign/cosign.pub")
		if err != nil {
			return &controller.Result{Message: "error reading public key"}, nil
		}

		pipeline, err := parsePipeline(r.Object.Raw)
		if err != nil {
			return &controller.Result{Message: err.Error()}, nil
		}

		annotations := pipeline.GetAnnotations()

		messageAnnotationFound := false
		signatureAnnotationFound := false

		var sig, msg string
		for key, val := range annotations {
			if strings.EqualFold(key, signatureAnnotation) {
				if val != "" {
					signatureAnnotationFound = true
					sig = val
				}
			}
			if strings.EqualFold(key, messageAnnotation) {
				if val != "" {
					messageAnnotationFound = true
					msg = val
				}
			}
		}

		if !signatureAnnotationFound && !messageAnnotationFound {
			return &controller.Result{Message: "signature or message annotation not found"}, nil
		}

		gzipMsg, _ := base64.StdEncoding.DecodeString(msg)
		rawMsg := gzipDecompress(gzipMsg)
		rawSig, _ := base64.StdEncoding.DecodeString(sig)

		publicKey, err := cryptoutils.UnmarshalPEMToPublicKey(cosignPub)

		verifier, err := signature.LoadECDSAVerifier(publicKey.(*ecdsa.PublicKey), crypto.SHA256)
		if err != nil {
			klog.Error(err)
		}
		err = verifier.VerifySignature(bytes.NewReader(rawSig), bytes.NewReader(rawMsg))

		if err != nil {
			return &controller.Result{Message: "Signature validation failed"}, nil
		}
		return &controller.Result{Allowed: true}, nil
	}
}

func gzipDecompress(in []byte) []byte {
	reader := bytes.NewReader(in)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return in
	}
	output, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return in
	}
	return output
}
