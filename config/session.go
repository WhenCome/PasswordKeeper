package config

import (
	"../utils/randutil"
	"../utils/fileutil"
	"time"
	"encoding/xml"
	"fmt"
)

// 会话
type Session struct {
	Token 		string
	CreateTime	int64
	UpdateTime	int64
	ExpireTime	int64
}

// 创建新会话
func NewSession() *Session {
	token := randutil.GetRandAlphaDigitString(64)
	createTime := time.Now().Unix()
	expireTime := getSessionExpireTime(createTime)  // 暂时默认会话有效期为12个小时
	sess := &Session{
		Token:token,
		CreateTime:createTime,
		UpdateTime:createTime,
		ExpireTime:expireTime,
	}
	// ignore the error of save session info
	sess.save()
	return sess
}

// 加载会话信息
func LoadSession() *Session {
	if !fileutil.IsFileExists(SessionTokenFile) {
		return nil
	}
	xmlBytes,err := fileutil.GetContents(SessionTokenFile)
	if err != nil {
		fmt.Printf("Load session failed : %s \n", err)
		return nil
	}
	sess := &Session{}
	err = xml.Unmarshal(xmlBytes, sess)
	if err != nil {
		fmt.Printf("Parse session failed : %s \n", err)
		return nil
	}
	return sess
}

// 获取过期时间
func getSessionExpireTime(t int64) int64 {
	return t + 12 * 3600
}

// 检查会话是否有效
func (sess *Session) IsValid() bool {
	if sess == nil {
		return false
	}
	now := time.Now().Unix()
	if now > sess.ExpireTime {
		return false
	}
	return true
}

// 重新激活一次session
func (sess *Session) Revalid() {
	if sess == nil {
		return
	}
	now := time.Now().Unix()
	sess.UpdateTime = now
	sess.ExpireTime = getSessionExpireTime(now)
	sess.save()
}

// 保存会话信息
func (sess *Session) save() error {
	byteXml, err := xml.Marshal(sess)
	if err != nil {
		return err
	}
	strXml := string(byteXml)
	_,err = fileutil.WriteFile(SessionTokenFile, strXml)
	return err
}

// 销毁会话
func (sess *Session) Destroy() error {
	if !fileutil.IsFileExists(SessionTokenFile) {
		return nil
	}
	sess.ExpireTime = time.Now().Unix()
	return sess.save()
}
