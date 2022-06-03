package main

import (
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func IsArnValid(arn string) bool {
	var validArn = regexp.MustCompile(`^arn:aws:cloudformation:(?P<Region>[^:\n]*):\d{12}:stack/?/.*$`)

	return validArn.MatchString(arn)
}

func GetStackResources(arn string) []string {
	region := strings.Split(arn, ":")[3]
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		ErrorLogger.Fatalf("failed to load aws config, %v", err.Error())
	}

	var resourceIds []string

	client := cloudformation.NewFromConfig(cfg)
	pages := cloudformation.NewListStackResourcesPaginator(client, &cloudformation.ListStackResourcesInput{
		StackName: &arn,
	})

	for pages.HasMorePages() {
		page, err := pages.NextPage(context.TODO())
		if err != nil {
			ErrorLogger.Fatalf("failed to get stack resources, %v", err.Error())
		}
		for _, resource := range page.StackResourceSummaries {
			resourceIds = append(resourceIds, *resource.PhysicalResourceId)
		}
	}

	return resourceIds
}
