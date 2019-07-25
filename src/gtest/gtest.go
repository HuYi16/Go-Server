package main

import(
    //R "common/hredis"
    L "common/hlog"
   S "common/hsocket"
   "fmt"
   "time"
)

func cbRead(iId int,buf []byte,iSize int)bool{
    L.W(fmt.Sprintf("%d read,size[%d],content[%s]!!",iId,iSize,string(buf)),L.Level_Trace)
    return true
}

func cbDis(iId int){
    L.W(fmt.Sprintf("%d close!!",iId),L.Level_Trace)
    return 
}

func main(){
    buf := []byte("                                                  ")
    buf2 := []byte("this is a test!!!")
    copy(buf[2:],buf2[:10])
    L.W(string(buf),L.Level_Trace)
    L.W("",L.Level_Trace)
    
    iId,ok,err :=  S.DialC(3564,"127.0.0.1",cbRead,cbDis)
    if ok{
        L.W(fmt.Sprintf("%d",iId),L.Level_Trace)
    }
    L.W(err,L.Level_Trace)
    time.Sleep(time.Second * 60)
    
}
