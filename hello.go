package hello

import (
  "net/http"

  "html/template"

  "appengine"
  "appengine/user"
)

func init() {
  http.HandleFunc("/", handler)
}


var templates = template.Must(template.ParseGlob("templates/*"))

func handler(writer http.ResponseWriter, request *http.Request) {
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
  print(currentUser.Email)
  templates.ExecuteTemplate(writer, "indexPage", currentUser)
}

