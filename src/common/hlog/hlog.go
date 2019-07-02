package hlog

import(
    "fmt"
    "path/filepath"
    "os"
    "time"
    "strings"
    "io"
    "sync"
    "bufio"
)

const (
    Level_Trace = iota
    Level_Warning
    Level_Normal
    Level_Error
)
// the lock of write log
var m sync.Mutex
// the defult log file name
var default_file_name string

var ibufferwritblank int

var inowwriteblank int

var szLogBuf string

//make a default_file_name
func init(){
    ibufferwritblank = 0 
    inowwriteblank = 0
    path,err :=  filepath.Abs(filepath.Dir(os.Args[0]))
    if nil != err{
        panic(err)
    }
    val := strings.Split(os.Args[0],"/")
    if len(val) == 0{
        fmt.Print("file name not find")
        panic(1)
    }
    default_file_name = fmt.Sprintf("%s/%s%s.txt",path,val[len(val)-1],time.Now().Format("2006-01-02 15:04:05"))
    /*
    path += "/"
    path += val[len(val)-1]
    path += Time.Now().Format("2006-01-02 15:04:05")
    */
    fmt.Print("path is ",default_file_name,"\n")
}
func SetDelay(iBlank int){
    if iBlank < 0 {
        return
    }
    ibufferwritblank = iBlank
}
//set the log file path and name 
func SetLogName(path string)bool{
    if path != ""{
        default_file_name = path
        return true
    }
    return false 
}
// check the arg ok
func check(file string,arg string) bool{
    if file == "" || arg == ""{
        fmt.Print("log warning file or content is null")
        return false
    }
    return true
}

//write log by direct io
func write(file string,arg string){
    m.Lock()
    defer m.Unlock()
    arg = fmt.Sprintf("%s %s \n",time.Now().Format("2006-01-02 15:04:05"),arg)
    f,err := os.OpenFile(file,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)    
        if err != nil{
            panic(err)
        }
    _,err = io.WriteString(f,arg)
    if nil != err{
        panic(err)
    }
    f.Close()
}
func writebuff(file,arg string){
   m.Lock()
   defer m.Unlock()
   arg = fmt.Sprintf("%s %s \n",time.Now().Format("2006-01-02 15:04:05"),arg)
   szLogBuf += arg
   inowwriteblank++
   if inowwriteblank >= ibufferwritblank{
      inowwriteblank = 0
    f,err := os.OpenFile(file,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)    
    if err != nil{
         panic(err)
     }
     w := bufio.NewWriter(f)
     _,errw := w.WriteString(szLogBuf)
     if nil != errw{
          f.Close()
          panic(errw)
    }
     w.Flush()
    f.Close()
   }
   return 
}
func wLog(file,arg string,bBuf bool,iLv int){
    if file == "" || arg == "" {
        return 
    }
    switch iLv{
    case Level_Warning:
        arg = "warning  " + arg
   case Level_Normal:
       arg = "normal  " + arg
   case Level_Error:
       arg = "error  " + arg
    default:
        fmt.Print(time.Now().Format("2006-01-02 15:04:05")," ",arg,"\n")
        return 
    }
    if bBuf{
        writebuff(file,arg)
    }else{
        write(file,arg) 
    }
    return 
}
//writ the log by direct io and self define the log file
func WF(file string,arg string,iLv int){
    if !check(file,arg){
        return 
    }
    wLog(file,arg,false,iLv)
}
//write log by direct io default file 
func W(arg string,iLv int){
    if !check(default_file_name,arg){
        return 
    }
    wLog(default_file_name,arg,false,iLv)
}

//write log by buffer and defult file
func WB(arg string,iLv int){
    if !check(default_file_name,arg){
        return 
    }
    wLog(default_file_name,arg,true,iLv)
}
//write log by buffer and self define file
func WBF(file,arg string,iLv int){
    if !check(file,arg){
        return 
    }
    wLog(file,arg,true,iLv)
}

