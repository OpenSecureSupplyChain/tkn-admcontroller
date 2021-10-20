package keys

import (
	"crypto"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
)

//GetPublicKey :
func GetPublicKey(pubKeyfp string) (crypto.PublicKey, error) {
	// Retrieve the cosign pub key, this must be loaded as a secret
	// kubectl create secret generic cosign-pub --from-file=./cosign.pub
	cosignPub, err := ioutil.ReadFile("/etc/cosign/cosign.pub")
	if err != nil {
		return nil, errors.Wrapf(err, "error reading public key")
	}
	return cryptoutils.UnmarshalPEMToPublicKey(cosignPub)
}

//TODO: Add functions to fetch keys from KMS
