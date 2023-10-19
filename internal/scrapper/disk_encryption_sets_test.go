package scrapper_test

import (
	. "azure-scrapper/internal/scrapper"
	"context"
	"errors"

	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScrapper_ListDiskEncryptionSets(t *testing.T) {
	tests := []struct {
		name                  string
		diskEncryptionFactory DiskEncryptionSetClientFactory
		handlerError          error
		expect                func(t *testing.T, err error)
	}{
		{
			name: "disk encryption set iteration succeeds",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{item: &compute.DiskEncryptionSetsClientListResponse{
					DiskEncryptionSetList: compute.DiskEncryptionSetList{
						Value: []*compute.DiskEncryptionSet{{}},
					},
				}}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "disk encryption handler fails",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return NewPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{item: &compute.DiskEncryptionSetsClientListResponse{
					DiskEncryptionSetList: compute.DiskEncryptionSetList{
						Value: []*compute.DiskEncryptionSet{{}},
					},
				}}, nil
			},
			handlerError: errors.New("failed to Handle disk encryption set"),
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "disk encryption set iteration fails",
			diskEncryptionFactory: func(sub string, cred az.TokenCredential, opts *arm.ClientOptions) (DiskEncryptionSetPager, error) {
				return FailPager[compute.DiskEncryptionSetsClientListOptions, compute.DiskEncryptionSetsClientListResponse]{}, nil
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := NewScrapper(testCred(), "not-important", WithDiskEncryptionSetFactory(tt.diskEncryptionFactory))
			require.NoError(t, err)
			tt.expect(t, scraper.ListDiskEncryptionSets(context.Background(), func(r *compute.DiskEncryptionSet) error { return tt.handlerError }))
		})
	}
}
