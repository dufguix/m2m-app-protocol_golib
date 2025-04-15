package applayerprotocol

import (
	"encoding/binary"
	"math"
)

var Converter converter

type converter struct{}

func (converter) BoolToBytes(val bool, buff []uint8) []uint8 {
	if len(buff) < 1 {
		return buff
	}
	if val {
		buff[0] = 1
	} else {
		buff[0] = 0
	}
	return buff[:1]
}

func (converter) BoolFromBytes(bytes []uint8) bool {
	if len(bytes) < 1 {
		return false
	}
	if bytes[0] == 0 {
		return false
	} else {
		return true
	}
}

func (converter) Uint8ToBytes(val uint8, buff []uint8) []uint8 {
	if len(buff) < 1 {
		return buff
	}
	buff[0] = val
	return buff[:1]
}

func (converter) Uint8FromBytes(bytes []uint8) uint8 {
	if len(bytes) < 1 {
		return 0
	}
	return bytes[0]
}

func (converter) Uint16ToBytes(val uint16, buff []uint8) []uint8 {
	if len(buff) < 2 {
		return buff
	}
	buff[0] = uint8(val)
	buff[1] = uint8(val >> 8)
	return buff[:2]
}

func (converter) Uint16FromBytes(bytes []uint8) uint16 {
	if len(bytes) < 2 {
		return 0
	}
	return binary.LittleEndian.Uint16(bytes)
}

func (converter) Uint32ToBytes(val uint32, buff []uint8) []uint8 {
	if len(buff) < 4 {
		return buff
	}
	buff[0] = uint8(val)
	buff[1] = uint8(val >> 8)
	buff[2] = uint8(val >> 16)
	buff[3] = uint8(val >> 24)
	return buff[:4]
}

func (converter) Uint32FromBytes(bytes []uint8) uint32 {
	if len(bytes) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(bytes)
}

func (converter) Float32ToBytes(val float32, buff []uint8) []uint8 {
	if len(buff) < 4 {
		return buff
	}
	bits := math.Float32bits(val)
	buff[0] = uint8(bits)
	buff[1] = uint8(bits >> 8)
	buff[2] = uint8(bits >> 16)
	buff[3] = uint8(bits >> 24)
	return buff[:4]
}

func (converter) Float32FromBytes(bytes []uint8) float32 {
	if len(bytes) < 4 {
		return 0.0
	}
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func (converter) Float64ToBytes(val float64, buff []uint8) []uint8 {
	if len(buff) < 8 {
		return buff
	}
	bits := math.Float64bits(val)
	buff[0] = byte(bits)
	buff[1] = byte(bits >> 8)
	buff[2] = byte(bits >> 16)
	buff[3] = byte(bits >> 24)
	buff[4] = byte(bits >> 32)
	buff[5] = byte(bits >> 40)
	buff[6] = byte(bits >> 48)
	buff[7] = byte(bits >> 56)
	return buff[:8]
}

func (converter) Float64FromBytes(bytes []uint8) float64 {
	if len(bytes) < 8 {
		return 0.0
	}
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}
