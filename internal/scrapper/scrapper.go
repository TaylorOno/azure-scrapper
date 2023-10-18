package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"golang.org/x/sync/errgroup"
)

type pageHandler[T any] func(r *T) error

type Scrapper struct {
	resourceGroupClient     ResourceGroupsPager
	providersClient         ProvidersPager
	networksClient          VirtualNetworkPager
	diskEncryptionSetClient DiskEncryptionSetPager
	clusterClient           ClusterPager
	nodePoolClient          NodePoolPager
}

// NewScrapper initialize the scrapper using the provided credentials for a single subscription.
// By default, the scrapper is initialized with the clients provided by the azure-go-sdk.
// clients can be overwritten by passing in option functions.
func NewScrapper(cred az.TokenCredential, sub string, opts ...OptionsFunc) (*Scrapper, error) {
	o := DefaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	rgc, err := o.resourceGroupClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	pc, err := o.providersClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	nc, err := o.virtualNetworkClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	desc, err := o.diskEncryptionSetClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	cc, err := o.clusterClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	npc, err := o.nodePoolClientFactory(sub, cred, nil)
	if err != nil {
		return nil, err
	}

	return &Scrapper{
		resourceGroupClient:     rgc,
		providersClient:         pc,
		networksClient:          nc,
		diskEncryptionSetClient: desc,
		clusterClient:           cc,
		nodePoolClient:          npc,
	}, nil
}

func (s *Scrapper) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.ListResourceGroups(ctx, consoleHandler[resource.ResourceGroup])
	})
	g.Go(func() error {
		return s.ListProviders(ctx, consoleHandler[resource.Provider])
	})
	g.Go(func() error {
		return s.ListVirtualNetworks(ctx, consoleHandler[network.VirtualNetwork])
	})
	g.Go(func() error {
		return s.ListDiskEncryptionSets(ctx, consoleHandler[compute.DiskEncryptionSet])
	})
	g.Go(func() error {
		return s.ListClusters(ctx, consoleHandler[container.ManagedCluster])
	})

	return g.Wait()
}

func processPage[T any](page []*T, err error, pageHandler pageHandler[T]) error {
	for _, v := range page {
		if err = pageHandler(v); err != nil {
			return fmt.Errorf("failed to process page: %w", err)
		}
	}
	return nil
}

func consoleHandler[T any](t *T) error {
	return json.NewEncoder(os.Stdout).Encode(t)
}
