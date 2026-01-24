// Reference data types
export interface FilterKey {
    id: number
    name: string
}

export interface SubscriptionScope {
    id: number
    name: string
}

export interface FilterOperator {
    id: number
    name: string
}

export interface NotificationType {
    id: number
    name: string
}

export interface FeedFilterSetType {
    id: number
    name: string
}

export interface TorrentCategory {
    id: number
    name: string
    scope_id: number
    scope_name: string
    media_type: string
}

export interface BookType {
    id: number
    name: string
}

export interface Library {
    id: number
    name: string
    path: string
    book_type_id: number
    book_type_name: string
}

export interface LibraryRequest {
    name: string
    path: string
    book_type_name: string
}

// Hardcover search
export interface HardcoverAuthorSearchResult {
    slug: string
    name: string
}

// Feedwatcher2 types
export interface Author {
    id: number
    name: string
    hardcover_ref: string | null
}

export interface AuthorAlias {
    id: number
    author_id: number
    name: string
}

export interface AuthorSubscription {
    id: number
    author_id: number
    author_name: string
    scope_id: number
    scope_name: string
    notifier_id: number | null
    notifier_name: string | null
    ebook_library_id: number
    ebook_library_name: string
    audiobook_library_id: number
    audiobook_library_name: string
}

export interface AuthorSubscriptionItem {
    id: number
    author_subscription_id: number
    torrent_hash: string
    booksearch_id: string
    torrent_url: string
    title: string
    downloaded_at: string
}

// Main resource types
export interface Feed {
    id: number
    name: string
    url: string
}

export interface Notifier {
    id: number
    name: string
    type_id: number
    type_name: string
    url: string
}

export interface FeedFilter {
    id: number
    name: string
    feed_id: number
    feed_name: string
    category_id: number
    category_name: string
    notifier_id: number
    notifier_name: string
}

export interface FeedAuthorFilter {
    id: number
    author: string
    feed_id: number
    feed_name: string
    category_id: number
    category_name: string
    notifier_id: number
    notifier_name: string
}

export interface FeedFilterSet {
    id: number
    feed_filter_id: number
    type_id: number
    type_name: string
}

export interface FeedFilterSetEntry {
    id: number
    feed_filter_set_id: number
    key_id: number
    key_name: string
    operator_id: number
    operator_name: string
    value: string
}

// Request types (for create/update operations)
// Reference data types are read-only - no request types needed
// (FilterKey, FilterOperator, NotificationType, FeedFilterSetType, TorrentCategory)

export interface FeedRequest {
    name: string
    url: string
}

export interface NotifierRequest {
    name: string
    type_id: number
    url?: string
}

export interface FeedFilterRequest {
    name: string
    feed_id: number
    category_id: number
    notifier_id: number
}

export interface FeedAuthorFilterRequest {
    author: string
    feed_id: number
    category_id: number
    notifier_id: number
}

export interface FeedFilterSetRequest {
    feed_filter_id: number
    type_id: number
}

export interface FeedFilterSetEntryRequest {
    feed_filter_set_id: number
    key_id: number
    operator_id: number
    value: string
}

export interface AuthorRequest {
    name: string
    hardcover_ref?: string | null
}

export interface AuthorAliasRequest {
    name: string
}

export interface AuthorSubscriptionRequest {
    scope_name: string
    notifier_id?: number | null
    ebook_library_name: string
    audiobook_library_name: string
}

export interface Torrent {
    hash: string
    name: string
    category: string
    state: string
    tags: string
}

export interface TorrentChangeCategoryRequest {
    category: string
}

export interface TorrentChangeTagsRequest {
    tags: string
}

// Audiobook Wizard API types
export interface Person {
    name: string
    asin?: string
}

export interface Series {
    name: string
    asin?: string
    position?: string
}

export interface Genre {
    name: string
    asin: string
    type: string
}

export interface BookMetadata {
    asin: string
    title: string
    subtitle?: string
    authors: Person[]
    narrators: Person[]
    description: string
    summary: string
    publisherName: string
    releaseDate: string
    runtimeLengthMin: number
    language: string
    isbn?: string
    rating: string
    genres: Genre[]
    seriesPrimary?: Series
    seriesSecondary?: Series
    image?: string
    copyright?: number
    isAdult?: boolean
    literatureType?: string
    region: string
    formatType: string
}

export interface TorrentImportInfo {
    hash: string
    name: string
    category: string
    tags: string
    asin?: string
    title?: string
    author?: string
    suggested_library?: string
    local_path: string
}

export interface SearchASINRequest {
    title: string
    author?: string
}

export interface PreviewDirectoryRequest {
    metadata: BookMetadata
}

export interface PreviewDirectoryResponse {
    directory_name: string
}

export interface ExecuteImportRequest {
    hash: string
    metadata: BookMetadata
    library_name: string
}

export interface ExecuteImportResponse {
    success: boolean
    destination_path: string
    message?: string
}
