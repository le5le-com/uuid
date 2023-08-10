package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

type UUID [16]byte

var InvalidFormat = "Invalid uuid format"

// NewV7 returns a UUID Version 7
// 生成一个uuidv7
//
//	 0                   1                   2                   3
//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                           unix_ts_ms                          |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|          unix_ts_ms           |  ver  |       rand_a          |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|var|                        rand_b                             |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                            rand_b                             |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func NewV7() (UUID, error) {
	var uuid UUID

	ms := time.Now().UTC().UnixMilli()

	// 1. 设置前面6个字节的时间
	// uint64为8个字节，uuidv7只需要6个字节。舍弃高位2个字节
	// 等效于：
	// var b8 [8]byte
	// binary.BigEndian.PutUint64(b8[:], uint64(ms))
	// copy(uuid[0:6], b8[2:8])
	//
	// 下面写法避免申请新的内存空间，参考 https://github.com/cmackenzie1/go-uuid/blob/main/uuid.go
	binary.BigEndian.PutUint64(uuid[:], uint64(ms&((1<<48)-1)<<16))

	// 2. 设置剩余随机数
	_, err := io.ReadFull(rand.Reader, uuid[6:])
	if err == nil {

		// 3. 设置规定标志位
		uuid[6] = (uuid[6] & 0x0f) | 0x70 // Version 7 [0111]
		uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant [10]
	}

	return uuid, err
}

// Parse a UUID string and returns the UUID
func Parse(s string) (UUID, error) {
	var uuid UUID
	var err error

	data := []byte(s)

	switch len(s) {
	case 32:
		// No dash
		_, err = hex.Decode(uuid[0:16], data[0:32])
		if err != nil {
			return UUID{}, err
		}
	case 36:
		if data[8] != '-' || data[13] != '-' || data[18] != '-' || data[23] != '-' {
			return UUID{}, errors.New(InvalidFormat)
		}

		_, err = hex.Decode(uuid[0:4], data[0:8])
		if err != nil {
			return UUID{}, err
		}
		_, err = hex.Decode(uuid[4:6], data[9:13])
		if err != nil {
			return UUID{}, err
		}
		_, err = hex.Decode(uuid[6:8], data[14:18])
		if err != nil {
			return UUID{}, err
		}
		_, err = hex.Decode(uuid[8:10], data[19:23])
		if err != nil {
			return UUID{}, err
		}
		_, err = hex.Decode(uuid[10:16], data[24:36])
		if err != nil {
			return UUID{}, err
		}
	default:
		return UUID{}, errors.New(InvalidFormat)
	}

	return uuid, nil
}

// UUIDV7FromObjectID returns a UUIDv7 from mongodb objectId string
func UUIDV7FromObjectID(s string) (UUID, error) {
	var uuid UUID

	objectId, err := hex.DecodeString(s)
	if err != nil {
		return uuid, err
	}

	seconds := int64(binary.BigEndian.Uint32(objectId[0:4]))
	ms := time.Unix(seconds, 0).UTC().UnixMilli()
	binary.BigEndian.PutUint64(uuid[:], uint64(ms&((1<<48)-1)<<16))
	uuid[6] = (uuid[6] & 0x0f) | 0x70 // Version 7 [0111]
	uuid[7] = objectId[4]
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant [10]
	copy(uuid[9:16], objectId[5:])

	return uuid, err
}

// String returns a hexadecimal string of a UUIDv7.
func (uuid UUID) String() string {
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return string(buf[:])
}

// ShortString returns a no dash hexadecimal string of a UUIDv7.
func (uuid UUID) ShortString() string {
	buf := make([]byte, 32)
	hex.Encode(buf[:], uuid[:])
	return string(buf[:])
}

// TimeFromV7 returns a time using time.Now hight bytes (now is 0 0).
func (uuid UUID) TimeFromV7() time.Time {
	var b8 [8]byte
	copy(b8[2:8], uuid[0:6])

	unixMilli := int64(binary.BigEndian.Uint64(b8[0:8]))
	unixSec := unixMilli / 1000
	nanoSec := (unixMilli % 1000) * 1e6
	return time.Unix(unixSec, nanoSec)
}

// ObjectID returns a mongodb objectId
func (uuid UUID) ObjectID() [12]byte {
	var b [12]byte
	binary.BigEndian.PutUint32(b[0:4], uint32(uuid.TimeFromV7().Unix()))
	b[4] = uuid[7]
	copy(b[5:12], uuid[9:16])

	return b
}

// ObjectIDHex returns a mongodb objectId string
func (uuid UUID) ObjectIDHex() string {
	bytes := uuid.ObjectID()
	return hex.EncodeToString(bytes[:])
}
