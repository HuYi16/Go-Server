package  redissql

import(
    "common/log/plog"
    "common/redigo/redis"
)

var redis_addr_map map[string][]string
var sql_addr_map map[string]map[]string
func init(){
    plog.Plog("start redis and sql dirver")
    redis_addr_map = make(map[string][]string)
    sql_addr_map = make(map[string][]string)
}

func initRedis() bool{
    return true
}
func SetReidsAddr(readonlyaddr []string,writeonlyaddr []string) bool{
    if len(readonlyaddr) == 0 ||len(writeonlyaddr) == 0{
        plog.Plog("SetReidsAddr param len is 0!!!")
        return false
    }
    valueread := addarg[:]
    valuewrite := writeonlyaddr[:]
    redis_addr_map["r"] = valueread
    redis_addr_map["w"] = valuewrite
    return true
}

func SetSqlAddr(readaddr []string,writeaddr []string) bool{
    if len(readaddr) == 0 || len(writeaddr) == 0{
        plog.Plog("SetSqlAddr param len is 0!!!")
        return false
    } 
    valread := readaddr[:]
    valwrite := writeaddr[:]
    sql_addr_map["r"] = valread
    sql_addr_map["w"] = valwrite
}
