package main

import (
    "github.com/AbdoulSy/content"
    "github.com/AbdoulSy/userDescriptor"
    "github.com/AbdoulSy/codeDescriptor"
    "github.com/AbdoulSy/pageBuilder"
	"github.com/AbdoulSy/adspgo"
	"github.com/AbdoulSy/layout"
  	"github.com/AbdoulSy/commitHistoryReader"
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/projects", projects)
	mux.HandleFunc("/visualisation", visualisation)
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	log.Println("ADSP Web server listening on the port 8080")
	http.ListenAndServe(":8080", mux)
}

var tpl *template.Template

//Docss Holds the Content of the JSON String
var Docss content.T

var User userDescriptor.User

var Commits commitHistoryReader.History

func init() {
	tpl = template.Must(template.ParseGlob("views/templates/*"))
}

func index(w http.ResponseWriter, req *http.Request) {

	//indev var assignments
    docName := "FileWalkAndDescription"
    host := "http://localhost:3465"
    fileReader := &codeDescriptor.T {
    	Name: docName,
    	Host: host,
    }

    userReader := &userDescriptor.T {
    	Name: "CurrentUser",
    	Host: "http://localhost:3465/current-user",
    }

    commitsReader := &commitHistoryReader.T {
    	Name: "Last 9 Commits",
    	Host: "http://localhost:3465/commit-history",
    }

    //fileReader injects content into Docss structure
    fileReader.GetBodyAsTextSync(&Docss)

    userReader.GetBodyAsTextSync(&User)

    commitsReader.GetBodyAsTextSync(&Commits)

	log.Printf("%+v", User)

	pageBuilder := &pageBuilder.T{
		Config: adspgo.Configuration("HOME"),
	}

	myPage, err := pageBuilder.Build(Docss)

	if err != nil {
		log.Println(err);
	}

	c, er := layout.BuildBasicLayoutWithPage(myPage, User, Commits)

	errtmpl := tpl.ExecuteTemplate(w, "layout", c)
	if errtmpl != nil || er != nil {
		log.Println(errtmpl)
	}
}

func projects(w http.ResponseWriter, req *http.Request) {


    userReader := &userDescriptor.T {
    	Name: "CurrentUser",
    	Host: "http://localhost:3465/current-user",
    }
    userReader.GetBodyAsTextSync(&User)


	pageBuilder := &pageBuilder.T{
		Config: adspgo.Configuration("PROJECTS"),
	}
	projectPage, err := pageBuilder.Build(Docss)
	c, er := layout.BuildBasicLayoutWithPage(projectPage, User, Commits)
	err = tpl.ExecuteTemplate(w, "layout", c)
	if err != nil || er != nil {
		log.Println(err)
	}

}

func visualisation(w http.ResponseWriter, req *http.Request) {

	pageBuilder := &pageBuilder.T{
		Config: adspgo.Configuration("VISUALISATION"),
	}
	userReader := &userDescriptor.T {
    	Name: "CurrentUser",
    	Host: "http://localhost:3465/current-user",
    }
    commitsReader := &commitHistoryReader.T {
        Name: "Last 9 Commits",
        Host: "http://localhost:3465/commit-history",
    }

    commitsReader.GetBodyAsTextSync(&Commits)


	userReader.GetBodyAsTextSync(&User)

	visualisationPage, err := pageBuilder.Build(Docss)
	c, er := layout.BuildBasicLayoutWithPage(visualisationPage, User, Commits)
	err = tpl.ExecuteTemplate(w, "visualisation_layout", c)
	if err != nil || er != nil {
		log.Println(err)
	}
}
