package validation

import (
	"fmt"
	"regexp"
	"strings"
)

type service struct {
}

type Service interface {
	ValidateAndSimplify(value string) (string, error)
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) isValid(expr string) bool {
	pattern := `^[-+]?\d*\.?\d+([\/*+-]\d*\.?\d+)*$`
	re := regexp.MustCompile(pattern)
	if re.MatchString(expr) && s.hasOperationSign(expr) {
		return true
	}
	return false
}

func (s *service) hasOperationSign(input string) bool {
	re := regexp.MustCompile(`[/*+-]`)
	return re.MatchString(strings.TrimLeft(input, "-"))
}

func (s *service) ValidateAndSimplify(value string) (string, error) {
	simplified := strings.ReplaceAll(value, " ", "")
	if len(simplified) > 50 {
		return value, fmt.Errorf("invalid expression value. max length must be less than %d", 50)
	}

	if !s.isValid(simplified) {
		return value, fmt.Errorf("invalid expression value format only +-/* operators and numbers allowed")
	}
	return simplified, nil
}
