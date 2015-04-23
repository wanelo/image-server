package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/golang/glog"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func homeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

// SignRequest signs the 'date' field of req.
func (c *Client) SignRequest(req *http.Request) error {
	start := time.Now()
	if c.signer == nil {
		var err error
		c.signer, err = c.getSigner(SDC_IDENTITY)
		if err != nil {
			return err
		}
	}
	err := signRequest(req, fmt.Sprintf("/%s/keys/%s", MANTA_USER, MANTA_KEY_ID), c.signer)
	elapsed := time.Since(start)
	glog.Infof("Took %s to sign request", elapsed)
	return err
}

func signRequest(req *http.Request, keyid string, signer ssh.Signer) error {
	now := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", now)
	signed, err := signer.Sign(rand.Reader, []byte(fmt.Sprintf("date: %s", now)))
	if err != nil {
		return fmt.Errorf("could not sign request: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(signed.Blob)
	authz := fmt.Sprintf("Signature keyId=%q,algorithm=%q,signature=%q", keyid, "rsa-sha1", sig)
	req.Header.Set("Authorization", authz)
	return nil
}

func (c *Client) getSigner(keyPath string) (signer ssh.Signer, err error) {
	if keySigner != nil {
		return keySigner, nil
	}

	// The keySigner will be nil if it was impossible to load the key, for reasons
	// like the key needing a passphrase. In that case, we ask the SSH Agent for
	// the key, just in case it has it.
	return c.getSignerFromSSHAgent(keyPath)
}

func getSignerFromPrivateKey(keyPath string) (signer ssh.Signer, err error) {
	var encodedKey []byte
	encodedKey, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	signer, err = ssh.ParsePrivateKey(encodedKey)
	return signer, err
}

func (c *Client) getSignerFromSSHAgent(keyPath string) (ssh.Signer, error) {
	if c.agentConn == nil {
		var err error
		c.agentConn, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			log.Fatal(err)
		}
	}
	ag := agent.NewClient(c.agentConn)

	signers, err := ag.Signers()
	if err != nil {
		log.Panic(err)
	}

	for _, signer := range signers {
		if signer.PublicKey().(*agent.Key).Comment == keyPath {
			return signer, nil
		}
	}
	return nil, fmt.Errorf("Signer was not found")
}
