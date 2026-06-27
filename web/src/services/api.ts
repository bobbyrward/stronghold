import type {
    NotificationType,
    TorrentCategory,
    SubscriptionScope,
    BookType,
    Library,
    LibraryRequest,
    Feed,
    FeedRequest,
    Notifier,
    NotifierRequest,
    Torrent,
    // Audiobook Wizard types
    BookMetadata,
    TorrentImportInfo,
    SearchASINRequest,
    PreviewDirectoryRequest,
    PreviewDirectoryResponse,
    ExecuteImportRequest,
    ExecuteImportResponse,
    // Feedwatcher2 types
    Author,
    AuthorRequest,
    AuthorAlias,
    AuthorAliasRequest,
    AuthorSubscription,
    AuthorSubscriptionRequest,
    AuthorSubscriptionItem,
    HardcoverAuthorSearchResult,
    PaginatedEventLogResponse,
    EventLog,
    VersionInfo
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
    // Notification Types (read-only reference data)
    notificationTypes: {
        list: () => request<NotificationType[]>('/notification-types'),
        get: (id: number) => request<NotificationType>(`/notification-types/${id}`)
    },

    // Torrent Categories (read-only reference data)
    torrentCategories: {
        list: () => request<TorrentCategory[]>('/torrent-categories'),
        get: (id: number) => request<TorrentCategory>(`/torrent-categories/${id}`)
    },

    // Subscription Scopes (read-only reference data)
    subscriptionScopes: {
        list: () => request<SubscriptionScope[]>('/subscription-scopes')
    },

    // Book Types (read-only reference data)
    bookTypes: {
        list: () => request<BookType[]>('/book-types'),
        get: (id: number) => request<BookType>(`/book-types/${id}`)
    },

    // Libraries
    libraries: {
        list: (bookTypeId?: number) => {
            const params = bookTypeId ? `?book_type_id=${bookTypeId}` : ''
            return request<Library[]>(`/libraries${params}`)
        },
        get: (id: number) => request<Library>(`/libraries/${id}`),
        create: (data: LibraryRequest) =>
            request<Library>('/libraries', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: LibraryRequest) =>
            request<Library>(`/libraries/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/libraries/${id}`, { method: 'DELETE' })
    },

    // Hardcover (external search)
    hardcover: {
        searchAuthors: (query: string) =>
            request<HardcoverAuthorSearchResult[]>(`/hardcover/authors/search?q=${encodeURIComponent(query)}`)
    },

    // Authors (feedwatcher2)
    authors: {
        list: (query?: string) => {
            const params = query ? `?q=${encodeURIComponent(query)}` : ''
            return request<Author[]>(`/authors${params}`)
        },
        get: (id: number) => request<Author>(`/authors/${id}`),
        create: (data: AuthorRequest) =>
            request<Author>('/authors', {
                method: 'POST',
                body: JSON.stringify(data)
            }),
        update: (id: number, data: AuthorRequest) =>
            request<Author>(`/authors/${id}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            }),
        delete: (id: number) =>
            request<void>(`/authors/${id}`, { method: 'DELETE' }),

        // Nested aliases
        aliases: {
            list: (authorId: number) =>
                request<AuthorAlias[]>(`/authors/${authorId}/aliases`),
            get: (authorId: number, id: number) =>
                request<AuthorAlias>(`/authors/${authorId}/aliases/${id}`),
            create: (authorId: number, data: AuthorAliasRequest) =>
                request<AuthorAlias>(`/authors/${authorId}/aliases`, {
                    method: 'POST',
                    body: JSON.stringify(data)
                }),
            update: (authorId: number, id: number, data: AuthorAliasRequest) =>
                request<AuthorAlias>(`/authors/${authorId}/aliases/${id}`, {
                    method: 'PUT',
                    body: JSON.stringify(data)
                }),
            delete: (authorId: number, id: number) =>
                request<void>(`/authors/${authorId}/aliases/${id}`, { method: 'DELETE' })
        },

        // Nested subscription (one per author)
        subscription: {
            get: (authorId: number) =>
                request<AuthorSubscription>(`/authors/${authorId}/subscription`),
            create: (authorId: number, data: AuthorSubscriptionRequest) =>
                request<AuthorSubscription>(`/authors/${authorId}/subscription`, {
                    method: 'POST',
                    body: JSON.stringify(data)
                }),
            update: (authorId: number, data: AuthorSubscriptionRequest) =>
                request<AuthorSubscription>(`/authors/${authorId}/subscription`, {
                    method: 'PUT',
                    body: JSON.stringify(data)
                }),
            delete: (authorId: number) =>
                request<void>(`/authors/${authorId}/subscription`, { method: 'DELETE' }),
            items: (authorId: number) =>
                request<AuthorSubscriptionItem[]>(`/authors/${authorId}/subscription/items`)
        }
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
    },

    // Event Logs (read-only, paginated)
    eventLogs: {
        list: (params: Record<string, string>) => {
            const query = new URLSearchParams(params).toString()
            return request<PaginatedEventLogResponse>(`/event-logs?${query}`)
        },
        get: (id: number) => request<EventLog>(`/event-logs/${id}`)
    },

    // Version info
    version: {
        get: () => request<VersionInfo>('/version')
    }
}
