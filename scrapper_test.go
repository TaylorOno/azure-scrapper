package main

import (
	"context"
	"errors"
	"testing"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			}
			tt.want(t, s.Run())
		})
	}
}

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
			handlerError: errors.New("failed to handle resource group"),
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

func TestScrapper_ListProviders(t *testing.T) {
	tests := []struct {
		name             string
		providersFactory ProvidersClientFactory
		handlerError     error
		want             func(t *testing.T, err error)
	}{
		{
			name: "provider iteration succeeds",
			providersFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ProvidersPager, error) {
				return NewPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{item: &resource.ProvidersClientListResponse{
					ProviderListResult: resource.ProviderListResult{
						Value: []*resource.Provider{{}},
					},
				}}, nil
			},
			want: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "provider handler fails",
			providersFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ProvidersPager, error) {
				return NewPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{item: &resource.ProvidersClientListResponse{
					ProviderListResult: resource.ProviderListResult{
						Value: []*resource.Provider{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to handle provider"),
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "provider iteration fails",
			providersFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ProvidersPager, error) {
				return FailPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{}, nil
			},
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithProvidersFactory(tt.providersFactory))
			require.NoError(t, err)
			tt.want(t, scraper.ListProviders(context.Background(), func(r *resource.Provider) error { return tt.handlerError }))
		})
	}
}

func TestScrapper_ListVirtualNetworks(t *testing.T) {
	tests := []struct {
		name                   string
		virtualNetworksFactory VirtualNetworkClientFactory
		handlerError           error
		want                   func(t *testing.T, err error)
	}{
		{
			name: "virtual network iteration succeeds",
			virtualNetworksFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (VirtualNetworkPager, error) {
				return NewPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{item: &network.VirtualNetworksClientListAllResponse{
					VirtualNetworkListResult: network.VirtualNetworkListResult{
						Value: []*network.VirtualNetwork{{}},
					},
				}}, nil
			},
			want: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "virtual network handler fails",
			virtualNetworksFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (VirtualNetworkPager, error) {
				return NewPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{item: &network.VirtualNetworksClientListAllResponse{
					VirtualNetworkListResult: network.VirtualNetworkListResult{
						Value: []*network.VirtualNetwork{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to handle virtual network"),
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "virtual network iteration fails",
			virtualNetworksFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (VirtualNetworkPager, error) {
				return FailPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{}, nil
			},
			want: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithVirtualNetworksFactory(tt.virtualNetworksFactory))
			require.NoError(t, err)
			tt.want(t, scraper.ListVirtualNetworks(context.Background(), func(r *network.VirtualNetwork) error { return tt.handlerError }))
		})
	}
}

func TestScrapper_ListDiskEncryptionSets(t *testing.T) {
	tests := []struct {
		name                  string
		diskEncryptionFactory DiskEncryptionSetClientFactory
		handlerError          error
		expect                func(t *testing.T, err error)
	}{
		{
			name: "disk encryption set iteration succeeds",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{item: &compute.DiskEncryptionSetsClientListResponse{
					DiskEncryptionSetList: compute.DiskEncryptionSetList{
						Value: []*compute.DiskEncryptionSet{{}},
					},
				}}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "disk encryption handler fails",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{item: &compute.DiskEncryptionSetsClientListResponse{
					DiskEncryptionSetList: compute.DiskEncryptionSetList{
						Value: []*compute.DiskEncryptionSet{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to handle disk encryption set"),
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "disk encryption set iteration fails",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return FailPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithDiskEncryptionSetFactory(tt.diskEncryptionFactory))
			require.NoError(t, err)
			tt.expect(t, scraper.ListDiskEncryptionSets(context.Background(), func(r *compute.DiskEncryptionSet) error { return tt.handlerError }))
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

func errorPager[T any]() *rt.Pager[T] {
	return rt.NewPager[T](rt.PagingHandler[T]{
		More:    func(t T) bool { return false },
		Fetcher: func(ctx context.Context, t *T) (T, error) { return *new(T), errors.New("failed to iterate") },
	})
}
