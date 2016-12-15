package main

import (
    "github.com/AbdoulSy/adspPageTodolist"
    "github.com/AbdoulSy/adspRouterNav"
    "github.com/AbdoulSy/adspUser"
    "github.com/AbdoulSy/adspLayoutBuilder"
	"github.com/AbdoulSy/adspgoConfig"
    "github.com/AbdoulSy/adspRdf"
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
var Docss adspPageTodolist.T
var User adspUser.User
var Commits commitHistoryReader.History

func init() {
	tpl = template.Must(template.ParseGlob("views/templates/*"))
}

func index(w http.ResponseWriter, req *http.Request) {

	//index var assignments

    fileReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/",
        },
    }
    //fileReader injects content into Docss structure
    fileReader.ReadRemoteJsonBodyInItem(&Docss)

    userReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/current-user",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/current-user",
        },
    }

    userReader.ReadRemoteJsonBodyInItem(&User)

    commitsReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/commit-history",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/commit-history",
        },
    }
    commitsReader.ReadRemoteJsonBodyInItem(&Commits)

	log.Printf("%+v", User)

	pageBuilder := &adspLayoutBuilder.Builder {
		Config: adspgoConfig.Configuration("HOME"),
	}

	myPage, err := pageBuilder.Build(Docss)

	if err != nil {
		log.Println(err);
	}

	c, er := adspLayoutBuilder.BuildBasicLayoutWithPage(myPage, User, Commits)

	errtmpl := tpl.ExecuteTemplate(w, "layout", c)
	if errtmpl != nil || er != nil {
		log.Println(errtmpl)
	}
}

func projects(w http.ResponseWriter, req *http.Request) {
    userReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/current-user",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/current-user",
        },
    }

    userReader.ReadRemoteJsonBodyInItem(&User)


	pageBuilder := &adspLayoutBuilder.Builder {
		Config: adspgoConfig.Configuration("PROJECTS"),
	}
	projectPage, err := pageBuilder.Build(Docss)
	c, er := adspLayoutBuilder.BuildBasicLayoutWithPage(projectPage, User, Commits)
	err = tpl.ExecuteTemplate(w, "layout", c)
	if err != nil || er != nil {
		log.Println(err)
	}

}

func visualisation(w http.ResponseWriter, req *http.Request) {

	pageBuilder := &adspLayoutBuilder.Builder{
		Config: adspgoConfig.Configuration("VISUALISATION"),
	}
	userReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/current-user",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/current-user",
        },
    }

    userReader.ReadRemoteJsonBodyInItem(&User)

    commitsReader := adspRouterNav.RDFRouterDocumentBuilder {
        Route: "/commit-history",
        Uri: adspRdf.UriType {
            FullUri: "http://localhost:3465/commit-history",
        },
    }
    commitsReader.ReadRemoteJsonBodyInItem(&Commits)

	visualisationPage, err := pageBuilder.Build(Docss)
	c, er := adspLayoutBuilder.BuildBasicLayoutWithPage(
        visualisationPage, User, Commits)
	err = tpl.ExecuteTemplate(w, "visualisation_layout", c)
	if err != nil || er != nil {
		log.Println(err)
	}
}
