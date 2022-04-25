package middlewareprovider

import (
	"net/http"

	"Assignment/models"
	"Assignment/providers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)

type Middleware struct {
	DB *sqlx.DB
}

func NewMiddleware(arDB *sqlx.DB) providers.MiddlewareProvider {
	return Middleware{
		DB: arDB,
	}
}

func (m Middleware) GetUserContext(req *http.Request) *models.UserContext {
	return req.Context().Value(models.UserContextKey).(*models.UserContext)
}

func (m Middleware) Default() chi.Middlewares {
	return chi.Chain(
		corsOptions().Handler,
		middleware.RequestID,
		middleware.RequestLogger(NewStructuredLogger()),
		middleware.Recoverer,
	)
}

func (m Middleware) AUTH() chi.Middlewares {
	return chi.Chain(
		authMiddleware(m.DB),
		checkForEmptyContext(),
	)
}
