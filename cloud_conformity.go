package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const CLOUD_CONFORMITY_BASE_URL = "https://ap-southeast-2-api.cloudconformity.com/v1/"

type Data struct {
	Data []Entity `json:"data"`
}

// Entity can be Account or Check
type Entity struct {
	Type       string                 `json:"type"`
	Id         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

func GetRequest(url string) (entities []Entity, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", "ApiKey "+CLOUD_CONFORMITY_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	entities = data.Data
	return entities, nil
}

func PostRequest(url string, payload string) (entities []Entity, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", "ApiKey "+CLOUD_CONFORMITY_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	entities = data.Data
	return entities, nil
}

func GetAllAccounts() []Entity {
	accounts, err := GetRequest(CLOUD_CONFORMITY_BASE_URL + "accounts")
	if err != nil {
		ErrorLogger.Fatalf("faile to get all accounts, %v", err.Error())
	}

	return accounts
}

func GetAwsAccounts(accounts []Entity) []Entity {
	var awsAccounts []Entity

	for _, account := range accounts {
		if _, exist := account.Attributes["awsaccount-id"]; exist {
			awsAccounts = append(awsAccounts, account)
		}
	}

	return awsAccounts
}

func GetAwsAccountCcId(awsAccounts []Entity, accountId string) string {
	for _, account := range awsAccounts {
		if account.Attributes["awsaccount-id"] == accountId {
			return account.Id
		}
	}

	return ""
}

func GetResourcesFailedChecks(ccId string, resourceId string) []Entity {
	path := "checks?accountIds=" + ccId + "&filter[resource]=" + resourceId + "&filter[statuses]=FAILURE"

	failedChecks, err := GetRequest(CLOUD_CONFORMITY_BASE_URL + path)
	if err != nil {
		ErrorLogger.Fatalf("failed to get failed checks, %v", err.Error())
	}

	return failedChecks
}

func ScanCfnTemplate(contents string) []Entity {
	payload := `{"data":{"attributes":{"type":"cloudformation-template","contents":` + contents + "}}}}"

	results, err := PostRequest(CLOUD_CONFORMITY_BASE_URL+"template-scanner/scan", payload)
	if err != nil {
		ErrorLogger.Fatalf("failed to scan the template, %v", err.Error())
	}

	return results
}
