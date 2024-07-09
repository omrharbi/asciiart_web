package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ascii "omar/Functions"
)

type data struct {
	Text    string
	Banner  string
	Message string
}

type messageErrors struct {
	Code        int
	Message     string
	Description string
}

type errors struct {
	Message string
}

var tmp1 *template.Template
var SourceDir="src"
func main() {
	port := ":8080"
	mux:=http.NewServeMux()
	fmt.Printf("http://localhost%s", port)
	tmp1 = template.Must(template.ParseGlob("templates/*.html"))
	mux.HandleFunc("/styles/", customAssetsHandler)
	mux.HandleFunc("/", ShowInfromations)
	mux.HandleFunc("/ascii-art", SendPost)
	
	log.Fatal(http.ListenAndServe(port, mux))

}

func customAssetsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/styles"):]
	fullPath := filepath.Join("src", path)
	fileInfo, err := os.Stat(fullPath)

	if !os.IsNotExist(err) && !fileInfo.IsDir() {
		http.StripPrefix("/styles/", http.FileServer(http.Dir("src"))).ServeHTTP(w, r)
	}else{
		ErrorHandler(w, http.StatusNotFound, "Status Not Found", ("This Page not available"))
		return
	}
}

func ShowInfromations(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, "Status Not Found", ("This Page not available"))
		return
	}
	err:=tmp1.ExecuteTemplate(w, "index.html", nil)
	if err!=nil{
		ErrorHandler(w, http.StatusNotFound, "Status Not Found", ("This Page not available"))
		return
	}

}

func SendPost(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed ", ("This Page not available"))
		return
	}
	text := r.PostFormValue("text")
	// if text == "" {
	// 	http.Redirect(w, r, "/?error=emptytext", http.StatusSeeOther)
	// 	return
	// }
	for _, v := range text {
		if v < 32 || v > 126 {
			data := errors{
				Message: "the Text Uncorect",
			}
			tmp1.ExecuteTemplate(w, "index.html", data)
			// ErrorHandler(w, http.StatusBadRequest, "Status Bad Request", ("The Text Is Not in The Ascii Art"))
			return
		}
	}
	if text == "" {
		data := errors{
			Message: "The Input It's Empty",
		}
		tmp1.ExecuteTemplate(w, "index.html", data)
		return
	}
	banner := r.PostFormValue("banner")
	err, valid := CheckFile(banner)
	if valid || err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", ("you don't Have this Page Or this Page Is Empty "))
		return
	}

	if banner == "shadow.txt" || banner == "standard.txt" || banner == "thinkertoy.txt" {
		slice, count := ascii.WriteTextFileAscii(banner)
		if count != 855 {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", ("The file does not contain all the letters "))
			return
		}
	 	result := ascii.PrintAsciArt(slice, text)
		datasend := data{
			Text:   result,
			Banner: banner,
		}
		asciiArtPage := "templates/ascii-art.html"
		_, valid := CheckFile(asciiArtPage)
		if valid {
			ErrorHandler(w, 500, "Bad Request", "this Page Is Empty  ")
			return
		} else {
			err := tmp1.ExecuteTemplate(w, "ascii-art.html", datasend)
			if err != nil {
				ErrorHandler(w, 500, "Bad Request", "you don't")
				return
			}
		}
	}
}

func ErrorHandler(w http.ResponseWriter, code int, message string, description string) {
	errors := messageErrors{
		Code:        code,
		Message:     message,
		Description: description,
	}
	w.WriteHeader(code)
	tmp1.ExecuteTemplate(w, "pageErrors.html", errors)
}

func CheckFile(path string) (error, bool) {
	file := filepath.Join(path)
	filepath, err := os.Stat(file)
	if os.IsNotExist(err) {
		return err, false
	}
	return nil, filepath.Size() == 0
}

