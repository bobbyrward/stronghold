# stronghold

automation of all the things

## sub systems

* ~book importer~
* ~book feed watcher~
* audiobook import
* manga downloader
* book requests

### Book Importer

Regularly polls qbit to see if any torrents in the configured categories are both
finished and missing the correct "imported" tag.

### Feed Watcher

Regularly polls a list of RSS feeds for new items, and adds them to qbit if they
match any of the filters
