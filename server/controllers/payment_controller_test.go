package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/dto"
	"server/models"
	"server/database"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandlePayment_ValidRequest(t *testing.T) {
    e := echo.New()

    p := models.Product{Name: "TestProduct", Price: 100.0}
    database.DB.Create(&p)

    items := []dto.CartItemDTO{
        {ProductID: p.ID, Quantity: 2},
    }
    body, _ := json.Marshal(items)

    req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := HandlePayment(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    var resp map[string]interface{}
    err = json.Unmarshal(rec.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "ok", resp["status"])
    assert.Greater(t, resp["total"].(float64), 0.0)
}

func TestHandlePayment_EmptyCart(t *testing.T) {
	e := echo.New()
	body, _ := json.Marshal([]dto.CartItemDTO{})
	req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandlePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "empty")
}

func TestHandlePayment_InvalidJSON(t *testing.T) {
	e := echo.New()
	invalid := []byte("{not json}")
	req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(invalid))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandlePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid")
}

func TestHandlePayment_ProductNotFound(t *testing.T) {
	e := echo.New()
	items := []dto.CartItemDTO{
		{ProductID: 99999, Quantity: 1},
	}
	body, _ := json.Marshal(items)
	req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandlePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, 0.0, resp["total"].(float64))
}

func TestHandlePayment_MixedValidAndInvalidProducts(t *testing.T) {
	e := echo.New()
	valid := models.Product{Name: "TestMix", Price: 10.0}
	database.DB.Create(&valid)

	items := []dto.CartItemDTO{
		{ProductID: valid.ID, Quantity: 2},
		{ProductID: 99999, Quantity: 1},
	}
	body, _ := json.Marshal(items)
	req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandlePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, 20.0, resp["total"])
	assert.IsType(t, 20.0, resp["total"])
	assert.NotNil(t, resp["status"])
	assert.NotEmpty(t, resp["total"])
	assert.True(t, resp["total"].(float64) > 0)
	assert.Contains(t, rec.Body.String(), "total")
}

func TestHandlePayment_ZeroQuantity(t *testing.T) {
	e := echo.New()
	p := models.Product{Name: "ZeroQty", Price: 99.9}
	database.DB.Create(&p)

	items := []dto.CartItemDTO{
		{ProductID: p.ID, Quantity: 0},
	}
	body, _ := json.Marshal(items)
	req := httptest.NewRequest(http.MethodPost, "/payment", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandlePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, 0.0, resp["total"].(float64))
}
