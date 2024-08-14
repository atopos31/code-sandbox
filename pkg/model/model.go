package model

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type CodeRequest struct {
	Code     string   `json:"code" binding:"required"`
	Language string   `json:"language" binding:"oneof=go c cpp java python"`
	MaxTime  float64  `json:"max_time" binding:"required"`
	MaxMem   int      `json:"max_mem" binding:"required"`
	Stdin    []string `json:"stdin"`
}

type CodeResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Meta []*CodeMETA `json:"meta"`
}

type CodeMETA struct {
	Time     float64 `json:"time"      metric:"time"`
	TimeWall float64 `json:"time_wall" metric:"time-wall"`
	MaxRss   int     `json:"max_rss"   metric:"max-rss"`
	Killed   bool    `json:"killed"    metric:"killed"`
	Message  string  `json:"message"   metric:"message"`
	Status   string  `json:"status"    metric:"status"`
	Exitsig  int     `json:"exitsig"   metric:"exitsig"`
	Stderr   string  `json:"stderr"`
	StdOut   string  `json:"stdout"`
}

type BuildMeta struct {
	Message string `json:"message"   metric:"message"`
	Status  string `json:"status"    metric:"status"`
	Exitsig int    `json:"exitsig"   metric:"exitsig"`
	Stderr  string `json:"stderr"`
}

type RunMeta struct {
	Time     float64 `json:"time"      metric:"time"`
	TimeWall float64 `json:"time_wall" metric:"time-wall"`
	MaxRss   int     `json:"max_rss"   metric:"max-rss"`
	Killed   bool    `json:"killed"    metric:"killed"`
	Message  string  `json:"message"   metric:"message"`
	Status   string  `json:"status"    metric:"status"`
	Exitsig  int     `json:"exitsig"   metric:"exitsig"`
	Stdin    string  `json:"stdin"`
	Stderr   string  `json:"stderr"`
	StdOut   string  `json:"stdout"`
}

func NewBuildMeta(errPath string, metaPath string) *BuildMeta {
	buildMeta := new(BuildMeta)
	data, err := os.ReadFile(errPath)
	if err != nil {
		buildMeta.Stderr = "err: read buildErrPath failed"
	} else {
		buildMeta.Stderr = string(data)
	}

	if err := MarshalMetaFrompath(buildMeta, metaPath); err != nil {
		return buildMeta
	}
	return buildMeta
}

func NewRunMeta(errPath string, outPath string, metaPath string) *RunMeta {
	runMeta := new(RunMeta)
	errdata, err := os.ReadFile(errPath)
	if err != nil {
		runMeta.Stderr = "err: read buildErrPath failed"
	} else {
		runMeta.Stderr = string(errdata)
	}
	outdata ,err := os.ReadFile(outPath)
	if err != nil {
		runMeta.StdOut = "err: read buildErrPath failed"
	} else {
		runMeta.StdOut = string(outdata)
	}

	if err := MarshalMetaFrompath(runMeta, metaPath); err != nil {
		return runMeta
	}
	return runMeta
}

func MarshalMetaFrompath[T *BuildMeta | *RunMeta](b T, filePath string) error {
	v := reflect.ValueOf(b).Elem()
	data, err := os.Open(filePath)
	defer data.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		key := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		} else {
			continue
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.Tag.Get("metric") == key {
				// 根据字段的类型设置值
				switch v.Field(i).Kind() {
				case reflect.Float64:
					floatValue, _ := strconv.ParseFloat(value, 64)
					v.Field(i).SetFloat(floatValue)
				case reflect.Int:
					intValue, _ := strconv.Atoi(value)
					v.Field(i).SetInt(int64(intValue))
				default:
					v.Field(i).SetString(value)
				}
			}
		}
	}

	return nil
}

func NewCodeMETA(stderrPath string, stdoutPath string, metaPath string) *CodeMETA {
	codeMETA := new(CodeMETA)
	if err := codeMETA.SetStderrFrompath(stderrPath); err != nil {
		return nil
	}
	if err := codeMETA.SetStdOutFrompath(stdoutPath); err != nil {
		return nil
	}
	if err := codeMETA.MarshalFrompath(metaPath); err != nil {
		return nil
	}
	return codeMETA
}

func (c *CodeMETA) SetStderrFrompath(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	c.Stderr = string(data)
	return nil
}

func (c *CodeMETA) SetStdOutFrompath(filepath string) error {
	data, err := os.Open(filepath)
	defer data.Close()
	if err != nil {
		return err
	}
	var lines []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查扫描器是否遇到错误
	if err := scanner.Err(); err != nil {
		return err
	}
	c.StdOut = strings.Join(lines, "\n")

	return nil
}

func (c *CodeMETA) MarshalFrompath(filepath string) error {
	v := reflect.ValueOf(c).Elem()
	data, err := os.Open(filepath)
	defer data.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		key := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		} else {
			continue
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.Tag.Get("metric") == key {
				// 根据字段的类型设置值
				switch v.Field(i).Kind() {
				case reflect.Float64:
					floatValue, _ := strconv.ParseFloat(value, 64)
					v.Field(i).SetFloat(floatValue)
				case reflect.Int:
					intValue, _ := strconv.Atoi(value)
					v.Field(i).SetInt(int64(intValue))
				default:
					v.Field(i).SetString(value)
				}
			}
		}
	}
	return nil
}
