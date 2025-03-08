package entry

// Index is the key directory entry in memory.
type Index struct {
	Fid    uint32 // file id
	Offset int64  // offset in file
	Size   uint64 // size of key
}

// IndexWithTimeout is the key directory entry with expiration time in memory.
type IndexWithTimeout struct {
	Index
	Tstamp uint64 // timestamp
	Exp    uint64 // expiration time
}
