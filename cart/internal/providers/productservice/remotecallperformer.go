package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var errUnmarshableBody = errors.New("error unmarshalling response body from Product Service")

type (
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	RequestWithSettableToken interface {
		SetToken(string)
	}

	errorResponse struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}

	// RemoteCallPerformer предназначен для осуществления вызова некого метода и получения результата
	RemoteCallPerformer[responseT any] struct {
		token   string
		baseURL *url.URL
		client  httpClient
	}
)

// NewRCPerformer создаёт новый RemoteCallPerformer
func NewRCPerformer[respT any](httpClient httpClient, baseURL *url.URL, token string) *RemoteCallPerformer[respT] {
	return &RemoteCallPerformer[respT]{
		baseURL: baseURL,
		token:   token,
		client:  httpClient,
	}
}

// Perform осуществляет вызов метода method, используя reqBody в качестве запроса и вовращая ответ.
// reqBody должен быть сериализуем в json.
func (rcp *RemoteCallPerformer[respT]) Perform(ctx context.Context, method string, reqBody RequestWithSettableToken) (*respT, error) {
	reqBody.SetToken(rcp.token)
	req, err := rcp.newPOSTRequest(ctx, method, reqBody)
	reqBody.SetToken("")
	if err != nil {
		return nil, err
	}
	response, err := rcp.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during request: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		errorResp, err := rcp.parseError(response)
		if err != nil {
			return nil, fmt.Errorf("recieved %d status code, but failed to decode error response: %w", response.StatusCode, err)
		}
		return nil, fmt.Errorf("HTTP request responded with: %d , message: %s", errorResp.Code, errorResp.Message)
	}
	return rcp.parseResponse(response)
}

func (rcp *RemoteCallPerformer[respT]) newPOSTRequest(ctx context.Context, method string, reqBody RequestWithSettableToken) (*http.Request, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	reqBodyReader := bytes.NewReader(bodyBytes)
	urlMethod := rcp.baseURL.JoinPath(method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlMethod.String(), reqBodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	return req, nil
}

func (rcp *RemoteCallPerformer[respT]) parseResponse(response *http.Response) (*respT, error) {
	var responseObj respT
	if err := rcp.parseToAny(response, &responseObj); err != nil {
		return nil, err
	}
	return &responseObj, nil
}

func (rcp *RemoteCallPerformer[respT]) parseError(response *http.Response) (*errorResponse, error) {
	var errorObj errorResponse
	if err := rcp.parseToAny(response, &errorObj); err != nil {
		return nil, err
	}
	return &errorObj, nil
}

func (rcp *RemoteCallPerformer[respT]) parseToAny(response *http.Response, toObj any) error {
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(resBody, toObj); err != nil {
		return errUnmarshableBody
	}
	return nil
}
