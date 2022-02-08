package xmysql

type (
	//MysqlConfig mysql 连接配置
	MYSQLConfig struct {
		DNS             string `json:"DNS" yaml:"DNS"`
		MaxIdleConns    int    `json:"MaxIdleConns" yaml:"MaxIdleConns"`
		MaxPoolSize     int    `json:"MaxPoolSize" yaml:"MaxPoolSize"`
		MaxIdleTime     int64  `json:"MaxIdleTime" yaml:"MaxIdleTime"`         //毫秒 1000=1s
		ConnMaxLifetime int    `json:"ConnMaxLifetime" yaml:"ConnMaxLifetime"` // 秒
		Log             int    `json:"Log" yaml:"Log"`                         //是否开启日志 0 开启
	}
)

var (
	enableLog = false
)

func EnableLog() {
	enableLog = true
}
func IsEnable() bool {
	return enableLog
}
