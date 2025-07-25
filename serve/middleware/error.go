package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/warjiang/page-spy-api/api/room"
	"github.com/warjiang/page-spy-api/serve/common"
)

func Error() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				res := common.NewErrorResponse(err)
				if res.Code == room.ServeError {
					return c.JSON(http.StatusInternalServerError, res)
				}

				return c.JSON(http.StatusBadRequest, res)
			}
			return nil
		}
	}
}
