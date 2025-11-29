# Book Database Schema Design

This document defines the database schema for the book downloading and library management system.

## Overview

A comprehensive system to track:

- **Downloads**: Active and completed downloads from torrents
- **Library Catalog**: Imported books with full metadata

Supports both **audiobooks** and **ebooks** with normalized relationships for authors, narrators, and series.

---

## Models

### Person

Base table for all people (can be author, narrator, or both).

```yaml
model: Person
fields:
  - name: Name
    type: string
    constraints: [not null, uniqueIndex]
  - name: SortName
    type: string
    description: "For sorting (e.g., 'King, Stephen')"
  - name: Description
    type: string
```

---

### Book

Main book entity representing a unique work (not a specific format).

```yaml
model: Book
fields:
  - name: Title
    type: string
    constraints: [not null, index]
  - name: Subtitle
    type: string
  - name: Description
    type: string
  - name: Language
    type: string
    default: "en"
  - name: PublishDate
    type: "*time.Time"
  - name: Publisher
    type: string
  - name: BookType
    type: string
    constraints: [not null]
    description: "audiobook, ebook, or both"
  - name: Duration
    type: int
    description: "Audiobook duration in minutes"
relationships:
  - name: Authors
    type: has_many
    model: BookAuthor
  - name: Narrators
    type: has_many
    model: BookNarrator
  - name: Series
    type: has_many
    model: BookSeries
  - name: Identifiers
    type: has_many
    model: BookIdentifier
  - name: Files
    type: has_many
    model: BookFile
```

---

### BookAuthor

Many-to-many relationship between Book and Person for authors.

```yaml
model: BookAuthor
fields:
  - name: BookID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_author"]
  - name: PersonID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_author"]
  - name: Book
    type: belongs_to
    model: Book
  - name: Person
    type: belongs_to
    model: Person
```

---

### BookNarrator

Many-to-many relationship between Book and Person for narrators.

```yaml
model: BookNarrator
fields:
  - name: BookID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_narrator"]
  - name: PersonID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_narrator"]
  - name: Book
    type: belongs_to
    model: Book
  - name: Person
    type: belongs_to
    model: Person
```

---

### Series

Series information for book collections.

```yaml
model: Series
fields:
  - name: Name
    type: string
    constraints: [not null, uniqueIndex]
  - name: Description
    type: string
```

---

### BookSeries

Many-to-many relationship between Book and Series with position.

```yaml
model: BookSeries
fields:
  - name: BookID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_series"]
  - name: SeriesID
    type: uint
    constraints: [not null, "uniqueIndex:idx_book_series"]
  - name: Position
    type: float64
    description: "Position in series (float for novellas like 1.5, 2.5)"
relationships:
  - name: Book
    type: belongs_to
    model: Book
  - name: Series
    type: belongs_to
    model: Series
```

---

### BookIdentifier

External identifiers for books (ISBN, ASIN, etc.).

```yaml
model: BookIdentifier
fields:
  - name: BookID
    type: uint
    constraints: [not null, index]
  - name: Type
    type: string
    constraints: [not null]
    description: "isbn, isbn13, asin, goodreads, librarything, audible"
  - name: Value
    type: string
    constraints: [not null]
relationships:
  - name: Book
    type: belongs_to
    model: Book
indexes:
  - name: idx_book_identifier_type_value
    fields: [BookID, Type]
    unique: true
```

---

### BookFile

Physical files in the library.

```yaml
model: BookFile
fields:
  - name: BookID
    type: uint
    constraints: [not null, index]
  - name: FilePath
    type: string
    constraints: [not null]
    description: "Full path to file"
  - name: FileName
    type: string
    constraints: [not null]
  - name: FileType
    type: string
    constraints: [not null]
    description: "epub, mobi, azw3, m4b, mp3"
  - name: FileSize
    type: int64
  - name: Checksum
    type: string
    description: "SHA256 for deduplication"
  - name: AddedAt
    type: time.Time
relationships:
  - name: Book
    type: belongs_to
    model: Book
```

---

### Download

Tracks active and completed downloads.

```yaml
model: Download
fields:
  # Source information
  - name: TorrentHash
    type: string
    constraints: [uniqueIndex]
  - name: TorrentName
    type: string
  - name: Category
    type: string
    description: "qBittorrent category"

  # Links
  - name: BookID
    type: "*uint"
    description: "Link to book once identified"

  # Status
  - name: Status
    type: string
    constraints: [not null]
    default: "downloading"
    description: "downloading, completed, importing, imported, failed, manual_intervention"

  - name: ImportedAt
    type: "*time.Time"

  # Error tracking
  - name: ErrorMessage
    type: string
  - name: RetryCount
    type: int
relationships:
  - name: SourceFeed
    type: belongs_to
    model: Feed
    optional: true
  - name: Book
    type: belongs_to
    model: Book
    optional: true
```

---

### ImportHistory

Audit trail of imports.

```yaml
model: ImportHistory
fields:
  - name: DownloadID
    type: uint
    constraints: [not null]
  - name: BookID
    type: uint
    constraints: [not null]
  - name: BookFileID
    type: uint
    constraints: [not null]
  - name: SourcePath
    type: string
  - name: DestPath
    type: string
  - name: ImportedAt
    type: time.Time
relationships:
  - name: Download
    type: belongs_to
    model: Download
  - name: Book
    type: belongs_to
    model: Book
  - name: BookFile
    type: belongs_to
    model: BookFile
```

---

## Design Notes

### Person/Author/Narrator Pattern

- A single `Person` can be both an author and narrator (e.g., memoir with author narration)
- The `BookAuthor` and `BookNarrator` tables allow tracking roles and ordering

### Series Position as Float

- Using `float64` for series position allows handling novellas (1.5, 2.5)
- Standard books use whole numbers (1.0, 2.0, 3.0)

### BookType

- Set on `Book` model to indicate what formats exist
- Values: `audiobook`, `ebook`, `both`
- A book marked as `both` can have both audiobook and ebook files

### Download Lifecycle

1. `downloading` - Torrent active in qBittorrent
2. `completed` - Download finished, awaiting import
3. `importing` - Import process running
4. `imported` - Successfully imported to library
5. `failed` - Import failed (check ErrorMessage)
6. `manual_intervention` - Needs manual review

### File Checksums

- `Checksum` field uses SHA256
- Enables duplicate detection across imports
- Can skip re-importing identical files

---

## Indexes

Key indexes for performance:

- `Book.Title` - Searching by title
- `Person.Name` - Searching by author/narrator name
- `BookIdentifier(BookID, Type)` - Unique constraint on identifier per book
- `Download.TorrentHash` - Quick lookup by hash
- `Download.Status` - Filtering active downloads

---

## Migration Notes

When implementing, add to `AutoMigrate()` in this order:

1. Person
2. Series
3. Book
4. BookAuthor
5. BookNarrator
6. BookSeries
7. BookIdentifier
8. BookFile
9. Download
10. ImportHistory

---

## Future Considerations

- **Tags/Genres**: Add BookTag model for categorization
- **Collections**: User-defined book collections
- **Reading Progress**: Track audiobook/ebook progress
- **Ratings/Reviews**: User ratings and notes
- **Cover Images**: Store cover art paths
- **Multiple Libraries**: Support multiple destination libraries per format
