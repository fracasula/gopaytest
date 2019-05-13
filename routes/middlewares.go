package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
)

const linksSelfKey = "CTX_LINKS_SELF"
const pageSizeKey = "CTX_PAGE_SIZE"
const pageNumberKey = "CTX_PAGE_NUMBER"

func With(r *http.Request) *ctxReader {
	return &ctxReader{request: r}
}

type ctxReader struct {
	request *http.Request
}

func (c *ctxReader) SelfLink() string {
	return c.request.Context().Value(linksSelfKey).(string)
}

func (c *ctxReader) PageValues() (pageSize, pageNumber, offset int) {
	pageSize = c.request.Context().Value(pageSizeKey).(int)
	pageNumber = c.request.Context().Value(pageNumberKey).(int)
	offset = (pageNumber - 1) * pageSize
	return
}

func linksSelfMiddleware(baseURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, linksSelfKey, baseURL+r.URL.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func paginationMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pageSize := "10"
		pageNumber := "1"

		pageSizeFromHeader := r.URL.Query().Get("$size")
		if pageSizeFromHeader != "" {
			pageSize = pageSizeFromHeader
		}
		pageNoFromHeader := r.URL.Query().Get("$page")
		if pageNoFromHeader != "" {
			pageNumber = pageNoFromHeader
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			RenderError(w, r, "Invalid page size", http.StatusBadRequest)
			return
		}

		pageNoInt, err := strconv.Atoi(pageNumber)
		if err != nil || pageNoInt < 1 {
			RenderError(w, r, "Invalid page number", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, pageSizeKey, pageSizeInt)
		ctx = context.WithValue(ctx, pageNumberKey, pageNoInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func loggerMiddleware(l *log.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: l, NoColor: false,
	})
}
