package hsocket

import(
   L "common/hlog"
   "net"
   "sync"
   "time"
   "math/rand"
   "unsafe"
   "fmt"
   "io"
   "strings"
)

type FakeSlice struct{
    addr uintptr
    len int
    cap int
}

const(
    iMaxBufSize = 8096
    IDEADTIME = 1000*time.Millisecond   //defual deadtime
    IHEARTTIMEOUT = 10  //heartbeat 10 seconds
)
//the function read callback
type cbRead func(iId int,buf []byte,iSize int) bool 

//the func disconnect callback
type cbDiscon func(iId int) 

type MSGHead struct{
    lens int //msg body lenth
    //ServerId int //to or from ServerId
    bHeartBeat bool 
}

type stConnInfo struct{
    szHost  string    //the host of the server
    iPort   int   //the port of the host
    bServerType bool      //true is server type false is client type
    con   net.Conn     //the connect handle
    cbD   cbDiscon      //disconnect callback function
}

type stBaseInfo struct{    
    mCanUseID map[int]int //the closed and can use socket id
    iNowMaxID int   // use to create a new socketid
    iInitID int     //the init ID
    iMaxConnectNumber int //the max connect number for server type
}

var BaseInfo stBaseInfo  //the base info for package
var mConInfo map[int]stConnInfo //the online info map
var m sync.Mutex      //the lock for online info
var bStop bool

func init(){
   L.W("init socket package",L.Level_Normal)
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   BaseInfo.iNowMaxID = r.Int()%100000
   if BaseInfo.iNowMaxID <= 10000 || BaseInfo.iNowMaxID >= 90000{
       BaseInfo.iNowMaxID = 76534
   }
   BaseInfo.iInitID = BaseInfo.iNowMaxID
   BaseInfo.iMaxConnectNumber = 65535
   mConInfo = make(map[int]stConnInfo,3000) //init 3000 size
   BaseInfo.mCanUseID = make(map[int]int,1000) //init 1000 size
   bStop = false
}

func StopTcp(){
    bStop = true
}

