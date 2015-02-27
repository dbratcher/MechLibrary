package mechlibrary

import (
  "net/http"

  "fmt"

  "appengine"
  "appengine/blobstore"

  "encoding/json"
)

func handleAndroidAdd(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  context.Infof("Made it here URL: %v", request.URL)
  uploadURL, err := blobstore.UploadURL(context, "/mechanics/add", nil)
  if err != nil {
    context.Infof("error %v", err.Error())
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }
  context.Infof("redirect URL: %v", uploadURL.RequestURI())
  fmt.Fprintf(writer, "%v", uploadURL.RequestURI())
}



func handleAndroidLatest(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)

  mechanicArr, _ := latestMechanicQuery(context)

  js, _ := json.Marshal(mechanicArr)
  writer.Header().Set("Content-Type", "application/json")
  writer.Write(js)
}


func handleAndroidVotes(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)

  mechanicArr, _ := votesMechanicQuery(context)

  js, _ := json.Marshal(mechanicArr)
  writer.Header().Set("Content-Type", "application/json")
  writer.Write(js)
}
