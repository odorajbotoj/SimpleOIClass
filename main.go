// zjyz oi lab
// work by 2207xuezihao (odorajbotoj)
// jf3 control system
// teacher
// version 1
// 23 09 25
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	SEND_DIR      string = "send/"  // 教师下发的文件
	UPLD_ROOT_DIR string = "upld/"  // 学生上传的文件根目录
	ID_MAP_DIR    string = "idmap/" // ip与id的对照规则
	PORT          string = ":8080"  // 服务使用的端口号
)

func getSend() string { // 获取教师下发的文件
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
	name, err := ioutil.ReadFile(ID_MAP_DIR + ip)
	if err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(ID_MAP_DIR+ip, []byte(""), 0644)
			if err != nil {
				log.Println("get id err: ", err)
			}
			return "已新建记录，请联系教师修改", 0
		}
		return "", 0
	}
	if string(name) == "" {
		return "请联系教师修改", 0
	}
	return string(name), len(name)
}

func uFunc(w http.ResponseWriter, r *http.Request) { // 处理上传文件的POST
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
		fn := f.Filename
		if f.Size > 102400 || strings.HasSuffix(fn, ".cpp") { // 非cpp文件不收取
			fr, _ := f.Open()
			fo, _ := os.Create(UPLD_ROOT_DIR + id + "/" + fn)
			io.Copy(fo, fr)
			fo.Close()
			fr.Close()
			log.Println(ip + " '" + id + "' uploaded '" + fn + "'")
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
	w.Write([]byte(fmt.Sprintf(ROOT_HTML, ip, id, getSend(), getUpld(ip))))
	return
}

func main() {
	log.Println("ZJYZIT LAB")
	log.Println("Simple OI Class")
	log.Println("teacher")
	log.Println("version 1")

	// http server
	http.HandleFunc("/", rootFunc)
	http.HandleFunc("/send", sendFunc)
	http.HandleFunc("/u", uFunc)
	http.HandleFunc("/del", delFunc)

	log.Fatalln(http.ListenAndServe(PORT, nil))
}

const ROOT_HTML = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Simple OI Class</title>
    <style type="text/css">
      a:link{color: black;}
      a:visited{color: red;}
      a:hover{color: orange;}
      a:active{color: brown;}
    </style>
  </head>
  <body>
    <font size="7">Simple OI Class</font>
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
        <p>· 提交源代码（如problem.cpp）, 其它文件不会被收取. 根据NOI相关规定，单个文件大小不得大于100KB. </p>
        <ul>
          <form method="post" action="/u" enctype="multipart/form-data">
            <input type="file" id="file" name="file" accept=".cpp" multiple style="background:white;border-radius:3px;" />
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
    <p>ZJYZIT LAB</p>
    <p>2023.09.26</p>
  </body>
</html>
`
