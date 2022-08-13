package wrapper

import (
	_ctx "bitbucket.org/bexstech/middleware/context"

	"github.com/labstack/echo"
)

func UserId(fn func(c echo.Context, userId string) error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return fn(c, _ctx.GetUserId(c))
	}
}
