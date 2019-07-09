package hsocket

import(
   L "/common/hlog"
   "net"
   "sync"
   "time"
   "math/rand"
   "unsafe"
   "fmt"
)

//the function read callback
type cbRead func(iId int,buf []byte,iSize int) bool 

//the func disconnect callback
type cbDiscon func(iId int) 


type stConnInfo struct{
    szHost  string    //the host of the server
    iPort   string   //the port of the host
    bServerType bool      //true is server type false is client type
    con   net.Conn     //the connect handle
}

type stBaseInfo struct{    
    mCanUseID map[int]int //the closed and can use socket id
    iNowMaxID int   // use to create a new socketid
    iInitID int     //the init ID
    iMaxConnectNumber int //the max connect number for server type
}

var BaseInfo stBaseInfo  //the base info for package
var mConInfo map[int]stConnInfo //the online info map
var m syncc.Mutex      //the lock for online info

func init(){
   L.W("inti socket package",L.Level_Normal)
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   BaseInfo.iNowMaxID = r.Int31()%100000
   if BaseInfo.iNowMaxID <= 10000 || BaseInfo.iNowMaxID >= 90000{
       BaseInfo.iNowMaxID = 76534
   }
   BaseInfo.iInitID = BaseInfo.iNowMaxID
   BaseInfo.iMaxConnect = 65535
   mConInfo = make(map[int]stConnInfo,3000) //init 3000 size 
}
//set the max connect number
func SetMaxConnct(iMaxNum int)bool{
    if iMaxNum <= 0{
        return false
    }
    BaseInfo.iMaxConnect = iMaxNum
    return true
}
//get a socket if for call  function
func getNewID() int{
    for k,_ := range mCanUseID{
        delete(mCanUseID,k)
        return k
    }
    if ++BaseInfo.iNowMaxID <= BaseInfo.iInitID + BaseInfo.iMaxConnect{
        return BaseInfo.iNowMaxID
    }
    L.W("socket max ID used!!!",L.Level_Error)
    return -1
}
//close socket  update online info and id info
func closeCon(iId int){
    m.Lock()
    defer m.Unlock()
    v,ok := mConInfo[iId]
    if ok{
        v.con.Close()
        delete(mConInfo,iId)
        mCanUseID[iId] = 0
    }
}
//update the online info
func updateConInfo(iId int,iPort int,szHost string,conn net.Conn,bServer bool,bAdd bool) bool{
    if bAdd{
         m.Lock()
         defer m.Unlock()
         v,ok := mConInfo[iId]
         if !ok{
             arg := stConnInfo{ szHost, iPort, bServer,conn}
             mConInfo[iId] = arg
         }else{
             return false
         }
     }else{
         closeCon(iId)
     }
    return true
}
//socket start read 
func conRead(cbR cbRead,cbD cbDiscon){
    if nil == cbR || nil == cbD{
        return 
    }
}
// function tcp start clent mode
func DailC(iPort int,szHost string,cbR cbRead,cbD cbDiscon)(int,bool,string){
    if iPort <= 0 || iPort >= 65535{
        return -1,false,"port is invalued!"
    }
    if szHost==""{
        return -1,false,"host is invalued!"
    }
    if nil == cbR || nil == cbD{
        return -1,false,"callback function is invalued!!"
    }
    conn,err := net.Dail("tcp",fmt.Sprintf("%s:%d",szHost,iPort))
    if nil == err{
        L.W(fmt.Sprintf("connect to %s:%d fail,err:",szHost,iPort,err),L.Level_Error)
        return -1,false,err
    }
    iId := getNewID()
    if iId == -1{
        conn.Close()
        return -1,false,"socket id use up!!"
    }
    if !updateConInfo(iId,iPort,szHost,conn,false,true){
        conn.Close()
        return -1,false,"socket id exsist!!"
    }
    go conRead(cbR,cbD)
    return iId,true,""
}

