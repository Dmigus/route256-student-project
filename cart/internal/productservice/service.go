package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/lister"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ProductService возвращает информацию о товарах в "специальном сервисе"
type ProductService struct {
	token   string
	baseURL *url.URL
	client  HTTPClient
}

func New(httpClient HTTPClient, baseURL *url.URL, token string) *ProductService {
	return &ProductService{
		baseURL: baseURL,
		token:   token,
		client:  httpClient,
	}
}

// IsItemPresent принимает ИД товара и возращает true, если он существует в "специальном сервисе"
func (p *ProductService) IsItemPresent(ctx context.Context, skuId service.SkuId) (bool, error) {
	reqBody := listSkusRequest{
		Token:         p.token,
		StartAfterSku: int64(skuId) - 1,
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
	if len(respDTO.Skus) > 0 && respDTO.Skus[0] == int64(skuId) {
		return true, nil
	}
	return false, nil
}

// GetProductsInfo принимает ИД товаров и возвращет их название и цену в том же порядке, как было в skuIds.
func (p *ProductService) GetProductsInfo(ctx context.Context, skuIds []service.SkuId) ([]lister.ProductInfo, error) {
	prodInfos := make([]lister.ProductInfo, 0, len(skuIds))
	for _, skuId := range skuIds {
		prodInfo, err := p.getProductInfo(ctx, skuId)
		if err != nil {
			return nil, err
		}
		prodInfos = append(prodInfos, prodInfo)
	}
	return prodInfos, nil
}

func (p *ProductService) getProductInfo(ctx context.Context, skuId service.SkuId) (lister.ProductInfo, error) {
	reqBody := getProductRequest{
		Token: p.token,
		Sku:   skuId,
	}
	req, err := p.newPOSTRequest(ctx, "get_product", reqBody)
	if err != nil {
		return lister.ProductInfo{}, err
	}
	response, err := p.client.Do(req)
	if err != nil {
		return lister.ProductInfo{}, fmt.Errorf("error during request: %s\n", err)
	}
	var respDTO getProductResponse
	if err = p.parseResponse(response, &respDTO); err != nil {
		return lister.ProductInfo{}, err
	}
	return lister.ProductInfo{
		Name:  respDTO.Name,
		Price: service.Price(respDTO.Price),
	}, nil
}

func (p *ProductService) newPOSTRequest(ctx context.Context, method string, reqBody any) (*http.Request, error) {
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

func (p *ProductService) parseResponse(response *http.Response, toObj any) error {
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(resBody, toObj); err != nil {
		return err
	}
	return nil
}
