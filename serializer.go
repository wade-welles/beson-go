package beson

import (
    "bytes"
    "encoding/binary"
    "math"
    "strings"

    "beson/types"
)

func Serialize(data interface{}) []byte {
    return serializeContent(data)
}

func serializeContent(data interface{}) []byte {
    t := getType(data)
    typeBuffer := serializeType(t)
    dataBuffers := serializeData(t, data)

    bytesBuffer := bytes.NewBuffer(make([]byte, 0))
    bytesBuffer.Write(typeBuffer)
    bytesBuffer.Write(dataBuffers)

    length := len(typeBuffer) + len(dataBuffers)
    serialContent := make([]byte, length)
    bytesBuffer.Read(serialContent)

    return serialContent
}

func getType(data interface{}) string {
    var t string

    if data == nil {
        t = DATA_TYPE["NULL"]
        return t
    }

    switch data.(type) {
    case *types.Bool:
        if data.(*types.Bool).Get() {
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

func serializeType(t string) []byte {
    typeHeader := make([]byte, 0)
    if t != "" {
        t = strings.ToUpper(t)
        typeHeader = TYPE_HEADER[t]
    }
    return typeHeader
}

func serializeData(t string, data interface{}) []byte {
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
        s := data.(*types.String)
        buffers = serializeString(s)
    case DATA_TYPE["ARRAY"]:
        slice := data.(*types.Slice)
        buffers = serializeSlice(slice)
    case DATA_TYPE["MAP"]:
        m := data.(*types.Map)
        buffers = serializeMap(m)
    case DATA_TYPE["BINARY"]:
        b := data.(*types.Binary)
        buffers = serializeBinary(b)
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

func serializeString(value *types.String) []byte {
    str := value.Get()
    length := len(str)
    lengthBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(lengthBytes, uint32(length))
    
    dataBytes := []byte(str)
    buf := concatBytesArray(lengthBytes, dataBytes)
    return buf
}

func serializeShortString(value *types.String) []byte {
    str := value.Get()
    length := len(str)
    lengthBytes := make([]byte, 2)
    binary.LittleEndian.PutUint16(lengthBytes, uint16(length))
    
    dataBytes := []byte(str)
    buf := concatBytesArray(lengthBytes, dataBytes)
    return buf
}

func serializeSlice(value *types.Slice) []byte {
    slice := value.Get()
    subBytesBuffer := bytes.NewBuffer(make([]byte, 0))
    for _, element := range slice {
        subType := getType(element)
        subTypeBytes := serializeType(subType)
        subDataBytes := serializeData(subType, element)
        subBytesBuffer.Write(subTypeBytes)
        subBytesBuffer.Write(subDataBytes)
    }

    length := subBytesBuffer.Len()
    lengthBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

    dataBytes := make([]byte, length)
    subBytesBuffer.Read(dataBytes)

    buf := concatBytesArray(lengthBytes, dataBytes)
    return buf
}

func serializeMap(value *types.Map) []byte {
    subBytesBuffer := bytes.NewBuffer(make([]byte, 0))
    m := value.Get()
    for key, value := range m {
        // serialize key
        k := types.NewString(key).(*types.String)
        keyBytes := serializeShortString(k)

        // serialize value
        subType := getType(value)
        subTypeBytes := serializeType(subType)
        subDataBytes := serializeData(subType, value)

        subBytesBuffer.Write(subTypeBytes)
        subBytesBuffer.Write(keyBytes)
        subBytesBuffer.Write(subDataBytes)
    }

    length := subBytesBuffer.Len()
    lengthBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

    dataBytes := make([]byte, length)
    subBytesBuffer.Read(dataBytes)

    buf := concatBytesArray(lengthBytes, dataBytes)
    return buf
}

func serializeBinary(value *types.Binary) []byte {
    dataBytes := value.ToBytes()
    length := len(dataBytes)
    lengthBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

    buf := concatBytesArray(lengthBytes, dataBytes)
    return buf
}

func concatBytesArray(b1 []byte, b2 ...[]byte) []byte {
    buf := bytes.NewBuffer(make([]byte, 0))
    
    buf.Write(b1)
    for _, element := range b2 {
        buf.Write(element)
    }

    newBytes := make([]byte, buf.Len())
    buf.Read(newBytes)

    return newBytes
}
