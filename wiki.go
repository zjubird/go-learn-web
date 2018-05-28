package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"html/template"
	"regexp"
)

type Page struct{
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error{
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error){
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body:body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil{
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, templ string, p *Page){
	err := templates.ExecuteTemplate(w, templ + ".html", p)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(rsp http.ResponseWriter, r *http.Request, title string){
	//rsp is repsonse, r is request
	p, err := loadPage(title)
	if err != nil{
		http.Redirect(rsp, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(rsp,"view",p)
	//t, _ := template.ParseFiles("view.html")
	//t.Execute(rsp, p)
	//fmt.Fprintf(rsp, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(rsp http.ResponseWriter, r *http.Request, title string){
	p, err := loadPage(title)
	if err != nil{
		p = &Page{Title: title}
	}
	renderTemplate(rsp, "edit", p)
//	t, _ := template.ParseFiles("edit.html")
//	t.Excecute(rsp, p)
}

func saveHandler(rsp http.ResponseWriter, req *http.Request, title string){
	body := req.FormValue("body")
	p := &Page{Title:title, Body:[]byte(body)}
	err := p.save()
	if err != nil{
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rsp,req,"/view/"+title, http.StatusFound)
}

func main(){
//	p1 := &Page{Title :"TestPage", Body:[]byte("This is a sample Page")}
//	p1.save()
//
//	p2, _ := loadPage("TestPage")
//
//	fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", makeHandler(viewHandler)) //handle all requests with view to handler viewHandler
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))//listen on port 8080
	fmt.Println("...")
}
