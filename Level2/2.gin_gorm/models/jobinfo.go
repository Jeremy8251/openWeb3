package models

import "time"

type Jobinfo struct {
	JobId      string `gorm:"column:jobId"`
	Jobhref    string
	Jobname    string
	Jobprice   string
	Infotag    string
	Jobdate    time.Time `gorm:"type:date"`
	Jobaddress string
	Detail     string
	DetailEl   string `gorm:"column:detailEl"`
	Code       string
	Number     string
	Gender     string
	Age        string
	Language   string
	Education  string
	Worktime   string
	Restinfo   string `gorm:"column:restinfo"`
	Restday    string
}

// 操作数据库的名称
func (Jobinfo) TableName() string {
	return "jobinfo"
}
