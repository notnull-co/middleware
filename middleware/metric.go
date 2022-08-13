package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/bexstech/middleware/context"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metrics *prometheus.HistogramVec
	once    sync.Once
)

func Monitor(serviceName string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := time.Now()

			var err error
			var errMsg string
			if err = next(c); err != nil {
				c.Error(err)
				errMsg = err.Error()
			}

			once.Do(func() {
				metrics = promauto.NewHistogramVec(prometheus.HistogramOpts{
					Name: fmt.Sprintf("%s_rest_interface", serviceName),
					Help: fmt.Sprintf("%s rest interface requests", serviceName),
				}, []string{"endpoint", "status_code", "user_id", "method", "err"})
			})

			metrics.With(prometheus.Labels{
				"endpoint":    c.Request().RequestURI,
				"status_code": strconv.Itoa(c.Response().Status),
				"user_id":     context.GetUserId(c),
				"method":      c.Request().Method,
				"err":         errMsg,
			}).Observe(time.Since(t).Seconds())

			return err
		}
	}
}
