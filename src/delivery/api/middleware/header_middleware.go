package middleware

import (
	"encoding/json"
	"peramalan-stok-be/src/helper/logger"
	"peramalan-stok-be/src/helper/response"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// HeaderCheck ...
func HeaderCheck(headerList map[string]interface{}, uriSkipper []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, uri := range uriSkipper {
				if strings.Contains(c.Request().RequestURI, uri) {
					return next(c)
				}
			}

			resp := response.NewResponse()
			for index, val := range headerList {
				isExist := c.Request().Header.Get(index)
				if isExist == "" {
					return resp.SendUnauthorized(c, index+" is required", nil)
				}

				headerVal := val.([]string)
				if len(headerVal) > 0 {
					check := func(slice []string, val string) bool {
						for _, item := range slice {
							if item == val {
								return true
							}
						}
						return false
					}(headerVal, isExist)
					if !check {
						return resp.SendBadRequest(c, index+" unknown header value", nil)
					}
				}

				if index == "Company-ID" || index == "Personal-ID" || index == "User-ID" || index == "Territory-ID" || index == "Territory-Area-ID" {
					val, _ := strconv.Atoi(isExist)
					c.Set(index, uint(val))
				} else {
					logger.Default().Println(index, isExist)
					c.Set(index, isExist)
				}
			}

			// roles
			roles := c.Request().Header.Get("Roles")
			var tempRoles []map[string]interface{}
			err := json.Unmarshal([]byte(roles), &tempRoles)
			if err != nil {
				logger.Default().Println(err)
			}
			c.Set("Roles", tempRoles)

			parentUserIDString := c.Request().Header.Get("Parent-User-ID")
			parentPersonalIDString := c.Request().Header.Get("Parent-Personal-ID")
			isAsSalesString := c.Request().Header.Get("Is-As-Sales")

			parentUserID, _ := strconv.Atoi(parentUserIDString)
			c.Set("Parent-User-ID", uint(parentUserID))

			parentPersonalID, _ := strconv.Atoi(parentPersonalIDString)
			c.Set("Parent-Personal-ID", uint(parentPersonalID))

			isAsSales, _ := strconv.ParseBool(isAsSalesString)
			c.Set("Is-As-Sales", isAsSales)

			return next(c)
		}
	}
}
