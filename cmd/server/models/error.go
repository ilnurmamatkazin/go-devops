package models

// RequestError структура, расширяющая стандартный error
type RequestError struct {
	StatusCode int
	Err        error
}

// Error реалзизация метода
func (r *RequestError) Error() string {
	return r.Err.Error()
}
