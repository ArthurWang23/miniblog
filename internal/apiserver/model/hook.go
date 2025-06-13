package model

import (
	"github.com/ArthurWang23/miniblog/internal/pkg/auth"
	"github.com/ArthurWang23/miniblog/internal/pkg/rid"
	"gorm.io/gorm"
)

// 数据库禁止保存明文密码
// 用户密码在入库前需要进行加密处理，为了加密明文密码字符串
// 通过BeforeCreate钩子实现在入库前对密码进行加密处理

func (m *UserM) BeforeCreate(tx *gorm.DB) error {
	var err error
	m.Password, err = auth.Encrypt(m.Password)
	if err != nil {
		return err
	}
	return nil
}

// 添加数据库表userID和postID字段的自动生成钩子，用来生成并保存记录的唯一标识符
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.New(uint64(m.ID))
	return tx.Save(m).Error
}

func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))
	return tx.Save(m).Error
}
