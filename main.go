package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

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

var (
	tmp1      *template.Template
	SourceDir = "src"
	datasend  data
)

func main() {
	port := ":8080"
	// ux := http.NewServeMux()
	fmt.Printf("http://localhost%s", port)
	tmp1 = template.Must(template.ParseGlob("templates/*.html"))

	// http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("srr"))))
	http.HandleFunc("/styles/", customAssetsHandler)
	http.HandleFunc("/", ShowInfromations)
	http.HandleFunc("/ascii-art", SendPost)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/download", downloadHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}

// this func to save cookie in header and redirect url to /Download
func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	content := url.QueryEscape(datasend.Text)
	// Set a cookie with the content
	cookie := &http.Cookie{
		Name:     "fileContent",
		Value:    content,
		Expires:  time.Now().Add(5 * time.Second),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/download", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("fileContent") // this to get name cookie and return name plus value ==>filecontent="value"
	fmt.Println(cookie)
	if err != nil {
		http.Error(w, "No content found", http.StatusNotFound)
		return
	}
	content, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		fmt.Fprintf(w, "Error decoding cookie: %v", err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(content))

	// Clear the cookie after download
	http.SetCookie(w, &http.Cookie{
		Name:     "fileContent",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
}

func customAssetsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/styles"):]
	fullPath := filepath.Join("src", path)
	fileInfo, err := os.Stat(fullPath)

	if !os.IsNotExist(err) && !fileInfo.IsDir() {
		http.StripPrefix("/styles/", http.FileServer(http.Dir("src"))).ServeHTTP(w, r)
	} else {
		ErrorHandler(w, http.StatusNotFound, "Status Not Found", ("This Page not available"))
		return
	}
}

func ShowInfromations(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, "Status Not Found", ("This Page not available"))
		return
	}
	err := tmp1.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
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
	if text == "" {
		data := errors{
			Message: "The Input It's Empty",
		}
		tmp1.ExecuteTemplate(w, "index.html", data)
		return
	}
	for _, v := range text {
		if v < 32 || v > 126 {
			data := errors{
				Message: "the Text Uncorect",
			}
			tmp1.ExecuteTemplate(w, "index.html", data)
			return
		}
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
		datasend = data{
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

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attechment; filename=ascii_art_web.txt")
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(datasend.Text))
}
