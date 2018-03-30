package config

// 定义数据库相关配置信息
// 数据库基本信息
var DbDriver string = "sqlite3"
var DbName string = "pwdkeeper.db"
var DbFile string = ""
// 数据库初始信息文件(此文件存在表示数据库已经创建)
var DbInitFlag = "db_init_flag"
var DbInitFlagFile string = ""
// 数据表初始化SQL
var DbInitSqls map[string]string = map[string]string{
	"pwd_items" : `create table if not exists pwd_items (
        _id integer primary key autoincrement,
        item varchar(255),
        pwd_encrypt varchar(3000),
        salt varchar(64),
        enabled varchar(10),
		description varchar(3000),
		create_time varchar(32),
		disable_time varchar(32),
		update_time varchar(32)
    )`,
}