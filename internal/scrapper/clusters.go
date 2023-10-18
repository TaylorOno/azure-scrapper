package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
)

// ClusterPager used to scrape aks cluster information
type ClusterPager interface {
	NewListPager(options *container.ManagedClustersClientListOptions) *rt.Pager[container.ManagedClustersClientListResponse]
}

type ClusterClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ClusterPager, error)

func defaultClusterClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ClusterPager, error) {
	return container.NewManagedClustersClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListClusters(ctx context.Context, pageHandler pageHandler[container.ManagedCluster]) error {
	pager := s.clusterClient.NewListPager(nil)
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
