firebase
========

A Go package for the Firebase API #golang

## Installation

Install the package with go:

```sh
go get github.com/melvinmt/firebase
```

And add it to your go file:

```go
package name

import (
    "github.com/melvinmt/firebase"
)
```

## Usage

```go
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

    personRef := firebase.NewReference(personUrl).Export(false)

    fred := Person{}

    if err = personRef.Value(fred); err != nil {
        panic(err)
    }

    fmt.Println(fred.Name.First, fred.Name.Last) // prints: Fred Swanson
}
```

## Docs

http://godoc.org/github.com/melvinmt/firebase
