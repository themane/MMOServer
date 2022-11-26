package mongoRepository

type NoSuchCombinationError struct{}

func (m *NoSuchCombinationError) Error() string {
	return "no such combination exists"
}
