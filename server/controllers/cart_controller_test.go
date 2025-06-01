package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/database"
	"server/dto"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const saveCartPath = "/save-cart"

func init() {
	database.Connect()
}

func TestSaveCartValid(t *testing.T) {
	e := echo.New()
	items := []dto.CartItemDTO{
		{ProductID: 1, Quantity: 2},
	}
	body, _ := json.Marshal(items)

	req := httptest.NewRequest(http.MethodPost, saveCartPath, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := SaveCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "cart_id")
}

func TestSaveCartEmpty(t *testing.T) {
	e := echo.New()
	body, _ := json.Marshal([]dto.CartItemDTO{})

	req := httptest.NewRequest(http.MethodPost, saveCartPath, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := SaveCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Empty")
}

func TestSaveCartInvalidJSON(t *testing.T) {
	e := echo.New()
	invalidBody := []byte("{not json}")

	req := httptest.NewRequest(http.MethodPost, saveCartPath, bytes.NewReader(invalidBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := SaveCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid")
}
