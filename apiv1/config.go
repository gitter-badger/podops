package apiv1

import (
	"github.com/fupas/commons/pkg/env"
)

const (
	defaultBucketProduction = "production.podops.dev"
	defaultBucketCDN        = "cdn.podops.dev"

	defaultPortalEndpoint  = "https://podops.dev"
	defaultAPIEndpoint     = "https://api.podops.dev"
	defaultCDNEndpoint     = "https://cdn.podops.dev"
	defaultStorageEndpoint = "https://storage.googleapis.com/cdn.podops.dev"
)

var (
	// DefaultPortalEndpoint points to the portal
	DefaultPortalEndpoint string = env.GetString("BASE_URL", defaultPortalEndpoint)

	// DefaultAPIEndpoint points to the API
	DefaultAPIEndpoint string = env.GetString("API_ENDPOINT", defaultAPIEndpoint)

	// DefaultCDNEndpoint points to the CDN
	DefaultCDNEndpoint string = env.GetString("CDN_URL", defaultCDNEndpoint)

	// BucketUpload is the canonical name of the upload bucket
	//BucketUpload string = env.GetString("BUCKET_UPLOAD", defaultBucketUpload)

	// BucketProduction is the canonical name of the production bucket
	BucketProduction string = env.GetString("BUCKET_PRODUCTION", defaultBucketProduction)

	// BucketCDN is the canonical name of the CDN bucket
	BucketCDN string = env.GetString("BUCKET_CDN", defaultBucketCDN)

	// StorageEndpoint is the direct link to assets in Google Storage
	StorageEndpoint string = env.GetString("STORAGE_ENDPOINT", defaultStorageEndpoint)
)
