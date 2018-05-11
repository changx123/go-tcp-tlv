package route

import (
	"go-tcp-tlv"
	"bytes"
	"encoding/binary"
	"fmt"
)

type Module struct {
	//模块
	Module uint16
	//方法
	Action uint16
	//Tlv 数据结构
	Tlv *go_tcp_tlv.Tlv
}

//路由结构
type Route struct {
	//打包结构指针
	Conn *go_tcp_tlv.TlvConn
	//路由指向
	Module Module
	//路由列表函数指针
	routeFun *RouteFun
}

type RouteFun struct {
	//uint16 对应路由函数
	routesFun map[uint16]map[uint16]func(r *Route) error
	//中间件路由函数列表
	useFun []func(r *Route) error
}

//声明新的路由 传入路由列表 生产路由列表
func NewRoute() (*RouteFun) {
	var f RouteFun
	f.routesFun = make(map[uint16]map[uint16]func(r *Route) error)
	return &f
}

//添加中间件
func (routeFun *RouteFun) Use(f func(r *Route) error) {
	routeFun.useFun = append(routeFun.useFun, f)
}

//添加路由 对应回调函数
func (routeFun *RouteFun) Route(module uint16, action uint16, f func(r *Route) error) {
	_, ok := routeFun.routesFun[module]
	if !ok {
		routeFun.routesFun[module] = make(map[uint16]func(r *Route) error)
	}
	routeFun.routesFun[module][action] = f
}

//初始一个新的连接
func (routeFun *RouteFun) NewConn(conn *go_tcp_tlv.TlvConn) *Route {
	var r Route
	r.Conn = conn
	r.routeFun = routeFun
	return &r
}

//开启路由监听
func (route *Route) Listen() {
	listen := route.Conn.Reader()
LISTEN:
	for v := range listen {
		route.analysisMA(v)
		//路由结构
		var route_p Route
		route_p.Conn = route.Conn
		route_p.Module = route.Module
		for _, vFun := range route.routeFun.useFun {
			err := vFun(&route_p)
			if err != nil {
				if err == USE_FUN_SKIP {
					continue LISTEN
				} else if err == BREAK_OFF_CLIENT {
					route.Conn.CloseConn()
					break LISTEN
				}
			}
		}
		actions, err := route.readerModule(route.Module.Module)
		if err != nil {
			if err == ERR_NOT_ACTIONS {
				fmt.Println(ERR_NOT_ACTIONS)
				continue
			}
		}
		pFun, ok := actions[route.Module.Action]
		if !ok {
			fmt.Println(ERR_NOT_ACTION)
			continue
		}
		go pFun(&route_p)
	}
}

func (route *Route) Write(module uint16, action uint16, b []byte) (uint32, error) {
	b = route.WriteModule(module, action, b)
	return route.Conn.Write(b)
}

func (route *Route) WriteModule(module uint16, action uint16, b []byte) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, route.Conn.Endian, module)
	binary.Write(bytesBuffer, route.Conn.Endian, action)
	binary.Write(bytesBuffer, route.Conn.Endian, b)
	return bytesBuffer.Bytes()
}

func (route *Route) analysisMA(tlv go_tcp_tlv.Tlv) {
	var module Module
	bytesBuffer := bytes.NewBuffer(tlv.Value)
	binary.Read(bytesBuffer, route.Conn.Endian, &module.Module)
	binary.Read(bytesBuffer, route.Conn.Endian, &module.Action)
	tlv.Value = bytesBuffer.Bytes()
	module.Tlv = &tlv
	route.Module = module
}

func (route *Route) readerModule(module uint16) (map[uint16]func(r *Route) error, error) {
	actionS, ok := route.routeFun.routesFun[module]
	if !ok {
		return nil, ERR_NOT_ACTIONS
	}
	return actionS, nil
}
