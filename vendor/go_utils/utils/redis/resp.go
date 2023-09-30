package redis

import (
	"bytes"
	"errors"
	"fmt"
	"go_utils/utils"
	"strconv"
	"strings"
)

const (
	RespTypeSimpleString byte = 43
	RespTypeError        byte = 45
	RespTypeInt          byte = 58
	RespTypeBulkString   byte = 36
	RespTypeArray        byte = 42
	RespTypeCR           byte = 13
	RespTypeLF           byte = 10
)

var RespArraySep []byte = []byte("/////")

func generateCommandBytesFromRawCommand(cmd string) []byte {
	var rawCommandBytes []byte
	commandParts := strings.Split(cmd, " ")
	totalCommandParts := len(commandParts)
	rawCommandBytes = append(rawCommandBytes, []byte(fmt.Sprintf("*%d\r\n", totalCommandParts))...)
	for i := 0; i < totalCommandParts; i++ {
		rawCommandBytes = append(rawCommandBytes, []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(commandParts[i]), commandParts[i]))...)
	}
	return rawCommandBytes
}

// Quick command functions
func (r *RedisClient) SimpleCommandSingleReturn(cmd ...string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		if err := r.DialRedis(); err != nil {
			return "", err
		}
	}
	if ret, err := r.respSendRawCommand(strings.Join(cmd, " ")); err == nil {
		utils.LogPrintDebug4(string(ret))
		return string(ret), nil
	} else {
		return "", err
	}
}

func (r *RedisClient) SimpleCommandMultiReturn(cmd ...string) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		if err := r.DialRedis(); err != nil {
			return []string{}, err
		}
	}
	var retArray []string
	if ret, err := r.respSendRawCommand(strings.Join(cmd, " ")); err == nil {
		utils.LogPrintDebug4(string(ret))
		retByteArray := bytes.Split(ret, RespArraySep)
		for _, v := range retByteArray {
			retArray = append(retArray, string(v))
		}
		return retArray, nil
	} else {
		return []string{}, err
	}
}

// RESP Send-Recv
func (r *RedisClient) respSendRawCommand(cmd string) ([]byte, error) {
	rawCommandBytes := generateCommandBytesFromRawCommand(cmd)
	writtenBytes, err := r.conn.Write(rawCommandBytes)
	if err != nil || writtenBytes != len(rawCommandBytes) {
		r.conn.Close()
		r.conn = nil
		return nil, errors.New("Failed to write command to redis")
	}
	return r.respReadLine()
}

// RESP Read recv type
func (r *RedisClient) respReadLine() ([]byte, error) {
	typeByte, err := r.rdBuf.ReadByte()
	if err != nil {
		return nil, err
	}
	utils.LogPrintDebug4("TypeByte is: ", typeByte)
	switch typeByte {
	case RespTypeSimpleString, RespTypeInt:
		if ret, err := r.respReadSimple(); err == nil {
			return ret, nil
		} else {
			return nil, err
		}
	case RespTypeError:
		if ret, err := r.respReadSimple(); err == nil {
			return ret, errors.New(string(ret))
		} else {
			return nil, err
		}
	case RespTypeBulkString:
		if ret, err := r.respReadBulkString(); err == nil {
			return ret, nil
		} else {
			return nil, err
		}
	case RespTypeArray:
		if ret, err := r.respReadArray(); err == nil {
			return bytes.Join(ret, RespArraySep), nil
		} else {
			return nil, err
		}
	default:
		utils.LogPrintError("Unsupported type:", typeByte)
		return nil, nil
	}
}

// RESP Read single line
func (r *RedisClient) respReadToCRLF() ([]byte, error) {
	var readBuf []byte
	for {
		if rdTemp, err := r.rdBuf.ReadBytes(RespTypeLF); err == nil {
			readBuf = append(readBuf, rdTemp...)
		} else {
			return nil, err
		}
		if readBuf[len(readBuf)-2] == RespTypeCR {
			// utils.LogPrintDebug4(tempBuf)
			break
		}
	}
	return readBuf, nil
}

// RESP Parse types
func (r *RedisClient) respReadSimple() ([]byte, error) {
	if ret, err := r.respReadToCRLF(); err == nil {
		utils.LogPrintDebug4("Simple ret: ", ret)
		return ret[:len(ret)-2], nil
	} else {
		return nil, err
	}
}

func (r *RedisClient) respReadBulkString() ([]byte, error) {
	var bulkStringLength int
	if ret, err := r.respReadToCRLF(); err == nil {
		utils.LogPrintDebug4("Bulk string length: ", ret)
		if length, err := strconv.Atoi(string(ret[:len(ret)-2])); err == nil {
			bulkStringLength = length
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	if bulkStringLength == -1 {
		return []byte{}, nil
	}
	if ret, err := r.respReadToCRLF(); err == nil {
		utils.LogPrintDebug4("Bulk string: ", ret)
		if len(ret) != bulkStringLength+2 {
			return nil, errors.New("Incorrect bulk string length detected")
		}
		return ret[:len(ret)-2], nil
	} else {
		return nil, err
	}
}

func (r *RedisClient) respReadArray() ([][]byte, error) {
	var arrayLength int
	var retArray = [][]byte{}
	if ret, err := r.respReadToCRLF(); err == nil {
		utils.LogPrintDebug4("Bulk string length: ", ret)
		if length, err := strconv.Atoi(string(ret[:len(ret)-2])); err == nil {
			arrayLength = length
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	if arrayLength == 0 {
		return retArray, nil
	}
	for i := 0; i < arrayLength; i++ {
		ret, err := r.respReadLine()
		if err != nil {
			return retArray, err
		}
		retArray = append(retArray, ret)
	}
	return retArray, nil
}
