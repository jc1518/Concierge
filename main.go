package main

import (
	"flag"
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

func main() {
	stackArns := flag.String("stacks-arn", "", "CloudFormation stacks ARN, use comma to seperate if more than one")
	flag.Parse()

	InfoLogger.Println("Validate CloudFormation stack ARN")
	for _, stackArn := range strings.Split(*stackArns, ",") {
		if !IsArnValid(stackArn) {
			ErrorLogger.Fatalf("%v is not a valid CloudFormation stack ARN", stackArn)
		}
	}
	InfoLogger.Println("Passed validation")

	InfoLogger.Println("Retrieve CloudConformity account information")
	allAccounts := GetAllAccounts()
	allAwsAccounts := GetAwsAccounts(allAccounts)

	for _, stackArn := range strings.Split(*stackArns, ",") {
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
