package validator

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"encoding/base64"

	"github.com/OpenSecureSupplyChain/tkn-admcontroller/pkg/validator/keys"
	"github.com/pkg/errors"
	"github.com/sigstore/sigstore/pkg/signature"
	"k8s.io/klog"
)

//ValidateYAMLObject : validates yaml object
func ValidateYAMLObject(msg, sig string) (bool, error) {
	pubKeyfp := ""
	publicKey, err := keys.GetPublicKey(pubKeyfp)
	if err != nil {
		return false, err
	}
	gzipMsg, _ := base64.StdEncoding.DecodeString(msg)
	rawMsg := gzipDecompress(gzipMsg)
	rawSig, _ := base64.StdEncoding.DecodeString(sig)
	verifier, err := signature.LoadECDSAVerifier(publicKey.(*ecdsa.PublicKey), crypto.SHA256)
	if err != nil {
		klog.Error(err)
	}
	err = verifier.VerifySignature(bytes.NewReader(rawSig), bytes.NewReader(rawMsg))
	if err != nil {
		return false, errors.Wrapf(err, "error verifying signature")
	}
	return true, nil
}
