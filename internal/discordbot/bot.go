package discordbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/autobrr/go-qbittorrent"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

type Bot struct {
	session    *discordgo.Session
	config     *config.DiscordBotConfig
	bookSearch *booksearch.BookSearchService
	qbitClient *qbittorrent.Client
	db         *gorm.DB
}

func Run() error {
	cfg := config.Config.DiscordBot

	if cfg.Token == "" {
		return fmt.Errorf("discord bot token is required")
	}

	bot, err := NewBot(&cfg)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	return bot.Start()
}

func NewBot(cfg *config.DiscordBotConfig) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to create discord session"), err)
	}

	qbitClient, err := qbit.CreateClient()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to create qbittorrent client"), err)
	}

	db, err := models.ConnectDB()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to connect to database"), err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to automigrate database"), err)
	}

	bot := &Bot{
		session:    session,
		config:     cfg,
		bookSearch: booksearch.NewBookSearchService(),
		qbitClient: qbitClient,
		db:         db,
	}

	bot.setupHandlers()

	return bot, nil
}

func (b *Bot) setupHandlers() {
	b.session.AddHandler(b.ready)
	b.session.AddHandler(b.interactionCreate)
}

func (b *Bot) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord session: %w", err)
	}
	defer func() {
		if err := b.session.Close(); err != nil {
			slog.ErrorContext(ctx, "Failed to close discord session", slog.Any("error", err))
		}
	}()

	slog.InfoContext(ctx, "Discord bot is running. Press CTRL+C to exit.")

	<-ctx.Done()
	slog.InfoContext(ctx, "Shutting down Discord bot...")

	return nil
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	ctx := context.Background()
	slog.InfoContext(ctx, "Discord bot is ready", slog.String("username", event.User.Username))

	err := b.registerCommands()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to register commands", slog.Any("error", err))
	}
}

func (b *Bot) registerCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "requestbook",
			Description: "Search and request books to be added to qBittorrent",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Search query for books",
					Required:    true,
				},
			},
		},
	}

	for _, command := range commands {
		guildID := b.config.GuildID
		if guildID == "" {
			_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
			if err != nil {
				return fmt.Errorf("failed to register global command %s: %w", command.Name, err)
			}
		} else {
			_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, command)
			if err != nil {
				return fmt.Errorf("failed to register guild command %s: %w", command.Name, err)
			}
		}
	}

	ctx := context.Background()
	if b.config.GuildID == "" {
		slog.InfoContext(ctx, "Registered global commands")
	} else {
		slog.InfoContext(ctx, "Registered guild commands", slog.String("guildId", b.config.GuildID))
	}

	return nil
}
