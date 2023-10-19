package scrapper_test

import (
	. "azure-scrapper/internal/scrapper"
	"context"
	"errors"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScrapper_ListClusters(t *testing.T) {
	tests := []struct {
		name                 string
		clusterClientFactory ClusterClientFactory
		handlerError         error
		expect               func(t *testing.T, err error)
	}{
		{
			name: "cluster iteration succeeds",
			clusterClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ClusterPager, error) {
				return NewPager[container.ManagedClustersClientListOptions, container.ManagedClustersClientListResponse]{item: &container.ManagedClustersClientListResponse{
					ManagedClusterListResult: container.ManagedClusterListResult{
						Value: []*container.ManagedCluster{{}},
					},
				}}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "cluster handler fails",
			clusterClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ClusterPager, error) {
				return NewPager[container.ManagedClustersClientListOptions, container.ManagedClustersClientListResponse]{item: &container.ManagedClustersClientListResponse{
					ManagedClusterListResult: container.ManagedClusterListResult{
						Value: []*container.ManagedCluster{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to Handle cluster"),
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "cluster iteration fails",
			clusterClientFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ClusterPager, error) {
				return FailPager[container.ManagedClustersClientListOptions, container.ManagedClustersClientListResponse]{}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithClusterFactory(tt.clusterClientFactory))
			require.NoError(t, err)
			tt.expect(t, scraper.ListClusters(context.Background(), func(r *container.ManagedCluster) error { return tt.handlerError }))
		})
	}
}
