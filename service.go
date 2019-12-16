package awesomeProject

type IService interface {
	AddCar(car *Car) error
	GetCar() ([]Car, error)
	UpdateCar(car *Car) error
	DeleteCar(id int) error
}

type Service struct {
}

func (s *Service) AddCar(car *Car) error {
	return nil
}
