package cookieservice

import (
	"fmt"
	"net/http"
)

type CookieService interface {
	GetCookie(w http.ResponseWriter, r *http.Request) (string, error)
}

type cookieService struct{}

func NewCookieService() CookieService {
	return &cookieService{}
}

func (c *cookieService) GetCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", fmt.Errorf("get cookie fail: %w", err)
	}

	return cookie.Value, nil
}