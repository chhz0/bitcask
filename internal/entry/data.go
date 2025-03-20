package entry

const (
	DataDeleted byte = iota
	DataNormal
	DataT
)

const (
	// HeaderSize 文件日志记录头大小
	// crc[4] + keySize[4] + valueSize[4] + type[1]
	HeaderSize = 4 + 4 + 4 + 1

	// MaxKeySize 文件日志记录最大key大小 2^32 - 1
	MaxKeySize = 1<<32 - 1
)

// Header 数据文件日志记录头结构, 大小为13字节
// crc[4] + keySize[4] + valueSize[4] + type[1]
type Header struct {
	CRC uint32 // crc校验码
	Ksz uint32 // keySize
	Vsz uint32 // valueSize -1 代表删除
}

// HeaderWithExp 数据文件日志记录头结构, 大小为29字节
// crc[4] + keySize[4] + valueSize[4] + type[1] + tstamp[8] + exp[8]
// todo: support timeout
type HeaderWithTimeout struct {
	*Header
	Tstamp uint64 // 时间戳
	Exp    uint64 // 过期时间
}

// DataEntry 文件日志记录结构
type DataEntry struct {
	*Header
	Key   []byte
	Value []byte
}
