package model

type TaskList struct {
	Name          string `json:"name" form:"name" gorm:"name"`
	Level         int8   `json:"level" form:"level" gorm:"level"`
	DependId      int    `json:"depend_id" form:"depend_id" gorm:"depend_id"`
	DependStatus  int8   `json:"depend_status" form:"depend_status" gorm:"depend_status"`
	Spec          string `json:"spec" form:"spec" gorm:"spec"`
	Protocol      int8   `json:"protocol" form:"protocol" gorm:"protocol"`
	Command       string `json:"command" form:"command" gorm:"command"`
	HttpMethod    string `json:"http_method" form:"http_method" gorm:"http_method"`
	Timeout       int    `json:"timeout" form:"timeout" gorm:"timeout"`
	RetryTimes    int8   `json:"retry_times" form:"retry_times" gorm:"retry_times"`
	RetryInterval int    `json:"retry_interval" form:"retry_interval" gorm:"retry_interval"`
	Remark        string `json:"remark" form:"remark" gorm:"remark"`
	Status        int8   `json:"status" form:"status" gorm:"status"`
	CreatedBy     int    `json:"created_by" form:"created_by" gorm:"created_by"`
	BaseModel
}

func (TaskList) TableName() string {
	return "task_list"
}
