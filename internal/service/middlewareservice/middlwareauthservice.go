package middlewareservice

import (
	"net/http"

	"github.com/LoaltyProgramm/to-do-app/internal/config"
	"github.com/LoaltyProgramm/to-do-app/internal/service/authservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/cookieservice"
)

type MiddlewareService interface {
	MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc
}

type middlewareService struct {
	jwtService authservice.JWTService
	cookieservice cookieservice.CookieService
	cfg config.Config
}

func NewMiddlewareService(jwtService authservice.JWTService, cookieservice cookieservice.CookieService, cfg config.Config) MiddlewareService {
	return &middlewareService{
		jwtService: jwtService,
		cookieservice: cookieservice,
	}
}

func (m *middlewareService) MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(m.cfg.TodoPassword) != 0 {
			token, err := m.cookieservice.GetCookie(w, r)
			if err != nil || token == "" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			valid, err := m.jwtService.ParseToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if !valid {
				http.Error(w, "Token invalid", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}