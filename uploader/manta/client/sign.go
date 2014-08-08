package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"time"
)

// Client is a Manta client. Client is not safe for concurrent use.
type Client struct {
	User   string
	KeyId  string
	Key    string
	Url    string
	signer Signer
}

func ensureHomedir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("manta: could not determine home directory: %v", err)
	}
	return user.HomeDir
}

// SignRequest signs the 'date' field of req.
func (c *Client) SignRequest(req *http.Request) error {
	if c.signer == nil {
		var err error
		c.signer, err = loadPrivateKey(c.Key)
		if err != nil {
			return fmt.Errorf("could not load private key %q: %v", c.Key, err)
		}
	}
	return signRequest(req, fmt.Sprintf("/%s/keys/%s", MANTA_USER, MANTA_KEY_ID), c.signer)
}

func signRequest(req *http.Request, keyid string, priv Signer) error {
	now := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", now)
	signed, err := priv.Sign([]byte(fmt.Sprintf("date: %s", now)))
	if err != nil {
		return fmt.Errorf("could not sign request: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(signed)
	authz := fmt.Sprintf("Signature keyId=%q,algorithm=%q,signature=%q", keyid, "rsa-sha256", sig)
	req.Header.Set("Authorization", authz)
	return nil
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}
	rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("An error occurred while parsing the key: %s", err)
	}
	return newSignerFromKey(rsa)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPublicKey rsa.PublicKey

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
