package models

import "time"

type UserFile struct {
	Id                 int
	Identity           string
	UserIdentity       string
	ParentId           int
	RepositoryIdentity string
	Ext                string
	Name               string
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated_at"`
	DeletedAt          time.Time `xorm:"deleted_at"`
}

func (r *UserFile) TableName() string {
	return "user_file"
}
