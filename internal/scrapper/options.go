package scrapper

// Options holds factory methods used to create clients that are passed to the scrapper.
type Options struct {
	resourceGroupClientFactory     ResourceGroupClientFactory
	providersClientFactory         ProvidersClientFactory
	virtualNetworkClientFactory    VirtualNetworkClientFactory
	diskEncryptionSetClientFactory DiskEncryptionSetClientFactory
	clusterClientFactory           ClusterClientFactory
	nodePoolClientFactory          NodePoolClientFactory
}

// DefaultOptions initialize scrapper to user the default client factories from the azure-go-sdk.
func DefaultOptions() *Options {
	return &Options{
		resourceGroupClientFactory:     defaultResourceGroupFactory,
		providersClientFactory:         defaultProvidersClientFactory,
		virtualNetworkClientFactory:    defaultNetworkClientFactory,
		diskEncryptionSetClientFactory: defaultDiskEncryptionSetClientFactory,
		clusterClientFactory:           defaultClusterClientFactory,
		nodePoolClientFactory:          defaultNodePoolClientFactory,
	}
}

type OptionsFunc func(opt *Options)

func WithResourceGroupsFactory(f ResourceGroupClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.resourceGroupClientFactory = f
	}
}

func WithProvidersFactory(f ProvidersClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.providersClientFactory = f
	}
}

func WithVirtualNetworksFactory(f VirtualNetworkClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.virtualNetworkClientFactory = f
	}
}

func WithDiskEncryptionSetFactory(f DiskEncryptionSetClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.diskEncryptionSetClientFactory = f
	}
}

func WithClusterFactory(f ClusterClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.clusterClientFactory = f
	}
}

func WithNodePoolFactory(f NodePoolClientFactory) OptionsFunc {
	return func(opt *Options) {
		opt.nodePoolClientFactory = f
	}
}
