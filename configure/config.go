package configure

import (
	"github.com/xxlixin1993/easyGo/utils"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var appConfig *Config

var (
	DefaultSelection  = "local"
	DefaultComment    = []byte{'#'}
	DefaultCommentSem = []byte{';'}
)

type Config struct {
	// 并发时限制
	sync.RWMutex

	// Section:key=value
	data map[string]map[string]string
}

// 错误码
const (
	KInitConfigError = iota + 1
	KInitLogError
)

// 错误信息
const (
	KUnknownTypeMsg = "unknown type"
)

// 初始化配置
func InitConfig(filePath string, mod string) error {
	if !utils.FileExists(filePath) {
		return errors.New("no such file or dir")
	}

	appConfig = &Config{}
	err := appConfig.parse(filePath, mod)
	if err != nil {
		return err
	}

	return nil
}

// 解析配置文件
func (c *Config) parse(fileName string, mod string) error {
	c.Lock()
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer c.Unlock()
	defer f.Close()

	buf := bufio.NewReader(f)

	var section string
	var lineNum int

	for {
		lineNum++
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if bytes.Equal(line, []byte{}) {
			continue
		} else if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		switch {
		case bytes.HasPrefix(line, DefaultComment):
			continue
		case bytes.HasPrefix(line, DefaultCommentSem):
			continue
		case bytes.HasPrefix(line, []byte{'['}) && bytes.HasSuffix(line, []byte{']'}):
			section = string(line[1 : len(line)-1])
			continue
		default:
			if section == mod {
				optionVal := bytes.SplitN(line, []byte{'='}, 2)
				if len(optionVal) != 2 {
					return fmt.Errorf("parse %s the content error : line %d , %s = ? ", fileName, lineNum, optionVal[0])
				}
				option := bytes.TrimSpace(optionVal[0])
				value := bytes.TrimSpace(optionVal[1])
				c.AddConfig(section, string(option), string(value))
			}
		}
	}

	return nil
}

// 添加一个新配置
func (c *Config) AddConfig(section string, option string, value string) bool {
	if section == "" {
		section = DefaultSelection
	}

	if len(c.data) == 0 {
		c.data = make(map[string]map[string]string)
	}

	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
	}

	_, ok := c.data[section][option]
	c.data[section][option] = value

	return !ok
}

// Get section.key or key
func (c *Config) get(key string) string {
	var (
		section string
		option  string
	)

	keys := strings.Split(strings.ToLower(key), "::")

	if len(keys) >= 2 {
		section = keys[0]
		option = keys[1]
	} else {
		section = DefaultSelection
		option = keys[0]
	}

	if value, ok := c.data[section][option]; ok {
		return value
	}

	return ""
}

func DefaultString(key string, defaultVal string) string {
	if v := appConfig.String(key); v != "" {
		return v
	}
	return defaultVal
}

func DefaultStrings(key string, defaultVal []string) []string {
	if v := appConfig.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func DefaultBool(key string, defaultVal bool) bool {
	if b, err := appConfig.Bool(key); err == nil {
		return b
	}
	return defaultVal
}

func DefaultInt(key string, defaultVal int) int {
	if b, err := appConfig.Int(key); err == nil {
		return b
	}
	return defaultVal
}

func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.get(key))
}

func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.get(key))
}

func (c *Config) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.get(key), 10, 64)
}

func (c *Config) Float64(key string) (float64, error) {
	return strconv.ParseFloat(c.get(key), 64)
}

func (c *Config) String(key string) string {
	return c.get(key)
}

func (c *Config) Strings(key string) []string {
	v := c.get(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ",")
}
