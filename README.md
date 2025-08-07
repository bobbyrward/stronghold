# stronghold

automation of all the things

## sub systems

* ~book importer~
* ~book feed watcher~
* ~discord bot~
* audiobook import
* manga downloader

### Book Importer

Regularly polls qbit to see if any torrents in the configured categories are both
finished and missing the correct "imported" tag.

### Feed Watcher

Regularly polls a list of RSS feeds for new items, and adds them to qbit if they
match any of the filters

### Discord Bot

Interactive Discord bot that provides book request functionality through slash commands.

**Features:**
- `/requestbook <query>` - Search for books and add them to qBittorrent
- Interactive book selection using Discord buttons
- Mock book search API with 10 sample books (Rothfuss, Herbert, Sanderson)
- Mock qBittorrent integration for adding torrents
