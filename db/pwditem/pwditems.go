package pwditem

import (
	"database/sql"
	"../../db"
)

// 定义密码项信息
type PwdItem struct {
	Id              int64
	Item            string
	Password        string
	Description     string
	CreateTime      string
	UpdateTime      string
}

// 创建一个空的Trade struct
func NewPwdItem() *PwdItem {
	return &PwdItem{}
}

// 根据item key获取一条密码项
func GetByItem(item string) (*PwdItem, error) {
	querySql := "select _id,item,password,description,create_time,update_time from pwd_items where item = ?"
	stmt, err := db.Db.Connection.Prepare(querySql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	// 开始查询
	row := stmt.QueryRow(item)
	pwdItem := NewPwdItem()
	err = row.Scan(
		&pwdItem.Id,
		&pwdItem.Item,
		&pwdItem.Password,
		&pwdItem.Description,
		&pwdItem.CreateTime,
		&pwdItem.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pwdItem, nil
}

// 将交易信息插入数据库
func (pwdItem *PwdItem) InsertToDb()  (int64, error) {
	insertSql := "insert into pwd_items(item,password,description,create_time,update_time) values(?,?,?,?,?)"
	stmt, err := db.Db.Connection.Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	// 执行
	result, err := stmt.Exec(
		pwdItem.Item,
		pwdItem.Password,
		pwdItem.Description,
		pwdItem.CreateTime,
		pwdItem.UpdateTime)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	pwdItem.Id = id
	return id, nil
}

// 将交易信息更新到数据库中
func (pwdItem *PwdItem) UpdateToDb()  (int64, error) {
	insertSql := "update pwd_items set password = ?, description = ?, create_time = ?, update_time = ? where _id = ?"
	stmt, err := db.Db.Connection.Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	// 执行
	result, err := stmt.Exec(
		pwdItem.Password,
		pwdItem.Description,
		pwdItem.CreateTime,
		pwdItem.UpdateTime,
		pwdItem.Id)
	if err != nil {
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affectedRows, nil
}

// 根据item key获取一条密码项
func GetItems() ([]*PwdItem, error) {
	querySql := "select _id,item,update_time from pwd_items"
	stmt, err := db.Db.Connection.Prepare(querySql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	// 开始查询
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	items := make([]*PwdItem, 0)
	for rows.Next() {
		pItem := NewPwdItem()
		err = rows.Scan(
			&pItem.Id,
			&pItem.Item,
			&pItem.UpdateTime)
		if err != nil {
			return nil, err
		}
		items = append(items, pItem)
	}
	return items, nil
}

// 根据key删除对应的项
func DeleteByItem(item string) (int64, error) {
	sql := "delete from pwd_items where item = ?"
	stmt, err := db.Db.Connection.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(item)
	if err != nil {
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affectedRows, nil
}