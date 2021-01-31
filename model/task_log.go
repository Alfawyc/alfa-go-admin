package model

import "time"

type TaskLog struct {
	Id         int       `json:"id" form:"id" gorm:"id"`
	TaskId     int       `json:"task_id" form:"task_id" gorm:"task_id"`
	Result     string    `json:"result" form:"result" gorm:"result"`
	RetryTimes int8      `json:"retry_times" form:"retry_times" gorm:"retry_times"`
	StartTime  time.Time `json:"start_time" form:"start_time" gorm:"start_time"`
	EndTime    time.Time `json:"end_time" form:"end_time" gorm:"end_time"`
}

func (TaskLog) TableName() string {
	return "task_log"
}
