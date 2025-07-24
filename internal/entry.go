package internal

type Entry struct {
	CRC    uint32
	Tstamp int64
	Key    []byte
	Val    []byte
}

type KV struct {
	K []byte
	V []byte
}
