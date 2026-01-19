package client

import (
	"context"
	"log/slog"

	"github.com/carlmjohnson/requests"
)

type GenericClient[Request any, Response any] struct {
	BaseUrl  string
	TypeName string
}

// ReadOnlyGenericClient is a client for read-only resources that only supports List and Get operations.
type ReadOnlyGenericClient[Response any] struct {
	BaseUrl  string
	TypeName string
}

func (c *ReadOnlyGenericClient[Response]) List(ctx context.Context) ([]Response, error) {
	var response []Response

	slog.InfoContext(ctx, "Fetching list", slog.String("type", c.TypeName))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s", c.TypeName).
		Method("GET").
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ReadOnlyGenericClient[Response]) Get(ctx context.Context, id uint) (Response, error) {
	var response Response

	slog.InfoContext(ctx, "Fetching item", slog.String("type", c.TypeName), slog.Any("id", id))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s/%d", c.TypeName, id).
		Method("GET").
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *GenericClient[Request, Response]) List(ctx context.Context) ([]Response, error) {
	var response []Response

	slog.InfoContext(ctx, "Fetching list", slog.String("type", c.TypeName))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s", c.TypeName).
		Method("GET").
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *GenericClient[Request, Response]) Get(ctx context.Context, id uint) (Response, error) {
	var response Response

	slog.InfoContext(ctx, "Fetching item", slog.String("type", c.TypeName), slog.Any("id", id))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s/%d", c.TypeName, id).
		Method("GET").
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *GenericClient[Request, Response]) Create(ctx context.Context, request Request) (Response, error) {
	var response Response

	slog.InfoContext(ctx, "Creating item", slog.String("type", c.TypeName))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s", c.TypeName).
		Method("POST").
		BodyJSON(&request).
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *GenericClient[Request, Response]) Delete(ctx context.Context, id uint) error {
	slog.InfoContext(ctx, "Delete item", slog.String("type", c.TypeName), slog.Any("id", id))

	err := requests.
		URL(c.BaseUrl).
		Pathf("/api/%s/%d", c.TypeName, id).
		Method("DELETE").
		Fetch(ctx)
	if err != nil {
		return err
	}

	return nil
}
