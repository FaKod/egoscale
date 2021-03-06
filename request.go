package egoscale

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Error formats a CloudStack error into a standard error
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("API error %s %d (%s %d): %s", e.ErrorCode, e.ErrorCode, e.CSErrorCode, e.CSErrorCode, e.ErrorText)
}

// Success computes the values based on the RawMessage, either string or bool
func (e *booleanResponse) IsSuccess() (bool, error) {
	if e.Success == nil {
		return false, errors.New("not a valid booleanResponse, Success is missing")
	}

	str := ""
	if err := json.Unmarshal(e.Success, &str); err != nil {
		boolean := false
		if e := json.Unmarshal(e.Success, &boolean); e != nil {
			return false, e
		}
		return boolean, nil
	}
	return str == "true", nil
}

// Error formats a CloudStack job response into a standard error
func (e *booleanResponse) Error() error {
	success, err := e.IsSuccess()

	if err != nil {
		return err
	}

	if success {
		return nil
	}

	return fmt.Errorf("API error: %s", e.DisplayText)
}

func (client *Client) parseResponse(resp *http.Response, key string) (json.RawMessage, error) {
	contentType := resp.Header.Get("content-type")

	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("body content-type response expected \"application/json\", got %q", contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	response, ok := m[key]
	if !ok {
		if resp.StatusCode >= 400 {
			response, ok = m["errorresponse"]
		}
		if !ok {
			for k := range m {
				return nil, fmt.Errorf("malformed JSON response, %q was expected, got %q", key, k)
			}
		}
	}

	if resp.StatusCode >= 400 {
		errorResponse := new(ErrorResponse)
		if e := json.Unmarshal(response, errorResponse); e != nil && errorResponse.ErrorCode <= 0 {
			return nil, fmt.Errorf("%d %s", resp.StatusCode, b)
		}
		return nil, errorResponse
	}

	n := map[string]json.RawMessage{}
	if err := json.Unmarshal(response, &n); err != nil {
		return nil, err
	}

	if len(n) > 1 {
		return response, nil
	}

	if len(n) == 1 {
		for k := range n {
			// boolean response and asyncjob result may also contain
			// only one key
			if k == "success" || k == "jobid" {
				return response, nil
			}
			return n[k], nil
		}
	}

	return response, nil
}

// asyncRequest perform an asynchronous job with a context
func (client *Client) asyncRequest(ctx context.Context, request AsyncCommand) (interface{}, error) {
	var err error

	res := request.asyncResponse()
	client.AsyncRequestWithContext(ctx, request, func(j *AsyncJobResult, er error) bool {
		if er != nil {
			err = er
			return false
		}
		if j.JobStatus == Success {
			if r := j.Response(res); err != nil {
				err = r
			}
			return false
		}
		return true
	})
	return res, err
}

// syncRequest performs a sync request with a context
func (client *Client) syncRequest(ctx context.Context, request syncCommand) (interface{}, error) {
	body, err := client.request(ctx, request)
	if err != nil {
		return nil, err
	}

	response := request.response()
	err = json.Unmarshal(body, response)

	// booleanResponse will alway be valid...
	if err == nil {
		if br, ok := response.(*booleanResponse); ok {
			success, e := br.IsSuccess()
			if e != nil {
				return nil, e
			}
			if !success {
				err = errors.New("not a valid booleanResponse")
			}
		}
	}

	if err != nil {
		errResponse := new(ErrorResponse)
		if e := json.Unmarshal(body, errResponse); e == nil && errResponse.ErrorCode > 0 {
			return errResponse, nil
		}
		return nil, err
	}

	return response, nil
}

// BooleanRequest performs the given boolean command
func (client *Client) BooleanRequest(req Command) error {
	resp, err := client.Request(req)
	if err != nil {
		return err
	}

	if b, ok := resp.(*booleanResponse); ok {
		return b.Error()
	}

	panic(fmt.Errorf("command %q is not a proper boolean response. %#v", req.name(), resp))
}

// BooleanRequestWithContext performs the given boolean command
func (client *Client) BooleanRequestWithContext(ctx context.Context, req Command) error {
	resp, err := client.RequestWithContext(ctx, req)
	if err != nil {
		return err
	}

	if b, ok := resp.(*booleanResponse); ok {
		return b.Error()
	}

	panic(fmt.Errorf("command %q is not a proper boolean response. %#v", req.name(), resp))
}

// Request performs the given command
func (client *Client) Request(request Command) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	switch request.(type) {
	case syncCommand:
		return client.syncRequest(ctx, request.(syncCommand))
	case AsyncCommand:
		return client.asyncRequest(ctx, request.(AsyncCommand))
	default:
		panic(fmt.Errorf("command %q is not a proper Sync or Async command", request.name()))
	}
}

