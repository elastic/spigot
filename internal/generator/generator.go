package generator

type Generator interface {
	Next() ([]byte, error)
}
