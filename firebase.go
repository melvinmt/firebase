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

type Reference struct {
	url     string
	postfix string
	client  *http.Client
	token   string
	export  bool
}

// Retrieve a new Firebase reference for given url.
func NewReference(url string) *Reference {

	// Initialize Reference struct.
	r := &Reference{
		url:     url,
		postfix: ".json",
		client:  &http.Client{},
		export:  false,
	}

	return r
}

func (r *Reference) Auth(token string) {
	r.token = token
}

func (r *Reference) Export(toggle bool) {
	r.export = toggle
}

// Execute a new HTTP Request.
func (r *Reference) executeRequest(method string, body io.Reader) ([]byte, error) {

	apiUrl := r.url + r.postfix

	v := url.Values{}
	if r.token != "" {
		v.Set("auth", r.token)
	}
	if r.export == true {
		v.Set("format", "export")
	}
	q := v.Encode()
	if len(q) > 0 {
		apiUrl = apiUrl + "?" + q
	}

	// Prepare HTTP Request.
	req, err := http.NewRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Make actual HTTP request.
	resp, err := r.client.Do(req)
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

	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = r.executeRequest("PUT", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Push(v interface{}) error {

	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = r.executeRequest("POST", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Update(v interface{}) error {

	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = r.executeRequest("PATCH", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (r *Reference) Delete() error {

	_, err := r.executeRequest("DELETE", nil)
	if err != nil {
		return err
	}

	return nil
}
