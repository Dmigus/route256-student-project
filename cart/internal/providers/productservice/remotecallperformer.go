package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var errUnmarshableBody = fmt.Errorf("error unmarshalling response body from Product Service")

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
	RemoteCallPerformer struct {
		token   string
		baseURL *url.URL
		client  httpClient
	}
)

func NewRCPerformer(httpClient httpClient, baseURL *url.URL, token string) *RemoteCallPerformer {
	return &RemoteCallPerformer{
		baseURL: baseURL,
		token:   token,
		client:  httpClient,
	}
}

// Perform осуществляет вызов метода method, используя reqBody в качестве запроса и записывая ответ в respBody.
// reqBody должен быть сериализуем в json.
func (rcp *RemoteCallPerformer) Perform(ctx context.Context, method string, reqBody RequestWithSettableToken, respBody any) error {
	reqBody.SetToken(rcp.token)
	req, err := rcp.newPOSTRequest(ctx, method, reqBody)
	if err != nil {
		return err
	}
	response, err := rcp.client.Do(req)
	if err != nil {
		return fmt.Errorf("error during request: %s\n", err)
	}
	if response.StatusCode != http.StatusOK {
		var errorResp *errorResponse
		err = rcp.parseResponse(response, errorResp)
		if err != nil {
			return fmt.Errorf("recieved %d status code, but failed to decode error response: %w", response.StatusCode, err)
		}
		return fmt.Errorf("HTTP request responded with: %d , message: %s", errorResp.Code, errorResp.Message)
	}
	return rcp.parseResponse(response, respBody)
}

func (rcp *RemoteCallPerformer) newPOSTRequest(ctx context.Context, method string, reqBody any) (*http.Request, error) {
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

func (rcp *RemoteCallPerformer) parseResponse(response *http.Response, toObj any) error {
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(resBody, toObj); err != nil {
		return errUnmarshableBody
	}
	return nil
}
