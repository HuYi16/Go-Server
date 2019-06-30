package main
import (
    R "common/hredis"
    L "common/hlog"
    "fmt"
)

var R_handle R.Redis_info_st
func init(){
    L.W("LB server init!!",L.Level_Warning)
}

func main(){
    L.W("this is test",L.Level_Trace)
    R_handle.Zerost()
    R_handle.SetRedisInfo("127.0.0.1",6379,"")
    if R_handle.Connect(){
        res,err :=R_handle.Con.Do("set","test","1")
        L.W(fmt.Sprintf("%s,%s",res,err),L.Level_Normal)
        R_handle.CloseRedis()
    }
    L.W("LB server quit!!!",L.Level_Error)
}
