package firebase

import (
	"testing"
)

func TestThatValueWorks(t *testing.T) {
	url := "https://go-test-firebase.firebaseio.com/"

	var r *Reference
	r = NewReference(url)

	type List struct {
		MessageList struct {
			Text    string `json:"text"`
			User_Id string `json:"user_id"`
		} `json:"message_list"`
	}

	var list List
	err := r.Value(&list)
	if err != nil {
		t.Error(err)
	}

}
