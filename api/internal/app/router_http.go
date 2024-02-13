package app

import "github.com/VadimGossip/calculator/api/internal/api/server/calculatorapi"

func initHttpRouter(app *App) {
	s := app.apiHttpServer

	c := calculatorapi.NewController(app.managerService)
	s.GET("/expression", c.GetAllExpressions)
	s.POST("/expression", c.CreateExpression)

	s.POST("/sub_expression/start", c.StartSubExpressionEval)

	s.GET("/agent", c.GetAllAgents)
	s.POST("/agent", c.AgentHeartbeat)

	s.GET("/duration", c.GetAllOperationDurations)
	s.POST("/duration", c.SaveOperationDurations)
}
