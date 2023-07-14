package support

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Env 获取配置文件(未获取到则使用第二个参数作为默认值)
func Env(name string, value ...any) any {
	var env any
	if len(value) == 1 {
		env = value[0]
	}
	if localEnv := os.Getenv(name); localEnv != "" {
		//处理Bool类型
		if strings.ToUpper(localEnv) == "TRUE" {
			env = true
			return env
		} else if strings.ToUpper(localEnv) == "FALSE" {
			env = false
			return env
		}
		//处理int类型
		reg := regexp.MustCompile("^(?:[0-9]|[1-9][0-9]{1,3}|10000)$") // 匹配1-10000以内的全部数字
		match := reg.MatchString(localEnv)
		if match {
			num, _ := strconv.Atoi(localEnv)
			env = num
			return env
		}
		env = localEnv
	}
	return env
}
