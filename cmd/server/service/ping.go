package service

func (s *Service) Ping() (err error) {
	return s.repository.Ping()
}
