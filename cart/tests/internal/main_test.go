package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"route256.ozon.ru/project/cart/internal/app"
	addPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/add"
	deletePkg "route256.ozon.ru/project/cart/internal/controllers/handlers/delete"
	listPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/list"
	"strconv"
	"testing"
)

func TestAddCheckDeleteCheck(t *testing.T) {
	t.Parallel()
	config, err := app.NewConfig("../../configs/config.json")
	if err != nil {
		log.Fatal(err)
	}
	appl := app.NewApp(config)
	appl.InitHandlers()

	ctx := context.Background()
	// добавляем товар в корзину для пользователя
	userId := 123
	skuId := 773297411
	var urlPath string
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	body, _ := json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{1},
	)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	req.SetPathValue(addPkg.UserIdSegment, strconv.Itoa(userId))
	req.SetPathValue(addPkg.SkuIdSegment, strconv.Itoa(skuId))
	respRec := httptest.NewRecorder()
	appl.AddHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusOK, respRec.Code, "add item to cart failed")

	// проверим, что он появился и проверим цену
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	req.SetPathValue(listPkg.UserIdSegment, strconv.Itoa(userId))
	respRec = httptest.NewRecorder()
	appl.ListHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusOK, respRec.Code, "list items from cart failed")
	respBody := struct {
		Items []struct {
			SkuId int64  `json:"sku_id"`
			Name  string `json:"name"`
			Count uint16 `json:"count"`
			Price uint32 `json:"price"`
		} `json:"items"`
		TotalPrice uint32 `json:"total_price"`
	}{}
	err = json.Unmarshal(respRec.Body.Bytes(), &respBody)
	require.NoError(t, err, "wrong response body")
	assert.NotZero(t, respBody.TotalPrice, "total price = 0")
	require.Len(t, respBody.Items, 1, "there was not 1 position in cart")
	assert.Equal(t, "Кроссовки Nike JORDAN", respBody.Items[0].Name)

	// удалим этот товар
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, urlPath, nil)
	req.SetPathValue(deletePkg.UserIdSegment, strconv.Itoa(userId))
	req.SetPathValue(deletePkg.SkuIdSegment, strconv.Itoa(skuId))
	respRec = httptest.NewRecorder()
	appl.DeleteHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusNoContent, respRec.Code, "delete item from cart failed")

	// запросим товары снова и проверим статус код
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	req.SetPathValue(listPkg.UserIdSegment, strconv.Itoa(userId))
	respRec = httptest.NewRecorder()
	appl.ListHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusNotFound, respRec.Code, "wrong status code for empty cart")
}

func TestAddAddCheckClearCheck(t *testing.T) {
	t.Parallel()
	config, err := app.NewConfig("../../configs/config.json")
	if err != nil {
		log.Fatal(err)
	}
	appl := app.NewApp(config)
	appl.InitHandlers()

	ctx := context.Background()
	// добавляем товар в корзину для пользователя два раза
	userId := 123
	skuId := 773297411
	var urlPath string
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	body, _ := json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{1},
	)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	req.SetPathValue(addPkg.UserIdSegment, strconv.Itoa(userId))
	req.SetPathValue(addPkg.SkuIdSegment, strconv.Itoa(skuId))
	respRec := httptest.NewRecorder()
	appl.AddHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusOK, respRec.Code, "add item to cart failed")
	body, _ = json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{5},
	)
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	req.SetPathValue(addPkg.UserIdSegment, strconv.Itoa(userId))
	req.SetPathValue(addPkg.SkuIdSegment, strconv.Itoa(skuId))
	respRec = httptest.NewRecorder()
	appl.AddHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusOK, respRec.Code, "add item to cart failed")

	// проверим, что он сложился
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	req.SetPathValue(listPkg.UserIdSegment, strconv.Itoa(userId))
	respRec = httptest.NewRecorder()
	appl.ListHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusOK, respRec.Code, "list items from cart failed")
	respBody := struct {
		Items []struct {
			SkuId int64  `json:"sku_id"`
			Name  string `json:"name"`
			Count uint16 `json:"count"`
			Price uint32 `json:"price"`
		} `json:"items"`
		TotalPrice uint32 `json:"total_price"`
	}{}
	err = json.Unmarshal(respRec.Body.Bytes(), &respBody)
	require.NoError(t, err, "wrong response body")
	require.Len(t, respBody.Items, 1, "there was not 1 position in cart")
	assert.Equal(t, "Кроссовки Nike JORDAN", respBody.Items[0].Name)
	assert.Equal(t, 6, int(respBody.Items[0].Count), "item count mismatch")

	// очистим корзину
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, urlPath, nil)
	req.SetPathValue(deletePkg.UserIdSegment, strconv.Itoa(userId))
	respRec = httptest.NewRecorder()
	appl.ClearHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusNoContent, respRec.Code, "clear cart failed")

	// запросим товары снова и проверим статус код
	urlPath, _ = url.JoinPath("somehost", "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	req.SetPathValue(listPkg.UserIdSegment, strconv.Itoa(userId))
	respRec = httptest.NewRecorder()
	appl.ListHandler.ServeHTTP(respRec, req)
	require.Equal(t, http.StatusNotFound, respRec.Code, "wrong status code for empty cart")
}
