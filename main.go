package main

import (
	"Screenshotweb/utils"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"
)

func main() {
	var sleep int64
	flag.Int64Var(&sleep, "s", 10, "延时时间,分钟")
	flag.Parse()
	for i := 0; ; i++ {
		netime := time.Now()
		lines := utils.LoadFile("./url.txt")

		fmt.Println(netime.Format("2006-01-02 15:04:05") + "开始测试")
		for _, line := range lines {
			parse, err := url.Parse(line)
			if err != nil {
				fmt.Println("解析网址错误")
				return
			}
			fmt.Println("开始获取: " + line)
			isOK, file := utils.Screenshot(line)
			if isOK {
				utils.Whitemark("./"+file, parse.Host, netime)
				os.Remove("./" + file)
			} else {
				fmt.Println("网站链接失败: " + line)
			}

		}
		fmt.Println("开始下一轮等待")
		time.Sleep(time.Duration(sleep) * time.Second)
	}

}
