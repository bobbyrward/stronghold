# Stronghold

**Stronghold** is an automated feed monitoring and media management system designed to integrate with qBittorrent. It monitors RSS/torrent feeds, applies configurable filters, automatically imports media files, and sends notifications through various channels.  In reality, this is all just an excuse to test out Claude AI's code generation capabilities on a large project.

## Overview

Stronghold automates the workflow of monitoring content feeds, downloading torrents, organizing media files, and managing your media library. It's particularly useful for managing audiobook and ebook collections with metadata enrichment from Audible and automatic library organization.

## Key Features

- 📡 **RSS Feed Monitoring** - Monitor multiple RSS/torrent feeds with custom filters
- 📚 **Automated Media Import** - Automatic import of audiobooks and ebooks from qBittorrent
- 🎯 **Smart Filtering** - Complex filter sets with AND/OR logic for precise content matching
- 🔍 **Metadata Enrichment** - Automatic metadata lookup from Audible for audiobooks
- 🔔 **Multi-Channel Notifications** - Discord webhook support with extensible notification system
- 🌐 **REST API** - Full-featured API for programmatic control
- 🖥️ **Web UI** - Modern Vue.js interface for management and manual imports
- 🤖 **Discord Bot** - Interactive book search and information lookup

## TODO

- Migrate to database backed configuration and state
- Use hardcover for metadata
- Better feed watching
- More author based feed watching and ui to support it
- Managing manga

## Architecture

```
┌─────────────────┐
│   Web UI (Vue)  │
└────────┬────────┘
         │
         ↓
┌─────────────────┐      ┌──────────────────┐
│   API Server    │◄─────┤  PostgreSQL DB   │
│   (Echo/GORM)   │      └──────────────────┘
└────────┬────────┘
         │
    ┌────┴────┬──────────┬────────────┬──────────┐
    ↓         ↓          ↓            ↓          ↓
┌─────────┐ ┌──────┐ ┌────────┐ ┌─────────┐ ┌────────┐
│  Feed   │ │ Book │ │Audiobook│ │ Discord │ │ Book   │
│ Watcher │ │Import│ │ Import  │ │   Bot   │ │ Search │
└────┬────┘ └──┬───┘ └────┬────┘ └────┬────┘ └────────┘
     │         │          │           │
     ↓         ↓          ↓           ↓
┌────────────────────────────────────────┐
│           qBittorrent                  │
└────────────────────────────────────────┘
```

## Subsystems

### 1. Feed Watcher

Monitors RSS/torrent feeds at configurable intervals and automatically downloads matching torrents.

**Features:**

- Multiple feed support with per-feed configuration
- Complex filter sets with AND/OR logic
- Author-based filtering
- Torrent category assignment
- Download client integration (qBittorrent)
- Discord notifications on match

**Use Cases:**

- Monitor audiobook feeds for specific authors
- Track ebook releases matching title patterns
- Auto-download content meeting specific criteria

### 2. Audiobook Importer

Automatically imports audiobook torrents from qBittorrent with metadata enrichment.

**Features:**

- Automatic metadata extraction from M4B/MP3 files
- ASIN lookup and Audible integration
- OPF metadata file generation
- Series detection and organization
- Directory naming from metadata templates
- Manual intervention workflow for edge cases

**Workflow:**

1. Detects completed audiobook torrents
2. Extracts metadata from files (ASIN, title, author)
3. Enriches metadata from Audible API
4. Generates directory structure and OPF files
5. Tags torrents as imported
6. Sends completion notifications

### 3. Book Importer

Imports ebook files (EPUB, MOBI, AZW3) from qBittorrent torrents.

**Features:**

- Multiple format support (EPUB, MOBI, AZW3)
- Configurable library destinations
- Hard-link or copy fallback
- Discord notifications
- Manual intervention tagging

**Workflow:**

1. Scans completed torrents for ebook files
2. Copies/links files to configured library
3. Tags torrents as imported
4. Notifies via Discord

### 4. API Server

RESTful API server providing programmatic access to all Stronghold features.

**Endpoints:**

- **Feeds** - Manage RSS feed sources
- **Feed Filters** - Configure content filters
- **Feed Author Filters** - Author-based filtering
- **Feed Filter Sets** - Complex filter combinations
- **Notifiers** - Notification channel configuration
- **Torrents** - View unimported and manual intervention torrents
- **Audiobook Wizard** - Manual import workflow for audiobooks

**Technology Stack:**

- Echo v4 (HTTP framework)
- GORM (ORM)
- PostgreSQL (production) / SQLite (testing)
- Structured logging (slog)

See [API Documentation](docs/openapi.yaml) for complete OpenAPI specification.

### 5. Web UI

Modern single-page application for managing Stronghold.

**Features:**

- Feed and filter management
- Torrent monitoring
- Audiobook import wizard
- Real-time updates
- Responsive design

**Technology:**

- Vue.js 3
- Vue Router
- TypeScript
- Tailwind CSS (assumed from modern Vue setup)

### 6. Discord Bot

Interactive Discord bot for book search and information lookup.

**Commands:**

- Book search by title/author
- Series information
- Metadata display

### 7. Book Search Service

Standalone service for searching and retrieving book metadata.

**Features:**

- Audible integration
- Metadata caching
- REST API endpoint

## License

[MIT](https://opensource.org/license/mit)

## Acknowledgments

- Built with [Echo](https://echo.labstack.com/) web framework
- ORM by [GORM](https://gorm.io/)
- qBittorrent integration via [go-qbittorrent](https://github.com/autobrr/go-qbittorrent)
- RSS parsing by [gofeed](https://github.com/mmcdole/gofeed)
- Frontend powered by [Vue.js](https://vuejs.org/)
