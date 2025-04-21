package qbit

import (
	"strings"

	"github.com/autobrr/go-qbittorrent"

	"github.com/bobbyrward/stronghold/internal/config"
)

func CreateClient() (*qbittorrent.Client, error) {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     config.Config.BookImporter.QbitURL,
		Username: config.Config.BookImporter.QbitUsername,
		Password: config.Config.BookImporter.QbitPassword,
	})

	return client, nil
}

func TagsFromTagList(tagList string) map[string]bool {
	tags := make(map[string]bool)

	for _, tag := range strings.Split(tagList, ",") {
		tags[tag] = true
	}

	return tags
}
