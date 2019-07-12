package main
import (
    R "common/hredis"
    L "common/hlog"
    S "common/hsocket"
    "fmt"
)

var R_handle R.Redis_info_st
func init(){
    L.W("LB server init!!",L.Level_Warning)
}

func cbRead(iId int,buf []byte,iSize int)bool{
    return true
}
func cbDis(iId int){
    return 
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
    ok,err := S.DailS(3564,cbRead,cbDis)
    if ok{
        L.W("start server suc!!",L.Level_Trace)
    }else{
        L.W(err,L.Level_Error)
    }

    L.W("LB server quit!!!",L.Level_Error)
}
