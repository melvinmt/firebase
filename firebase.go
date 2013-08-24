package firebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Reference struct {
	url          string
	postfix      string
	client       *http.Client
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

func (r *Reference) Auth(token string) *Reference {
	r.token = token

	return r
}

func (r *Reference) Export(toggle bool) *Reference {
	r.export = toggle

	return r
}

// Execute a new HTTP Request.
func (r *Reference) executeRequest(method string, body io.Reader) ([]byte, error) {

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

	fmt.Println("ApiURL:", apiUrl)

	// Prepare HTTP Request.
	req, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		fmt.Println("err1:", err)
		return nil, err
	}

	client := &http.Client{}

	// Make actual HTTP request.
	if r.response, err = client.Do(req); err != nil {
		fmt.Println("err2:", err)
		return nil, err
	}
	defer r.response.Body.Close()

	// Check status code for errors.
	status := r.response.Status
	fmt.Println("status:", status)
	if strings.HasPrefix(status, "2") == false {
		fmt.Println("err3:", err)
		return nil, errors.New(status)
	}

	// Read body.
	if r.responseBody, err = ioutil.ReadAll(r.response.Body); err != nil {
		fmt.Println("err4:", err)
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

	fmt.Println("v:", v)
	fmt.Println("body:", string(jsonData))

	// PUT the data to Firebase.
	_, err = r.executeRequest("PUT", bytes.NewReader(jsonData))
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
	_, err = r.executeRequest("POST", bytes.NewReader(jsonData))
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

	fmt.Println("v:", v)
	fmt.Println("body:", string(jsonData))

	// PATCH the data on Firebase.
	_, err = r.executeRequest("PATCH", bytes.NewReader(jsonData))
	fmt.Println("err:", err)
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
