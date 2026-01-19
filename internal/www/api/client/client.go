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

	FeedFilterSetTypeResponse = api.FeedFilterSetTypeResponse
	FeedFilterSetTypeClient   = ReadOnlyGenericClient[FeedFilterSetTypeResponse]

	FeedFilterSetRequest  = api.FeedFilterSetRequest
	FeedFilterSetResponse = api.FeedFilterSetResponse
	FeedFilterSetClient   = GenericClient[FeedFilterSetRequest, FeedFilterSetResponse]

	FeedFilterRequest  = api.FeedFilterRequest
	FeedFilterResponse = api.FeedFilterResponse
	FeedFilterClient   = GenericClient[FeedFilterRequest, FeedFilterResponse]

	FilterKeyResponse = api.FilterKeyResponse
	FilterKeyClient   = ReadOnlyGenericClient[FilterKeyResponse]

	FilterOperatorResponse = api.FilterOperatorResponse
	FilterOperatorClient   = ReadOnlyGenericClient[FilterOperatorResponse]

	NotificationTypeResponse = api.NotificationTypeResponse
	NotificationTypeClient   = ReadOnlyGenericClient[NotificationTypeResponse]

	NotifierRequest  = api.NotifierRequest
	NotifierResponse = api.NotifierResponse
	NotifierClient   = GenericClient[NotifierRequest, NotifierResponse]

	TorrentCategoryResponse = api.TorrentCategoryResponse
	TorrentCategoryClient   = ReadOnlyGenericClient[TorrentCategoryResponse]

	SubscriptionScopeResponse = api.SubscriptionScopeResponse
	SubscriptionScopeClient   = ReadOnlyGenericClient[SubscriptionScopeResponse]
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
	SubscriptionScopes   *SubscriptionScopeClient
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
		SubscriptionScopes:   &SubscriptionScopeClient{BaseUrl: baseUrl, TypeName: "subscription-scopes"},
	}
}
