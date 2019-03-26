package beson

import (
    "encoding/binary"
    "fmt"
    "math"
    "strings"

    "beson/types"
)

func aaa() {
    fmt.Println(".")
}

func Deserialize(buffer []byte, anchor uint32)(uint32, types.RootType) {
    return deserializeContent(buffer, anchor)
}

func deserializeContent(buffer []byte, start uint32)(uint32, types.RootType) {
    var anchor uint32
    var t string
    var value types.RootType

    anchor, t = deserializeType(buffer, start)
    anchor, value = deserializeData(t, buffer, anchor)

    return anchor, value
}

func deserializeType(buffer []byte, start uint32)(uint32, string) {
    var length uint32 = 2
    end := start + length
    typeData := buffer[start:end]
    
    t := getTypeHeaderKey(typeData)
    return end, t
}

func deserializeData(t string , buffer []byte, start uint32)(uint32, types.RootType) {
    var anchor uint32
    var value types.RootType

    switch t {
    case DATA_TYPE["NULL"]:
        anchor, value = deserializeNull(start)
    case DATA_TYPE["TRUE"], DATA_TYPE["FALSE"]:
        anchor, value = deserializeBoolean(t, start)
    case DATA_TYPE["INT8"]:
        anchor, value = deserializeInt8(buffer, start)
    case DATA_TYPE["INT16"]:
        anchor, value = deserializeInt16(buffer, start)
    case DATA_TYPE["INT32"]:
        anchor, value = deserializeInt32(buffer, start)
    case DATA_TYPE["INT64"]:
        anchor, value = deserializeInt64(buffer, start)
    case DATA_TYPE["INT128"]:
        anchor, value = deserializeInt128(buffer, start)
    case DATA_TYPE["UINT8"]:
        anchor, value = deserializeUInt8(buffer, start)
    case DATA_TYPE["UINT16"]:
        anchor, value = deserializeUInt16(buffer, start)
    case DATA_TYPE["UINT32"]:
        anchor, value = deserializeUInt32(buffer, start)
    case DATA_TYPE["UINT64"]:
        anchor, value = deserializeUInt64(buffer, start)
    case DATA_TYPE["UINT128"]:
        anchor, value = deserializeUInt128(buffer, start)
    case DATA_TYPE["FLOAT32"]:
        anchor, value = deserializeFloat32(buffer, start)
    case DATA_TYPE["FLOAT64"]:
        anchor, value = deserializeFloat64(buffer, start)
    case DATA_TYPE["STRING"]:
        anchor, value = deserializeString(buffer, start)
    }
    return anchor, value
}

func getTypeHeaderKey(typeData []uint8) string {
    var t string
    for key, value := range TYPE_HEADER {
        if value[0] == typeData[0] && value[1] == typeData[1] {
            t = strings.ToLower(key)
            break
        }
    }
    return t
}

func deserializeNull(start uint32)(uint32, types.RootType) {
    return start, nil
}

func deserializeBoolean(t string, start uint32)(uint32, types.RootType) {
    var value *types.Bool
    if t == DATA_TYPE["TRUE"] {
        value = types.NewBool(true).(*types.Bool)
    } else {
        value = types.NewBool(false).(*types.Bool)
    }
    return start, value
}

func deserializeInt8(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 1
    num := int8(buffer[start])
    value := types.NewInt8(num).(*types.Int8)

    return end, value
}

func deserializeInt16(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 2
    num := binary.LittleEndian.Uint16(buffer[start:end])
    value := types.NewInt16(int16(num)).(*types.Int16)

    return end, value
}

func deserializeInt32(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 4
    num := binary.LittleEndian.Uint32(buffer[start:end])
    value := types.NewInt32(int32(num)).(*types.Int32)

    return end, value
}

func deserializeInt64(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 8
    num := binary.LittleEndian.Uint64(buffer[start:end])
    value := types.NewInt64(int64(num)).(*types.Int64)

    return end, value
}

func deserializeInt128(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 16
    numLow := binary.LittleEndian.Uint64(buffer[start:start + 8])
    numHigh := binary.LittleEndian.Uint64(buffer[start + 8:end])
    
    
    value := types.NewInt128("0", 2).(*types.Int128)
    value.SetLow(numLow)
    value.SetHigh(numHigh)

    return end, value
}

func deserializeUInt8(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 1
    num := buffer[start]
    value := types.NewUInt8(num).(*types.UInt8)

    return end, value
}

func deserializeUInt16(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 2
    num := binary.LittleEndian.Uint16(buffer[start:end])
    value := types.NewUInt16(num).(*types.UInt16)

    return end, value
}

func deserializeUInt32(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 4
    num := binary.LittleEndian.Uint32(buffer[start:end])
    value := types.NewUInt32(num).(*types.UInt32)

    return end, value
}

func deserializeUInt64(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 8
    num := binary.LittleEndian.Uint64(buffer[start:end])
    value := types.NewUInt64(num).(*types.UInt64)

    return end, value
}

func deserializeUInt128(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 16
    numLow := binary.LittleEndian.Uint64(buffer[start:start + 8])
    numHigh := binary.LittleEndian.Uint64(buffer[start + 8:end])
    
    
    value := types.NewUInt128("0", 2).(*types.UInt128)
    value.SetLow(numLow)
    value.SetHigh(numHigh)

    return end, value
}

func deserializeFloat32(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 4
    numUint32 := binary.LittleEndian.Uint32(buffer[start:end])
    num := math.Float32frombits(numUint32)
    value := types.NewFloat32(num).(*types.Float32)

    return end, value
}

func deserializeFloat64(buffer []byte, start uint32)(uint32, types.RootType) {
    end := start + 8
    numUint64 := binary.LittleEndian.Uint64(buffer[start:end])
    num := math.Float64frombits(numUint64)
    value := types.NewFloat64(num).(*types.Float64)

    return end, value
}

func deserializeString(buffer []byte, start uint32)(uint32, types.RootType) {
    length := binary.LittleEndian.Uint32(buffer[start:start + 4])
    end := start + 4 + length
    str := string(buffer[start + 4:end])
    value := types.NewString(str).(*types.String)

    return end, value
}
