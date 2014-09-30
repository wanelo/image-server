package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.google.com/p/go.crypto/ssh"
)

var MANTA_URL, MANTA_USER, MANTA_KEY_ID, SDC_IDENTITY string
var keySigner ssh.Signer

// Client is a Manta client. Client is not safe for concurrent use.
type Client struct {
	User        string
	KeyId       string
	Key         string
	Url         string
	HTTPTimeout time.Duration
	signer      ssh.Signer
	agentConn   io.ReadWriter
}

// Entry represents an object stored in Manta, either a file or a directory
type Entry struct {
	Name       string `json:"name"`           // Entry name
	Etag       string `json:"etag,omitempty"` // If type is 'object', object UUID
	Size       int    `json:"size,omitempty"` // If type is 'object', object size (content-length)
	Type       string `json:"type"`           // Entry type, one of 'directory' or 'object'
	Mtime      string `json:"mtime"`          // ISO8601 timestamp of the last update
	Durability int    `json:"durability"`
}

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
		User:        MANTA_USER,
		KeyId:       MANTA_KEY_ID,
		Key:         SDC_IDENTITY,
		Url:         MANTA_URL,
		HTTPTimeout: 5 * time.Second,
	}
}

// PutObject creates or overwrites an object
// https://apidocs.joyent.com/manta/api.html#PutObject
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

// GetObject retrieves an object
// https://apidocs.joyent.com/manta/api.html#GetObject
func (c *Client) GetObject(path string) (io.Reader, error) {
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	// FIXME: Must close resp.Body, and we're not currently.
	// defer resp.Body.Close()

	err = c.ensureStatus(resp, 200)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// PutDirectory in the Joyent Manta Storage Service is an idempotent create-or-update operation
// https://apidocs.joyent.com/manta/api.html#PutDirectory
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

// DeleteDirectory deletes a directory. The directory must be empty. There is no response data from this request. On success an HTTP 204 is returned
func (c *Client) DeleteDirectory(path string) error {
	resp, err := c.Delete(path)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.ensureStatus(resp, 204)
}

// ListDirectory lists the contents of a directory
// https://apidocs.joyent.com/manta/api.html#ListDirectory
func (c *Client) ListDirectory(path string) ([]Entry, error) {
	headers := make(http.Header)
	headers.Add("content-type", "application/json; type=directory")

	resp, err := c.Get(path, headers)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var entries []Entry

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		entry := new(Entry)
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), entry)

		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, c.ensureStatus(resp, 200)
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

// Delete executes a DELETE request and returns the response.
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Do("DELETE", path, nil, nil)
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

	tr := &http.Transport{ResponseHeaderTimeout: c.HTTPTimeout}
	client := &http.Client{Transport: tr}
	return client.Do(req)
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
