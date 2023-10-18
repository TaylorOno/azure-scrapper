package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// ResourceGroupsPager used to scrape resource group information
type ResourceGroupsPager interface {
	NewListPager(options *resource.ResourceGroupsClientListOptions) *rt.Pager[resource.ResourceGroupsClientListResponse]
}

type ResourceGroupClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ResourceGroupsPager, error)

func defaultResourceGroupFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ResourceGroupsPager, error) {
	return resource.NewResourceGroupsClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListResourceGroups(ctx context.Context, pageHandler pageHandler[resource.ResourceGroup]) error {
	pager := s.resourceGroupClient.NewListPager(nil)
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
