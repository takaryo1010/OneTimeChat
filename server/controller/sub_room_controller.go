package controller

import "github.com/labstack/echo"

func GetCookie(c echo.Context, key string) string {
	cookie, err := c.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
