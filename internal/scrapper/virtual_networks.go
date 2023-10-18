package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

// VirtualNetworkPager used to scrape virtual network information
type VirtualNetworkPager interface {
	NewListAllPager(options *network.VirtualNetworksClientListAllOptions) *rt.Pager[network.VirtualNetworksClientListAllResponse]
}

type VirtualNetworkClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (VirtualNetworkPager, error)

func defaultNetworkClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (VirtualNetworkPager, error) {
	return network.NewVirtualNetworksClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListVirtualNetworks(ctx context.Context, pageHandler pageHandler[network.VirtualNetwork]) error {
	pager := s.networksClient.NewListAllPager(nil)
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
