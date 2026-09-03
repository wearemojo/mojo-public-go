package mrpc

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/wearemojo/mojo-public-go/lib/authenforce"
	"github.com/wearemojo/mojo-public-go/lib/cher"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
	"github.com/xeipuuv/gojsonschema"
)

const (
	// VersionPreview is used for experimental endpoints in development which
	// are coming but a version identifier has not been decided yet or may be
	// withdrawn at any point.
	VersionPreview = "preview"

	// VersionLatest is used by engineers only to call the latest version of an
	// endpoint in utilities like cURL and Paw.
	VersionLatest = "latest"
)

// expVersion matches valid method versions.
var expVersion = regexp.MustCompile(`^(?:preview|20\d{2}-\d{2}-\d{2})$`)

// expMethod matches valid method names.
var expMethod = regexp.MustCompile(`^[a-z][a-z\d]*(?:_[a-z\d]+)*$`)

// Server is an HTTP-compatible RPC handler. Methods are registered explicitly
// against each version they serve; there is no implied inheritance between
// versions.
type Server struct {
	baseAuth Middleware
	mw       []Middleware

	// versions maps version -> method -> endpoint. Preview lives under the
	// "preview" key and never participates in the latest alias.
	versions      map[string]map[string]HandlerFunc
	latestVersion string
}

// NewServer returns a Server whose endpoints are all guarded by the given
// base authorization enforcers, with request logging installed as the
// outermost middleware.
func NewServer(baseEnforcers authenforce.Enforcers) *Server {
	return &Server{
		baseAuth: enforcerMiddleware(baseEnforcers),
		mw:       []Middleware{Logger()},
		versions: map[string]map[string]HandlerFunc{},
	}
}

func (s *Server) register(method, version string, schema gojsonschema.JSONLoader, enforcers authenforce.Enforcers, handler HandlerFunc) {
	if !expMethod.MatchString(method) || !expVersion.MatchString(version) {
		panic(fmt.Sprintf("invalid method/version: %s %s", method, version))
	}

	if _, ok := s.versions[version][method]; ok {
		panic(fmt.Sprintf("cannot set '%s' on version '%s', it's already defined", method, version))
	}

	if len(enforcers) == 0 {
		panic(fmt.Sprintf("enforcers cannot be empty for %s %s: pass authenforce.Enforcers{authenforce.UnsafeNoAuthentication} to explicitly allow unauthenticated access", method, version))
	}

	if handler == nil {
		panic(fmt.Sprintf("handler cannot be nil for %s %s", method, version))
	}

	chain := make([]Middleware, 0, len(s.mw)+3)
	chain = append(chain, s.mw...)
	chain = append(chain, s.baseAuth)

	if schema != nil {
		compiled, err := gojsonschema.NewSchemaLoader().Compile(schema)
		if err != nil {
			panic(fmt.Sprintf("json schema error in %s: %s", method, err))
		}

		chain = append(chain, validate(compiled))
	}

	chain = append(chain, enforcerMiddleware(enforcers))

	if s.versions[version] == nil {
		s.versions[version] = map[string]HandlerFunc{}
	}

	s.versions[version][method] = compose(chain, handler)

	if version != VersionPreview && version > s.latestVersion {
		s.latestVersion = version
	}
}

func compose(mw []Middleware, h HandlerFunc) HandlerFunc {
	for _, m := range slices.Backward(mw) {
		h = m(h)
	}

	return h
}

func (s *Server) methodSet(version string) (methods map[string]HandlerFunc, ok bool) {
	if version == VersionLatest {
		if s.latestVersion == "" {
			return nil, false
		}

		return s.versions[s.latestVersion], true
	}

	methods, ok = s.versions[version]

	return methods, ok
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request, info *requestInfo) error {
	methods, ok := s.methodSet(info.Version)
	if !ok {
		return cher.New(cher.NotFound, cher.M{"version": info.Version})
	}

	matched, ok := methods[info.Method]
	if !ok {
		return cher.New(cher.NotFound, cher.M{"method": info.Method, "version": info.Version})
	}

	return matched(w, r)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if strings.ToUpper(r.Method) != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.RawQuery != "" {
		s.writeError(ctx, w, cher.New("unexpected_input", nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	method, version, ok := requestPath(r.URL.Path)
	if !ok {
		s.writeError(ctx, w, cher.New(cher.NotFound, nil))
		return
	}

	info := &requestInfo{Version: version, Method: method}
	ctx = setInfo(ctx, info)
	r = r.WithContext(ctx)

	s.writeError(ctx, w, s.serve(w, r, info))
}

// expRequestPath only matches HTTP paths formed of /<version>/<method>.
var expRequestPath = regexp.MustCompile(`^/(preview|latest|20\d{2}-\d{2}-\d{2})/([a-z0-9\_]+)$`)

func requestPath(path string) (method, version string, ok bool) {
	match := expRequestPath.FindStringSubmatch(path)
	if len(match) != 3 {
		return "", "", false
	}

	return match[2], match[1], true
}

func (s *Server) writeError(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	body := errorResponseBody(err)

	w.WriteHeader(body.StatusCode())

	if werr := json.MarshalEncode(jsontext.NewEncoder(w), body, deterministic); werr != nil {
		mlog.Warn(ctx, merr.New(ctx, "mrpc_write_error_failed", nil, werr))
	}
}

// errorResponseBody maps an error onto the cher.E sent to the client.
func errorResponseBody(err error) cher.E {
	if cErr, ok := errors.AsType[cher.E](err); ok {
		return cErr
	}

	if syntaxErr, ok := errors.AsType[*jsontext.SyntacticError](err); ok {
		return cher.New("invalid_json", cher.M{
			"error":  syntaxErr.Error(),
			"offset": syntaxErr.ByteOffset,
		})
	}

	semanticErr, ok := errors.AsType[*json.SemanticError](err)
	if !ok {
		return cher.New(cher.Unknown, nil)
	}

	// `name` is a JSON Pointer, where v1 used a dot-separated path
	meta := cher.M{"name": string(semanticErr.JSONPointer)}

	// both are documented as potentially unknown
	if semanticErr.GoType != nil {
		meta["expected"] = semanticErr.GoType.Kind().String()
	}
	if semanticErr.JSONKind != 0 {
		meta["actual"] = semanticErr.JSONKind.String()
	}

	return cher.New("invalid_json", meta)
}
