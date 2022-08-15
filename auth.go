package middleware

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Validate func(echo.Context, map[string]interface{}) error

func Authenticator(validationFn Validate, validators ...*auth0.JWTValidator) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var token *jwt.JSONWebToken
			var err error

			for _, validator := range validators {
				token, err = validator.ValidateRequest(c.Request())

				if err == nil {
					break
				}
			}

			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			var claims map[string]interface{}
			if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			if validationFn != nil {
				if err := validationFn(c, claims); err != nil {
					return c.NoContent(http.StatusForbidden)
				}
			}

			return next(c)
		}
	}
}
