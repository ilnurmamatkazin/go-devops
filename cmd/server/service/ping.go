package service

func (s *ServiceMetric) Ping() (err error) {
	return s.storage.Ping()
}
