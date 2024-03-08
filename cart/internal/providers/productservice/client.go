package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"route256.ozon.ru/project/cart/internal/models"
)

var (
	errUnmarshableBody    = fmt.Errorf("error unmarshalling response body from Product Service")
	errInvalidSkusArray   = fmt.Errorf("no list sku in response")
	errInvalidPrice       = fmt.Errorf("returned price is not valid")
	errInvalidProductName = fmt.Errorf("returned name is not valid")
	errSkuIdIsNotUInt32   = fmt.Errorf("SkuId is not in range UInt32")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client возвращает информацию о товарах в "специальном сервисе"
type Client struct {
	token   string
	baseURL *url.URL
	client  HTTPClient
}

func New(httpClient HTTPClient, baseURL *url.URL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		client:  httpClient,
	}
}

// IsItemPresent принимает ИД товара и возращает true, если он существует в "специальном сервисе"
func (p *Client) IsItemPresent(ctx context.Context, skuId int64) (bool, error) {
	if err := p.checkSkuId(skuId); err != nil {
		return false, errSkuIdIsNotUInt32
	}
	reqBody := listSkusRequest{
		Token:         p.token,
		StartAfterSku: uint32(skuId - 1),
		Count:         1,
	}
	req, err := p.newPOSTRequest(ctx, "list_skus", reqBody)
	if err != nil {
		return false, err
	}
	response, err := p.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error during request: %s\n", err)
	}
	var respDTO listSkusResponse
	err = p.parseResponse(response, &respDTO)
	if err != nil {
		return false, err
	}
	err = validateListSkusResponse(respDTO)
	if err != nil {
		return false, err
	}
	if len(*respDTO.Skus) > 0 && (*respDTO.Skus)[0] == uint32(skuId) {
		return true, nil
	}
	return false, nil
}

// GetProductsInfo принимает ИД товаров и возвращет их название и цену в том же порядке, как было в skuIds.
func (p *Client) GetProductsInfo(ctx context.Context, skuIds []int64) ([]models.ProductInfo, error) {
	prodInfos := make([]models.ProductInfo, 0, len(skuIds))
	for _, skuId := range skuIds {
		prodInfo, err := p.getProductInfo(ctx, skuId)
		if err != nil {
			return nil, err
		}
		prodInfos = append(prodInfos, prodInfo)
	}
	return prodInfos, nil
}

func (p *Client) getProductInfo(ctx context.Context, skuId int64) (models.ProductInfo, error) {
	if err := p.checkSkuId(skuId); err != nil {
		return models.ProductInfo{}, errSkuIdIsNotUInt32
	}
	reqBody := getProductRequest{
		Token: p.token,
		Sku:   uint32(skuId),
	}
	req, err := p.newPOSTRequest(ctx, "get_product", reqBody)
	if err != nil {
		return models.ProductInfo{}, err
	}
	response, err := p.client.Do(req)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error during request: %s\n", err)
	}
	var respDTO getProductResponse
	if err = p.parseResponse(response, &respDTO); err != nil {
		return models.ProductInfo{}, err
	}
	err = validateGetProductResponse(respDTO)
	if err != nil {
		return models.ProductInfo{}, err
	}
	return models.ProductInfo{
		Name:  *respDTO.Name,
		Price: *respDTO.Price,
	}, nil
}

func (p *Client) newPOSTRequest(ctx context.Context, method string, reqBody any) (*http.Request, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	reqBodyReader := bytes.NewReader(bodyBytes)
	urlMethod := p.baseURL.JoinPath(method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlMethod.String(), reqBodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	return req, nil
}

func (p *Client) parseResponse(response *http.Response, toObj any) error {
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(resBody, toObj); err != nil {
		return errUnmarshableBody
	}
	return nil
}

func (p *Client) checkSkuId(skuId int64) error {
	if skuId < 0 || skuId > math.MaxUint32 {
		return errSkuIdIsNotUInt32
	}
	return nil
}

func validateListSkusResponse(resp listSkusResponse) error {
	if resp.Skus == nil {
		return errInvalidSkusArray
	}
	return nil
}

func validateGetProductResponse(resp getProductResponse) error {
	var err error
	if resp.Name == nil || !models.IsStringValidName(*resp.Name) {
		err = errors.Join(err, errInvalidProductName)
	}
	if resp.Price == nil {
		err = errors.Join(err, errInvalidPrice)
	}
	return err
}
