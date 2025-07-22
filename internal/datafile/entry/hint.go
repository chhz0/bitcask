package entry

type HintHeader struct {
	Version uint32
}

type HintEntry struct {
	FileID    int
	Timestamp uint64
	KeySize   uint32
	ValueSize uint32
	ValuePos  int64
	Key       []byte
}
