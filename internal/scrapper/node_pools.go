package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
)

// NodePoolPager used to scrape node pool information
type NodePoolPager interface {
	NewListPager(resourceGroupName string, resourceName string, options *container.AgentPoolsClientListOptions) *rt.Pager[container.AgentPoolsClientListResponse]
}

type NodePoolClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (NodePoolPager, error)

func defaultNodePoolClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (NodePoolPager, error) {
	return container.NewAgentPoolsClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListNodePool(ctx context.Context, rg string, name string, pageHandler pageHandler[container.AgentPool]) error {
	pager := s.nodePoolClient.NewListPager(rg, name, nil)
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
