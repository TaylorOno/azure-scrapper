package scrapper_test

import (
	. "azure-scrapper/internal/scrapper"
	"context"
	"errors"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScrapper_ListResourceGroups(t *testing.T) {
	tests := []struct {
		name                  string
		resourceGroupsFactory ResourceGroupClientFactory
		handlerError          error
		want                  func(t *testing.T, err error)
	}{
		{
			name: "resource group iteration succeeds",
			resourceGroupsFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ResourceGroupsPager, error) {
				return NewPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{item: &resource.ResourceGroupsClientListResponse{
					ResourceGroupListResult: resource.ResourceGroupListResult{
						Value: []*resource.ResourceGroup{{}},
					},
				}}, nil
			},
			want: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "resource group handler fails",
			resourceGroupsFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ResourceGroupsPager, error) {
				return NewPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{item: &resource.ResourceGroupsClientListResponse{
					ResourceGroupListResult: resource.ResourceGroupListResult{
						Value: []*resource.ResourceGroup{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to Handle resource group"),
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "resource group iteration fails",
			resourceGroupsFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ResourceGroupsPager, error) {
				return FailPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{}, nil
			},
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithResourceGroupsFactory(tt.resourceGroupsFactory))
			require.NoError(t, err)
			tt.want(t, scraper.ListResourceGroups(context.Background(), func(r *resource.ResourceGroup) error { return tt.handlerError }))
		})
	}
}
