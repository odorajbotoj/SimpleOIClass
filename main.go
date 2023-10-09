// zjyz oi lab
// work by 2207xuezihao (odorajbotoj)
// Simple OI Class
// server
// version 4.0.1
// 23 10 09
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var ( // 默认配置
	SEND_DIR      string = "send/"           // 教师下发的文件
	UPLD_ROOT_DIR string = "upld/"           // 学生上传的文件根目录
	ID_MAP_DIR    string = "idmap/"          // ip与id的对照规则
	PORT          string = ":8080"           // 服务使用的端口号
	ACCEPT        string = ".cpp"            // 允许上传的文件后缀
	TITLE         string = "Simple OI Class" // 网页标题
	REG           bool   = false             // 是否为注册模式
)

func readKVMap(filename string) map[string]string {
	// 读取文件
	// 使用K=V形式
	// https://www.cnblogs.com/rickiyang/p/11074169.html
	m := make(map[string]string)
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Println("无法读取" + filename)
		return m
	}
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("在读取" + filename + "的过程中出错. Skip.")
			continue
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		m[key] = value
	}
	return m
}

func getSend() string { // 获取教师下发的文件
	if REG {
		return ""
	}
	rd, err := ioutil.ReadDir(SEND_DIR)
	if err != nil {
		log.Println("readSend: ", err)
		return fmt.Sprintln("readSend: ", err)
	}
	ret := ""
	for _, fi := range rd {
		if !fi.IsDir() {
			ret = fmt.Sprintf("%s<li><a href=\"send?fn=%s\" target=\"_blank\">%s</a>&nbsp;&nbsp;(%d字节)</li>", ret, url.QueryEscape(fi.Name()), fi.Name(), fi.Size())
		}
	}
	return ret
}

func sendFunc(w http.ResponseWriter, r *http.Request) { // 下载教师下发的文件
	if REG {
		return
	}
	fn, err := url.QueryUnescape(r.URL.Query().Get("fn"))
	if fn == "" && err == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404. File not found."))
	} else if err != nil {
		log.Println("sendFunc: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint("Err: ", err)))
	} else {
		b, err := ioutil.ReadFile(SEND_DIR + fn)
		if err != nil {
			log.Println("sendFunc: ", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprint("Err: ", err)))
			return
		}
		finfo, _ := os.Stat(SEND_DIR + fn)
		w.Header().Set("Content-Disposition", "attachment; filename="+fn)
		w.Header().Set("Content-Type", http.DetectContentType(b))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", finfo.Size()))
		w.Write(b)
		return
	}
}

func getIP(r *http.Request) string { // 获取ip地址
	forwarded := r.Header.Get("X-FORWARDED_FOR")
	if forwarded != "" {
		return strings.Split(forwarded, ":")[0]
	} else {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
}

func getID(ip string) (string, int) { // 获取学生姓名
	name, err := ioutil.ReadFile(ID_MAP_DIR + ip + ".txt")
	if len(name) >= 3 {
		if name[0] == 0xEF || name[1] == 0xBB || name[2] == 0xBF { // 去除UTF-8 BOM以适配Windows记事本
			name = name[3:]
		}
	}
	r := strings.NewReplacer("\r", "", "\n", "") // 去除CR, LF
	if err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(ID_MAP_DIR+ip+".txt", []byte(""), 0644)
			if err != nil {
				log.Println("get id err: ", err)
			}
			return "已新建记录，未设置姓名", 0
		}
		return "", 0
	}
	if string(name) == "" {
		return "未设置姓名", 0
	}
	return r.Replace(string(name)), len(name)
}

