package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var allAccounts = []Entity{
	{
		Type: "accounts",
		Id:   "12345",
		Attributes: map[string]interface{}{
			"name":       "azure account 01",
			"cloud-type": "azure",
			"cloud-data": map[string]interface{}{
				"azure": map[string]interface{}{
					"subscriptionId": "123-456-789",
				},
			},
		},
	},
	{
		Type: "accounts",
		Id:   "67890",
		Attributes: map[string]interface{}{
			"name":          "aws account 01",
			"awsaccount-id": "123456789012",
			"cloud-type":    "aws",
		},
	},
}

var awsAccounts = []Entity{
	{
		Type: "accounts",
		Id:   "67890",
		Attributes: map[string]interface{}{
			"name":          "aws account 01",
			"awsaccount-id": "123456789012",
			"cloud-type":    "aws",
		},
	},
}

func TestGetRequest(t *testing.T) {
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse []Entity
		expectedErr      error
	}{
		{
			name: "Get_all_accounts",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data": [
					{"type": "accounts", "id": "12345", "attributes": {"name": "azure account 01", "cloud-type": "azure", "cloud-data": {"azure" :{"subscriptionId": "123-456-789"}}}}, 
					{"type": "accounts", "id": "67890","attributes": {"name":"aws account 01", "awsaccount-id": "123456789012", "cloud-type":"aws"}}]}`))
			})),
			expectedResponse: allAccounts,
			expectedErr:      nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.server.Close()
			resp, err := MakeRequest(tc.server.URL, nil)
			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("expect %v, got %v", tc.expectedResponse, resp)
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expect %v, got %v", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestGetAwsAccounts(t *testing.T) {
	expect := awsAccounts
	actual := GetAwsAccounts(allAccounts)
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}

func TestGetAwsAccountCcId(t *testing.T) {
	expect := "67890"
	actual := GetAwsAccountCcId(awsAccounts, "123456789012")
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}
