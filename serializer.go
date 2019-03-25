package beson

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math"
    "strings"

    "beson/types"
)

func Serialize(data interface{}) []byte {
    fmt.Println(".")
    return SerializeContent(data)
}

func SerializeContent(data interface{}) []byte {
    t := GetType(data)
    typeBuffer := SerializeType(t)
    dataBuffers := SerializeData(t, data)

    bytesBuffer := bytes.NewBuffer(make([]byte, 0))
    bytesBuffer.Write(typeBuffer)
    bytesBuffer.Write(dataBuffers)

    length := len(typeBuffer) + len(dataBuffers)
    serialContent := make([]byte, length)
    bytesBuffer.Read(serialContent)

    return serialContent
}

func GetType(data interface{}) string {
    var t string

    if data == nil {
        t = DATA_TYPE["NULL"]
        return t
    }

    switch data.(type) {
    case *types.Bool:
        if data.(bool) {
            t = DATA_TYPE["TRUE"]
        } else {
            t = DATA_TYPE["FALSE"]
        }
    case *types.Float32:
        t = DATA_TYPE["FLOAT32"]
    case *types.Float64:
        t = DATA_TYPE["FLOAT64"]
    case *types.Int8:
        t = DATA_TYPE["INT8"]
    case *types.Int16:
        t = DATA_TYPE["INT16"]
    case *types.Int32:
        t = DATA_TYPE["INT32"]
    case *types.Int64:
        t = DATA_TYPE["INT64"]
    case *types.Int128:
        t = DATA_TYPE["INT128"]
    case *types.UInt8:
        t = DATA_TYPE["UINT8"]
    case *types.UInt16:
        t = DATA_TYPE["UINT16"]
    case *types.UInt32:
        t = DATA_TYPE["UINT32"]
    case *types.UInt64:
        t = DATA_TYPE["UINT64"]
    case *types.UInt128:
        t = DATA_TYPE["UINT128"]
    case *types.Binary:
        t = DATA_TYPE["BINARY"]
    case *types.String:
        t = DATA_TYPE["STRING"]
    case *types.Slice:
        t = DATA_TYPE["ARRAY"]
    case *types.Map:
        t = DATA_TYPE["MAP"]
    default:
        t = ""
    }

    return t
}

func SerializeType(t string) []byte {
    typeHeader := make([]byte, 0)
    if t != "" {
        t = strings.ToUpper(t)
        typeHeader = TYPE_HEADER[t]
    }
    return typeHeader
}

func SerializeData(t string, data interface{}) []byte {
    var buffers []byte

    switch t {
    case DATA_TYPE["NULL"]:
        buffers = serializeNull()
    case DATA_TYPE["TRUE"], DATA_TYPE["FALSE"]:
        buffers = serializeBoolean()
    case DATA_TYPE["UINT8"]:
        buffers = make([]byte, 1)
        buffers[0] = data.(*types.UInt8).Get()
    case DATA_TYPE["UINT16"]:
        buffers = make([]byte, 2)
        binary.LittleEndian.PutUint16(buffers, data.(*types.UInt16).Get())
    case DATA_TYPE["UINT32"]:
        buffers = make([]byte, 4)
        binary.LittleEndian.PutUint32(buffers, data.(*types.UInt32).Get())
    case DATA_TYPE["UINT64"]:
        buffers = make([]byte, 8)
        binary.LittleEndian.PutUint64(buffers, data.(*types.UInt64).Get())
    case DATA_TYPE["UINT128"]:
        buffers = serializeUInt128(data.(*types.UInt128))
    case DATA_TYPE["INT8"]:
        buffers = make([]byte, 1)
        buffers[0] = uint8(data.(*types.Int8).Get())
    case DATA_TYPE["INT16"]:
        buffers = make([]byte, 2)
        binary.LittleEndian.PutUint16(buffers, uint16(data.(*types.Int16).Get()))
    case DATA_TYPE["INT32"]:
        buffers = make([]byte, 4)
        binary.LittleEndian.PutUint32(buffers, uint32(data.(*types.Int32).Get()))
    case DATA_TYPE["INT64"]:
        buffers = make([]byte, 8)
        binary.LittleEndian.PutUint64(buffers, uint64(data.(*types.Int64).Get()))
    case DATA_TYPE["INT128"]:
        buffers = serializeInt128(data.(*types.Int128))
    case DATA_TYPE["FLOAT32"]:
        bits := math.Float32bits(data.(*types.Float32).Get())
        buffers = make([]byte, 4)
        binary.LittleEndian.PutUint32(buffers, bits)
    case DATA_TYPE["FLOAT64"]:
        bits := math.Float64bits(data.(*types.Float64).Get())
        buffers = make([]byte, 8)
        binary.LittleEndian.PutUint64(buffers, bits)
    case DATA_TYPE["STRING"]:
        buffers = []byte(data.(*types.String).Get())
    case DATA_TYPE["ARRAY"]:
        // TODO
        
    }

    return buffers
}

func serializeNull() []byte {
    buf := make([]byte, 0)
    return buf
}

func serializeBoolean() []byte {
    buf := make([]byte, 0)
    return buf
}

func serializeUInt128(value *types.UInt128) []byte {
    buf := value.ToBytes()
    return buf
}

func serializeInt128(value *types.Int128) []byte {
    buf := value.ToBytes()
    return buf
}
