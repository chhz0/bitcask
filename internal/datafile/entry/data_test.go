package entry

import (
	"encoding/binary"
	"hash/crc32"
	"testing"

	"errors"

	errs "github.com/chhz0/go-bitcask/internal/errors"
)

func TestDataEntry_Encode(t *testing.T) {
	tests := []struct {
		name    string
		entry   *DataEntry
		want    []byte
		wantErr bool
	}{
		{
			name: "normal case",
			entry: &DataEntry{
				DataHeader: &DataHeader{
					Tstamp: 1234567890,
					Flag:   1,
				},
				Key:   []byte("test"),
				Value: []byte("value"),
			},
			want: func() []byte {
				ks, vs := 4, 5
				buf := make([]byte, HeaderSize+ks+vs)
				withoutCRC := buf[4:]
				binary.BigEndian.PutUint64(withoutCRC[:8], 1234567890)
				binary.BigEndian.PutUint32(withoutCRC[8:12], uint32(ks))
				binary.BigEndian.PutUint32(withoutCRC[12:16], uint32(vs))
				withoutCRC[16] = 1
				copy(withoutCRC[17:], "test")
				copy(withoutCRC[17+ks:], "value")
				crc := crc32.ChecksumIEEE(withoutCRC)
				binary.LittleEndian.PutUint32(buf[:4], crc)
				return buf
			}(),
			wantErr: false,
		},
		{
			name: "empty key and value",
			entry: &DataEntry{
				DataHeader: &DataHeader{
					Tstamp: 0,
					Flag:   0,
				},
				Key:   []byte{},
				Value: []byte{},
			},
			want: func() []byte {
				ks, vs := 0, 0
				buf := make([]byte, HeaderSize+ks+vs)
				withoutCRC := buf[4:]
				binary.BigEndian.PutUint64(withoutCRC[:8], 0)
				binary.BigEndian.PutUint32(withoutCRC[8:12], uint32(ks))
				binary.BigEndian.PutUint32(withoutCRC[12:16], uint32(vs))
				withoutCRC[16] = 0
				crc := crc32.ChecksumIEEE(withoutCRC)
				binary.LittleEndian.PutUint32(buf[:4], crc)
				return buf
			}(),
			wantErr: false,
		},
		{
			name: "large key and value",
			entry: &DataEntry{
				DataHeader: &DataHeader{
					Tstamp: 987654321,
					Flag:   2,
				},
				Key:   make([]byte, 1000),
				Value: make([]byte, 2000),
			},
			want: func() []byte {
				ks, vs := 1000, 2000
				buf := make([]byte, HeaderSize+ks+vs)
				withoutCRC := buf[4:]
				binary.BigEndian.PutUint64(withoutCRC[:8], 987654321)
				binary.BigEndian.PutUint32(withoutCRC[8:12], uint32(ks))
				binary.BigEndian.PutUint32(withoutCRC[12:16], uint32(vs))
				withoutCRC[16] = 2
				copy(withoutCRC[17:], make([]byte, 1000))
				copy(withoutCRC[17+ks:], make([]byte, 2000))
				crc := crc32.ChecksumIEEE(withoutCRC)
				binary.LittleEndian.PutUint32(buf[:4], crc)
				return buf
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.entry.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataEntry.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareBytes(got, tt.want) {
				t.Errorf("DataEntry.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *DataEntry
		wantErr error
	}{
		{
			name: "valid record",
			data: func() []byte {
				entry := &DataEntry{
					DataHeader: &DataHeader{
						Tstamp: 1234567890,
						Ksz:    4,
						Vsz:    5,
						Flag:   1,
					},
					Key:   []byte("test"),
					Value: []byte("value"),
				}
				entry.CRC = calculateCRC(entry)
				buf := make([]byte, HeaderSize+4+5)
				binary.LittleEndian.PutUint32(buf[:4], entry.CRC)
				binary.BigEndian.PutUint64(buf[4:12], entry.Tstamp)
				binary.BigEndian.PutUint32(buf[12:16], entry.Ksz)
				binary.BigEndian.PutUint32(buf[16:20], entry.Vsz)
				buf[20] = entry.Flag
				copy(buf[HeaderSize:], entry.Key)
				copy(buf[HeaderSize+4:], entry.Value)
				return buf
			}(),
			want: &DataEntry{
				DataHeader: &DataHeader{
					Tstamp: 1234567890,
					Ksz:    4,
					Vsz:    5,
					Flag:   1,
				},
				Key:   []byte("test"),
				Value: []byte("value"),
			},
			wantErr: nil,
		},
		{
			name:    "data too short",
			data:    make([]byte, HeaderSize-1),
			want:    nil,
			wantErr: errs.ErrInvalidRecord,
		},
		{
			name: "invalid key size",
			data: func() []byte {
				buf := make([]byte, HeaderSize)
				binary.BigEndian.PutUint32(buf[12:16], 100) // Ksz = 100
				return buf
			}(),
			want:    nil,
			wantErr: errs.ErrInvalidRecord,
		},
		{
			name: "crc validation failed",
			data: func() []byte {
				entry := &DataEntry{
					DataHeader: &DataHeader{
						Tstamp: 1234567890,
						Ksz:    4,
						Vsz:    5,
						Flag:   1,
					},
					Key:   []byte("test"),
					Value: []byte("value"),
				}
				buf := make([]byte, HeaderSize+4+5)
				binary.LittleEndian.PutUint32(buf[:4], 0) // Wrong CRC
				binary.BigEndian.PutUint64(buf[4:12], entry.Tstamp)
				binary.BigEndian.PutUint32(buf[12:16], entry.Ksz)
				binary.BigEndian.PutUint32(buf[16:20], entry.Vsz)
				buf[20] = entry.Flag
				copy(buf[HeaderSize:], entry.Key)
				copy(buf[HeaderSize+4:], entry.Value)
				return buf
			}(),
			want:    nil,
			wantErr: errs.ErrCRCValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.data)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				// Clear CRC for comparison
				got.CRC = 0
				tt.want.CRC = 0
				if !compareDataEntry(got, tt.want) {
					t.Errorf("Decode() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
func compareDataEntry(a, b *DataEntry) bool {
	if a.DataHeader.Tstamp != b.DataHeader.Tstamp ||
		a.DataHeader.Ksz != b.DataHeader.Ksz ||
		a.DataHeader.Vsz != b.DataHeader.Vsz ||
		a.DataHeader.Flag != b.DataHeader.Flag {
		return false
	}
	if !compareBytes(a.Key, b.Key) || !compareBytes(a.Value, b.Value) {
		return false
	}
	return true
}