func GetStop()bool{
    return bStop
}
//set the max connect number
func SetMaxConnct(iMaxNum int)bool{
    if iMaxNum <= 0{
        return false
    }
    BaseInfo.iMaxConnectNumber = iMaxNum
    return true
}
//get a socket if for call  function
func getNewID() int{
    for k,_ := range BaseInfo.mCanUseID{
        delete(BaseInfo.mCanUseID,k)
        return k
    }
    if BaseInfo.iNowMaxID  < BaseInfo.iInitID + BaseInfo.iMaxConnectNumber{
        BaseInfo.iNowMaxID++
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
        v.cbD(iId)
        v.con.Close()
        delete(mConInfo,iId)
        BaseInfo.mCanUseID[iId] = 0
    }
}
//update the online info
func updateConInfo(iId int,iPort int,szHost string,conn net.Conn,cbD cbDiscon,bServer,bAdd bool) bool{
    if nil == conn || nil == cbD || iId < 0 || iPort < 0{
        return false
    }
    if bAdd{
         m.Lock()
         defer m.Unlock()
         _,ok := mConInfo[iId]
         if !ok{
             arg := stConnInfo{ szHost, iPort, bServer,conn,cbD}
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
func conRead(cbR cbRead,iId int,conn net.Conn,bServer bool){
    defer  closeCon(iId)
    if nil == cbR || iId <= 0{
        return 
    }
    //set INIFTtime 50ms
    var stHead MSGHead
    iHeadSize := int(unsafe.Sizeof(stHead))
   // iHeadSize := *(*int)(unsafe.Pointer(unsafe.Sizeof(stHead)))
    buf := make([]byte,iMaxBufSize)
    bufHead := make([]byte,iHeadSize)
    var iReadByte int
    var err error
    iNowBuf := 0
    iStartIndex := 0
    llNowTime := time.Now().Unix() + IHEARTTIMEOUT
   // L.W("start read loop....",L.Level_Trace)
    for !GetStop(){
        errread := conn.SetReadDeadline(time.Now().Add(IDEADTIME))
        if nil != errread{
            L.W(fmt.Sprintf("SetReadDeadline err:%s",errread),L.Level_Error)
        }
        
        //check heart beat or send heart beat
        if bServer{
            if llNowTime < (time.Now().Unix()-2){//delay 2 seconds
               // L.W(fmt.Sprintf("[%d] - [%d]",llNowTime,time.Now().Unix() + IHEARTTIMEOUT+2),L.Level_Trace)
                return 
            }
        }else {
               // L.W(fmt.Sprintf("client [%d] - [%d]",llNowTime,time.Now().Unix() + IHEARTTIMEOUT+2),L.Level_Trace)
            if llNowTime -1 <=time.Now().Unix(){//1 seconds beforhand
                if !doHeartBeat(conn){
                    return 
                }else{
                    llNowTime = time.Now().Unix() + IHEARTTIMEOUT
                }
            }
        }
       // L.W("start read....",L.Level_Trace)
        iReadByte = 0
        err = nil
        iReadByte,err = conn.Read(buf[iNowBuf:])
       // L.W(fmt.Sprintf("read number %d,iNowBuf[%d],iHeadSize[%d]",iReadByte,iNowBuf,iHeadSize),L.Level_Trace)
        if err != nil{
            //L.W(fmt.Sprintf("[%d]read err...[%s] ",iId,err),L.Level_Trace)
            if strings.Contains(err.Error(),"timeout"){
                continue
            }
            L.W(fmt.Sprintf("[%d]read err,[%s] ",iId,err),L.Level_Error)
            return 
            if err != io.EOF{
                L.W(fmt.Sprintf("[%d]read err,[%s] ",iId,err),L.Level_Error)
                return 
            }
            continue
        }
        iNowBuf += iReadByte
         for iNowBuf-iStartIndex >= iHeadSize{//judge if head is complete!
             copy(bufHead,buf[iStartIndex:iStartIndex+iHeadSize])
             //stHead = MSGHead(bufHead)
             stHead = *(*(**MSGHead)(unsafe.Pointer(&bufHead)))
            // L.W(fmt.Sprintf("%d,%d,%d,%s,%d",iNowBuf,iReadByte,iStartIndex,string(bufHead),stHead.lens),L.Level_Trace)
             if stHead.bHeartBeat{
                    //do heartbeat back
                    if bServer && !doHeartBeat(conn){
                        return 
                    }else{
                        llNowTime = time.Now().Unix() + IHEARTTIMEOUT  //update next heart dead time
                        iStartIndex +=iHeadSize
                    }
                }else{
                if iHeadSize + stHead.lens <= iNowBuf{//get head body complete
                //copy body data to a new buf
                    iStartIndex += iHeadSize
                    body  := buf[iStartIndex:iStartIndex + stHead.lens]
                    if !cbR(iId,body,stHead.lens){
                        return
                    }
                    iStartIndex += stHead.lens
                }else{
                    break  //continue read 
                }
            }
        }
        copy(buf,buf[iStartIndex:iNowBuf])
        iNowBuf -= iStartIndex
        iStartIndex = 0
    }
}
// function tcp start clent mode 
//need auto reconnect?
func DialC(iPort int,szHost string,cbR cbRead,cbD cbDiscon)(int,bool,string) {
    if iPort <= 0 || iPort >= 65535{
        return -1,false,"port is invalued!"
    }
    if szHost==""{
        return -1,false,"host is invalued!"
    }
    if nil == cbR || nil == cbD{
        return -1,false,"callback function is invalued!!"
    }
    conn,err := net.Dial("tcp",fmt.Sprintf("%s:%d",szHost,iPort))
    if nil != err{
        L.W(fmt.Sprintf("connect to %s:%d fail,err:",szHost,iPort,err),L.Level_Error)
        return -1,false,fmt.Sprintf("%s",err)
    }
    iId := getNewID()
    if iId == -1{
        conn.Close()
        return -1,false,"socket id use up!!"
    }
    if !updateConInfo(iId,iPort,szHost,conn,cbD,false,true){
        conn.Close()
        return -1,false,"socket id exsist!!"
    }
    go conRead(cbR,iId,conn,false)
    return iId,true,"start client suc!!"
}

// handle the clent Accept
func handleAccp(cbR cbRead,cbD cbDiscon,l net.Listener,iPort int){
    if nil == cbR ||nil ==  cbD || iPort <= 0{
        L.W("call back fun is nil!! server accept start fail!!",L.Level_Error)
        return
    }
    var conn net.Conn
    var err error
    L.W("start Accept",L.Level_Trace)
    for !GetStop(){
        conn,err = l.Accept()
        L.W("some one connct!!",L.Level_Trace)
        if nil != err{
            L.W(fmt.Sprintf("accept fail:%s",err),L.Level_Error)
        }
        iId := getNewID()
        if iId == -1{
            conn.Close()
            L.W("get socketid fail",L.Level_Normal)
            return 
        }
         if !updateConInfo(iId,iPort,conn.RemoteAddr().String(),conn,cbD,true,true){
            conn.Close()
            L.W("insert into map socket fail!!",L.Level_Normal)
            return 
        }
        go conRead(cbR,iId,conn,true)
    }
    return 
}
//start net server mode
func DialS(iPort int,cbR cbRead,cbD cbDiscon)(bool,string){
    if iPort <= 0 {
        return false,"server port is invalued!!"
    }
    if nil == cbR || nil == cbD{
        return false ,"server callback func is nil!!"
    }
    l,err := net.Listen("tcp",fmt.Sprintf("%s:%d","127.0.0.1",iPort))
    if nil != err{
        L.W(fmt.Sprintf("start listen:port[%d],err:%s",iPort,err),L.Level_Error)
        return false,"start listen err!!"
    }
    go handleAccp(cbR,cbD,l,iPort)
    return true,""
}

func doHeartBeat(con net.Conn)bool{
    head := MSGHead{lens:100,
                    bHeartBeat:true}
    iHeadSize := int(unsafe.Sizeof(head))
    tempHead :=&FakeSlice{uintptr(unsafe.Pointer(&head)),iHeadSize,iHeadSize}
    bufSend := *(*[]byte)(unsafe.Pointer(tempHead))
    n,err := con.Write(bufSend)
    if err != nil{
        L.W(fmt.Sprintf("send fail,err:%s",err),L.Level_Error)
        return false
    }
    if n != len(bufSend){
        L.W("send data number is wrong!!",L.Level_Error)
        return false
    }
   // L.W(fmt.Sprintf("do heart beat,%s",string(bufSend)),L.Level_Trace)
    return true
 }

func Write(iId int,buf []byte,lens int)bool{
    head :=MSGHead{lens:lens,
                    bHeartBeat:false}
    iHeadSize := int(unsafe.Sizeof(head))
    //iHeadSize := *(*int)(unsafe.Pointer(unsafe.Sizeof(head)))
    if lens > iMaxBufSize - iHeadSize{
        return false
    }
    v,ok := mConInfo[iId]
    if !ok{
        return false
    }
    tempHead :=&FakeSlice{uintptr(unsafe.Pointer(&head)),iHeadSize,iHeadSize}
    bufSend := *(*[]byte)(unsafe.Pointer(tempHead))
    bufSend = append(bufSend,buf...)
   // bufSend := make([]byte,lens + iHeadSize)
   // copy(bufSend,[]byte(head))
    //copy(bufSend[unsafe.Sizeof(head):],buf)
    n,err := v.con.Write(bufSend)
    if err != nil{
        L.W(fmt.Sprintf("send fail,err:%s",err),L.Level_Error)
        closeCon(iId)
        return false
    }
    if n != len(bufSend){
        L.W("send data number is wrong!!",L.Level_Error)
        closeCon(iId)
        return false
    }
    return true
}
