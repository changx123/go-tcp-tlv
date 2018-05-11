package go_tcp_tlv

import (
	"net"
	"encoding/binary"
)

//tlv数据结构
type Tlv struct {
	//数据类型 [00000000] 1字节 [0默认 自己约定]
	Type uint8
	//数据长度 [00000000][...][...][...] 4字节
	Length uint32
	//详细数据
	Value []byte
}

//tcp连接和配置
type TlvConn struct {
	//socket TCP连接
	conn net.Conn
	//每次读取Buffer长度 默认1024
	ReadBufferSize uint32
	//端序设置 (binary.BigEndian 大端序 ， binary.LittleEndian 小端序)默认大端序
	Endian binary.ByteOrder
}

//初始化连接
func (conn *TlvConn) NewConn(tcp_conn net.Conn) {
	//设置读取长度默认值
	if conn.ReadBufferSize == 0 {
		conn.ReadBufferSize = 1024
	}
	//设置端序默认值
	if conn.Endian == nil {
		conn.Endian = binary.BigEndian
	}
	conn.conn = tcp_conn
}

func (conn *TlvConn) ReaderConn() net.Conn {
	return conn.conn
}

func (conn *TlvConn) CloseConn() {
	ConnCloseChan[conn] = make(chan string, 1)
	ConnCloseChan[conn] <- "close"
	conn.conn.Close()
}