// RequestWithContext preforms a request with a context
func (client *Client) RequestWithContext(ctx context.Context, request Command) (interface{}, error) {
	switch request.(type) {
	case syncCommand:
		return client.syncRequest(ctx, request.(syncCommand))
	case AsyncCommand:
		return client.asyncRequest(ctx, request.(AsyncCommand))
	default:
		panic(fmt.Errorf("command %q is not a proper Sync or Async command", request.name()))
	}
}

// AsyncRequest performs the given command
func (client *Client) AsyncRequest(request AsyncCommand, callback WaitAsyncJobResultFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	client.AsyncRequestWithContext(ctx, request, callback)
}

// AsyncRequestWithContext preforms a request with a context
func (client *Client) AsyncRequestWithContext(ctx context.Context, request AsyncCommand, callback WaitAsyncJobResultFunc) {
	body, err := client.request(ctx, request)
	if err != nil {
		callback(nil, err)
		return
	}

	jobResult := new(AsyncJobResult)
	if err := json.Unmarshal(body, jobResult); err != nil {
		r := new(ErrorResponse)
		if e := json.Unmarshal(body, r); e != nil && r.ErrorCode > 0 {
			if !callback(nil, r) {
				return
			}
		}
		if !callback(nil, err) {
			return
		}
	}

	// Successful response
	if jobResult.JobID == "" || jobResult.JobStatus != Pending {
		callback(jobResult, nil)
		// without a JobID, the next requests will only fail
		return
	}

	for iteration := 0; ; iteration++ {
		time.Sleep(client.RetryStrategy(int64(iteration)))

		req := &QueryAsyncJobResult{JobID: jobResult.JobID}
		resp, err := client.syncRequest(ctx, req)
		if err != nil && !callback(nil, err) {
			return
		}

		result, ok := resp.(*AsyncJobResult)
		if !ok {
			if !callback(nil, fmt.Errorf("wrong type. AsyncJobResult expected, got %T", resp)) {
				return
			}
		}

		if result.JobStatus == Failure {
			if !callback(nil, result.Error()) {
				return
			}
		} else {
			if !callback(result, nil) {
				return
			}
		}
	}
}

// Payload builds the HTTP request from the given command
func (client *Client) Payload(request Command) (string, error) {
	params := url.Values{}
	err := prepareValues("", &params, request)
	if err != nil {
		return "", err
	}
	if hookReq, ok := request.(onBeforeHook); ok {
		if err := hookReq.onBeforeSend(&params); err != nil {
			return "", err
		}
	}
	params.Set("apikey", client.APIKey)
	params.Set("command", request.name())
	params.Set("response", "json")

	// This code is borrowed from net/url/url.go
	// The way it's encoded by net/url doesn't match
	// how CloudStack works.
	var buf bytes.Buffer
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		prefix := csEncode(k) + "="
		for _, v := range params[k] {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(csEncode(v))
		}
	}

	return buf.String(), nil
}

// Sign signs the HTTP request and return it
func (client *Client) Sign(query string) (string, error) {
	mac := hmac.New(sha1.New, []byte(client.apiSecret))
	_, err := mac.Write([]byte(strings.ToLower(query)))
	if err != nil {
		return "", err
	}

	signature := csEncode(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	return fmt.Sprintf("%s&signature=%s", csQuotePlus(query), signature), nil
}

// request makes a Request while being close to the metal
func (client *Client) request(ctx context.Context, req Command) (json.RawMessage, error) {
	payload, err := client.Payload(req)
	if err != nil {
		return nil, err
	}
	query, err := client.Sign(payload)
	if err != nil {
		return nil, err
	}

	method := "GET"
	url := fmt.Sprintf("%s?%s", client.Endpoint, query)

	var body io.Reader
	// respect Internet Explorer limit of 2048
	if len(url) > 1<<11 {
		url = client.Endpoint
		method = "POST"
		body = strings.NewReader(query)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Add("User-Agent", fmt.Sprintf("exoscale/egoscale (%v)", Version))

	if method == "POST" {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Length", strconv.Itoa(len(query)))
	}

	resp, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	// XXX: addIpToNic is kind of special
	key := fmt.Sprintf("%sresponse", strings.ToLower(req.name()))
	if key == "addiptonicresponse" {
		key = "addiptovmnicresponse"
	}

	text, err := client.parseResponse(resp, key)
	if err != nil {
		return nil, err
	}

	return text, nil
}
