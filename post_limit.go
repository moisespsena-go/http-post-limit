package post_limit

import (
	"context"
	"errors"
	"net/http"
)

var (
	DefaultMaxPostSize            int64 = 1024 * 1024 * 10 // 10Mb
	DefaultPostLimitFailedHandler       = PostLimitFailedHandler(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "max post size exceeded", http.StatusBadRequest)
	})
)

const (
	PostSizeKey contextKey = "post-size"
)

type (
	contextKey             string
	FormParser             func(r *http.Request) (err error)
	PostLimitFailedHandler = http.HandlerFunc

	Middleware struct {
		http.Handler
		*Opts
	}

	Opts struct {
		MaxPostSize  int64
		ErrorHandler http.HandlerFunc
	}
)

func New(handler http.Handler, opt ...*Opts) *Middleware {
	var opts *Opts
	for _, opts = range opt {
	}
	if opts == nil {
		opts = &Opts{}
	}
	if opts.MaxPostSize == 0 {
		opts.MaxPostSize = DefaultMaxPostSize
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = DefaultPostLimitFailedHandler
	}
	return &Middleware{
		handler,
		opts,
	}
}

func ParseForm(r *http.Request) (err error) {
	var maxPostSize = MaxPostSizeOf(r)
	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		err = r.ParseForm()
	case "multipart/form-data":
		err = r.ParseMultipartForm(maxPostSize)
	default:
		err = errors.New("bad content-type")
	}
	return nil
}

func (m *Middleware) ServeHttp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		if r.ContentLength > m.MaxPostSize {
			m.ErrorHandler(w, r)
			return
		}
	}
	r = r.WithContext(context.WithValue(r.Context(), PostSizeKey, m.MaxPostSize))
	m.Handler.ServeHTTP(w, r)
}

func MaxPostSizeOf(r *http.Request) int64 {
	var maxPostSize = DefaultMaxPostSize
	if v := r.Context().Value(PostSizeKey); v != nil {
		if vi := v.(int64); vi > 0 {
			maxPostSize = vi
		}
	}
	return maxPostSize
}
