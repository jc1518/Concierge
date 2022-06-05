package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func TestValidArn(t *testing.T) {
	stackArn := "arn:aws:cloudformation:ap-southeast-2:123456789012:stack/my-stack/69d48220-010d-11ec-982a-06dd10360dfc"
	isValid := IsArnValid(stackArn)
	if !isValid {
		t.Errorf(`IsArnValid(%s) = %t; want true`, stackArn, isValid)
	}
}

func TestInvalidValidArn(t *testing.T) {
	stackArn := "arn:aws:cloudformation:ap-southeast-2:123456789:stack/my-stack/69d48220-010d-11ec-982a-06dd10360dfc"
	isValid := IsArnValid(stackArn)
	if isValid {
		t.Errorf(`IsArnValid(%s) = %t; want false`, stackArn, isValid)
	}
}

type mockListStackResourcesPager struct {
	PageNum int
	Pages   []*cloudformation.ListStackResourcesOutput
}

func (m *mockListStackResourcesPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockListStackResourcesPager) NextPage(ctx context.Context, f ...func(*cloudformation.Options)) (output *cloudformation.ListStackResourcesOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestListStackResources(t *testing.T) {
	r01 := "resource-01"
	r02 := "resource-02"
	r03 := "resource-03"
	r04 := "resource-04"
	r05 := "resource-05"

	pager := &mockListStackResourcesPager{
		Pages: []*cloudformation.ListStackResourcesOutput{
			{
				StackResourceSummaries: []types.StackResourceSummary{
					{
						PhysicalResourceId: &r01,
					},
					{
						PhysicalResourceId: &r02,
					},
				},
			},
			{
				StackResourceSummaries: []types.StackResourceSummary{
					{
						PhysicalResourceId: &r03,
					},
					{
						PhysicalResourceId: &r04,
					},
				},
			},
			{
				StackResourceSummaries: []types.StackResourceSummary{
					{
						PhysicalResourceId: &r05,
					},
				},
			},
		},
	}

	resourceIds, err := ListStackResources(context.TODO(), pager)

	if err != nil {
		t.Fatalf("expect no error, but got %v", err.Error())
	}

	if expect, actual := []string{"resource-01", "resource-02", "resource-03", "resource-04", "resource-05"}, resourceIds; !reflect.DeepEqual(expect, actual) {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}
