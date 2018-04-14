package route

import (
	"errors"
)

//找不到方法列表
var ERR_NOT_ACTIONS = errors.New("ERROR: module => not find actions")

//找不到方法列表
var ERR_NOT_ACTION = errors.New("ERROR: actions => not find action")

//中间件跳出本次路由指令 不继续执行下面方法
var USE_FUN_SKIP = errors.New("CONTROL: use fun This skip")

var BREAK_OFF_CLIENT = errors.New("CONTROL: use fun This BREAK_OFF_CLIENT")

