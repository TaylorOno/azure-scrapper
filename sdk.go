package main

import (
	az "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	rt "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	container "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	resource "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// ResourceGroupsPager Interfaces for query functions use by our applications.
// These interfaces are satisfied by the azure clients return from factory methods.
type ResourceGroupsPager interface {
	NewListPager(options *resource.ResourceGroupsClientListOptions) *rt.Pager[resource.ResourceGroupsClientListResponse]
}

type ProvidersPager interface {
	NewListPager(options *resource.ProvidersClientListOptions) *rt.Pager[resource.ProvidersClientListResponse]
}

type VirtualNetworkPager interface {
	NewListAllPager(options *network.VirtualNetworksClientListAllOptions) *rt.Pager[network.VirtualNetworksClientListAllResponse]
}

type DiskEncryptionSetPager interface {
	NewListPager(options *compute.DiskEncryptionSetsClientListOptions) *rt.Pager[compute.DiskEncryptionSetsClientListResponse]
}

type ClusterPager interface {
	NewListPager(options *container.ManagedClustersClientListOptions) *rt.Pager[container.ManagedClustersClientListResponse]
}

type NodePoolPager interface {
	NewListPager(resourceGroupName string, resourceName string, options *container.AgentPoolsClientListOptions) *rt.Pager[container.AgentPoolsClientListResponse]
}

// Factory signatures from the azure-sdk-for-go used by our applications
type ResourceGroupClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ResourceGroupsPager, error)
type ProvidersClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ProvidersPager, error)
type VirtualNetworkClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (VirtualNetworkPager, error)
type DiskEncryptionSetClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (DiskEncryptionSetPager, error)
type ClusterClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ClusterPager, error)
type NodePoolClientFactory func(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (NodePoolPager, error)

// Default factories redirect to azure-sdk-for-go factories to meet interface requirements
func defaultResourceGroupFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ResourceGroupsPager, error) {
	return resource.NewResourceGroupsClient(subscriptionID, credential, options)
}

func defaultProvidersClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ProvidersPager, error) {
	return resource.NewProvidersClient(subscriptionID, credential, options)
}

func defaultNetworkClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (VirtualNetworkPager, error) {
	return network.NewVirtualNetworksClient(subscriptionID, credential, options)
}

func defaultDiskEncryptionSetClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (DiskEncryptionSetPager, error) {
	return compute.NewDiskEncryptionSetsClient(subscriptionID, credential, options)
}

func defaultClusterClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (ClusterPager, error) {
	return container.NewManagedClustersClient(subscriptionID, credential, options)
}

func defaultNodePoolClientFactory(subscriptionID string, credential az.TokenCredential, options *arm.ClientOptions) (NodePoolPager, error) {
	return container.NewAgentPoolsClient(subscriptionID, credential, options)
}
