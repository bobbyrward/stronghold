package qbit

import (
	"context"
	"log/slog"
	"strings"

	"github.com/autobrr/go-qbittorrent"

	"github.com/bobbyrward/stronghold/internal/config"
)

func CreateClient() (*qbittorrent.Client, error) {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Creating qBittorrent client", 
		slog.String("host", config.Config.Qbit.URL),
		slog.String("username", config.Config.Qbit.Username))
	
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     config.Config.Qbit.URL,
		Username: config.Config.Qbit.Username,
		Password: config.Config.Qbit.Password,
	})

	slog.InfoContext(ctx, "Successfully created qBittorrent client")
	
	return client, nil
}

func TagsFromTagList(tagList string) map[string]bool {
	tags := make(map[string]bool)

	for _, tag := range strings.Split(tagList, ",") {
		tags[tag] = true
	}

	return tags
}
