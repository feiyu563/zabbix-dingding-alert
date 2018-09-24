package main

import (
	"os"
	"bytes"
	"net/http"
	"io/ioutil"
	"github.com/Unknwon/goconfig"
	"log"
	"time"
)
func main()  {
	//取输入参数1和2
	user:=os.Args[1]
	text:=os.Args[2]
	fileName := time.Now().Format("20060102") +".dingding.log"
	logFile,_:= os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	defer logFile.Close()
	Log := log.New(logFile,"[Info]", log.Ldate|log.Ltime) // log.Ldate|log.Ltime|log.Lshortfile
	Log.Println("开始发送消息!")
	SendMsg(user,text,Log)
}

//发送消息到钉钉
func SendMsg(user,text string,Log *log.Logger)  {
	jsonstring:=`{
     "msgtype": "markdown",
     "markdown": {"title":"DF云告警信息",
"text":"`+text+`"
     },
    "at": {
        "atMobiles": [
            "`+user+`"
        ], 
        "isAtAll": true
    }
 }`
	reader:=bytes.NewReader([]byte(jsonstring))
	cfg, err :=goconfig.LoadConfigFile("dingding.conf")
	if err != nil {
		Log.SetPrefix("[Err]")
		Log.Println("读取配置文件失败[config.ini]")
		return
	}
	url,err:=cfg.GetValue("setup","url")
	if err != nil {
		Log.SetPrefix("[Err]")
		Log.Fatalf("无法获取键值（%s）：%s", "url", err)
		return
	}
	resp:=Post(url,reader)
	Log.SetPrefix("[Info]")
	Log.Fatalf("消息发送完成,服务器返回内容：%s", string(resp))
}

func Post(url string,reader *bytes.Reader) []byte  {
	request,_ := http.NewRequest("POST", url, reader)
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, _ := client.Do(request)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	return respBytes
}