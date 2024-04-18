package calculatorapi

import (
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/expression"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	CreateExpression(c *gin.Context)
	GetAllExpressions(c *gin.Context)
	GetAllAgents(c *gin.Context)
	SaveOperationDurations(c *gin.Context)
	GetAllOperationDurations(c *gin.Context)
}

type controller struct {
	expressionService expression.Service
}

var _ Controller = (*controller)(nil)

func NewController(expressionService expression.Service) *controller {
	return &controller{expressionService: expressionService}
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
	e := &domain.Expression{
		ReqUid: req.ReqUid,
		Value:  req.ExpressionValue,
	}

	simplifiedValue, err := ctrl.expressionService.ValidateAndSimplify(e.Value)
	if err != nil {
		e.State = domain.ExpressionStateError
		e.ErrorMsg = fmt.Sprintf("validate expression error: %s", err)
	} else {
		e.State = domain.ExpressionStateNew
		e.Value = simplifiedValue
	}

	if err = ctrl.expressionService.RegisterExpression(c.Request.Context(), e); err != nil {
		errMsg := fmt.Sprintf("create expression error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateExpression",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, CreateExpressionResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, CreateExpressionResponse{Expression: e, Error: e.ErrorMsg, Status: http.StatusOK})
}

func (ctrl *controller) GetAllExpressions(c *gin.Context) {
	expressions, err := ctrl.expressionService.GetExpressions(c.Request.Context())
	if err != nil {
		errMsg := fmt.Sprintf("get all expressions: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateAllExpressions",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, GetExpressionsResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, GetExpressionsResponse{Expressions: expressions, Status: http.StatusOK})
}

func (ctrl *controller) GetAllAgents(c *gin.Context) {
	agents, err := ctrl.expressionService.GetAgents(c.Request.Context())
	if err != nil {
		errMsg := fmt.Sprintf("get all agents: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateAllAgents",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, GetAgentsResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, GetAgentsResponse{Agents: agents, Status: http.StatusOK})
}

func (ctrl *controller) SaveOperationDurations(c *gin.Context) {
	operationDurations := make(map[string]uint32)
	if err := c.BindJSON(&operationDurations); err != nil {
		errMsg := fmt.Sprintf("Parse request error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "SaveOperationDurations",
		}).Error(errMsg)
		c.JSON(http.StatusBadRequest, CommonResponse{Error: errMsg, Status: http.StatusBadRequest})
		return
	}
	if err := ctrl.expressionService.SaveOperationDurations(c.Request.Context(), operationDurations); err != nil {
		errMsg := fmt.Sprintf("save operation durations error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "SaveOperationDurations",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, CommonResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{Status: http.StatusOK})
}

func (ctrl *controller) GetAllOperationDurations(c *gin.Context) {
	operationDurations, err := ctrl.expressionService.GetOperationDurations(c.Request.Context())
	if err != nil {
		errMsg := fmt.Sprintf("get all operation durations: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "GetAllOperationDurations",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, GetOperationDurationsResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, GetOperationDurationsResponse{OperationDuration: operationDurations, Status: http.StatusOK})
}
