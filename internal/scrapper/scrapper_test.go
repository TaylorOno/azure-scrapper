package scrapper

import (
	"context"
	"errors"
	"testing"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
)

func TestNewScraper(t *testing.T) {
	tests := []struct {
		name          string
		withFactories []OptionsFunc
		want          func(*testing.T, *Scrapper, error)
	}{
		{
			name:          "successful creation",
			withFactories: []OptionsFunc{},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:          "fails if resource group client factory fails",
			withFactories: []OptionsFunc{WithResourceGroupsFactory(brokenFactory[ResourceGroupsPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.Error(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if provider client factory fails",
			withFactories: []OptionsFunc{WithProvidersFactory(brokenFactory[ProvidersPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.Error(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if virtual network client factory fails",
			withFactories: []OptionsFunc{WithVirtualNetworksFactory(brokenFactory[VirtualNetworkPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if disk encryption set client factory fails",
			withFactories: []OptionsFunc{WithDiskEncryptionSetFactory(brokenFactory[DiskEncryptionSetPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if cluster client factory fails",
			withFactories: []OptionsFunc{WithClusterFactory(brokenFactory[ClusterPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if node pool client factory fails",
			withFactories: []OptionsFunc{WithNodePoolFactory(brokenFactory[NodePoolPager])},
			want: func(t *testing.T, scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scrapper, err := NewScrapper(testCred(), "not-important", tt.withFactories...)
			tt.want(t, scrapper, err)
		})
	}
}

func TestScrapper_Run(t *testing.T) {
	type clients struct {
		resourceGroupClient      ResourceGroupsPager
		providersClient          ProvidersPager
		networksClient           VirtualNetworkPager
		diskEncryptionSetsClient DiskEncryptionSetPager
		clusterClient            ClusterPager
	}
	tests := []struct {
		name    string
		clients clients
		want    func(t *testing.T, err error)
	}{
		{
			name: "Successful execution",
			clients: clients{
				resourceGroupClient: NewPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{
					item: &resource.ResourceGroupsClientListResponse{},
				},
				providersClient: NewPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{
					item: &resource.ProvidersClientListResponse{},
				},
				networksClient: NewPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{
					item: &network.VirtualNetworksClientListAllResponse{},
				},
				diskEncryptionSetsClient: NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{
					item: &compute.DiskEncryptionSetsClientListResponse{},
				},
				clusterClient: NewPager[container.ManagedClustersClientListOptions, container.ManagedClustersClientListResponse]{
					item: &container.ManagedClustersClientListResponse{},
				},
			},
			want: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scrapper{
				resourceGroupClient:     tt.clients.resourceGroupClient,
				providersClient:         tt.clients.providersClient,
				networksClient:          tt.clients.networksClient,
				diskEncryptionSetClient: tt.clients.diskEncryptionSetsClient,
				clusterClient:           tt.clients.clusterClient,
			}
			tt.want(t, s.Run())
		})
	}
}

func testCred() az.TokenCredential {
	return &azidentity.DefaultAzureCredential{}
}

func brokenFactory[T any](_ string, _ az.TokenCredential, _ *arm.ClientOptions) (T, error) {
	return *new(T), errors.New("failed to create client")
}

type NewPager[O any, T any] struct {
	item *T
}

func (n NewPager[O, T]) NewListPager(_ *O) *rt.Pager[T] {
	return singleItemPager[T](n.item)
}

func (n NewPager[O, T]) NewListAllPager(opts *O) *rt.Pager[T] {
	return n.NewListPager(opts)
}

type NewNodePager[O any, T any] struct {
	item *T
}

func (n NewNodePager[O, T]) NewListPager(_ string, _ string, _ *O) *rt.Pager[T] {
	return singleItemPager[T](n.item)
}

func singleItemPager[T any](item *T) *rt.Pager[T] {
	return rt.NewPager[T](rt.PagingHandler[T]{
		More:    func(t T) bool { return false },
		Fetcher: func(ctx context.Context, t *T) (T, error) { return *item, nil },
	})
}

type FailPager[O any, T any] struct {
}

func (f FailPager[O, T]) NewListPager(_ *O) *rt.Pager[T] {
	return errorPager[T]()
}

func (f FailPager[O, T]) NewListAllPager(opts *O) *rt.Pager[T] {
	return f.NewListPager(opts)
}

type FailNodePager[O any, T any] struct {
}

func (f FailNodePager[O, T]) NewListPager(_ string, _ string, _ *O) *rt.Pager[T] {
	return errorPager[T]()
}

func errorPager[T any]() *rt.Pager[T] {
	return rt.NewPager[T](rt.PagingHandler[T]{
		More:    func(t T) bool { return false },
		Fetcher: func(ctx context.Context, t *T) (T, error) { return *new(T), errors.New("failed to iterate") },
	})
}
