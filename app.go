package mechlibrary

import (
  "net/http"

  "html/template"

  "time"

  "appengine"
  "appengine/datastore"
  "appengine/user"
)

func init() {
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/mechanics/add", handleAddMechanic)
}

type Mechanic struct {
  Title  string
  Description string
  AddedTime time.Time
  AddedBy string
}

var templates = template.Must(template.ParseGlob("templates/*"))

func handleIndex(writer http.ResponseWriter, request *http.Request) {
  if request.Method != "GET" {
    http.Error(writer, "GET requests only", http.StatusMethodNotAllowed)
    return
  }
  if request.URL.Path != "/" {
    http.NotFound(writer, request)
    return
  }

  context := appengine.NewContext(request)
  currentUser := user.Current(context)
  if currentUser == nil {
    url, err := user.LoginURL(context, request.URL.String())
    if err != nil {
      http.Error(writer, err.Error(), http.StatusInternalServerError)
      return
    }
    writer.Header().Set("Location", url)
    writer.WriteHeader(http.StatusFound)
    return
  }


  templates.ExecuteTemplate(writer, "indexPage", currentUser)

}

 func mechanicKey(context appengine.Context) *datastore.Key {
   return datastore.NewKey(context, "Mechanic", "all_mechanics", 0, nil)
 }

func storeMechanic(context appengine.Context, writer http.ResponseWriter, request *http.Request) {
   mechanic := &Mechanic{
    Title: request.FormValue("title"),
    Description: request.FormValue("description"),
    AddedTime: time.Now(),
  }
  currentUser := user.Current(context)
  if currentUser != nil {
    mechanic.AddedBy = currentUser.String()
  }
  key := datastore.NewIncompleteKey(context, "Mechanic", mechanicKey(context))
  _, err := datastore.Put(context, key, mechanic)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(writer, request, "/mechanics/add", http.StatusSeeOther)
}

func showMechanicForm(context appengine.Context, writer http.ResponseWriter, request *http.Request) {
 query := datastore.NewQuery("Mechanic").Ancestor(mechanicKey(context)).Order("-AddedTime").Limit(10)
  var mechanicArr []*Mechanic
  _, err := query.GetAll(context, &mechanicArr)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  templates.ExecuteTemplate(writer, "addMechanic", mechanicArr)
}

func handleAddMechanic(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  if request.Method == "POST" {
    storeMechanic(context, writer, request)
  } else if request.Method == "GET" {
    showMechanicForm(context, writer, request)
  } else {
    http.Error(writer, "Must be a POST or GET request", http.StatusMethodNotAllowed)
  }
}
