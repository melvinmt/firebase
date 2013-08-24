/* A Go package for the Firebase API #golang
 */
package firebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type clientInterface interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type reference struct {
	url          string
	postfix      string
	client       clientInterface
	token        string
	export       bool
	response     *http.Response
	responseBody []byte
}

// Retrieve a new Firebase reference for given url.
func NewReference(url string) *reference {

	// Initialize reference struct.
	r := &reference{
		url:     url,
		postfix: ".json",
		export:  false,
		client:  &http.Client{},
	}

	return r
}

// Uses the Firebase secret or Auth Token to authenticate.
func (r *reference) Auth(token string) *reference {
	r.token = token

	return r
}

// Set to true if you want priority data to be returned.
func (r *reference) Export(toggle bool) *reference {
	r.export = toggle

	return r
}

// Execute a new HTTP Request.
func (r *reference) executeRequest(method string, body io.Reader) ([]byte, error) {

	apiUrl := r.url + r.postfix

	// Build query parameters (if any).
	v := url.Values{}
	if r.token != "" {
		v.Set("auth", r.token)
	}
	if r.export == true {
		v.Set("format", "export")
	}
	q := v.Encode()

	// Attach query parameters to apiUrl.
	if len(q) > 0 {
		apiUrl = apiUrl + "?" + q
	}

	// Prepare HTTP Request.
	req, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		return nil, err
	}

	// Make actual HTTP request.
	if r.response, err = r.client.Do(req); err != nil {
		return nil, err
	}
	defer r.response.Body.Close()

	// Check status code for errors.
	status := r.response.Status
	if strings.HasPrefix(status, "2") == false {
		return nil, errors.New(status)
	}

	// Read body.
	if r.responseBody, err = ioutil.ReadAll(r.response.Body); err != nil {
		return nil, err
	}

	return r.responseBody, nil
}

// Retrieve the current value for this reference.
func (r *reference) Value(v interface{}) error {

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

// Set the value for this reference (overwrites existing value).
func (r *reference) Write(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// PUT the data to Firebase.
	_, err = r.executeRequest("PUT", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

// Pushes a new object to this reference (effectively creates a list).
func (r *reference) Push(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// POST the data to Firebase.
	_, err = r.executeRequest("POST", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

// Update existing data.
func (r *reference) Update(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// PATCH the data on Firebase.
	_, err = r.executeRequest("PATCH", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

// Delete any values for this reference.
func (r *reference) Delete() error {

	// DELETE the data on Firebase.
	_, err := r.executeRequest("DELETE", nil)
	if err != nil {
		return err
	}

	return nil
}
