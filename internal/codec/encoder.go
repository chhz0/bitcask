package codec

import (
	"encoding/binary"
	"time"
)

const (
	crcSize    = 4
	tstampSize = 8
	keySize    = 4
	valueSize  = 4

	headerSize = crcSize + tstampSize + keySize + valueSize

	bufTstampEndIdx = crcSize + tstampSize
	bufKszEndIdx    = bufTstampEndIdx + keySize
	bufVszEndIdx    = bufKszEndIdx + valueSize
)

func Encode(key, val []byte, deleted bool) []byte {
	ksz := len(key)
	vsz := len(val)

	if deleted {
		vsz = 0
	}

	buf := make([]byte, headerSize+ksz+vsz)

	binary.BigEndian.PutUint64(buf[crcSize:bufTstampEndIdx], uint64(time.Now().Unix()))
	binary.BigEndian.PutUint32(buf[bufTstampEndIdx:bufKszEndIdx], uint32(ksz))
	binary.BigEndian.PutUint32(buf[bufKszEndIdx:bufVszEndIdx], uint32(vsz))

	copy(buf[bufVszEndIdx:bufVszEndIdx+ksz], key)

	if !deleted {
		copy(buf[bufVszEndIdx+ksz:bufVszEndIdx+ksz+vsz], val)
	}

	binary.BigEndian.PutUint32(buf[:crcSize], calculateCRC(buf[crcSize:]))

	return buf
}
