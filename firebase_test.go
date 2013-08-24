package firebase

import (
	// "fmt"
	"testing"
)

func TestThatValueWorks(t *testing.T) {
	url := "https://go-test-firebase.firebaseio.com/users/fred/name"

	var r *Reference
	r = NewReference(url)

	type Person struct {
		First string `json:"first"`
		Last  string `json:"last"`
	}

	person := Person{
		First: "Fred",
		Last:  "Swanson",
	}
	// err := r.Write(person)
	// fmt.Println("err:", err)
	// if err != nil {
	// 	t.Error(err)
	// }

	person.Last = "Johnson"

	err := r.Write(person)
	if err != nil {
		t.Error(err)
	}
	person.Last = "Johnson"

	err = r.Write(person)
	if err != nil {
		t.Error(err)
	}

	person.Last = "Tercan"

	err = r.Write(person)
	if err != nil {
		t.Error(err)
	}

	// err = r.Delete()
	// if err != nil {
	// 	t.Error(err)
	// }

}
