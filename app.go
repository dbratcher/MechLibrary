package mechlibrary

import (
  "net/http"
  "html/template"

  "appengine"
)

func init() {
  // Index
  http.HandleFunc("/", handleIndex)

  // Users
  http.HandleFunc("/user/signin", handleSignIn)
  http.HandleFunc("/user/signout", handleSignOut)

  // Mechanics
  http.HandleFunc("/mechanics/id", handleMechanic)
  http.HandleFunc("/mechanics/add", handleAddMechanic)
  http.HandleFunc("/mechanics/approve", handleApproveMechanic)
  http.HandleFunc("/mechanics/upvote", handleVotedMechanics)
  http.HandleFunc("/mechanics/list/votes", handleVotesMechanics)
  http.HandleFunc("/mechanics/list/latest", handleLatestMechanics)
  http.HandleFunc("/curate/mechanics", handleCurateMechanics)

  // Android Specific
  http.HandleFunc("/android/mechanics/add", handleAndroidAdd)
  http.HandleFunc("/android/mechanics/list/latest", handleAndroidLatest)
  http.HandleFunc("/android/mechanics/list/votes", handleAndroidVotes)

  // Screenshots
  http.HandleFunc("/serve/", handleServe)
}

var templates = template.Must(template.ParseGlob("templates/*"))

type IndexPage struct {
  SortByVotes []*Mechanic
  SortByTime []*Mechanic
}

func handleIndex(writer http.ResponseWriter, request *http.Request) {

  // error checking
  if request.Method != "GET" {
    http.Error(writer, "GET requests only", http.StatusMethodNotAllowed)
    return
  }
  if request.URL.Path != "/" {
    http.NotFound(writer, request)
    return
  }

  // load page
  context := appengine.NewContext(request)

  var page IndexPage
  page.SortByTime, _ = latestMechanicQuery(context)
  if len(page.SortByTime)>5 {
    page.SortByTime = page.SortByTime[:5]
  }
  page.SortByVotes, _ = votesMechanicQuery(context)
  if len(page.SortByVotes)>5 {
    page.SortByVotes = page.SortByVotes[:5]
  }
  templates.ExecuteTemplate(writer, "indexPage", page)
}


