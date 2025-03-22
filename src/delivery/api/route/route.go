package route

import (
	"encoding/json"
	"errors"
	"peramalan-stok-be/src/helper"
	"peramalan-stok-be/src/helper/logger"
	responseHelper "peramalan-stok-be/src/helper/response"
	viperHelper "peramalan-stok-be/src/helper/viper"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// NewRoute Handler
type NewRoute struct {
	Helper   helper.Helper
	Echo     *echo.Echo
	Config   viperHelper.Interface
	Response responseHelper.Interface
	DB       *gorm.DB
}

// Register ...
func (t *NewRoute) Register() {
	groupV1 := t.Echo.Group("v1")
	customMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeStarted := time.Now()
			err := next(c)
			status := c.Response().Status
			ip := c.RealIP()
			httpErr := new(echo.HTTPError)
			if errors.As(err, &httpErr) {
				status = httpErr.Code
			}
			// logger.Default().Println(echo.HeaderXRequestID)
			// manualLog := logfile.MainLog("req")
			fields := map[string]interface{}{
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				"user_agent": c.Request().UserAgent(),
				"remote":     ip,
				"method":     c.Request().Method,
				"path":       c.Request().URL.Path,
				"query":      c.Request().URL.RawQuery,
				"status":     status,
				"latency":    int64(time.Since(timeStarted) / time.Millisecond),
			}
			req, _ := json.Marshal(fields)
			logger.Default().Println(string(req))
			if err != nil {
				// logger.Default().Println("on error: ", err)

				return err
			}
			return nil

		}
	}

	t.Echo.Use(customMiddleware)

	t.PingRoute(groupV1.Group("/ping"))
	t.MainRoute(groupV1.Group("/main"))
}
