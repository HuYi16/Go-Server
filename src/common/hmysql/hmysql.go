package hmysql

import(
    _"common/mysql"
    "database/sql"
    "common/hlog/hlog"
    "common/hlog/hlog"
)
type Sql_info_st struct{
    Port         int
    Ip           string
    User         string
    Key          string
    DataBaseName string
}

func (s * Sql_info_st) do_init(){
    Ip = ""
    Port = 0
    User = ""
    Key = ""
    DataBaseName = ""
}
var 
func init(){
    hlog.Hlog("init sql moudle!!")

}
