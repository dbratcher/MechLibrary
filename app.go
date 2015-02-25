package mechlibrary

import (
  "net/http"
  "html/template"

  "appengine"
  "appengine/datastore"
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
  http.HandleFunc("/mechanics/list/voted", handleVotedMechanics)
  http.HandleFunc("/mechanics/list/latest", handleLatestMechanics)
  http.HandleFunc("/curate/mechanics", handleCurateMechanics)

  // Screenshots
  http.HandleFunc("/serve/", handleServe)
}

var templates = template.Must(template.ParseGlob("templates/*"))

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
  query := datastore.NewQuery("Mechanic").Ancestor(listKey(context, "approved_mechanics")).Order("-AddedTime").Limit(10)
  var mechanicArr []*Mechanic
  keys, err := query.GetAll(context, &mechanicArr)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  for i := 0; i < len(mechanicArr); i++ {
    mechanicArr[i].Id = keys[i].IntID()
  }

  templates.ExecuteTemplate(writer, "indexPage", mechanicArr)
}


