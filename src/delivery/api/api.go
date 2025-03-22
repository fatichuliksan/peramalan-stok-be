package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"peramalan-stok-be/src/delivery/api/route"
	"peramalan-stok-be/src/helper"
	responseHelper "peramalan-stok-be/src/helper/response"
	validatorHelper "peramalan-stok-be/src/helper/validator"
	viperHelper "peramalan-stok-be/src/helper/viper"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

// NewAPI struct ...
type NewAPI struct {
	Echo      *echo.Echo
	Config    viperHelper.Interface
	Validator validatorHelper.Interface
	Response  responseHelper.Interface
	Printer   *message.Printer
	DB        *gorm.DB
}

// Register ...
func (t *NewAPI) Register() *NewAPI {
	t.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		// AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderAccept,
			"App-ID",
			"App-Name",
			"App-Platform",
			"App-Version",
			"Operating-System",
			"Client-Local-IP",
			"Client-Public-IP",
			"Client-Timezone",
			"Client-Timestamp",
			"Client-Device-ID",
			"Client-Manufacture",
			"Client-Brand",
			"Client-Model",
			"Client-Operating-System",
		},
	}))

	t.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	t.Echo.HTTPErrorHandler = t.HTTPErrorHandlerCustom
	t.Echo.Validator = t.Validator.Validator()

	// t.Echo.Logger.SetLevel(lasbstackLog.DEBUG)
	t.Echo.Use(middleware.Recover())
	t.Echo.HideBanner = true
	t.Echo.Debug = true
	// t.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		if c.Request().UserAgent() == "ELB-HealthChecker/2.0" {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// 	Format: `{"time":"${time_rfc3339}","request_id":"${id}",` +
	// 		`"user_agent":"${user_agent}","remote":"${remote_ip}","method":"${method}",` +
	// 		`"uri":"${uri}","status":"${status}","latency":"${latency_human}","error":"${error}"}` + "\n",
	// }))

	t.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogRequestID: true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		Skipper: func(c echo.Context) bool {
			if c.Request().UserAgent() == "ELB-HealthChecker/2.0" {
				return true
			}
			return false
		},
		BeforeNextFunc: func(c echo.Context) {
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := map[string]interface{}{
				"time":        v.StartTime.Format("2006-01-02T15:04:05Z07:00"),
				"request_id":  v.RequestID,
				"personal_id": c.Request().Header.Get("Personal-ID"),
				"user_agent":  v.UserAgent,
				"remote_ip":   v.RemoteIP,
				"method":      v.Method,
				"uri":         v.URI,
				"status":      v.Status,
				"latency":     fmt.Sprintf("%vms", v.Latency.Milliseconds()),
				"error":       v.Error,
			}
			req, _ := json.Marshal(fields)
			fmt.Println(string(req))
			// logger.Default().Println(string(req))
			return nil
		},
	}))

	// t.Echo.Use(customMiddleware.HeaderCheck(map[string]interface{}{
	// 	"App-ID":             []string{},
	// 	"App-Name":           []string{"sfa", "dms", "wms"},
	// 	"App-Platform":       []string{"web", "mobile"},
	// 	"App-Version":        []string{},
	// 	"Company-ID":         []string{},
	// 	"Personal-ID":        []string{},
	// 	"User-ID":            []string{},
	// 	"Territory-ID":       []string{},
	// 	"Territory-Area-ID":  []string{},
	// 	"Device-ID":          []string{},
	// 	"Parent-User-ID":     []string{},
	// 	"Parent-Personal-ID": []string{},
	// 	"Is-As-Sales":        []string{},
	// }, []string{"authentication", "apikey", "ping", "/auth/callback"}))

	// t.Echo.Use(customMiddleware.UserDevice(t.DB.DBMaster, []string{"authentication", "apikey", "ping", "/auth/callback"}))

	t.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := c.Request().Header.Get("Accept-Language")
			if lang == "" {
				c.Set("Printer", t.Printer)
			} else {
				c.Set("Printer", message.NewPrinter(t.getLanguageTag(lang)))
			}
			return next(c)
		}
	})

	helper := helper.Helper{}
	helper.Response = t.Response
	helper.Config = t.Config

	route := route.NewRoute{
		Echo:     t.Echo,
		Config:   t.Config,
		Response: t.Response,
		DB:       t.DB,
		Helper:   helper,
	}
	route.Register()
	if helper.Config.GetBool("app.debug") {
		data, err := json.MarshalIndent(route.Echo.Routes(), "", "  ")
		if err != nil {
			log.Println(err)
		}
		os.WriteFile("routes.json", data, 0644)
	}
	return t
}

// HTTPErrorHandlerCustom ...
func (t *NewAPI) HTTPErrorHandlerCustom(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		if t.Echo.Debug {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			report = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}

	castedObject, ok := err.(validator.ValidationErrors)
	if ok {
		report.Code = http.StatusBadRequest
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				report.Message = fmt.Sprintf("%s is not valid email", err.Field())
			case "gte":
				report.Message = fmt.Sprintf("%s value must be greater than equal %s", err.Field(), err.Param())
			case "lte":
				report.Message = fmt.Sprintf("%s value must be lower than equal %s", err.Field(), err.Param())
			case "gt":
				report.Message = fmt.Sprintf("%s value must be greater than %s", err.Field(), err.Param())
			case "lt":
				report.Message = fmt.Sprintf("%s value must be lower than %s", err.Field(), err.Param())
			case "max":
				report.Message = fmt.Sprintf("%s maximum is %s digit", err.Field(), err.Param())
			case "min":
				report.Message = fmt.Sprintf("%s minimum is %s digit", err.Field(), err.Param())
			default:
				report.Message = fmt.Sprintf("%s field defined tag validation %s", err.Field(), err.Tag())
			}
			break
		}
	}
	t.Response.SendCustomResponse(c, report.Code, report.Message.(string), nil)
}

func (t *NewAPI) getLanguageTag(lang string) language.Tag {
	lt, _ := language.Parse(lang)
	tag, _, _ := message.DefaultCatalog.Matcher().Match(lt)
	return tag
}
