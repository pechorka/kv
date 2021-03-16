package web

// перечисление кодов ошибки сервиса
const (
	// ErrInternal внутренняя ошибка сервиса
	ErrInternal = 1
	// ErrNotFound - значение не найдено
	ErrNotFound = 2
	// ErrDecode - ошибка при чтении тела запроса
	ErrDecode = 3
)
