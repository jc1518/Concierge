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

type ListStackResourcesPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*cloudformation.Options)) (*cloudformation.ListStackResourcesOutput, error)
}

func ListStackResources(ctx context.Context, pager ListStackResourcesPager) (resourceIds []string, err error) {
	for pager.HasMorePages() {
		var page *cloudformation.ListStackResourcesOutput
		page, err := pager.NextPage(ctx)
		if err != nil {
			return resourceIds, err
		}
		for _, resource := range page.StackResourceSummaries {
			resourceIds = append(resourceIds, *resource.PhysicalResourceId)
		}
	}

	return resourceIds, nil
}

func GetStackResources(arn string) []string {
	region := strings.Split(arn, ":")[3]
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		ErrorLogger.Fatalf("failed to load aws config, %v", err.Error())
	}

	client := cloudformation.NewFromConfig(cfg)

	pager := cloudformation.NewListStackResourcesPaginator(client, &cloudformation.ListStackResourcesInput{
		StackName: &arn,
	})

	resp, err := ListStackResources(context.TODO(), pager)
	if err != nil {
		ErrorLogger.Fatalf("failed to get stack resources, %v", err.Error())
	}

	return resp

}
