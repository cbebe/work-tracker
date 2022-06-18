package worktracker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

type LayoutSpy struct {
	data any
}
type WithTitle struct {
	PageTitle string
}

func (l *LayoutSpy) Execute(wr io.Writer, data any) error {
	l.data = data
	return nil
}

type Test struct {
	p  string
	fn func(s *WorkServiceSpy) bool
}

type Handler func(t testing.TB, p string) *WorkServiceSpy

func testHandler(t *testing.T, handler Handler, tests []Test) {
	for _, tt := range tests {
		t.Run(tt.p, func(t *testing.T) {
			s := handler(t, tt.p)
			assert.True(t, tt.fn(s))
		})
	}
}

var tests = []Test{
	{"/", func(s *WorkServiceSpy) bool {
		return s.CalledGetWork
	}},
	{"/all", func(s *WorkServiceSpy) bool {
		return s.CalledGetWork
	}},
	{"/start", func(s *WorkServiceSpy) bool {
		return s.CalledStartWork
	}},
	{"/stop", func(s *WorkServiceSpy) bool {
		return s.CalledStopWork
	}},
}

func TestNewWorkHandler_OK(t *testing.T) {
	testHandler(t, mustBeOK, tests)
}

func TestNewWorkHandler_Error(t *testing.T) {
	testHandler(t, mustFail, tests)
}

func serveRequest(t testing.TB, p string, code int, err error) *WorkServiceSpy {
	t.Helper()
	s := &WorkServiceSpy{err: err}
	l := &LayoutSpy{}
	h := NewWorkHandler(s, l)
	response := httptest.NewRecorder()
	h.ServeHTTP(response, newGetRequest(p))
	assert.Equal(t, response.Code, code)
	return s
}

func mustBeOK(t testing.TB, p string) *WorkServiceSpy {
	return serveRequest(t, p, http.StatusOK, nil)
}

func mustFail(t testing.TB, p string) *WorkServiceSpy {
	return serveRequest(t, p, http.StatusInternalServerError, fmt.Errorf("oops"))
}

func newGetRequest(p string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, p, nil)
	return req
}
