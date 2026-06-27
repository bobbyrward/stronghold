package client

import "github.com/bobbyrward/stronghold/internal/www/api"

type (
	FeedRequest  = api.FeedRequest
	FeedResponse = api.FeedResponse
	FeedClient   = GenericClient[FeedRequest, FeedResponse]

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
	NotificationTypes    *NotificationTypeClient
	Notifiers            *NotifierClient
	TorrentCategories    *TorrentCategoryClient
	SubscriptionScopes   *SubscriptionScopeClient
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseUrl:              baseUrl,
		Feeds:                &FeedClient{BaseUrl: baseUrl, TypeName: "feeds"},
		NotificationTypes:    &NotificationTypeClient{BaseUrl: baseUrl, TypeName: "notification-types"},
		Notifiers:            &NotifierClient{BaseUrl: baseUrl, TypeName: "notifiers"},
		TorrentCategories:    &TorrentCategoryClient{BaseUrl: baseUrl, TypeName: "torrent-categories"},
		SubscriptionScopes:   &SubscriptionScopeClient{BaseUrl: baseUrl, TypeName: "subscription-scopes"},
	}
}
