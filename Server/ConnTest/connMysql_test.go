package ConnTest

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"testing"
)

func TestMysqlConn(t *testing.T) {
	db, err := sql.Open("mysql", "root:huhao123@tcp(127.0.0.1:3306)/game?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	var result sql.Result
	result, err = db.Exec("insert into user_table(user_id, name,password) values(?,?,?)", 1234, "huhaos", "huhao123")
	if err != nil {
		fmt.Println(err)
		return
	}

	lastId, _ := result.LastInsertId()
	fmt.Println("新插入记录的ID为", lastId)

	var row *sql.Row
	row = db.QueryRow("select * from  user_table")
	var name, password string
	var user_id int
	err = row.Scan(&user_id, &password, &name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user_id, "\t", name, "\t", password)

	result, err = db.Exec("insert into user_table(user_id,name,password) values(?,?,?)", 4444, "black", "huhao123456")

	var rows *sql.Rows
	rows, err = db.Query("select * from user_table")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var name, password string
		var user_id int
		rows.Scan(&user_id, &name, &password)
		fmt.Println(user_id, "\t", name, "\t", password)
	}
	rows.Close()

	db.Exec("truncate table user")
}

func TestSqLiteConn(t *testing.T) {
	DBConnector, err := sql.Open("sqlite3", "./../DB/sqllite.db")
	if err != nil || DBConnector == nil {
		fmt.Printf("连接失败")
		os.Exit(1)
	}
	DBConnector.SetMaxOpenConns(16)
	DBConnector.SetMaxIdleConns(8)
	p := DBConnector.Ping()
	fmt.Printf("连接成功 %v", p)

	sqlstring := string("replace into user_table(user_id,name,password) values (?,?,?)")
	stmt, err__ := DBConnector.Prepare(sqlstring)
	defer stmt.Close()
	if err__ != nil {
		fmt.Println("create err ", err__)
		return
	}
	result, err := stmt.Exec(
		33334,
		"huhao",
		"huhao123",
	)
	if err != nil {
		fmt.Println("create err ", err)
	}
	effectRow, err := result.RowsAffected()
	if err != nil {
		fmt.Println("create err ", err)
	}
	fmt.Println("create effect rows:", effectRow)

}

func TestRedisConn(t *testing.T) {

}
