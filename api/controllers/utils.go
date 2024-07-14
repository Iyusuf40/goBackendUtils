package controllers

import (
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
)

func GetBodyInMap(c echo.Context) map[string]any {
	return ReadFromReaderIntoMap(c.Request().Body)
}

func ReadFromReaderIntoMap(r io.Reader) map[string]any {
	body, _ := io.ReadAll(r)
	var bodyMap map[string]any
	json.Unmarshal(body, &bodyMap)
	return bodyMap
}
