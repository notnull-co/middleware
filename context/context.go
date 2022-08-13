package context

import (
	"github.com/labstack/echo"
)

func SetUserId(c echo.Context, merchantId string) {
	c.Set("user_id", merchantId)
}

func GetUserId(c echo.Context) string {
	if merchantId := c.Get("user_id"); merchantId != nil {
		if merchantId, ok := merchantId.(string); ok {
			return merchantId
		}
	}
	return ""
}
