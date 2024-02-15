package app

import "github.com/VadimGossip/calculator/api/internal/api/server/calculatorapi"

func initHttpRouter(app *App) {
	s := app.apiHttpServer

	c := calculatorapi.NewController(app.expressionService)
	s.GET("/expression", c.GetAllExpressions)
	s.POST("/expression", c.CreateExpression)

	s.POST("/sub_expression/start", c.StartSubExpressionEval)
	s.POST("/sub_expression/stop", c.StopSubExpressionEval)

	s.GET("/agent", c.GetAllAgents)
	s.POST("/agent", c.AgentHeartbeat)

	s.GET("/duration", c.GetAllOperationDurations)
	s.POST("/duration", c.SaveOperationDurations)
}
