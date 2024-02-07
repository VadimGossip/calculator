package app

import "github.com/VadimGossip/calculator/api/internal/api/server/calculatorapi"

func initHttpRouter(app *App) {
	s := app.apiHttpServer

	c := calculatorapi.NewController(app.managerService)
	s.GET("/expression", c.GetAllExpressions)
	s.POST("/expression", c.CreateExpression)
}
