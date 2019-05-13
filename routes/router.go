package routes

import (
	"gopaytest/container"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func NewRouter(c container.Container) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
		middleware.Timeout(30*time.Second),
		render.SetContentType(render.ContentTypeJSON),
		loggerMiddleware(c.Logger()),
		linksSelfMiddleware(c.BaseURL()),
		paginationMiddleware,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/payments", NewPaymentsRouter(
			c.PaymentsRepository(),
			c.Logger(),
		))
	})

	return router
}

func WritePaginationHeaders(w http.ResponseWriter, pageSize, pageNumber, count int) {
	pageCount := 0
	if count > 0 {
		pageCount = int(math.Ceil(float64(count) / float64(pageSize)))
	}

	w.Header().Set("X-Page-Count", strconv.Itoa(pageCount))
	w.Header().Set("X-Page-Number", strconv.Itoa(pageNumber))
	w.Header().Set("X-Page-Size", strconv.Itoa(pageSize))
	w.Header().Set("X-Total-Count", strconv.Itoa(count))
}
