package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

//go:embed static/*
var content embed.FS

func main() {
	port := flag.Int("p", 8090, "监听端口")
	flag.Parse()

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(content, "static/index.html")
		if err != nil {
			log.Println("parse tmpl failed: " + err.Error())
			return
		}
		tmpl.Execute(rw, nil)
	})

	http.Handle("/static/", http.FileServer(http.FS(content)))

	errorResult := func(msg string, rw http.ResponseWriter) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(500)
		rw.Write([]byte(msg))
		log.Println(msg)
	}

	http.HandleFunc("/upload", func(rw http.ResponseWriter, r *http.Request) {
		span, err := strconv.Atoi(r.URL.Query().Get("span"))
		if err != nil {
			span = 10
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			errorResult("get form file failed: "+err.Error(), rw)
			return
		}
		log.Printf("开始处理文件: %s 间隔: %d\n", header.Filename, span)
		outDir := filepath.Join(os.Getenv("USERPROFILE"), "Desktop", header.Filename+"--output")
		err = os.Mkdir(outDir, 0666)
		if err != nil && !os.IsExist(err) {
			errorResult(fmt.Sprintf("创建导出目录失败, %s", err), rw)
			return
		}
		err = api.Split(file.(io.ReadSeeker), outDir, header.Filename, span, nil)
		if err != nil {
			errorResult("split file failed: "+err.Error(), rw)
			return
		}
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(200)
		rw.Write([]byte("文件切割成功, 存放路径: " + outDir))
	})
	s := fmt.Sprintf("localhost:%d", *port)
	log.Printf("\n\n****************************************\n\n请打开浏览器访问 http://%s \n\n****************************************\n\n", s)
	http.ListenAndServe(s, nil)
}
