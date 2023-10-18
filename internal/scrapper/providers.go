package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// ProvidersPager used to scrape provider information
type ProvidersPager interface {
	NewListPager(options *resource.ProvidersClientListOptions) *rt.Pager[resource.ProvidersClientListResponse]
}

type ProvidersClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ProvidersPager, error)

func defaultProvidersClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ProvidersPager, error) {
	return resource.NewProvidersClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListProviders(ctx context.Context, pageHandler pageHandler[resource.Provider]) error {
	pager := s.providersClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to advance page: %w", err)
		}
		if err = processPage(page.Value, err, pageHandler); err != nil {
			return err
		}
	}
	return nil
}
