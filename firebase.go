/* A Go package for the Firebase API #golang
 */
package firebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type client interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type Reference struct {
	url          string
	postfix      string
	client       client
	token        string
	export       bool
	response     *http.Response
	responseBody []byte
}

// Retrieve a new Firebase reference for given url.
func NewReference(url string) *Reference {

	// Initialize Reference struct.
	r := &Reference{
		url:     url,
		postfix: ".json",
		export:  false,
		client:  &http.Client{},
	}

	return r
}

// Uses the Firebase secret or Auth Token to authenticate.
func (r *Reference) Auth(token string) *Reference {
	r.token = token

	return r
}

// Set to true if you want priority data to be returned.
func (r *Reference) Export(toggle bool) *Reference {
	r.export = toggle

	return r
}

// Execute a new HTTP Request.
func (r *Reference) executeRequest(method string, body []byte) ([]byte, error) {

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

	// Adding tiny sleep to prevent rate limited requests.
	time.Sleep(10 * time.Millisecond)

	// Prepare HTTP Request.
	req, err := http.NewRequest(method, apiUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Make actual HTTP request.
	if r.response, err = r.client.Do(req); err != nil {
		return nil, err
	}

	//Poor man's 307 handling.
	//This code will replay the previous request using the new url in the Location header
	//Hopefully https://code.google.com/p/go/issues/detail?id=7912 gets resolved so we can remove this.
	if strings.HasPrefix(r.response.Status, "307 ") {
		r.response.Body.Close() // close the original

		location := r.response.Header.Get("Location")
		req, err = http.NewRequest(method, location, bytes.NewReader(body))
		if r.response, err = r.client.Do(req); err != nil {
			return nil, err
		}
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

// Set the value for this Reference (overwrites existing value).
func (r *Reference) Write(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// PUT the data to Firebase.
	_, err = r.executeRequest("PUT", jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Pushes a new object to this Reference (effectively creates a list).
func (r *Reference) Push(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// POST the data to Firebase.
	_, err = r.executeRequest("POST", jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Update existing data.
func (r *Reference) Update(v interface{}) error {

	// JSON encode the data.
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// PATCH the data on Firebase.
	_, err = r.executeRequest("PATCH", jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Delete any values for this Reference.
func (r *Reference) Delete() error {

	// DELETE the data on Firebase.
	_, err := r.executeRequest("DELETE", nil)
	if err != nil {
		return err
	}

	return nil
}
