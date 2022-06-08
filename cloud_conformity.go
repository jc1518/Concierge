package main

import (
	"bytes"
	"encoding/json"
	"io"
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

func MakeRequest(url string, payload io.Reader) (entities []Entity, err error) {
	client := &http.Client{}

	method := "GET"
	if payload != nil {
		method = "POST"
	}

	req, err := http.NewRequest(method, url, payload)
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
	accounts, err := MakeRequest(CLOUD_CONFORMITY_BASE_URL+"accounts", nil)
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

	failedChecks, err := MakeRequest(CLOUD_CONFORMITY_BASE_URL+path, nil)
	if err != nil {
		ErrorLogger.Fatalf("failed to get failed checks, %v", err.Error())
	}

	return failedChecks
}

func ScanCfnTemplate(templateContents string) []Entity {
	payload := `{"data":{"attributes":{"type":"cloudformation-template","contents":` + templateContents + "}}}}"

	results, err := MakeRequest(CLOUD_CONFORMITY_BASE_URL+"template-scanner/scan", bytes.NewBufferString(payload))
	if err != nil {
		ErrorLogger.Fatalf("failed to scan the template, %v", err.Error())
	}

	return results
}
