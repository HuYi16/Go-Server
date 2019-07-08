package hsocket

import(
   L "/common/hlog"
   "net"
   "sync"
   "time"
   "unsafe"
)

type Read func() bool //read fun callback

type DisCon func() //disconnect fun callback

func init(){
   L.W("inti socket package",L.Level_Normal) 
}

