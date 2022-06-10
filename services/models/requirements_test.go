package models

import (
	"encoding/json"
	"errors"
	"testing"
)

//test various model functionality

func TestRequirementUnmarshal(t *testing.T) {
	type ReqtUser struct {
		Type RequirementType
	}
	tests := []struct {
		TestType      ReqtUser
		ExpectedError error
	}{
		{
			TestType: ReqtUser{
				Type: "standing",
			},
			ExpectedError: nil,
		},
		{
			TestType: ReqtUser{
				Type: "adhoc",
			},
			ExpectedError: nil,
		},
		{
			TestType: ReqtUser{
				Type: "somethingwrong",
			},
			ExpectedError: errors.New("don't care about message just that we have an error"),
		},
		{
			TestType: ReqtUser{
				Type: "",
			},
			ExpectedError: errors.New("still don't care about message"),
		},
	}
	for _, test := range tests {
		data, _ := json.Marshal(test.TestType)
		var unmarshallTestType ReqtUser
		err := json.Unmarshal(data, &unmarshallTestType)
		if err == nil && test.ExpectedError != nil {
			t.Errorf("expected non nil error but got nil error for test %v", test)
			t.FailNow()
		} else if err != nil && test.ExpectedError == nil {
			t.Errorf("expected nil error but got error %v, for test %v", err.Error(), test)
		}
	}
}
