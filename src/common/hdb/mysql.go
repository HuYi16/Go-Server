package db

import (
    L"common/hlog"
	"database/sql"
	"fmt"
	_ "common/mysql"
)

var (
	_mysql *sql.DB
)

type Mysql struct {
	User         string
	Password     string
	Host         string
	Db           string
	MaxOpenConns int
	MaxIdleConns int
}

type CallBack func() error

//start sql connect 
func NewMysql(arg Mysql) {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		arg.User, arg.Password, arg.Host, arg.Db)
	_mysql, err = sql.Open("mysql", dataSourceName)
	if err != nil {
        L.W(fmt.Sprintf("open sql:%s err:%s", dataSourceName, err),L.Level_Error)
		return
	}
	err = _mysql.Ping()
	if err != nil {
        L.W(fmt.Sprintf("mysql:%s ping err:%s", dataSourceName, err),L.Level_Error)
		return
	}
	_mysql.SetMaxIdleConns(arg.MaxIdleConns)
	_mysql.SetMaxOpenConns(arg.MaxOpenConns)
    L.W(fmt.Sprintf("sql run:%s", dataSourceName),L.Level_Normal)
	return
}

//close sql
func CloseMysql() {
	if _mysql == nil {
		return
	}
	_mysql.Close()
	L.W("mysql close",L.Level_Normal)
}


func GetMysql() *sql.DB {
	return _mysql
}

//sql execute fun
func execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return GetMysql().Exec(sqlStr, args...)
}

func Query(queryStr string, args ...interface{}) (map[int]map[string]string, error) {
	query, err := GetMysql().Query(queryStr, args...)
	results := make(map[int]map[string]string)
	if err != nil {
		return results, err
	}
	defer query.Close()
	cols, _ := query.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	i := 0
	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}
	return results, nil
}

//sql update
func Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// sql insert 
func Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

// sql delete
func Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}


type MysqlTransaction struct {
	SQLTX *sql.Tx
}

func (t *MysqlTransaction) execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return t.SQLTX.Exec(sqlStr, args...)
}


func Begin() (*MysqlTransaction, error) {
	var trans = &MysqlTransaction{}
	var err error
	if pingErr := GetMysql().Ping(); pingErr == nil {
		trans.SQLTX, err = GetMysql().Begin()
	}
	return trans, err
}

func (t *MysqlTransaction) Rollback() error {
	return t.SQLTX.Rollback()
}

func (t *MysqlTransaction) Commit() error {
	return t.SQLTX.Commit()
}

func (t *MysqlTransaction) Query(queryStr string, args ...interface{}) (map[int]map[string]string, error) {
	query, err := t.SQLTX.Query(queryStr, args...)
	results := make(map[int]map[string]string)
	if err != nil {
		return results, err
	}
	defer query.Close()
	cols, _ := query.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	i := 0
	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}
	return results, nil
}

func (t *MysqlTransaction) Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

func (t *MysqlTransaction) Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

func (t *MysqlTransaction) Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}
