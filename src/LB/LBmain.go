package main
import (
    log "common/hlog/hlog"
    "common/redis/hredis"
)

var redis_handle Redis_info_st
func init(){
    fmt.Println("LB server init!!")
}

func main(){
    redis_handle.Zerost()
    redis_handle.SetRedisInfo("127.0.0.1",6379,"")
    if redis_handle.Con2Redis(){
        res :=redis_handle.Con.Do("set","test","1")
        log.Hlog(res)
        redis_handle.CloseRedis()
    }
    log.Hlog("this is LB server!")
}
