package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
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

func ScanTemplate(templateFile string) {
	fmt.Print(templateFile)
	content := `{"AWSTemplateFormatVersion": "2010-09-09",  "Resources": {"AthenaS3Bucket": {"Type": "AWS::S3::Bucket"}}}`
	results := ScanCfnTemplate(content)
	for _, check := range results {
		if check.Attributes["status"] == "FAILURE" {
			WarningLogger.Printf("%v (%v): %v", check.Attributes["risk-level"], check.Attributes["descriptorType"], check.Attributes["message"])
		}
	}

}

func main() {
	stackArns := flag.String("stacks-arn", "", "CloudFormation stacks ARN, use comma to seperate if more than one")
	templateFile := flag.String("template-file", "", "CloudFormation template file")
	flag.Parse()

	if stackArns != nil {
		ScanStack(*stackArns)
	}

	if templateFile != nil {
		ScanTemplate(*templateFile)
	}

}
