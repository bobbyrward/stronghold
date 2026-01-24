package discordbot

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bwmarrin/discordgo"
	"github.com/cappuccinotm/slogx"
)

const (
	ChannelID_BobbysBookRequests = "1396454552961814608"
	ChannelID_BookRequests       = "1234593906209980478"
)

func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		b.handleSlashCommand(s, i)
	case discordgo.InteractionMessageComponent:
		b.handleComponentInteraction(s, i)
	}
}

func (b *Bot) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "requestbook":
		b.handleRequestBookCommand(s, i)
	}
}

func (b *Bot) handleRequestBookCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondWithError(s, i, "No search query provided")
		return
	}

	query := options[0].StringValue()
	if query == "" {
		b.respondWithError(s, i, "Search query cannot be empty")
		return
	}

	ctx := context.Background()
	slog.InfoContext(ctx, "Processing book request", slog.String("query", query), slog.String("userId", i.Member.User.ID))

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to defer interaction response", slog.Any("error", err))
		return
	}

	params := &booksearch.SearchParameters{
		Query:      query,
		MaxResults: 5,
	}

	searchResponse, err := b.bookSearch.Search(context.Background(), b.db, params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to search books", slog.Any("error", err), slog.String("query", query))
		b.editResponseWithError(s, i, "Failed to search for books")
		return
	}

	if len(searchResponse.Data) == 0 {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("No books found for query: **%s**", query),
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to send no results message", slog.Any("error", err))
		}
		return
	}

	dbResults, err := b.cacheSearchResults(searchResponse)
	if err != nil {
		b.editResponseWithError(s, i, "Failed to cache search results")
		return
	}

	b.sendBookSelectionMessage(s, i, dbResults, params)
}

func (b *Bot) cacheSearchResults(searchResponse *booksearch.SearchResponse) ([]models.SearchResponseItem, error) {
	ctx := context.Background()
	modelItems := make([]models.SearchResponseItem, len(searchResponse.Data))

	for idx, book := range searchResponse.Data {
		modelItems[idx] = book.ToModel()
	}

	result := b.db.Create(&modelItems)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to cache search results", slogx.Error(result.Error))
		return modelItems, result.Error
	}

	return modelItems, nil
}

func displayTitle(result *models.SearchResponseItem) string {
	category := ""

	switch result.MainCategory {
	case booksearch.MainCategoryEbooks:
		category = "ebook"
	case booksearch.MainCategoryAudiobooks:
		category = "audiobook"
	default:
		category = "unknown category"
	}

	return fmt.Sprintf("%s (%s)", result.Title, category)
}

func displayString(result *models.SearchResponseItem) string {
	parts := []string{}

	parts = append(parts, fmt.Sprintf("Authors: %s", result.Authors))

	if len(result.Series) > 0 {
		parts = append(parts, fmt.Sprintf("Series: %s", result.Series))
	}

	parts = append(parts, fmt.Sprintf("Type: %s", result.FileTypes))
	parts = append(parts, fmt.Sprintf("Tags: %s", result.Tags))

	return strings.Join(parts, "\n")
}

func (b *Bot) sendBookSelectionMessage(s *discordgo.Session, i *discordgo.InteractionCreate, searchResults []models.SearchResponseItem, searchParams *booksearch.SearchParameters) {
	ctx := context.Background()

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Books found for: %s", searchParams.Query),
		Description: "Select one or more books to add to qBittorrent:",
		Color:       0x00ff00,
		Fields:      make([]*discordgo.MessageEmbedField, len(searchResults)),
	}

	for idx, book := range searchResults {
		embed.Fields[idx] = &discordgo.MessageEmbedField{
			Name:   displayTitle(&book),
			Value:  displayString(&book),
			Inline: false,
		}
	}

	components := b.createBookSelectionComponents(searchResults)

	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send book selection message", slog.Any("error", err))
		b.editResponseWithError(s, i, "Failed to display book selection")
	}
}

