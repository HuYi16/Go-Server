package  redissql

import(
    "common/log/plog"
    "common/redigo/redis"
    "common/stdef/stdef"
    "database/sql"
    _ "common/mysql"
)

var redis_addr_map map[string][]stdef.SqlRedisinfodes
var sql_addr_map map[string][]stdef.SqlRedisinfodes

func init(){
    plog.Plog("start redis and sql dirver")
    redis_addr_map = make(map[string][]stdef.SqlRedisinfodes)
    sql_addr_map = make(map[string][]stdef.SqlRedisinfodes)
}

func ConRisOrSql(bRedis bool,addr string) interface{},error{
    var errres error
    bSuc = false 
    if addr != ""{
        if bRedis{
            c,err = redis.Dial("tcp",addr)
            errres = err
        else{
            db,err = sql.Open()
        }
    }
    return nil,err
}
func makestredissqlinfodes(arg []string) []stdef.SqlRedisinfodes{
    valinfo := make([]stdef.SqlRedisinfodes,0)
    if len(arg) == 0{
        return  valinfo
    }
    var stTemp stdef.SqlRedisinfodes
    for k,v := range arg{
        stTemp.AddrPort = v
        valinfo = append(valinfo,stTemp)
    }
    return valinfo
}

//set redis or sql addr and port  
//bRedis set redis addr or set sql addr
//readonlyaddr the ip and port arry or slice which just ueed to read data format like 127.0.0.1:6666
//writeonlyaddr the ip and port arry or slice which only used to write data
func SetReidsSqlAddr(bRedis bool,read []string,write []string) bool{
    if len(readonlyaddr) == 0 ||len(writeonlyaddr) == 0{
        plog.Plog("SetReidsSqlAddr param len is 0!!!")
        return false
    }
    if bRedis {
        redis_addr_map["r"] = makestredissqlinfodes(read)
        redis_addr_map["w"] = makestredissqlinfodes(write)
    }else{
        sql_addr_map["r"] = makestredissqlinfodes(read)
        sql_addr_map["w"] = makestredissqlinfodes(write)
    }
    return true
}
// add some new addr for redis or sql
func RedisSqlAddrAppend(bRedis bool,readaddr []string,writeaddr []string) int{
    if len(readaddr) == 0 && len(writeaddr) == 0{
        return 1 //readaddr and writeaddr both empty
    }
    return 0 //suc
}  
func connectsqlorredis(bRedis bool,addr string)  interface{}{
    if bRedis {
        plog.Plog("start connectsqlorredis")

    }
    return true
}
