package db

import (
    "os"
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "../utils/fileutil"
    "../config"
)

var Db *DBHelper

// 初始化，检查数据库是否存在，如果不存在则创建对应的表
func init() {
    Db = ConnectDB()
}

type DBHelper struct {
    Connection *sql.DB
}

// 连接数据库
func ConnectDB() *DBHelper {
    conn, err := sql.Open(config.DbDriver, config.DbFile)
    if err != nil {
        log.Fatalf("connect db failed : %s \n", err)
    }
    // 检查数据库是否已经创建
    if !fileutil.IsFileExists(config.DbInitFlagFile) {
        // initialize tables
        initTables(conn)
        // 创建标识文件
        f, err := os.Create(config.DbInitFlagFile)
        if err != nil {
            log.Fatalf("create db flag failed : %s \n", err)
        } else {
            f.Close()
        }
    }
    return &DBHelper{Connection:conn}
}

// 检查表是否存在
func IsTableExists(conn *sql.DB, tableName string) (bool, error) {
    // 获取数据表列表
    rows, err := conn.Query(".tables")
    if err != nil {
        return false, err
    }
    defer rows.Close()
    var tableCreated bool = false
    for rows.Next() {
        var tblName string
        err = rows.Scan(&tblName)
        if err != nil {
            break
        }
        // trades表为必须建立的表，如果此表存在，表示表已经创建
        if tblName == tableName {
            tableCreated = true
            break
        }
    }
    return tableCreated, nil
}

// 初始化表结构
func initTables(conn *sql.DB) {
    for table, sql := range config.DbInitSqls {
        _,err := conn.Exec(sql)
        if err != nil {
            log.Fatalf("initialize table %s failed : %s \n", table, err)
        }
    }
}
