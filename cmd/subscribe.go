package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/config"
)

type authorRequest struct {
	Name string `json:"name"`
}

type authorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type subscriptionRequest struct {
	ScopeName  string `json:"scope_name"`
	NotifierID *uint  `json:"notifier_id,omitempty"`
}

type subscriptionResponse struct {
	ID           uint    `json:"id"`
	AuthorID     uint    `json:"author_id"`
	AuthorName   string  `json:"author_name"`
	ScopeName    string  `json:"scope_name"`
	NotifierName *string `json:"notifier_name"`
}

type notifierResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func createSubscribeCmd() *cobra.Command {
	var scope string
	var notifier string

	subscribeCmd := &cobra.Command{
		Use:   "subscribe <author-name>",
		Short: "Create an author and subscription",
		Long:  `Creates an author with the given name and a subscription for that author with the specified scope.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSubscribeCmd(args[0], scope, notifier)
		},
	}

	subscribeCmd.Flags().StringVar(&scope, "scope", "personal", "Subscription scope (personal, family, kids, general)")
	subscribeCmd.Flags().StringVar(&notifier, "notifier", "", "Notifier name for subscription notifications")

	return subscribeCmd
}

func runSubscribeCmd(authorName, scope, notifierName string) error {
	ctx := context.Background()

	baseURL := config.Config.APIClient.URL
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	slog.DebugContext(ctx, "Creating author subscription", slog.String("author", authorName), slog.String("scope", scope))

	// Look up notifier if provided
	var notifierID *uint
	if notifierName != "" {
		notifier, err := getNotifierByName(ctx, baseURL, notifierName)
		if err != nil {
			return err
		}
		notifierID = &notifier.ID
		slog.DebugContext(ctx, "Found notifier", slog.String("name", notifier.Name), slog.Uint64("id", uint64(notifier.ID)))
	}

	// Create the author
	author, err := createAuthor(ctx, baseURL, authorName)
	if err != nil {
		return err
	}

	fmt.Printf("Created author: %s (ID: %d)\n", author.Name, author.ID)

	// Create the subscription
	sub, err := createSubscription(ctx, baseURL, author.ID, scope, notifierID)
	if err != nil {
		return err
	}

	if sub.NotifierName != nil {
		fmt.Printf("Created subscription: %s -> %s with notifier %s (ID: %d)\n", sub.AuthorName, sub.ScopeName, *sub.NotifierName, sub.ID)
	} else {
		fmt.Printf("Created subscription: %s -> %s (ID: %d)\n", sub.AuthorName, sub.ScopeName, sub.ID)
	}

	return nil
}

func createAuthor(ctx context.Context, baseURL, name string) (*authorResponse, error) {
	url := fmt.Sprintf("%s/api/authors", baseURL)

	reqBody := authorRequest{Name: name}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal author request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create author: " + string(body))
	}

	var author authorResponse
	if err := json.Unmarshal(body, &author); err != nil {
		return nil, fmt.Errorf("failed to parse author response: %w", err)
	}

	return &author, nil
}

func createSubscription(ctx context.Context, baseURL string, authorID uint, scope string, notifierID *uint) (*subscriptionResponse, error) {
	url := fmt.Sprintf("%s/api/authors/%d/subscription", baseURL, authorID)

	reqBody := subscriptionRequest{ScopeName: scope, NotifierID: notifierID}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create subscription: " + string(body))
	}

	var sub subscriptionResponse
	if err := json.Unmarshal(body, &sub); err != nil {
		return nil, fmt.Errorf("failed to parse subscription response: %w", err)
	}

	return &sub, nil
}

func getNotifierByName(ctx context.Context, baseURL, name string) (*notifierResponse, error) {
	url := fmt.Sprintf("%s/api/notifiers", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifiers: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to list notifiers: " + string(body))
	}

	var notifiers []notifierResponse
	if err := json.Unmarshal(body, &notifiers); err != nil {
		return nil, fmt.Errorf("failed to parse notifiers response: %w", err)
	}

	for _, n := range notifiers {
		if n.Name == name {
			return &n, nil
		}
	}

	return nil, fmt.Errorf("notifier not found: %s", name)
}
