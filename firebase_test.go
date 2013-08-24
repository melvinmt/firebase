package firebase

import (
    // "fmt"
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "testing"
)

func TestThatNewReferenceWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var r *Reference
    r = NewReference(url)
    _ = r
}

type MockClient struct {
    Stub   StubPerson
    Status string
}

type StubPerson struct {
    First string `json:"first"`
    Last  string `json:"last"`
}

type readerAndCloser struct {
    io.Reader
    io.Closer
}

func (s StubPerson) Close() error {
    return nil
}

func (m *MockClient) Do(req *http.Request) (resp *http.Response, err error) {

    jsonData, err := json.Marshal(m.Stub)
    if err != nil {
        return nil, err
    }

    resp = &http.Response{
        Status: m.Status,
        Body: &readerAndCloser{
            bytes.NewReader(jsonData),
            m.Stub,
        },
    }

    return resp, nil
}

func TestThatWriteWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var err error

    person := StubPerson{
        First: "Fred",
        Last:  "Swanson",
    }

    r := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "200 OK"},
    }

    err = r.Write(person)
    if err != nil {
        t.Error(err)
    }

    rn := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "404 Not Found"},
    }

    err = rn.Write(person)
    if err == nil {
        t.Error("Any status code other than 2XX should throw an error.")
    }
}

func TestThatUpdateWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var err error

    person := StubPerson{
        First: "Fred",
        Last:  "Swanson",
    }

    r := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "200 OK"},
    }

    err = r.Update(person)
    if err != nil {
        t.Error(err)
    }

    rn := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "404 Not Found"},
    }

    err = rn.Update(person)
    if err == nil {
        t.Error("Any status code other than 2XX should throw an error.")
    }
}

func TestThatPushWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var err error

    person := StubPerson{
        First: "Fred",
        Last:  "Swanson",
    }

    r := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "200 OK"},
    }

    err = r.Push(person)
    if err != nil {
        t.Error(err)
    }

    rn := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "404 Not Found"},
    }

    err = rn.Push(person)
    if err == nil {
        t.Error("Any status code other than 2XX should throw an error.")
    }
}

func TestThatDeleteWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var err error

    person := StubPerson{}

    r := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "200 OK"},
    }

    err = r.Delete()
    if err != nil {
        t.Error(err)
    }

    rn := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "404 Not Found"},
    }

    err = rn.Delete()
    if err == nil {
        t.Error("Any status code other than 2XX should throw an error.")
    }
}

func TestThatValueWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var err error

    person := StubPerson{
        First: "Fred",
        Last:  "Swanson",
    }

    r := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "200 OK"},
    }

    norsep := StubPerson{}

    err = r.Value(&norsep)
    if err != nil {
        t.Error(err)
    }

    if norsep.First != "Fred" || norsep.Last != "Swanson" {
        t.Error("Invalid values returned.")
    }

    rn := &Reference{
        url:     url,
        postfix: ".json",
        client:  &MockClient{person, "404 Not Found"},
    }

    norsep2 := StubPerson{}

    err = rn.Value(&norsep2)
    if err == nil {
        t.Error("Any status code other than 2XX should throw an error.")
    }
}

func TestThatAuthWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var r *Reference
    r = NewReference(url).Auth("token").Auth("overwrite_token")
    _ = r
}

func TestThatExportWorks(t *testing.T) {
    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    var r *Reference
    r = NewReference(url).Export(true).Export(false)
    _ = r
}

package main

import (
    "github.com/melvinmt/firebase"
    "fmt"
)

type PersonName struct {
    First string
    Last  string
}

type Person struct {
    Name PersonName
}

func main() {
    var err error

    url := "https://SampleChat.firebaseIO-demo.com/users/fred/name"

    // Can also be your Firebase secret:
    authToken := "MqL0c8tKCtheLSYcygYNtGhU8Z2hULOFs9OKPdEp"

    // Auth is optional:
    ref := firebase.NewReference(url).Auth(authToken)

    // Create the value.
    personName := PersonName{
        First: "Fred",
        Last:  "Swanson",
    }

    // Write the value to Firebase.
    if err = ref.Write(personName); err != nil {
        panic(err)
    }

    // Now, we're going to retrieve the person.
    personUrl := "https://SampleChat.firebaseIO-demo.com/users/fred"

    personRef := firebase.NewReference(url).Export(false)

    fred := Person{}

    if err = personRef.Value(fred); err != nil {
        panic(err)
    }

    fmt.Println(fred.Name.First, fred.Name.Last) // prints: Fred Swanson
}
