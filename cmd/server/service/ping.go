package service

// Ping функция проверки соединения с базой данных.
func (s *ServiceMetric) Ping() (err error) {
	return s.storage.Ping()
}
