package config

// 定义数据库相关配置信息
// 数据库基本信息
var DbDriver = "sqlite3"
var DbName = "pwdkeeper.db"
var DbFile = ""
// 数据库初始信息文件(此文件存在表示数据库已经创建)
var DbInitFlag = "db_init_flag"
var DbInitFlagFile = ""
// 数据表初始化SQL
var DbInitSqls map[string]string = map[string]string{
	"pwd_items" : `create table if not exists pwd_items (
        _id integer primary key autoincrement,
        item varchar(255),
        password varchar(3000),
		description varchar(3000),
		create_time varchar(32),
		update_time varchar(32)
    )`,
}
// 创建索引的sql
var DbIndexSqls []string = []string{
	"CREATE UNIQUE INDEX idx_item on pwd_items(item)",
}