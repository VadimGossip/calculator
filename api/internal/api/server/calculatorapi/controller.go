package calculatorapi

import (
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/manager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	CreateExpression(c *gin.Context)
}

type controller struct {
	managerService manager.Service
}

var _ Controller = (*controller)(nil)

func NewController(managerService manager.Service) *controller {
	return &controller{managerService: managerService}
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
	id, err := ctrl.managerService.RegisterExpression(c.Request.Context(), req.ExpressionValue)
	if err != nil {
		errMsg := fmt.Sprintf("create expression error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "CreateExpression",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, CreateExpressionResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, CreateExpressionResponse{Id: id, Status: http.StatusOK})
}

func (ctrl *controller) GetAllExpressions(c *gin.Context) {
	expressions, err := ctrl.managerService.GetExpressions(c.Request.Context())
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
