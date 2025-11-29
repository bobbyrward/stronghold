package client

import "github.com/bobbyrward/stronghold/internal/www/api"

type (
	FeedRequest  = api.FeedRequest
	FeedResponse = api.FeedResponse
	FeedClient   = GenericClient[FeedRequest, FeedResponse]

	FeedAuthorFilterRequest  = api.FeedAuthorFilterRequest
	FeedAuthorFilterResponse = api.FeedAuthorFilterResponse
	FeedAuthorFilterClient   = GenericClient[FeedAuthorFilterRequest, FeedAuthorFilterResponse]

	FeedFilterSetEntryRequest  = api.FeedFilterSetEntryRequest
	FeedFilterSetEntryResponse = api.FeedFilterSetEntryResponse
	FeedFilterSetEntryClient   = GenericClient[FeedFilterSetEntryRequest, FeedFilterSetEntryResponse]

	FeedFilterSetTypeRequest  = api.FeedFilterSetTypeRequest
	FeedFilterSetTypeResponse = api.FeedFilterSetTypeResponse
	FeedFilterSetTypeClient   = GenericClient[FeedFilterSetTypeRequest, FeedFilterSetTypeResponse]

	FeedFilterSetRequest  = api.FeedFilterSetRequest
	FeedFilterSetResponse = api.FeedFilterSetResponse
	FeedFilterSetClient   = GenericClient[FeedFilterSetRequest, FeedFilterSetResponse]

	FeedFilterRequest  = api.FeedFilterRequest
	FeedFilterResponse = api.FeedFilterResponse
	FeedFilterClient   = GenericClient[FeedFilterRequest, FeedFilterResponse]

	FilterKeyRequest  = api.FilterKeyRequest
	FilterKeyResponse = api.FilterKeyResponse
	FilterKeyClient   = GenericClient[FilterKeyRequest, FilterKeyResponse]

	FilterOperatorRequest  = api.FilterOperatorRequest
	FilterOperatorResponse = api.FilterOperatorResponse
	FilterOperatorClient   = GenericClient[FilterOperatorRequest, FilterOperatorResponse]

	NotificationTypeRequest  = api.NotificationTypeRequest
	NotificationTypeResponse = api.NotificationTypeResponse
	NotificationTypeClient   = GenericClient[NotificationTypeRequest, NotificationTypeResponse]

	NotifierRequest  = api.NotifierRequest
	NotifierResponse = api.NotifierResponse
	NotifierClient   = GenericClient[NotifierRequest, NotifierResponse]

	TorrentCategoryRequest  = api.TorrentCategoryRequest
	TorrentCategoryResponse = api.TorrentCategoryResponse
	TorrentCategoryClient   = GenericClient[TorrentCategoryRequest, TorrentCategoryResponse]
)

type Client struct {
	BaseUrl string

	Feeds                *FeedClient
	FeedAuthorFilters    *FeedAuthorFilterClient
	FeedFilterSetEntries *FeedFilterSetEntryClient
	FeedFilterSetTypes   *FeedFilterSetTypeClient
	FeedFilterSets       *FeedFilterSetClient
	FeedFilters          *FeedFilterClient
	FilterKeys           *FilterKeyClient
	FilterOperators      *FilterOperatorClient
	NotificationTypes    *NotificationTypeClient
	Notifiers            *NotifierClient
	TorrentCategories    *TorrentCategoryClient
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseUrl:              baseUrl,
		Feeds:                &FeedClient{BaseUrl: baseUrl, TypeName: "feeds"},
		FeedAuthorFilters:    &FeedAuthorFilterClient{BaseUrl: baseUrl, TypeName: "feed-author-filters"},
		FeedFilterSetEntries: &FeedFilterSetEntryClient{BaseUrl: baseUrl, TypeName: "feed-filter-set-entries"},
		FeedFilterSetTypes:   &FeedFilterSetTypeClient{BaseUrl: baseUrl, TypeName: "feed-filter-set-types"},
		FeedFilterSets:       &FeedFilterSetClient{BaseUrl: baseUrl, TypeName: "feed-filter-sets"},
		FeedFilters:          &FeedFilterClient{BaseUrl: baseUrl, TypeName: "feed-filters"},
		FilterKeys:           &FilterKeyClient{BaseUrl: baseUrl, TypeName: "filter-keys"},
		FilterOperators:      &FilterOperatorClient{BaseUrl: baseUrl, TypeName: "filter-operators"},
		NotificationTypes:    &NotificationTypeClient{BaseUrl: baseUrl, TypeName: "notification-types"},
		Notifiers:            &NotifierClient{BaseUrl: baseUrl, TypeName: "notifiers"},
		TorrentCategories:    &TorrentCategoryClient{BaseUrl: baseUrl, TypeName: "torrent-categories"},
	}
}
