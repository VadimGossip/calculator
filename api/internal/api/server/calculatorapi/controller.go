package calculatorapi

import (
	"errors"
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/auth"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/expression"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	ctxUserID           = "userID"
)

type Controller interface {
	CreateExpression(c *gin.Context)
	GetAllExpressions(c *gin.Context)
	GetAllAgents(c *gin.Context)
	SaveOperationDurations(c *gin.Context)
	GetAllOperationDurations(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc
}

type controller struct {
	expressionService expression.Service
	authService       auth.Service
}

var _ Controller = (*controller)(nil)

func NewController(expressionService expression.Service, authService auth.Service) *controller {
	return &controller{expressionService: expressionService, authService: authService}
}

func (ctrl *controller) getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", domain.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", domain.ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", domain.ErrEmptyToken
	}

	return headerParts[1], nil
}

func (ctrl *controller) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := ctrl.getTokenFromRequest(c)
		if err != nil {
			errMsg := fmt.Sprintf("Parse token error: %s", err)
			logrus.WithFields(logrus.Fields{
				"request": "authMiddleware",
			}).Error(errMsg)
			c.JSON(http.StatusUnauthorized, CommonResponse{Error: errMsg, Status: http.StatusUnauthorized})
			return
		}

		id, err := ctrl.authService.ParseToken(t)
		if err != nil {
			errMsg := fmt.Sprintf("AccessToken invalid or expired: %s", err)
			logrus.WithFields(logrus.Fields{
				"request": "authMiddleware",
			}).Error(errMsg)
			c.JSON(http.StatusUnauthorized, CommonResponse{Error: errMsg, Status: http.StatusUnauthorized})
			return
		}
		c.Set(ctxUserID, id)
		c.Next()
	}
}

func (ctrl *controller) getUserId(c *gin.Context) (int64, error) {
	id, ok := c.Get(ctxUserID)
	if !ok {
		return 0, domain.ErrUserNotFound
	}
	idInt, ok := id.(int64)
	if !ok {
		return 0, domain.ErrUserIdInvalid
	}

	return idInt, nil
}

func (ctrl *controller) Register(c *gin.Context) {
	var u domain.User
	if err := c.BindJSON(&u); err != nil {
		errMsg := fmt.Sprintf("Parse request error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "Register",
		}).Error(errMsg)
		c.JSON(http.StatusBadRequest, RegisterUserResponse{Error: errMsg, Status: http.StatusBadRequest})
		return
	}
	if err := ctrl.authService.Register(c.Request.Context(), &u); err != nil {
		errMsg := fmt.Sprintf("register user error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "Register",
		}).Error(errMsg)
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, RegisterUserResponse{Id: u.Id, Status: http.StatusOK})
}

func (ctrl *controller) Login(c *gin.Context) {
	var cr domain.Credentials
	if err := c.BindJSON(&cr); err != nil {
		errMsg := fmt.Sprintf("Parse request error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "Login",
		}).Error(errMsg)
		c.JSON(http.StatusBadRequest, CommonResponse{Error: errMsg, Status: http.StatusBadRequest})
		return
	}

	accessToken, refreshToken, err := ctrl.authService.Login(c.Request.Context(), cr)
	if err != nil {
		errMsg := fmt.Sprintf("Search user error: %s", err)
		logrus.WithFields(logrus.Fields{
			"request": "Login",
		}).Error(errMsg)
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, CommonResponse{Error: domain.ErrUserNotFound.Error(), Status: http.StatusBadRequest})
			return
		}
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{Error: errMsg, Status: http.StatusInternalServerError})
		return
	}
	refreshTokenTTL := ctrl.authService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.TokenResponse{Token: accessToken})
}

func (ctrl *controller) Refresh(c *gin.Context) {
	cookieRefresh, err := c.Cookie("refresh-token")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "Refresh",
		}).Error(err)
		c.JSON(http.StatusBadRequest, CommonResponse{Error: "parse refresh token error", Status: http.StatusBadRequest})
		return
	}

	accessToken, refreshToken, err := ctrl.authService.RefreshTokens(c.Request.Context(), cookieRefresh)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "Refresh",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, CommonResponse{Error: "refresh tokens error", Status: http.StatusBadRequest})
		return
	}
	refreshTokenTTL := ctrl.authService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.TokenResponse{Token: accessToken})
}

func (ctrl *controller) CreateExpression(c *gin.Context) {
	userId, err := ctrl.getUserId(c)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "CreateExpression",
		}).Error(err)
		c.JSON(http.StatusUnauthorized, CreateExpressionResponse{Error: err.Error(), Status: http.StatusUnauthorized})
		return
	}

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
		UserId: userId,
		ReqUid: req.ReqUid,
		Value:  req.ExpressionValue,
	}

	simplifiedValue, err := ctrl.expressionService.ValidateAndSimplify(e.Value)
	if err != nil {
		e.ErrorMsg = fmt.Sprintf("validate expression error: %s", err)
	} else {
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
	userId, err := ctrl.getUserId(c)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "GetAllExpressions",
		}).Error(err)
		c.JSON(http.StatusUnauthorized, CreateExpressionResponse{Error: err.Error(), Status: http.StatusUnauthorized})
		return
	}

	expressions, err := ctrl.expressionService.GetExpressions(c.Request.Context(), userId)
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
	_, err := ctrl.getUserId(c)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "GetAllAgents",
		}).Error(err)
		c.JSON(http.StatusUnauthorized, CreateExpressionResponse{Error: err.Error(), Status: http.StatusUnauthorized})
		return
	}

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
	_, err := ctrl.getUserId(c)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "SaveOperationDurations",
		}).Error(err)
		c.JSON(http.StatusUnauthorized, CreateExpressionResponse{Error: err.Error(), Status: http.StatusUnauthorized})
		return
	}

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
	_, err := ctrl.getUserId(c)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "GetAllOperationDurations",
		}).Error(err)
		c.JSON(http.StatusUnauthorized, CreateExpressionResponse{Error: err.Error(), Status: http.StatusUnauthorized})
		return
	}

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
