import type {
    FilterKey,
    FilterKeyRequest,
    FilterOperator,
    FilterOperatorRequest,
    NotificationType,
    NotificationTypeRequest,
    FeedFilterSetType,
    FeedFilterSetTypeRequest,
    TorrentCategory,
    TorrentCategoryRequest,
    Feed,
    FeedRequest,
    Notifier,
    NotifierRequest,
    FeedFilter,
    FeedFilterRequest,
    FeedAuthorFilter,
    FeedAuthorFilterRequest,
    FeedFilterSet,
    FeedFilterSetRequest,
    FeedFilterSetEntry,
    FeedFilterSetEntryRequest,
    Torrent,
    // Audiobook Wizard types
    BookMetadata,
    Library,
    TorrentImportInfo,
    SearchASINRequest,
    PreviewDirectoryRequest,
    PreviewDirectoryResponse,
    ExecuteImportRequest,
    ExecuteImportResponse
} from '@/types/api'

const BASE_URL = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${BASE_URL}${path}`, {
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers
        },
        ...options
    })

    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Request failed' }))
        throw new Error(error.error || 'Request failed')
    }

    if (response.status === 204) {
        return {} as T
    }

    return response.json()
}

export const api = {
    // Filter Keys
    filterKeys: {
        list: () => request<FilterKey[]>('/filter-keys'),
        get: (id: number) => request<FilterKey>(`/filter-keys/${id}`),
        create: (data: FilterKeyRequest) =>
            request<FilterKey>('/filter-keys', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FilterKeyRequest) =>
            request<FilterKey>(`/filter-keys/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/filter-keys/${id}`, { method: 'DELETE' })
    },

    // Filter Operators
    filterOperators: {
        list: () => request<FilterOperator[]>('/filter-operators'),
        get: (id: number) => request<FilterOperator>(`/filter-operators/${id}`),
        create: (data: FilterOperatorRequest) =>
            request<FilterOperator>('/filter-operators', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FilterOperatorRequest) =>
            request<FilterOperator>(`/filter-operators/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/filter-operators/${id}`, { method: 'DELETE' })
    },

    // Notification Types
    notificationTypes: {
        list: () => request<NotificationType[]>('/notification-types'),
        get: (id: number) => request<NotificationType>(`/notification-types/${id}`),
        create: (data: NotificationTypeRequest) =>
            request<NotificationType>('/notification-types', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: NotificationTypeRequest) =>
            request<NotificationType>(`/notification-types/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/notification-types/${id}`, { method: 'DELETE' })
    },

    // Feed Filter Set Types
    feedFilterSetTypes: {
        list: () => request<FeedFilterSetType[]>('/feed-filter-set-types'),
        get: (id: number) => request<FeedFilterSetType>(`/feed-filter-set-types/${id}`),
        create: (data: FeedFilterSetTypeRequest) =>
            request<FeedFilterSetType>('/feed-filter-set-types', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedFilterSetTypeRequest) =>
            request<FeedFilterSetType>(`/feed-filter-set-types/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feed-filter-set-types/${id}`, { method: 'DELETE' })
    },

    // Torrent Categories
    torrentCategories: {
        list: () => request<TorrentCategory[]>('/torrent-categories'),
        get: (id: number) => request<TorrentCategory>(`/torrent-categories/${id}`),
        create: (data: TorrentCategoryRequest) =>
            request<TorrentCategory>('/torrent-categories', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: TorrentCategoryRequest) =>
            request<TorrentCategory>(`/torrent-categories/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/torrent-categories/${id}`, { method: 'DELETE' })
    },

    // Feeds
    feeds: {
        list: () => request<Feed[]>('/feeds'),
        get: (id: number) => request<Feed>(`/feeds/${id}`),
        create: (data: FeedRequest) =>
            request<Feed>('/feeds', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedRequest) =>
            request<Feed>(`/feeds/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feeds/${id}`, { method: 'DELETE' })
    },

    // Notifiers
    notifiers: {
        list: () => request<Notifier[]>('/notifiers'),
        get: (id: number) => request<Notifier>(`/notifiers/${id}`),
        create: (data: NotifierRequest) =>
            request<Notifier>('/notifiers', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: NotifierRequest) =>
            request<Notifier>(`/notifiers/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/notifiers/${id}`, { method: 'DELETE' })
    },

    // Feed Filters
    feedFilters: {
        list: (feedId?: number) => {
            const params = feedId ? `?feed_id=${feedId}` : ''
            return request<FeedFilter[]>(`/feed-filters${params}`)
        },
        get: (id: number) => request<FeedFilter>(`/feed-filters/${id}`),
        create: (data: FeedFilterRequest) =>
            request<FeedFilter>('/feed-filters', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedFilterRequest) =>
            request<FeedFilter>(`/feed-filters/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feed-filters/${id}`, { method: 'DELETE' })
    },

    // Feed Author Filters
    feedAuthorFilters: {
        list: (feedId?: number) => {
            const params = feedId ? `?feed_id=${feedId}` : ''
            return request<FeedAuthorFilter[]>(`/feed-author-filters${params}`)
        },
        get: (id: number) => request<FeedAuthorFilter>(`/feed-author-filters/${id}`),
        create: (data: FeedAuthorFilterRequest) =>
            request<FeedAuthorFilter>('/feed-author-filters', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedAuthorFilterRequest) =>
            request<FeedAuthorFilter>(`/feed-author-filters/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feed-author-filters/${id}`, { method: 'DELETE' })
    },

    // Feed Filter Sets
    feedFilterSets: {
        list: (feedFilterId?: number) => {
            const params = feedFilterId ? `?feed_filter_id=${feedFilterId}` : ''
            return request<FeedFilterSet[]>(`/feed-filter-sets${params}`)
        },
        get: (id: number) => request<FeedFilterSet>(`/feed-filter-sets/${id}`),
        create: (data: FeedFilterSetRequest) =>
            request<FeedFilterSet>('/feed-filter-sets', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedFilterSetRequest) =>
            request<FeedFilterSet>(`/feed-filter-sets/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feed-filter-sets/${id}`, { method: 'DELETE' })
    },

    // Feed Filter Set Entries
    feedFilterSetEntries: {
        list: (feedFilterSetId?: number) => {
            const params = feedFilterSetId ? `?feed_filter_set_id=${feedFilterSetId}` : ''
            return request<FeedFilterSetEntry[]>(`/feed-filter-set-entries${params}`)
        },
        get: (id: number) => request<FeedFilterSetEntry>(`/feed-filter-set-entries/${id}`),
        create: (data: FeedFilterSetEntryRequest) =>
            request<FeedFilterSetEntry>('/feed-filter-set-entries', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: FeedFilterSetEntryRequest) =>
            request<FeedFilterSetEntry>(`/feed-filter-set-entries/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/feed-filter-set-entries/${id}`, { method: 'DELETE' })
    },

    // Torrents
    torrents: {
        unimported: () => request<Torrent[]>('/torrents/unimported'),
        manualIntervention: () => request<Torrent[]>('/torrents/manual'),
        changeCategory: (hash: string, category: string) => request<void>(`/torrents/${hash}/category`, {
            method: 'POST',
            body: JSON.stringify({ category: category })
        }),
        changeTags: (hash: string, tags: string) => request<void>(`/torrents/${hash}/tags`, {
            method: 'POST',
            body: JSON.stringify({ tags: tags })
        })
    },

    // Audiobook Wizard
    audiobookWizard: {
        getTorrentInfo: (hash: string) =>
            request<TorrentImportInfo>(`/audiobook-wizard/torrent/${hash}/info`),

        searchAsin: (data: SearchASINRequest) =>
            request<BookMetadata[]>('/audiobook-wizard/search-asin', {
                method: 'POST',
                body: JSON.stringify(data)
            }),

        getMetadataByAsin: (asin: string) =>
            request<BookMetadata>(`/audiobook-wizard/asin/${asin}/metadata`),

        previewDirectory: (data: PreviewDirectoryRequest) =>
            request<PreviewDirectoryResponse>('/audiobook-wizard/preview-directory', {
                method: 'POST',
                body: JSON.stringify(data)
            }),

        getLibraries: () =>
            request<Library[]>('/audiobook-wizard/libraries'),

        executeImport: (data: ExecuteImportRequest) =>
            request<ExecuteImportResponse>('/audiobook-wizard/execute-import', {
                method: 'POST',
                body: JSON.stringify(data)
            })
    }
}
