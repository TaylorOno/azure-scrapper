package scrapper

import (
	"context"
	"fmt"
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// DiskEncryptionSetPager used to scrape disk encryption set information
type DiskEncryptionSetPager interface {
	NewListPager(options *compute.DiskEncryptionSetsClientListOptions) *rt.Pager[compute.DiskEncryptionSetsClientListResponse]
}

type DiskEncryptionSetClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (DiskEncryptionSetPager, error)

func defaultDiskEncryptionSetClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (DiskEncryptionSetPager, error) {
	return compute.NewDiskEncryptionSetsClient(subscriptionID, credential, options)
}

func (s *Scrapper) ListDiskEncryptionSets(ctx context.Context, pageHandler pageHandler[compute.DiskEncryptionSet]) error {
	pager := s.diskEncryptionSetClient.NewListPager(nil)
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
