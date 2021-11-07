package output

type Output interface {
	Write(p []byte) (n int, err error)
	Close() error
}
