package hlog
import(
    "fmt"
    "time"
)

func Hlog(arg string){
    fmt.Println(time.Now.Format("2016-01-02 15:04:05"),"--",arg)
}
