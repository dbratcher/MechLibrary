package mechlibrary

import (
  "net/http"

  "appengine"
  "appengine/blobstore"
)

func handleServe(w http.ResponseWriter, r *http.Request) {
  blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}
