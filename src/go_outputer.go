package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var sqlFile string
var outputFile string
var format string
var separatorType string
var reset int
var confFile = "outputer.conf"

func main() {
	flag.StringVar(&format, "f", "csv", "导出类型 支持json和csv，默认csv")
	flag.StringVar(&separatorType, "t", "comma", "csv文件分隔符 comma tab,默认comma(逗号)")
	flag.StringVar(&sqlFile, "s", "qry.sql", "待执行的sql的名称")
	flag.StringVar(&outputFile, "o", "", "导出的文件名称")
	flag.IntVar(&reset, "r", 0, "-reset=1 重置数据库链接或者直接删除"+confFile)
	flag.Parse()

	separator := ","
	if separatorType == "tab" {
		separator = "\t"
	}

	//处理数据库连接输入
	var dbUrl string
	isInput := false
	if !PathExists(confFile) || reset == 1 {
		f := bufio.NewReader(os.Stdin)
		println("请输入数据库连接配置eg( user:password@tcp(ip:port)/dbname?charset=utf8 ):")
		dbUrl, _ = f.ReadString('\n')
		dbUrl = strings.ReplaceAll(dbUrl, "\n", "")
		if len(dbUrl) == 0 {
			println("没有输入数据库配置，将退出")
			return
		}
		isInput = true
	} else {
		data, err := ioutil.ReadFile(confFile)
		if err != nil {
			println("读取数据库配置错误", err.Error())
			return
		} else {
			dbUrl = AesDecrypt(string(data))
		}
	}

	//读取查询sql
	if !PathExists(sqlFile) {
		println("请在" + sqlFile + "里输入需要查询的sql语句")
		return
	}
	sqlData, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		println("读取sql错误！", err.Error())
		return
	}

	//连接数据库
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		println("数据库连接错误！", err.Error())
		return
	}

	rows, err := db.Query(string(sqlData))
	if err != nil {
		println("数据库查询错误", err.Error())
		println("请检查数据库连接配置，如要修改请删除" + confFile + "或者在启动参数加上--reset=true")
		return
	}

	//保存连接信息到配置文件
	if isInput {
		err = ioutil.WriteFile(confFile, []byte(AesEncrypt(dbUrl)), 0777)
		if err != nil {
			println("写入配置失败", err.Error())
		} else {
			println("数据库连接信息已被加密保存到" + confFile + ",下一次可直接使用")
		}
	}

	//创建导出文件
	var fileName string
	if (outputFile) != "" {
		if format == "csv" && !strings.Contains(outputFile, ".csv") {
			outputFile = outputFile + ".csv"
		} else if format == "json" && !strings.Contains(outputFile, ".json") {
			outputFile = outputFile + ".json"
		}
		fileName = outputFile
	} else {
		fileName = time.Now().Format("20060102-150405导出." + format)
	}
	_ = os.Remove(fileName)
	_, _ = os.Create(fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		println("init open " + err.Error())
		return
	}

	//写入header
	header, err := rows.Columns()
	var contentList []interface{}
	for range header {
		var val sql.NullString
		contentList = append(contentList, &val)
	}

	var resultList [][]string
	for rows.Next() {
		err := rows.Scan(contentList...)
		if err != nil {
			println("获取一行数据错误", err.Error())
			return
		} else {
			var strList []string
			for _, value := range contentList {
				sqlStr, ok := value.(*sql.NullString)
				if ok {
					strList = append(strList, (*sqlStr).String)
				}
			}
			resultList = append(resultList, strList)
		}
	}

	if format == "csv" {
		//写入header
		_, err = file.Write([]byte(strings.Join(header, separator) + "\n"))
		for _, rowList := range resultList {
			_, err = file.Write([]byte(strings.ReplaceAll(strings.Join(rowList, separator), "\n", " ") + "\n"))
			if err != nil {
				println("数据写入文件错误", err.Error())
				return
			}
		}
	} else if format == "json" {
		var jsonList []map[string]string
		for _, rowList := range resultList {
			tempMap := make(map[string]string)
			for headerIndex, headerValue := range header {
				tempMap[headerValue] = rowList[headerIndex]
			}
			jsonList = append(jsonList, tempMap)
		}

		writeData, err := json.Marshal(jsonList)
		if err != nil {
			println("json格式化错误", err.Error())
			return
		}

		_, err = file.Write(writeData)
		if err != nil {
			println("数据写入文件错误", err.Error())
			return
		}
	}

	println("数据导出完成,共导出" + strconv.Itoa(len(resultList)) + "条数据" + " 导出文件 " + fileName)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
