package main
import (
    "log"
    R "common/hredis"
)

var R_handle R.Redis_info_st
func init(){
    log.Print("LB server init!!")
}

func main(){
    R_handle.Zerost()
    R_handle.SetRedisInfo("127.0.0.1",6379,"")
    if R_handle.Connect(){
        res,err :=R_handle.Con.Do("set","test","1")
        log.Print(res,err)
        R_handle.CloseRedis()
    }
    log.Print("this is LB server!")
}
