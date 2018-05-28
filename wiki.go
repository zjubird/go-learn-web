package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"html/template"
)

type Page struct{
	Title string
	Body  []byte
}

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

func renderTemplate(w http.ResponseWriter, templ string, p *Page){
	t, _ := template.ParseFiles(templ + ".html")
	t.Execute(w, p)
}
func viewHandler(rsp http.ResponseWriter, r *http.Request){
	//rsp is repsonse, r is request
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(rsp,"view",p)
	//t, _ := template.ParseFiles("view.html")
	//t.Execute(rsp, p)
	//fmt.Fprintf(rsp, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}


func editHandler(rsp http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil{
		p = &Page{Title: title}
	}
	renderTemplate(rsp, "edit", p)
//	t, _ := template.ParseFiles("edit.html")
//	t.Excecute(rsp, p)
}

func saveHandler(rsp http.ResponseWriter, req *http.Request){
}

func main(){
//	p1 := &Page{Title :"TestPage", Body:[]byte("This is a sample Page")}
//	p1.save()
//
//	p2, _ := loadPage("TestPage")
//
//	fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", viewHandler) //handle all requests with view to handler viewHandler
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))//listen on port 8080
	fmt.Println("...")
}
