package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var config sync.Map

func init() {
	//fmt.Print("Env:\n")
}
func Get[T int | float64 | bool | string](key string, value T) T {
	if value, ok := config.Load(key); ok {
		return value.(T)
	}
	env := getConfigFormEnv(key)
	if len(env) != 0 {
		value = StringToType[T](env)
		//fmt.Printf("    %s:%v\n", key, value)
	}
	config.Store(key, value)
	return value
}

func getConfigFormEnv(key string) string {
	keys := []string{key, strings.ToLower(key), strings.ToUpper(key)}
	for _, key = range keys {
		if value := os.Getenv(key); len(value) > 0 {
			return value
		}
	}
	return ""
}

// 将字符串转换为对应的类型
func StringToType[T any](s string) T {
	var zero T
	switch any(zero).(type) {
	case int:
		parsed, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return any(parsed).(T)
	case float64:
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return any(parsed).(T)
	case bool:
		parsed, err := strconv.ParseBool(s)
		if err != nil {
			panic(err)
		}
		return any(parsed).(T)
	case string:
		return any(s).(T)
	default:
		panic(fmt.Errorf("unsupported type: %T", zero))
	}
}
