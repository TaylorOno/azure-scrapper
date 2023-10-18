package scrapper

import (
	"context"
	"errors"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScrapper_ListNodePools(t *testing.T) {
	tests := []struct {
		name                  string
		nodePoolClientFactory NodePoolClientFactory
		handlerError          error
		expect                func(t *testing.T, err error)
	}{
		{
			name: "node pool iteration succeeds",
			nodePoolClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (NodePoolPager, error) {
				return NewNodePager[container.AgentPoolsClientListOptions, container.AgentPoolsClientListResponse]{item: &container.AgentPoolsClientListResponse{
					AgentPoolListResult: container.AgentPoolListResult{
						Value: []*container.AgentPool{{}},
					},
				}}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "node pool handler fails",
			nodePoolClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (NodePoolPager, error) {
				return NewNodePager[container.AgentPoolsClientListOptions, container.AgentPoolsClientListResponse]{item: &container.AgentPoolsClientListResponse{
					AgentPoolListResult: container.AgentPoolListResult{
						Value: []*container.AgentPool{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to Handle node pool"),
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "node pool iteration fails",
			nodePoolClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (NodePoolPager, error) {
				return FailNodePager[container.AgentPoolsClientListOptions, container.AgentPoolsClientListResponse]{}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithNodePoolFactory(tt.nodePoolClientFactory))
			require.NoError(t, err)
			tt.expect(t, scraper.ListNodePool(context.Background(), "rg", "name", func(r *container.AgentPool) error { return tt.handlerError }))
		})
	}
}
