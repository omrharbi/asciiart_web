package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	ascii "omar/Functions"
	"path/filepath"
)

type DataSet struct{
	Text string
	Banner string
	Valid bool

}
type ErrorsMes struct{
	Code int
	MessageError string
	Description string
}

var tmp *template.Template
var port string

func main() {
	port =":8081"
	fmt.Printf("http://localhost%s \n", port)
	tmp=template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", GetIndexPage)
	http.HandleFunc("/asci_art", SendData)
	log.Fatal(http.ListenAndServe(port, nil))
}

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w,http.StatusNotFound,"Status Not Found","page Not Found")
		return
	}
	data:=DataSet{
		Valid: true,
	}
	tmp.ExecuteTemplate(w, "index.html", data)
}
func SendData(w http.ResponseWriter, r *http.Request){
	valid := true
	if r.Method!=http.MethodPost{
		ErrorHandler(w,http.StatusMethodNotAllowed,"Method Not Allowed","")
		return
	}
	text:=r.PostFormValue("text")
	if CheckWords(w,text) {
		return
	}
	banner:=r.PostFormValue("banner")
	if CheckExtension(w,banner){
		return
	}
	ascii_web:=ascii.WriteTextFileAscii(banner)

	result:=ascii.PrintAsciArt(ascii_web,text)
	data:=DataSet{
		Text:result,
		Banner: banner,
		Valid: valid,
	}
	tmp.ExecuteTemplate(w, "asci_art.html", data)
}

func CheckExtension(w http.ResponseWriter,banner string) bool{
	if filepath.Ext(banner)!=".txt"{
		ErrorHandler(w, http.StatusUnsupportedMediaType,"Invalid file extension",".txt")
		return true
	}
	return false
}
func CheckWords(w http.ResponseWriter, word string) bool{
	for _,n:=range word{
		if n<32 || n>126{
		ErrorHandler(w,http.StatusInternalServerError,"Internal Server Error","")
		return true
		}
	}
	if len(word)==0{
		ErrorHandler(w,http.StatusNoContent,"No Content","You Shold be Write Word ")
		return true
	}
 return	false
}

func ErrorHandler(w http.ResponseWriter, code int, errors string,des string){
	errorMes:=ErrorsMes{
		Code: code,
		MessageError: errors,
		Description: des,
	}
	w.WriteHeader(code)
	tmp.ExecuteTemplate(w,"pageErrors.html",errorMes)
}