package calculatorapi

import (
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/orchestrator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	CreateExpression(c *gin.Context)
}

type controller struct {
	orchestratorService orchestrator.Service
}

var _ Controller = (*controller)(nil)

func NewController(orchestratorService orchestrator.Service) *controller {
	return &controller{orchestratorService: orchestratorService}
}

func (ctrl *controller) CreateExpression(c *gin.Context) {
	var req CreateExpressionRequest
	if err := c.BindJSON(&req); err != nil {
		errMsg := fmt.Sprintf("Parse request error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateExpression",
		}).Error(errMsg)
		c.JSON(http.StatusBadRequest, CreateExpressionResponse{Error: errMsg, Status: http.StatusBadRequest})
		return
	}
	//validateService
	id, err := ctrl.orchestratorService.RegisterExpression(c.Request.Context(), req.ExpressionValue)
	if err != nil {
		errMsg := fmt.Sprintf("create expression error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateExpression",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, CreateExpressionResponse{Status: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, CreateExpressionResponse{Id: id, Status: http.StatusOK})
}
