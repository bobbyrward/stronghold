package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/bobbyrward/stronghold/internal/version"
)

func GetVersion() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version":    version.Version,
			"git_commit": version.GitCommit,
			"build_time": version.BuildTime,
		})
	}
}
