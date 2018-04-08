package go_tcp_tlv

import (
	"encoding/binary"
	"bytes"
)

func (conn *TlvConn) Write(b []byte) (uint32, error) {
	return conn.connWrite(b, 0)
}

func (conn *TlvConn) connWrite(b []byte, t uint8) (uint32, error) {
	//tlv 数据结构
	var tlv Tlv
	tlv.Type = uint8(t)
	tlv.Length = uint32(len(b))
	tlv.Value = b
	//tlv 转换
	b = tlv.writeToBt(conn.Endian)
	//发送数据
	l, err := conn.conn.Write(b)
	if err != nil {
		return uint32(l), err
	}
	return uint32(l), nil
}

//tlv转为byte
func (t *Tlv) writeToBt(endian binary.ByteOrder) []byte {
	newBuffer := bytes.NewBuffer([]byte{})
	//写入tag(T)
	binary.Write(newBuffer, endian, t.Type)
	//写入长度(L)
	binary.Write(newBuffer, endian, t.Length)
	//写入数据(V)
	binary.Write(newBuffer, endian, t.Value)
	return newBuffer.Bytes()
}
