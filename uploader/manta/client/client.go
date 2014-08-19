package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

var MANTA_URL, MANTA_USER, MANTA_KEY_ID, SDC_IDENTITY string
var keySigner ssh.Signer

// manta_url := os.Getenv("MANTA_URL")
func init() {
	MANTA_URL = os.Getenv("MANTA_URL")
	MANTA_USER = os.Getenv("MANTA_USER")
	MANTA_KEY_ID = os.Getenv("MANTA_KEY_ID")
	SDC_IDENTITY = os.Getenv("SDC_IDENTITY")
	if SDC_IDENTITY == "" {
		homeDir, err := homeDir()
		if err == nil {
			SDC_IDENTITY = filepath.Join(homeDir, ".ssh", "id_rsa")
		}
	}
	if requiresSSL() {
		keySigner, _ = getSignerFromPrivateKey(SDC_IDENTITY)
	}
}

// DefaultClient returns a Client instance configured from the
// default Manta environment variables.
func DefaultClient() *Client {
	return &Client{
		User:  MANTA_USER,
		KeyId: MANTA_KEY_ID,
		Key:   SDC_IDENTITY,
		Url:   MANTA_URL,
	}
}

func (c *Client) PutObject(destination string, contentType string, r io.Reader) error {
	headers := make(http.Header)
	headers.Add("content-type", "application/json")

	log.Println("filepath:", destination)
	resp, err := c.Put(destination, headers, r)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.ensureStatus(resp, 204)
}

func (c *Client) GetObject(path string) (io.Reader, error) {
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	err = c.ensureStatus(resp, 200)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) PutDirectory(path string) error {
	headers := make(http.Header)
	headers.Add("content-type", "application/json; type=directory")

	resp, err := c.Put(path, headers, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.ensureStatus(resp, 204)
}

func (c *Client) ensureStatus(resp *http.Response, expected int) error {
	if resp.StatusCode != expected {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return (err)
		}
		return fmt.Errorf("HTTP Response for %v was not %v got %v: %s", resp.Request, expected, resp.StatusCode, body)
	}

	return nil
}

// Get executes a GET request and returns the response.
func (c *Client) Get(path string, headers http.Header) (*http.Response, error) {
	return c.Do("GET", path, headers, nil)
}

// Put executes a PUT request and returns the response.
func (c *Client) Put(path string, headers http.Header, r io.Reader) (*http.Response, error) {
	return c.Do("PUT", path, headers, r)
}

// Post executes a POST request and returns the response.
func (c *Client) Post(path string, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.Do("POST", path, headers, body)
}

// Do executes a method request and returns the response.
func (c *Client) Do(method, path string, headers http.Header, r io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(method, path, r)
	if err != nil {
		return nil, err
	}
	req.Close = true

	for header, values := range headers {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	if requiresSSL() {
		if err := c.SignRequest(req); err != nil {
			return nil, err
		}
	}

	return http.DefaultClient.Do(req)
}

// NewRequest is similar to http.NewRequest except it appends path to
// the API endpoint this client is configured for.
func (c *Client) NewRequest(method, path string, r io.Reader) (*http.Request, error) {
	var url string
	if strings.HasPrefix(path, "/"+c.User) {
		url = fmt.Sprintf("%s%s", c.Url, path)
	} else {
		url = fmt.Sprintf("%s/%s/%s", c.Url, c.User, path)
	}
	log.Println("Making request to manta:", url)
	return http.NewRequest(method, url, r)
}

func requiresSSL() bool {
	// No need to use https if inside manta
	return MANTA_URL != "http://localhost:80/"
}