func (b *Bot) createBookSelectionComponents(searchResults []models.SearchResponseItem) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	for i := 0; i < len(searchResults); i += 5 {
		end := i + 5
		if end > len(searchResults) {
			end = len(searchResults)
		}

		var buttons []discordgo.MessageComponent
		for j := i; j < end; j++ {
			book := searchResults[j]
			buttons = append(buttons, discordgo.Button{
				Label:    fmt.Sprintf("%d", j+1),
				Style:    discordgo.SecondaryButton,
				CustomID: fmt.Sprintf("select_book_%d", book.ID),
			})
		}

		components = append(components, discordgo.ActionsRow{
			Components: buttons,
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Add Selected Books",
				Style:    discordgo.PrimaryButton,
				CustomID: "add_selected_books",
				Disabled: true,
			},
			discordgo.Button{
				Label:    "Cancel",
				Style:    discordgo.DangerButton,
				CustomID: "cancel_selection",
			},
		},
	})

	return components
}

func (b *Bot) handleComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	customID := i.MessageComponentData().CustomID

	slog.InfoContext(ctx, "")

	switch {
	case strings.HasPrefix(customID, "select_book_"):
		b.handleBookSelection(s, i)
	case customID == "add_selected_books":
		b.handleAddSelectedBooks(s, i)
	case customID == "cancel_selection":
		b.handleCancelSelection(s, i)
	}
}

func (b *Bot) handleBookSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	bookID := strings.TrimPrefix(i.MessageComponentData().CustomID, "select_book_")

	components := i.Message.Components
	var selectedBooks []string

	for _, component := range components {
		if component.Type() == discordgo.ActionsRowComponent {

			actionRow, ok := component.(*discordgo.ActionsRow)
			compType := reflect.TypeOf(component)

			slog.InfoContext(
				ctx,
				"Type() says actions row",
				slog.Bool("cast_result", ok),
				slog.String("typeName", compType.Name()),
				slog.Any("comp", component),
			)

			if ok {
				for compIdx, comp := range actionRow.Components {
					if button, ok := comp.(*discordgo.Button); ok {
						action := ""

						if strings.HasPrefix(button.CustomID, "select_book_") {
							slog.InfoContext(
								ctx,
								"HasPrefix",
								slog.String("customID", button.CustomID),
							)

							if button.CustomID == "select_book_"+bookID {
								if button.Style == discordgo.SecondaryButton {
									button.Style = discordgo.SuccessButton
									selectedBooks = append(selectedBooks, bookID)
									action = "Selecting"
								} else {
									button.Style = discordgo.SecondaryButton
									action = "Deselecting"
								}
								actionRow.Components[compIdx] = button
							} else if button.Style == discordgo.SuccessButton {
								selectedBooks = append(selectedBooks, strings.TrimPrefix(button.CustomID, "select_book_"))
								action = "wtf"
							}
						} else if button.CustomID == "add_selected_books" {
							button.Disabled = len(selectedBooks) == 0
							actionRow.Components[compIdx] = button
							action = "add selected"
						}

						slog.InfoContext(
							ctx,
							"book selection",
							slog.Int("handle button", compIdx),
							slog.String("label", button.Label),
							slog.Uint64("style", uint64(button.Style)),
							slog.Bool("disabled", button.Disabled),
							slog.String("customID", button.CustomID),
							slog.String("action", action),
						)

					} else {
						slog.InfoContext(
							ctx,
							"button not found?",
							slog.Int("handle button", compIdx),
							slog.String("label", button.Label),
							slog.Uint64("style", uint64(button.Style)),
							slog.Bool("disabled", button.Disabled),
							slog.String("customID", button.CustomID),
						)
					}
				}
			} else {
				slog.InfoContext(
					ctx,
					"Not action row",
					slog.Uint64("componentType", uint64(component.Type())),
				)
			}
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     i.Message.Embeds,
			Components: components,
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update button selection", slog.Any("error", err))
	}
}

