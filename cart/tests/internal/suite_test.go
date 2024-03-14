package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const addr = "http://0.0.0.0:8082"

type Suit struct {
	suite.Suite
}

func (s *Suit) SetupSuite() {
	if err := s.checkAppStartup(10 * time.Second); err != nil {
		log.Fatal("app is not started")
	}
}

func (s *Suit) checkAppStartup(dur time.Duration) error {
	probeUrl, _ := url.JoinPath(addr, "/healthz/alive")
	ctx, cancelFunc := context.WithTimeout(context.Background(), dur)
	defer cancelFunc()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, probeUrl, nil)
			resp, err := http.DefaultClient.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
			<-time.After(time.Second)
		}
	}
}

func (s *Suit) TestAddCheckDeleteCheck() {
	ctx := context.Background()
	// добавляем товар в корзину для пользователя
	userId := 123
	skuId := 773297411
	client := http.Client{}
	var urlPath string
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	body, _ := json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{1},
	)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	respRec, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respRec.StatusCode, "add item to cart failed")

	// проверим, что он появился и проверим цену
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respRec.StatusCode, "list items from cart failed")
	respBody := struct {
		Items []struct {
			SkuId int64  `json:"sku_id"`
			Name  string `json:"name"`
			Count uint16 `json:"count"`
			Price uint32 `json:"price"`
		} `json:"items"`
		TotalPrice uint32 `json:"total_price"`
	}{}
	body, err = io.ReadAll(respRec.Body)
	s.Require().NoError(err)
	err = json.Unmarshal(body, &respBody)
	s.Require().NoError(err)
	s.Require().NoError(err, "wrong response body")
	s.Assert().NotZero(respBody.TotalPrice, "total price = 0")
	s.Require().Len(respBody.Items, 1, "there was not 1 position in cart")
	s.Assert().Equal("Кроссовки Nike JORDAN", respBody.Items[0].Name)

	// удалим этот товар
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, respRec.StatusCode, "delete item from cart failed")

	// запросим товары снова и проверим статус код
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, respRec.StatusCode, "wrong status code for empty cart")
}

func (s *Suit) TestAddAddCheckClearCheck() {
	ctx := context.Background()
	// добавляем товар в корзину для пользователя два раза
	userId := 1234
	skuId := 773297411
	client := http.Client{}
	var urlPath string
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart", strconv.Itoa(skuId))
	body, _ := json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{1},
	)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	respRec, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respRec.StatusCode, "add item to cart failed")
	body, _ = json.Marshal(
		struct {
			Count uint16 `json:"count"`
		}{5},
	)
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respRec.StatusCode, "add item to cart failed")

	// проверим, что он сложился
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, respRec.StatusCode, "list items from cart failed")
	respBody := struct {
		Items []struct {
			SkuId int64  `json:"sku_id"`
			Name  string `json:"name"`
			Count uint16 `json:"count"`
			Price uint32 `json:"price"`
		} `json:"items"`
		TotalPrice uint32 `json:"total_price"`
	}{}
	body, err = io.ReadAll(respRec.Body)
	s.Require().NoError(err)
	err = json.Unmarshal(body, &respBody)
	s.Require().NoError(err)
	s.Require().NoError(err, "wrong response body")
	s.Require().Len(respBody.Items, 1, "there was not 1 position in cart")
	s.Assert().Equal("Кроссовки Nike JORDAN", respBody.Items[0].Name)
	s.Assert().Equal(6, int(respBody.Items[0].Count), "item count mismatch")

	// очистим корзину
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodDelete, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, respRec.StatusCode, "clear cart failed")

	// запросим товары снова и проверим статус код
	urlPath, _ = url.JoinPath(addr, "user", strconv.Itoa(userId), "cart")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	respRec, err = client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, respRec.StatusCode, "wrong status code for empty cart")
}
