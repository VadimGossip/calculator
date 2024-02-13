package parser

type storage struct {
	tokens         []string
	current        int
	subExpressions []subExpression
}

var _ Service = (*service)(nil)

func NewStorage(tokens []string) *storage {
	return &storage{tokens: tokens}
}

func (s *storage) getNextToken() string {
	if s.current < len(s.tokens) {
		token := s.tokens[s.current]
		s.current++
		return token
	}
	return ""
}

func (s *storage) getNewSubExpressionId() int64 {
	return int64(len(s.subExpressions) + 1)
}

func (s *storage) appendSubExpression(e subExpression) {
	s.subExpressions = append(s.subExpressions, e)
}

func (s *storage) getSubExpressions() []subExpression {
	return s.subExpressions
}

func (s *storage) decrementCurrent() {
	s.current--
}
