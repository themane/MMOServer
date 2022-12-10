package exceptions

type NoSuchCombinationError struct{}

func (m *NoSuchCombinationError) Error() string {
	return "no such combination exists"
}

type AlreadyExistsError struct {
	Message string
}

func (m *AlreadyExistsError) Error() string {
	return m.Message
}
