package firebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Reference struct {
	Url     string
	Postfix string
	Client  *http.Client
}

// Retrieve a new Firebase reference for given url.
func NewReference(url string) *Reference {

	// Initialize Reference struct.
	r := &Reference{
		Url:     url,
		Postfix: ".json",
		Client:  &http.Client{},
	}

	return r
}

// Execute a new HTTP Request.
func (r *Reference) executeRequest(method string, body io.Reader) ([]byte, error) {

	// Prepare HTTP Request.
	req, err := http.NewRequest(method, r.Url+r.Postfix, nil)
	if err != nil {
		return nil, err
	}

	// Make actual HTTP request.
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// Make sure to close Body reader eventually.
	defer resp.Body.Close()

	// Check status code for errors.
	status := resp.Header.Get("Status Code")
	if strings.HasPrefix(status, "2") == false {
		return nil, errors.New(status)
	}

	// Read body.
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// Retrieve the current value for this Reference.
func (r *Reference) Value(v interface{}) error {

	// GET the data from Firebase.
	resp, err := r.executeRequest("GET", nil)
	if err != nil {
		return err
	}

	// JSON decode the data into given interface.
	err = json.Unmarshal(resp, v)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Write(v interface{}) error {

	jsonData, err := json.Marshal(t.Mockup)
	if err != nil {
		return err
	}

	_, err := r.newRequest("PUT", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Push(v interface{}) error {

	jsonData, err := json.Marshal(t.Mockup)
	if err != nil {
		return err
	}

	_, err := r.newRequest("POST", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Update(v interface{}) error {

	jsonData, err := json.Marshal(t.Mockup)
	if err != nil {
		return err
	}

	_, err := r.newRequest("PATCH", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Delete() error {

	_, err := r.newRequest("DELETE", nil)
	if err != nil {
		return err
	}

	return nil
}
