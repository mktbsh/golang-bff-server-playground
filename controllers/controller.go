package controllers

import "github.com/labstack/echo/v4"

type Controller struct{}

type LoginUser struct {
	ID string
}

func (c Controller) UseLoginUser(ctx echo.Context) *LoginUser {
	user := ctx.Get("login-user")
	// todo: validate
	if user != nil {
		return user.(*LoginUser)
	}
	return nil
}
