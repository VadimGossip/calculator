package parser

import (
	"github.com/VadimGossip/calculator/api/internal/domain"
	"strconv"
	"strings"
	"unicode"
)

type subExpression struct {
	id        int64
	first     *float64
	second    *float64
	operation string
	prevSe1   *subExpression
	prevSe2   *subExpression
}

type service struct {
	tokens         []string
	current        int
	subExpressions []subExpression
}

type Service interface {
	ParseExpression(e domain.Expression) []domain.SubExpression
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) tokenize(input string) []string {
	var tokens []string
	var currentToken string

	for idx, char := range input {
		if unicode.IsDigit(char) || char == '.' || char == '-' && idx == 0 {
			currentToken += string(char)
		} else if strings.ContainsAny(string(char), "+-*/") {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		}
	}
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}
	return tokens
}

func (s *service) parseNumber(storage *storage) *float64 {
	token := storage.getNextToken()
	number, _ := strconv.ParseFloat(token, 64)
	return &number
}

func (s *service) parseExpression(storage *storage) {
	se1, first := s.parseHighCostOperations(storage)
	for {
		operator := storage.getNextToken()
		if operator != "+" && operator != "-" {
			storage.decrementCurrent()
			break
		}
		se2, second := s.parseHighCostOperations(storage)

		current := subExpression{id: storage.getNewSubExpressionId(),
			first:     first,
			second:    second,
			operation: operator}
		if se1 != nil {
			current.first = nil
			current.prevSe1 = se1

		}
		se1 = &current
		if se2 != nil {
			current.second = nil
			current.prevSe2 = se2
		}
		storage.appendSubExpression(current)
	}
}

func (s *service) parseHighCostOperations(storage *storage) (*subExpression, *float64) {
	first := s.parseNumber(storage)
	var counter int
	var prev *subExpression
	for {
		operator := storage.getNextToken()
		if operator != "*" && operator != "/" {
			storage.current--
			return prev, first
		}
		second := s.parseNumber(storage)
		current := subExpression{id: storage.getNewSubExpressionId(),
			first:     first,
			second:    second,
			operation: operator,
		}
		if counter > 0 {
			current.first = nil
			current.prevSe1 = prev
		}
		prev = &current
		storage.appendSubExpression(current)
		counter++
	}
}

func (s *service) ParseExpression(e domain.Expression) []domain.SubExpression {
	var result []domain.SubExpression
	st := NewStorage(s.tokenize(e.Value))
	s.parseExpression(st)
	seItems := st.getSubExpressions()
	for idx, item := range seItems {
		se := domain.SubExpression{
			Id:            item.id,
			ExpressionsId: 0,
			Val1:          item.first,
			Val2:          item.second,
			Operation:     item.operation,
			AgentId:       0,
			IsLast:        idx == len(seItems),
		}
		if item.prevSe1 != nil {
			se.SubExpressionId1 = &item.prevSe1.id
		}
		if item.prevSe2 != nil {
			se.SubExpressionId2 = &item.prevSe2.id
		}
		result = append(result, se)
	}
	return result
}
