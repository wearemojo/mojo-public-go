package mrpc_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
	"github.com/wearemojo/mojo-public-go/lib/authenforce"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/mrpc"
	"github.com/xeipuuv/gojsonschema"
)

type greetReq struct {
	Name string `json:"name"`
}

type greetRes struct {
	Message string `json:"message"`
}

func greet(_ context.Context, req *greetReq) (*greetRes, error) {
	if req.Name == "boom" {
		return nil, cher.New("name_not_allowed", cher.M{"name": req.Name})
	}

	return &greetRes{Message: "hello " + req.Name}, nil
}

func listNames(_ context.Context) ([]string, error) {
	return []string{"a", "b"}, nil
}

func touch(_ context.Context, _ *greetReq) error {
	return nil
}

func ping(_ context.Context) error {
	return nil
}

var greetSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,
	"required": ["name"],
	"properties": {"name": {"type": "string", "minLength": 1}}
}`)

func buildServer() http.Handler {
	rpc := mrpc.NewServer(authenforce.Enforcers{authenforce.UnsafeNoAuthentication})

	noAuth := authenforce.Enforcers{authenforce.UnsafeNoAuthentication}

	mrpc.Register(rpc, "greet", "2020-01-01", greetSchema, noAuth, greet)
	mrpc.RegisterNoReq(rpc, "list_names", "2020-01-01", noAuth, listNames)
	mrpc.RegisterNoRes(rpc, "touch", "2020-01-01", greetSchema, noAuth, touch)
	mrpc.RegisterNoReqRes(rpc, "ping", "2020-01-01", noAuth, ping)

	// 2021-06-01 serves everything 2020-01-01 did, but touch is dropped
	mrpc.Register(rpc, "greet", "2021-06-01", greetSchema, noAuth, greet)
	mrpc.RegisterNoReq(rpc, "list_names", "2021-06-01", noAuth, listNames)
	mrpc.RegisterNoReqRes(rpc, "ping", "2021-06-01", noAuth, ping)

	return rpc
}

type wireResult struct {
	status int
	header http.Header
	body   string
}

func fire(handler http.Handler, method, path, body string) wireResult {
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	ctx := clog.Set(context.Background(), logrus.NewEntry(logrus.New()))
	r := httptest.NewRequestWithContext(ctx, method, path, reqBody)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)

	return wireResult{
		status: res.StatusCode,
		header: res.Header,
		body:   string(raw),
	}
}

func TestWireContract(t *testing.T) {
	handler := buildServer()

	const jsonCT = "application/json; charset=utf-8"

	cases := []struct {
		name         string
		method, path string
		body         string

		status int
		ct     string
		wire   string
	}{
		{
			name: "struct response", method: "POST", path: "/2020-01-01/greet", body: `{"name":"world"}`,
			status: 200, ct: jsonCT, wire: "{\"message\":\"hello world\"}\n",
		},
		{
			name: "slice response", method: "POST", path: "/2020-01-01/list_names",
			status: 200, ct: jsonCT, wire: "[\"a\",\"b\"]\n",
		},
		{
			name: "no response 204", method: "POST", path: "/2020-01-01/touch", body: `{"name":"x"}`,
			status: 204, ct: jsonCT, wire: "",
		},
		{
			name: "bare 204", method: "POST", path: "/2020-01-01/ping",
			status: 204, ct: jsonCT, wire: "",
		},
		{
			name: "handler cher error", method: "POST", path: "/2020-01-01/greet", body: `{"name":"boom"}`,
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"name_not_allowed\",\"meta\":{\"name\":\"boom\"}}\n",
		},
		{
			name: "schema failure min length", method: "POST", path: "/2020-01-01/greet", body: `{"name":""}`,
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"bad_request\",\"reasons\":[{\"code\":\"schema_failure\",\"meta\":{\"field\":\"name\",\"message\":\"String length must be greater than or equal to 1\",\"type\":\"string_gte\"}}]}\n",
		},
		{
			name: "schema failure required", method: "POST", path: "/2020-01-01/greet", body: `{}`,
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"bad_request\",\"reasons\":[{\"code\":\"schema_failure\",\"meta\":{\"field\":\"(root)\",\"message\":\"name is required\",\"type\":\"required\"}}]}\n",
		},
		{
			name: "schema failure invalid type", method: "POST", path: "/2020-01-01/greet", body: `{"name":123}`,
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"bad_request\",\"reasons\":[{\"code\":\"schema_failure\",\"meta\":{\"field\":\"name\",\"message\":\"Invalid type. Expected: string, given: integer\",\"type\":\"invalid_type\"}}]}\n",
		},
		{
			name: "missing request body", method: "POST", path: "/2020-01-01/greet",
			status: 400, ct: jsonCT, wire: "{\"code\":\"invalid_json\"}\n",
		},
		{
			name: "unexpected body on no-input", method: "POST", path: "/2020-01-01/ping", body: `{"a":1}`,
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"bad_request\",\"reasons\":[{\"code\":\"unexpected_request_body\"}]}\n",
		},
		{
			name: "null byte now rejected", method: "POST", path: "/2020-01-01/ping", body: "\x00",
			status: 400, ct: jsonCT,
			wire: "{\"code\":\"bad_request\",\"reasons\":[{\"code\":\"unexpected_request_body\"}]}\n",
		},
		{
			name: "unknown method", method: "POST", path: "/2020-01-01/nope", body: `{}`,
			status: 404, ct: jsonCT,
			wire: "{\"code\":\"not_found\",\"meta\":{\"method\":\"nope\",\"version\":\"2020-01-01\"}}\n",
		},
		{
			name: "unknown version", method: "POST", path: "/2019-01-01/greet", body: `{"name":"x"}`,
			status: 404, ct: jsonCT,
			wire: "{\"code\":\"not_found\",\"meta\":{\"version\":\"2019-01-01\"}}\n",
		},
		{
			name: "inherited method later version", method: "POST", path: "/2021-06-01/greet", body: `{"name":"x"}`,
			status: 200, ct: jsonCT, wire: "{\"message\":\"hello x\"}\n",
		},
		{
			name: "method dropped in later version", method: "POST", path: "/2021-06-01/touch", body: `{"name":"x"}`,
			status: 404, ct: jsonCT,
			wire: "{\"code\":\"not_found\",\"meta\":{\"method\":\"touch\",\"version\":\"2021-06-01\"}}\n",
		},
		{
			name: "latest resolves", method: "POST", path: "/latest/greet", body: `{"name":"x"}`,
			status: 200, ct: jsonCT, wire: "{\"message\":\"hello x\"}\n",
		},
		{
			name: "preview not found", method: "POST", path: "/preview/greet", body: `{"name":"x"}`,
			status: 404, ct: jsonCT, wire: "{\"code\":\"not_found\",\"meta\":{\"version\":\"preview\"}}\n",
		},
		{
			name: "malformed path", method: "POST", path: "/greet", body: `{}`,
			status: 404, ct: jsonCT, wire: "{\"code\":\"not_found\"}\n",
		},
		{
			name: "get method not allowed", method: "GET", path: "/2020-01-01/greet",
			status: 405, ct: "", wire: "",
		},
		{
			name: "get on odd path", method: "GET", path: "/anything",
			status: 405, ct: "", wire: "",
		},
		{
			name: "query string rejected", method: "POST", path: "/2020-01-01/greet?a=1", body: `{"name":"x"}`,
			status: 400, ct: "", wire: "{\"code\":\"unexpected_input\"}\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)

			got := fire(handler, testCase.method, testCase.path, testCase.body)

			is.Equal(got.status, testCase.status)
			is.Equal(got.body, testCase.wire)
			is.Equal(got.header.Get("Content-Type"), testCase.ct)
		})
	}
}