func regFunc(w http.ResponseWriter, r *http.Request) {
	if !REG {
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("ParseForm err: ", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	ip := getIP(r)
	err = ioutil.WriteFile(ID_MAP_DIR+ip+".txt", []byte(r.Form["username"][0]), 0644)
	if err != nil {
		log.Println("set id err: ", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func uFunc(w http.ResponseWriter, r *http.Request) { // 处理上传文件的POST
	if REG {
		return
	}
	ip := getIP(r)
	id, idn := getID(ip)
	if idn == 0 { // 未填写姓名则不做处理
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	_, err := os.Stat(UPLD_ROOT_DIR + id)
	if os.IsNotExist(err) { // 未建立文件夹则建立
		os.Mkdir(UPLD_ROOT_DIR+id, os.ModePerm)
	}
	r.ParseMultipartForm(409600)
	files, ok := r.MultipartForm.File["file"]
	if !ok { // 出错则取消
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	for _, f := range files {
		if f.Size < 102400 {
			if find := strings.Contains(ACCEPT, path.Ext(f.Filename)); find {
				fr, _ := f.Open()
				fo, _ := os.Create(UPLD_ROOT_DIR + id + "/" + f.Filename)
				io.Copy(fo, fr)
				fo.Close()
				fr.Close()
				log.Println(ip + " '" + id + "' uploaded '" + f.Filename + "'")
			}
		} else {
			continue
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func getUpld(ip string) string { // 列出上传的文件
	id, idn := getID(ip)
	if idn == 0 {
		return ""
	}
	rd, err := ioutil.ReadDir(UPLD_ROOT_DIR + id)
	if err != nil {
		log.Println("readUpld: ", err)
		return ""
	}
	ret := ""
	for _, fi := range rd {
		if !fi.IsDir() {
			ret = fmt.Sprintf("%s<li><span>%s</span>&nbsp;&nbsp(%d字节)&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"del?fn=%s\" onclick=\"return confirm('确认删除 %s ?');\">[删除]</a></li>", ret, fi.Name(), fi.Size(), url.QueryEscape(fi.Name()), fi.Name())
		}
	}
	return ret
}

func delFunc(w http.ResponseWriter, r *http.Request) { // 删除上传的文件
	if REG {
		return
	}
	fn, err := url.QueryUnescape(r.URL.Query().Get("fn"))
	ip := getIP(r)
	id, idn := getID(ip)
	if fn != "" && err == nil && idn != 0 {
		err = os.Remove(UPLD_ROOT_DIR + id + "/" + fn)
		if err != nil {
			log.Println("del file err: ", err)
		}
		log.Println(ip + " '" + id + "' deleted '" + fn + "'")
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	id, _ := getID(ip)
	log.Println(ip + " '" + id + "' connected.")
	if !REG {
		w.Write([]byte(fmt.Sprintf(ROOT_HTML, TITLE, TITLE, ip, id, getSend(), ACCEPT, ACCEPT, getUpld(ip))))
	} else {
		w.Write([]byte(fmt.Sprintf(REG_HTML, TITLE, TITLE, ip, id)))
	}
	return
}

func init() {
	log.Println("ZJYZIT LAB")
	log.Println("Simple OI Class")
	log.Println("server")
	log.Println("version 4.0.1")

	config := readKVMap("config.txt")
	if _, ok := config["SEND"]; ok {
		SEND_DIR = config["SEND"]
		log.Println("set SEND_DIR = ", SEND_DIR)
	}
	if _, ok := config["UPLD"]; ok {
		UPLD_ROOT_DIR = config["UPLD"]
		log.Println("set UPLD_ROOT_DIR = ", UPLD_ROOT_DIR)
	}
	if _, ok := config["IDMAP"]; ok {
		ID_MAP_DIR = config["IDMAP"]
		log.Println("set ID_MAP_DIR = ", ID_MAP_DIR)
	}
	if _, ok := config["PORT"]; ok {
		PORT = config["PORT"]
		log.Println("set PORT = ", PORT)
	}
	if _, ok := config["ACCEPT"]; ok {
		ACCEPT = config["ACCEPT"]
		log.Println("set ACCEPT = ", ACCEPT)
	}
	if _, ok := config["TITLE"]; ok {
		TITLE = config["TITLE"]
		log.Println("set TITLE = ", TITLE)
	}
	if v, ok := config["REG"]; ok {
		if v == "ON" {
			REG = true
			log.Println("set REG =  true")
		}
	}
}

func main() {
	// http server
	http.HandleFunc("/", rootFunc)
	if !REG {
		http.HandleFunc("/send", sendFunc)
		http.HandleFunc("/u", uFunc)
		http.HandleFunc("/del", delFunc)
	} else {
		http.HandleFunc("/reg", regFunc)
	}
	log.Fatalln(http.ListenAndServe(PORT, nil))
}

const ROOT_HTML = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>%s</title>
    <style type="text/css">
      a:link{color: black;}
      a:visited{color: red;}
      a:hover{color: orange;}
      a:active{color: brown;}
    </style>
  </head>
  <body>
    <font size="7">%s</font>
    <input type=button value=点击刷新 onclick="location=location" style="background:#DDFCFA;border-radius:30px;height:48px;width:256px;color:dodgerblue;" />
    <hr/>
    <div id="id">
      <ul>
        <li>IP:&nbsp&nbsp&nbsp&nbsp%s</li>
        <li>Name:&nbsp&nbsp%s</li>
      </ul>
    </div>
    <hr/>
    <div id="sendFile">
      <h2>文件</h2>
      <div style="background:#DDFCFA;border-radius:15px;">
        <ul>
          %s
        </ul>
      </div>
    </div>
    <hr/>
    <div id="upload">
      <h2>提交的文件</h2>
      <div style="background:#DDFCFA;border-radius:15px;">
        <p>· 提交文件（如problem.cpp）, 单个文件大小不得大于100KB, 一次性提交文件大小总和不得大于400KB. 允许的类型：%s</p>
        <ul>
          <form method="post" action="/u" enctype="multipart/form-data">
            <input type="file" id="file" name="file" accept="%s" multiple style="background:white;border-radius:3px;" />
            <input type="submit" value="提交" />
          </form>
        </ul>
        <ul>
          %s
        </ul>
      </div>
    </div>
    <hr/>
    <p>不要哀求 学会进取 若是如此 终有所获</p>
    <p>物来顺应 未来不迎 当时不杂 既过不恋</p>
    <p>SimpleOIClass 4.0.1</p>
    <p>ZJYZIT LAB</p>
    <p>2023.10.09</p>
  </body>
</html>
`

const REG_HTML = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>%s</title>
  </head>
  <body>
    <font size="7">%s</font>
    <input type=button value=点击刷新 onclick="location=location" style="background:#DDFCFA;border-radius:30px;height:48px;width:256px;color:dodgerblue;" />
    <hr/>
    <div id="id">
      <ul>
        <li>IP:&nbsp&nbsp&nbsp&nbsp%s</li>
        <li>Name:&nbsp&nbsp%s</li>
      </ul>
    </div>
    <hr/>
    <div id="reg">
      <h2>提交姓名</h2>
      <div style="background:#DDFCFA;border-radius:15px;">
        <p>· 提交姓名，该名称将绑定此IP地址。</p>
        <form method="post" action="/reg">
          <input type="text" id="username" name="username" style="background:white;border-radius:3px;" />
          <input type="submit" value="提交" />
        </form>
      </div>
    </div>
    <hr/>
    <p>不要哀求 学会进取 若是如此 终有所获</p>
    <p>物来顺应 未来不迎 当时不杂 既过不恋</p>
    <p>SimpleOIClass 4.0.1</p>
    <p>ZJYZIT LAB</p>
    <p>2023.10.09</p>
  </body>
</html>
`