func (b *Bot) handleAddSelectedBooks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to defer message update", slog.Any("error", err))
		return
	}

	var selectedBookIDs []int
	components := i.Message.Components

	for _, component := range components {
		if actionRow, ok := component.(*discordgo.ActionsRow); ok {
			for _, comp := range actionRow.Components {
				if button, ok := comp.(*discordgo.Button); ok {
					if strings.HasPrefix(button.CustomID, "select_book_") && button.Style == discordgo.SuccessButton {
						idString := strings.TrimPrefix(button.CustomID, "select_book_")
						id, _ := strconv.Atoi(idString)
						selectedBookIDs = append(selectedBookIDs, id)
					}
				}
			}
		}
	}

	if len(selectedBookIDs) == 0 {
		slog.WarnContext(ctx, "No books selected for adding")
		return
	}

	var successCount int
	var failedBooks []string

	for _, bookID := range selectedBookIDs {
		var book models.SearchResponseItem

		result := b.db.First(&book, bookID)

		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to find cached search result",
				slog.Int("id", bookID),
			)
		}

		torrentURL := fmt.Sprintf("%s/tor/download.php/%s", config.Config.BookSearch.BaseURL, book.DlHash)

		qbitCategory := ""
		switch book.MainCategory {
		case booksearch.MainCategoryEbooks:
			if i.ChannelID == ChannelID_BobbysBookRequests {
				qbitCategory = "personal-books"
			} else {
				qbitCategory = "general-books"
			}
		case booksearch.MainCategoryAudiobooks:
			if i.ChannelID == ChannelID_BobbysBookRequests {
				qbitCategory = "audiobooks"
			} else {
				qbitCategory = "general-audiobooks"
			}
		default:
			qbitCategory = "unknown"
		}

		err = b.qbitClient.AddTorrentFromUrlCtx(
			ctx,
			torrentURL,
			map[string]string{
				"autoTMM":  "true",
				"category": qbitCategory,
			},
		)

		if err != nil {
			slog.ErrorContext(ctx, "Failed to add torrent to qBittorrent",
				slog.Any("error", err),
				slog.String("dlHash", book.DlHash),
				slog.String("torrentURL", torrentURL),
			)
			failedBooks = append(failedBooks, book.Title)
		} else {
			successCount++
			slog.InfoContext(ctx, "Successfully added torrent to qBittorrent",
				slog.String("dlHash", book.DlHash),
				slog.String("torrentURL", torrentURL),
				slog.String("channelid", i.ChannelID),
			)
		}
	}

	var embed *discordgo.MessageEmbed
	if successCount > 0 && len(failedBooks) == 0 {
		embed = &discordgo.MessageEmbed{
			Title:       "Books Added Successfully! 🎉",
			Description: fmt.Sprintf("Successfully added %d book(s) to qBittorrent", successCount),
			Color:       0x00ff00,
		}
	} else if successCount > 0 && len(failedBooks) > 0 {
		embed = &discordgo.MessageEmbed{
			Title:       "Partially Successful ⚠️",
			Description: fmt.Sprintf("Added %d book(s) successfully, failed to add %d book(s)", successCount, len(failedBooks)),
			Color:       0xffaa00,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       "Failed to Add Books ❌",
			Description: "Failed to add any books to qBittorrent. Check logs for details.",
			Color:       0xff0000,
		}
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update success message", slog.Any("error", err))
	}
}

func (b *Bot) handleCancelSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Selection Cancelled",
					Description: "Book selection has been cancelled.",
					Color:       0xff0000,
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to cancel selection", slog.Any("error", err))
	}
}

func (b *Bot) respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	ctx := context.Background()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send error response", slog.Any("error", err))
	}
}

func (b *Bot) editResponseWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	ctx := context.Background()
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: String("❌ " + message),
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to edit response with error", slog.Any("error", err))
	}
}

func String(s string) *string {
	return &s
}
