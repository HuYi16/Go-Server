package main
import (
    DB "common/hdb"
    L "common/hlog"
    S "common/hsocket"
    "fmt"
    "time"
)

var DBRedis  *DB.Cache
func init(){
    L.W("LB server init!!",L.Level_Warning)
}

func cbRead(iId int,buf []byte,iSize int)bool{
    szSend := fmt.Sprintf("sever send data :[%d]",time.Now().UnixNano())
    L.W(fmt.Sprintf("recv:%s,send[%s]",string(buf),szSend),L.Level_Trace)
    S.Write(iId,[]byte(szSend),len([]byte(szSend)))
    return true
}
func cbDis(iId int){
    return 
}

func main(){
    L.W("this is test",L.Level_Trace)
    redisCfg := DB.Redis{
        MaxIdle :                8,
        MaxActive :              64,
        IdleTimeout :            300,
        RedisServer :            "127.0.0.1:6379",
        DialReadTimeout :        3,
        DialWriteTimeout :       3,
        Auth            :        "",
        DbNum           :         0,
    }
    DBRedis = DB.NewRedis(redisCfg)
    if nil == DBRedis{
        L.W("init redis fail",L.Level_Error)
        return 
    }
    timeBlank := 3 * time.Second
    errredis := DBRedis.Put("htest",2,timeBlank)
    if errredis != nil{
        L.W("put err",L.Level_Trace)
    }
    v := DBRedis.Get("htest")
    L.W(fmt.Sprintf("key:htest  val:[%s]",v),L.Level_Trace)
    ok,err := S.DialS(3564,cbRead,cbDis)
    if ok{
        L.W("start server suc!!",L.Level_Trace)
    }else{
        L.W(err,L.Level_Error)
    }
    time.Sleep(50*time.Second)
    L.W("LB server quit!!!",L.Level_Error)
}
