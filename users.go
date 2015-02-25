package mechlibrary

import (
  "net/http"

  "appengine"
  "appengine/user"
)

func handleSignIn(writer http.ResponseWriter, request *http.Request) {
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
  writer.Header().Set("Location", "/")
}

func handleSignOut(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  currentUser := user.Current(context)
  if currentUser != nil {
    url, err := user.LogoutURL(context, request.URL.String())
    if err != nil {
      http.Error(writer, err.Error(), http.StatusInternalServerError)
      return
    }
    writer.Header().Set("Location", url)
    writer.WriteHeader(http.StatusFound)
    return
  }
  writer.Header().Set("Location", "/")
}
