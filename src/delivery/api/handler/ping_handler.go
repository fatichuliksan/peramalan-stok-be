package handler

import (
	"peramalan-stok-be/src/helper/response"
	"peramalan-stok-be/src/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PingHandler ...
type PingHandler struct {
	Response response.Interface
	DB       *gorm.DB
}

// Get ...
func (t *PingHandler) Get(c echo.Context) error {
	pingUsecase := usecase.NewPing(t.DB)
	res, err := pingUsecase.Ping("pong")

	if err != nil {
		return t.Response.SendError(c, err.Error(), res)
	}

	return t.Response.SendSuccess(c, res, nil)

	// return t.Response.SendSuccess(c, c.Get("Printer").(*message.Printer).Sprintf("%v, %d task(s) remaining!", "fatiq", 2), res)
	// return t.Response.SendSuccess(c, c.Get("Printer").(*message.Printer).Sprintf("%d task(s) remaining!", 3), res)
	// return t.Response.SendSuccess(c, c.Get("Printer").(*message.Printer).Sprintf("Hello Word"), res)
}
