package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	VERSION                  = "0.3.0-beta"
	CLOUD_CONFORMITY_API_KEY = os.Getenv("CLOUD_CONFORMITY_API_KEY")
	InfoLogger               *log.Logger
	WarningLogger            *log.Logger
	ErrorLogger              *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

	if CLOUD_CONFORMITY_API_KEY == "" {
		ErrorLogger.Fatalf("CLOUD_CONFORMITY_API_KEY variable is missing")
	}

}

func ScanStack(stackArns string) {
	InfoLogger.Println("Validate CloudFormation stack ARN")
	for _, stackArn := range strings.Split(stackArns, ",") {
		if !IsArnValid(stackArn) {
			ErrorLogger.Fatalf("%v is not a valid CloudFormation stack ARN", stackArn)
		}
	}
	InfoLogger.Println("Passed validation")

	InfoLogger.Println("Retrieve CloudConformity account information")
	allAccounts := GetAllAccounts()
	allAwsAccounts := GetAwsAccounts(allAccounts)

	for _, stackArn := range strings.Split(stackArns, ",") {
		account := strings.Split(stackArn, ":")[4]
		ccId := GetAwsAccountCcId(allAwsAccounts, account)
		InfoLogger.Printf("List resources in stack %v", stackArn)
		resourceIds := GetStackResources(stackArn)
		for _, resourceId := range resourceIds {
			InfoLogger.Printf("Retrieve check results of %v", resourceId)
			for _, check := range GetResourcesFailedChecks(ccId, resourceId) {
				WarningLogger.Printf("%v (%v): %v", check.Attributes["risk-level"], check.Attributes["resourceName"], check.Attributes["message"])
			}
		}
	}
}

func ReadTemplateFile(templateFile string) string {
	lower := strings.ToLower(templateFile)
	IsYamOrJson := strings.Contains(lower, ".yaml") || strings.Contains(lower, ".yml") || strings.Contains(lower, ".json")
	if !IsYamOrJson {
		ErrorLogger.Fatalf("CloudFormation template has to be json or yaml file")
	}

	InfoLogger.Printf("Load CloudFormation template %s", templateFile)
	contentsBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		ErrorLogger.Fatalf("Failed to load template file - %v", err.Error())
	}
	contents := string(contentsBytes)
	if !IsTemplateValid(contents) {
		ErrorLogger.Fatalf("Not a valid CloudFormation template")
	}

	if strings.Contains(lower, ".yaml") || strings.Contains(lower, ".yml") {
		var temp interface{}
		err := yaml.Unmarshal(contentsBytes, &temp)
		if err != nil {
			ErrorLogger.Fatalf("Failed to unmarshall yaml - %v", err.Error())
		}

		contentsBytes, err = json.Marshal(&temp)
		if err != nil {
			ErrorLogger.Fatalf("Failed to marshall json - %v", err.Error())
		}

		contents = string(contentsBytes)
	}
	return contents
}

func ScanTemplate(templateContents string) {
	InfoLogger.Printf("Scan template contents")
	results := ScanCfnTemplate(templateContents)
	for _, check := range results {
		if check.Attributes["status"] == "FAILURE" {
			WarningLogger.Printf("%v (%v): %v", check.Attributes["risk-level"], check.Attributes["descriptorType"], check.Attributes["message"])
		}
	}

}

func help() {
	fmt.Println("Version:", VERSION)
}

func main() {
	stackArns := flag.String("stacks-arn", "", "CloudFormation stacks ARN, use comma to seperate if more than one")
	templateFile := flag.String("template-file", "", "CloudFormation template file (json or yaml)")
	flag.Parse()

	if *stackArns == "" && *templateFile == "" {
		help()
	}

	if *stackArns != "" {
		ScanStack(*stackArns)
	}

	if *templateFile != "" {
		templateContents := ReadTemplateFile(*templateFile)
		ScanTemplate(templateContents)
	}
}
