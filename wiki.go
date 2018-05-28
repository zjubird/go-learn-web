package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
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

func viewHandler(rsp http.ResponseWriter, r *http.Request){
	//rsp is repsonse, r is request
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(rsp, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}


func main(){
//	p1 := &Page{Title :"TestPage", Body:[]byte("This is a sample Page")}
//	p1.save()
//
//	p2, _ := loadPage("TestPage")
//
//	fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", viewHandler) //handle all requests with view to handler viewHandler
	log.Fatal(http.ListenAndServe(":8080", nil))//listen on port 8080
}
