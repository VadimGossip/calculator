package app

import "github.com/VadimGossip/calculator/api/internal/api/server/calculatorapi"

func initHttpRouter(app *App) {
	s := app.apiHttpServer

	c := calculatorapi.NewController(app.expressionService, app.authService)

	authApi := s.Group("/api/v1")
	{
		authApi.POST("/register", c.Register)
		authApi.POST("/login", c.Login)
		authApi.POST("/refresh", c.Refresh)
	}

	expressionApi := s.Group("/expression")
	expressionApi.Use(c.AuthMiddleware())
	{
		expressionApi.GET("", c.GetAllExpressions)
		expressionApi.POST("", c.CreateExpression)
	}

	agentApi := s.Group("/agent")
	agentApi.Use(c.AuthMiddleware())
	{
		agentApi.GET("", c.GetAllAgents)
	}

	durationApi := s.Group("/duration")
	durationApi.Use(c.AuthMiddleware())
	{
		durationApi.GET("", c.GetAllOperationDurations)
		durationApi.POST("", c.SaveOperationDurations)
	}
}
