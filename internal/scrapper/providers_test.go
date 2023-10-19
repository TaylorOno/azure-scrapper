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
			handlerError: errors.New("failed to Handle provider"),
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
