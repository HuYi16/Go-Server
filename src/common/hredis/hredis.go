package hredis
//redis no need connect pool
import{
    "fmt"
    "common/hlog/hlog"    
    "common/redigo/redis"
}

type Redis_info_st struct{
    cCon              redis.Conn
    szAddr            string
    szUser            string
    szPsw             string
    iPort             int
    iNowDataBaseId    int
    bConnect          bool
}

func (st * Redis_info_st) CloseRedis(){
    if nil != st.cCon
        st.cCon.Close()
        st.bConnct = false
        st.cCon = nil
}

func (st * Redis_info_st) Zerost(){
    st.CloseRedis()
    st.cCon = nil
    szAddrPort = ""
    st.iNowDataBaseId = 0
    st.iPort = 0
    st.szUser = ""
    st.szPsw = ""
    bConnect = false
}

func (st * redis_info_st) SetRedisInfo(szAddr string,iPort int,szUser string,szPsw string){
    st.szAddr = szAddr
    st.szUser = szUser
    st.szPsw = szPsw
    st.iPort = iPort
}


func (st * Redis_info_st) Conn2Redis() bool{
    if bConnect{
        return true
    }
    if 0 == st.iPort || "" == st.szAddr{
        return false
    }
    s := fmt.Sprintf("%s:%d",st.szAddr,st.iPort)
    c,err := redis.Dial("tcp",s)
    if err != nil{
        plog.Plog(err)
        return false
    }
    if st.szUser != "" &&  st.szPsw != ""{
    err = c.Send(st.szUser,st.szPsw)
        if nil != c{
            st.C;pseRedis(c)
            return false
       }
    }
    st.cCon = c
    st.bConnct = true;
    return true
}
/*
var redis_info_map map[string][]hdef.Redis_info_st

func initmap(){
    redis_info_map = make(map[string][]hdef.Redis_info_st)
    redis_info_map["r"] = make([]hdef.Redis_info_st,0)
    redis_info_map["w"] = make([]hdef.Redis_info_st,0)
}

func init(){
    initmap()
    hlog.Hlog("start redis moudle!!")
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
func connect2redis(ipport string) redis.Conn{
    c,err := redis.Dial("tcp",ipport)
    if err != nil{
        plog.Plog(err)
        return nil
    }
    _,errdo :=  c.Do("SET","connctcheck","1")
    if errdo == nil{
        return nil
    }
    return c
}

func checkcon(c redis.Conn) bool{
    _,err:= c.Do("GET","connctcheck")
    if err == nil{
        return true
    }
    return false
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
        hlog.Hlog("Set Redis addr is empty!!")
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
//get a useable redis conn and you can ues connction to change dbid
func GetCanUseCon(bRead bool) redis.Conn{
    if len(redis_info_map) != 2{
        return nil
    }
    c:=make(hdef.Redis_info_st)
    c.Con = nil
    if bRead{
        for k,v := range redis_info_map["r"]{
            c =  v
            break
        }
    }else{
        for k,v := range redis_info_map["w"]{
            c = v
            break
        }
    }
    for k,v := redis_info_map{
        for kk,vv := range v{
            c = v
            break
        }
    }
    if c.Con == nil{
        return nil
    }
    if !checkcon(c.Con) && len(c.AddrPort) != 0{
        cc,err := redis.Dial("TCP",c.AddrPort)
        if err != nil{
            c.Con = cc
            return cc
        }
    }
    return nil
}
*/
