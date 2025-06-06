package middleware

import (
	"log/slog"
	"net/http"
)

type Catcher struct {
	http.ResponseWriter

	status int
	len    int
}

func NewCatcher(w http.ResponseWriter) *Catcher {
	return &Catcher{
		ResponseWriter: w,
		status:         http.StatusOK, // Default status code
		len:            0,
	}
}

func (c *Catcher) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

func (c *Catcher) Write(b []byte) (int, error) {
	n, err := c.ResponseWriter.Write(b)
	c.len += n
	return n, err
}

func (c *Catcher) Header() http.Header {
	return c.ResponseWriter.Header()
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		catcher := NewCatcher(w)

		// Call the next handler in the chain
		next.ServeHTTP(catcher, r)

		slog.Info(
			"http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", catcher.status,
			"length", catcher.len,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"referer", r.Referer(),
		)
	})
}
