package index

type Iterator interface {
	Rewind()
	Valid() bool
	Next()
	Key() []byte
	Value() *Entry
	Close()
}
