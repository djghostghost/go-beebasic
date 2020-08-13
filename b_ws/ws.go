package b_ws


import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WsHandler struct {
	client *fasthttp.Client

	Params map[string]string

	Headers map[string]string

	ContentType string

	Body []byte

	Timeout time.Duration
}

func Default() *WsHandler {
	return New(5*time.Second, 5*time.Second, 3)
}

func New(diaTimeout time.Duration, timeout time.Duration, attempts int) *WsHandler {
	return &WsHandler{
		client: &fasthttp.Client{
			MaxIdemponentCallAttempts: attempts,
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(addr, diaTimeout)
			},
		},
		Params: make(map[string]string),
		Headers: make(map[string]string),
		Timeout: timeout,
	}
}

func (ws *WsHandler) SetHeader(k, v string) *WsHandler {
	ws.Headers[k] = v
	return ws
}
func (ws *WsHandler) SetParam(k, v string) *WsHandler {
	ws.Params[k] = v
	return ws
}
func (ws *WsHandler) SetTimeout(timeout time.Duration) *WsHandler {
	ws.Timeout = timeout
	return ws
}

func (ws *WsHandler) DoRequest(method string, urlStr string) ([]byte, int, error) {
	var req = fasthttp.AcquireRequest()
	var resp = fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, 0, err
	}
	if method == http.MethodPost {
		for k, v := range ws.Params {
			req.PostArgs().Set(k, v)
		}
	} else {
		queryParams := u.Query()
		for k, v := range ws.Params {
			queryParams.Set(k, v)
		}
		u.RawQuery = queryParams.Encode()
	}
	fullUrl := u.String()
	logs.Debug("[WS] url: ", fullUrl)
	logs.Debug("[WS] params: ", ws.Params)
	req.Header.SetMethod(method)
	if ws.ContentType != "" {
		req.Header.SetContentType(ws.ContentType)
	}
	for k, v := range ws.Headers {
		req.Header.Set(k, v)
	}
	req.SetRequestURI(fullUrl)
	if ws.Body != nil {
		req.SetBody(ws.Body)
	}
	// DO request
	err = ws.client.DoTimeout(req, resp, ws.Timeout)
	if err != nil {
		return nil, 0, err
	}
	respBody := resp.Body()
	statusCode := resp.Header.StatusCode()
	logs.Debug("[WS] resp body: ", string(respBody), " status code: ", strconv.Itoa(statusCode))
	if statusCode >= 400 {
		return nil, statusCode, errors.New("[WS] status code is " + strconv.Itoa(statusCode))
	}
	return respBody, statusCode, nil
}

func (ws *WsHandler) FetchGetJson(urlStr string) (*gjson.Result, error) {
	body, _, err := ws.DoRequest(http.MethodGet, urlStr)
	if err != nil {
		return nil, err
	}
	r := gjson.ParseBytes(body)
	return &r, nil
}

func (ws *WsHandler) FetchPostJson(urlStr string) (*gjson.Result, error) {
	body, _, err := ws.DoRequest(http.MethodPost, urlStr)
	if err != nil {
		return nil, err
	}
	r := gjson.ParseBytes(body)
	return &r, nil
}

func (ws *WsHandler) FetchGetString(urlStr string) (string, error) {
	body, _, err := ws.DoRequest(http.MethodGet, urlStr)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (ws *WsHandler) FetchPostString(urlStr string) (string, error) {
	body, _, err := ws.DoRequest(http.MethodPost, urlStr)
	if err != nil {
		return "", err
	}
	return string(body), nil
}


