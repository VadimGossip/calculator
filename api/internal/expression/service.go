package expression

type service struct {
}

type Service interface {
	//Split
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}
