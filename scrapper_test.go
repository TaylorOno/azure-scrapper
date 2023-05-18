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
		want          func(*Scrapper, error)
	}{
		{
			name:          "successful creation",
			withFactories: []OptionsFunc{},
			want: func(scrapper *Scrapper, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:          "fails if resource group client factory fails",
			withFactories: []OptionsFunc{WithResourceGroupsFactory(brokenFactory[ResourceGroupsPager])},
			want: func(scrapper *Scrapper, err error) {
				assert.Error(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if provider client factory fails",
			withFactories: []OptionsFunc{WithProvidersFactory(brokenFactory[ProvidersPager])},
			want: func(scrapper *Scrapper, err error) {
				assert.Error(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if virtual network client factory fails",
			withFactories: []OptionsFunc{WithVirtualNetworksFactory(brokenFactory[VirtualNetworkPager])},
			want: func(scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
		{
			name:          "fails if disk encryption set client factory fails",
			withFactories: []OptionsFunc{WithDiskEncryptionSetFactory(brokenFactory[DiskEncryptionSetPager])},
			want: func(scrapper *Scrapper, err error) {
				assert.EqualError(t, err, "failed to create client")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want(NewScrapper(testCred(), "not-important", tt.withFactories...))
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
		want    func(error)
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
			want: func(err error) {
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
			tt.want(s.Run())
		})
	}
}

func TestScrapper_ListResourceGroups(t *testing.T) {
	tests := []struct {
		name          string
		withFactories []OptionsFunc
		want          func(error)
	}{
		{
			name: "resource group iteration succeeds",
			withFactories: []OptionsFunc{WithResourceGroupsFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ResourceGroupsPager, error) {
				return NewPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{item: &resource.ResourceGroupsClientListResponse{}}, nil
			})},
			want: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "resource group iteration fails",
			withFactories: []OptionsFunc{WithResourceGroupsFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ResourceGroupsPager, error) {
				return FailPager[resource.ResourceGroupsClientListOptions, resource.ResourceGroupsClientListResponse]{}, nil
			})},
			want: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", tt.withFactories...)
			require.NoError(t, err)
			tt.want(scraper.ListResourceGroups(context.Background(), func(r *resource.ResourceGroup) error { return nil }))
		})
	}
}

func TestScrapper_ListProviders(t *testing.T) {
	tests := []struct {
		name          string
		withFactories []OptionsFunc
		want          func(error)
	}{
		{
			name: "provider iteration succeeds",
			withFactories: []OptionsFunc{WithProvidersFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ProvidersPager, error) {
				return NewPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{item: &resource.ProvidersClientListResponse{}}, nil
			})},
			want: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "provider iteration fails",
			withFactories: []OptionsFunc{WithProvidersFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (ProvidersPager, error) {
				return FailPager[resource.ProvidersClientListOptions, resource.ProvidersClientListResponse]{}, nil
			})},
			want: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", tt.withFactories...)
			require.NoError(t, err)
			tt.want(scraper.ListProviders(context.Background(), func(r *resource.Provider) error { return nil }))
		})
	}
}

func TestScrapper_ListVirtualNetworks(t *testing.T) {
	tests := []struct {
		name          string
		withFactories []OptionsFunc
		want          func(error)
	}{
		{
			name: "virtual network iteration succeeds",
			withFactories: []OptionsFunc{WithVirtualNetworksFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (VirtualNetworkPager, error) {
				return NewPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{item: &network.VirtualNetworksClientListAllResponse{}}, nil
			})},
			want: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "virtual network iteration fails",
			withFactories: []OptionsFunc{WithVirtualNetworksFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (VirtualNetworkPager, error) {
				return FailPager[network.VirtualNetworksClientListAllOptions, network.VirtualNetworksClientListAllResponse]{}, nil
			})},
			want: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", tt.withFactories...)
			require.NoError(t, err)
			tt.want(scraper.ListVirtualNetworks(context.Background(), func(r *network.VirtualNetwork) error { return nil }))
		})
	}
}

func TestScrapper_ListDiskEncryptionSets(t *testing.T) {
	tests := []struct {
		name          string
		withFactories []OptionsFunc
		want          func(error)
	}{
		{
			name: "disk encryption set iteration succeeds",
			withFactories: []OptionsFunc{WithDiskEncryptionSetFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{item: &compute.DiskEncryptionSetsClientListResponse{}}, nil
			})},
			want: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "disk encryption set iteration fails",
			withFactories: []OptionsFunc{WithDiskEncryptionSetFactory(func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return FailPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{}, nil
			})},
			want: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", tt.withFactories...)
			require.NoError(t, err)
			tt.want(scraper.ListDiskEncryptionSets(context.Background(), func(r *compute.DiskEncryptionSet) error { return nil }))
		})
	}
}

func testCred() az.TokenCredential {
	return &azidentity.DefaultAzureCredential{}
}

func brokenFactory[T any](sub string, cred az.TokenCredential, opts *arm.ClientOptions) (T, error) {
	return *new(T), errors.New("failed to create client")
}

type NewPager[O any, T any] struct {
	item *T
}

func (n NewPager[O, T]) NewListPager(opts *O) *rt.Pager[T] {
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

func (f FailPager[O, T]) NewListPager(opts *O) *rt.Pager[T] {
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
