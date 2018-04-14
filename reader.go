package go_tcp_tlv

import (
	"bytes"
	"encoding/binary"
)

//读取数据
func (conn *TlvConn) Reader() <-chan Tlv {
	//创建输出通道
	out := make(chan Tlv)
	go func() {
		b := make([]byte, conn.ReadBufferSize)
		//接收数据缓冲区
		bio := bytes.NewBuffer([]byte{})
		//接收数据缓冲区长度
		biol := uint32(0)
		var t Tlv
		//EOF标识
		var EOF = false
		for {
			//判断字节是否大于5(T.L 协调头5字节)
			if t.Length == 0 && biol >= 5 {
				//读取tag
				binary.Read(bio, conn.Endian, &t.Type)
				//读取length
				binary.Read(bio, conn.Endian, &t.Length)
				//减去长度5
				biol -= 5
			}
			//判断协议头存在数据(V)长度并且当前缓冲区数据长度">="(L)收到的数据长度
			if t.Length != 0 && biol >= t.Length {
				//开辟需要取得数据长度
				tmpBt := make([]byte, t.Length)
				//获取数据
				binary.Read(bio, conn.Endian, &tmpBt)
				t.Value = tmpBt
				//传输数据到通道
				out <- t
				//减去(V)的长度
				biol -= t.Length
				//初始结构清空上次数据
				t = Tlv{}
				continue
			}
			if t.Length <= 0 && biol <= 0 && EOF {
				break
			}
			//接收数据
			l, err := conn.conn.Read(b)
			if err != nil || l < 1 {
				EOF = true
				continue
			}
			//取到的数据写入缓冲区
			binary.Write(bio, conn.Endian, b[:l])
			//增加缓冲区长度
			biol += uint32(l)
		}
		//关闭通道
		close(out)
	}()
	return out
}
