package scrapper

import (
	"context"
	"errors"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
			handlerError: errors.New("failed to Handle virtual network"),
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
