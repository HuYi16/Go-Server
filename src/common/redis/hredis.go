package hredis

import{
    "common/log/plog"    
    "common/redigo/redis"
    "common/hdef/hdef"
}

var redis_info_map map[string][]hdef.Redis_info_st

func initmap(){
    redis_info_map = make(map[string][]hdef.Redis_info_st)
    redis_info_map["r"] = make([]hdef.Redis_info_st,0)
    redis_info_map["w"] = make([]hdef.Redis_info_st,0)
}

func init(){
    initmap()
    plog.Plog("start redis moudle!!")
}

func checkifcon(ipport string,bread bool) bool{
    var con []hdef.Redis_info_st
    if bread {
        con = redis_info_map["r"]
    }else{
        con = redis_info_map["w"]
    }
    for k,v := con{
        if v.AddrPort == ipport{
            return true
        }
    }
    return false
}
//connct to redis
func connect2redis(ipport string) interface{}{
    c,err := redis.Dial("tcp",ipport)
    if err != nil{
        plog.Plog(err)
        return nil
    }
    return c
}

//close all redis connct
func Shutdownredis(){
    for k,v := range redis_info_map{
        for kk,vv := range v{
            vv.Con.Close()
        }
    }
    initmap()
}
//set redis addr and connect to redis
func SetRedis(bRead bool,addr []string) bool{
    if len(addr) == 0 {
        plog.Plog("Set Redis addr is empty!!")
        return false
    }
    conarry :=make([]hdef.Redis_info_st,0)
    var arg hdef.Redis_info_st

    for k,v := range addr{
        if checkifcon(v,bRead){     
            c:= connect2redis(v)
            if nil == c{
                return false
            }
            arg.AddrPort = v
            arg.Con = c
           conarry =  append(conarry,arg)
        }
    }
    if bRead{
        redis_info_map["r"] = append(redis_info_map["r"],conarry)
    }else{
        redis_info_map["w"] = append(redis_info_map["w"],conarry)
    }
    return true
}
//commond like type,key,value..... string arry and moudle will choose read redis or write redis
//if there has just one type connect ,func will no choose
//if there has no connect ,func will return nil,false
func DoRedis(command []string) interface{},bool{
    return nil,true
}
