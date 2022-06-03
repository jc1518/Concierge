package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const CLOUD_CONFORMITY_BASE_URL = "https://ap-southeast-2-api.cloudconformity.com/v1/"

func getRequest(path string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", CLOUD_CONFORMITY_BASE_URL+path, nil)
	if err != nil {
		ErrorLogger.Fatalf("failed to form the request, %v", err.Error())
	}

	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", "ApiKey "+CLOUD_CONFORMITY_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		ErrorLogger.Fatalf("failed to send request to url, %v", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLogger.Fatalf("failed to get response from url, %v", err.Error())
	}

	return body
}

// Entity can be Account or Check
type Entity struct {
	Type       string                 `json:"type"`
	Id         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Data struct {
	Data []Entity `json:"data"`
}

func GetAllAccounts() []Entity {
	var accounts Data

	err := json.Unmarshal(getRequest("accounts"), &accounts)
	if err != nil {
		ErrorLogger.Fatalf("faile to get all accounts, %v", err.Error())
	}

	return accounts.Data
}

func GetAwsAccounts(allAccounts []Entity) []Entity {
	var awsAccounts []Entity

	for _, account := range allAccounts {
		if _, exist := account.Attributes["awsaccount-id"]; exist {
			awsAccounts = append(awsAccounts, account)
		}
	}

	return awsAccounts
}

func GetAwsAccountCcId(allAwsAccounts []Entity, accountId string) string {
	for _, account := range allAwsAccounts {
		if account.Attributes["awsaccount-id"] == accountId {
			return account.Id
		}
	}

	return ""
}

func GetResourcesFailedChecks(ccId string, resourceId string) []Entity {
	var failedChecks Data

	path := "checks?accountIds=" + ccId + "&filter[resource]=" + resourceId + "&filter[statuses]=FAILURE"
	err := json.Unmarshal(getRequest(path), &failedChecks)
	if err != nil {
		ErrorLogger.Fatalf("faile to get failed checks, %v", err.Error())
	}

	return failedChecks.Data
}
