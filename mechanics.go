package mechlibrary

import (
  "net/http"

  "time"
  "strconv"

  "appengine"
  "appengine/blobstore"
  "appengine/datastore"
  "appengine/user"
)

type Mechanic struct {
  Id      int64
  Title  string
  Description string
  AddedTime time.Time
  AddedBy string
  Approved bool
  ScreenshotURL string
  Votes int
}

func listKey(context appengine.Context, book string) *datastore.Key {
  return datastore.NewKey(context, "Library", book, 0, nil)
}

func handleMechanic(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  int_id, _ := strconv.Atoi(request.URL.Query().Get("key"))
  key := datastore.NewKey(context, "Mechanic", "", int64(int_id), listKey(context, "all_mechanics"))

  mech := new(Mechanic)
  if err := datastore.Get(context, key, mech); err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }
  mech.Votes = fetchMechanicVote(context, key)
  mech.Id = int64(int_id)
  templates.ExecuteTemplate(writer, "mechanic", mech)
}

func handleVotedMechanics(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  //context.Errorf("handle vote request %v", request.URL.Query().Get("key"));
  Increment(context, request.URL.Query().Get("key"))
}

func fetchMechanicVote(context appengine.Context, key *datastore.Key) (int) {
  vote, update, _ := Count(context,  strconv.FormatInt(key.IntID(), 10))
  if update {
    mechKey := datastore.NewKey(context, "Mechanic", "", key.IntID(), listKey(context, "all_mechanics"))
    mech := new(Mechanic)
    if err := datastore.Get(context, mechKey, mech); err != nil {
      return 0
    }
    mech.Votes = vote
    datastore.Put(context, mechKey, mech)
  }
  return vote
}

func latestMechanicQuery(context appengine.Context) ([]*Mechanic, error) {
  var mechanicArr []*Mechanic
  query := datastore.NewQuery("Mechanic").Ancestor(listKey(context, "all_mechanics")).Filter("Approved=", true).Order("-AddedTime")
  keys, err := query.GetAll(context, &mechanicArr)
  if err != nil {
    return nil, err
  }

  for i := 0; i < len(mechanicArr); i++ {
    if len(mechanicArr[i].Description) > 69 {
      mechanicArr[i].Description = mechanicArr[i].Description[:69] + "..."
    }
    mechanicArr[i].Id = keys[i].IntID()
    mechanicArr[i].Votes = fetchMechanicVote(context, keys[i])
  }
  return mechanicArr, nil
}


func votesMechanicQuery(context appengine.Context) ([]*Mechanic, error) {
  var mechanicArr []*Mechanic
  query := datastore.NewQuery("Mechanic").Ancestor(listKey(context, "all_mechanics")).Filter("Approved=", true).Order("-Votes")
  keys, err := query.GetAll(context, &mechanicArr)
  if err != nil {
    return nil, err
  }

  for i := 0; i < len(mechanicArr); i++ {
    if len(mechanicArr[i].Description) > 69 {
      mechanicArr[i].Description = mechanicArr[i].Description[:69] + "..."
    }
    mechanicArr[i].Id = keys[i].IntID()
    mechanicArr[i].Votes = fetchMechanicVote(context, keys[i])
  }
  return mechanicArr, nil
}

func handleLatestMechanics(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)

  mechanicArr, _ := latestMechanicQuery(context)

  templates.ExecuteTemplate(writer, "mechanicLatestList", mechanicArr)
}

func handleVotesMechanics(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)

  mechanicArr, _ := votesMechanicQuery(context)

  templates.ExecuteTemplate(writer, "mechanicLatestList", mechanicArr)
}

func handleApproveMechanic(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  int_id, _ := strconv.Atoi(request.URL.Query().Get("id"))
  key := datastore.NewKey(context, "Mechanic", "", int64(int_id), listKey(context, "all_mechanics"))

  mech := new(Mechanic)
  if err := datastore.Get(context, key, mech); err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  mech.Approved = true

  newKey := datastore.NewKey(context, "Mechanic", "", int64(int_id), listKey(context, "all_mechanics"))

  datastore.Put(context, key, mech)
  newKey, err := datastore.Put(context, newKey, mech)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(writer, request, "/mechanics/list/latest", http.StatusFound)
}

func handleCurateMechanics(writer http.ResponseWriter, request *http.Request) {
   context := appengine.NewContext(request)

   query := datastore.NewQuery("Mechanic").Ancestor(listKey(context, "all_mechanics")).Order("-AddedTime")
  var mechanicArr []*Mechanic
  keys, err := query.GetAll(context, &mechanicArr)
  if err != nil {
    http.Error(writer, err.Error(), http.StatusInternalServerError)
    return
  }

  for i := 0; i < len(mechanicArr); i++ {
    mechanicArr[i].Id = keys[i].IntID()
    mechanicArr[i].Votes = fetchMechanicVote(context, keys[i])
  }

  templates.ExecuteTemplate(writer, "mechanicCurateList", mechanicArr)

}

func handleAddMechanic(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  if request.Method == "POST" {
    storeMechanic(context, writer, request)
  } else if request.Method == "GET" {
    uploadURL, err := blobstore.UploadURL(context, "/mechanics/add", nil)
    if err != nil {
      http.Error(writer, err.Error(), http.StatusInternalServerError)
      return
    }
    templates.ExecuteTemplate(writer, "addMechanic", uploadURL)
  } else {
    http.Error(writer, "Must be a POST or GET request", http.StatusMethodNotAllowed)
  }
}

func storeMechanic(context appengine.Context, writer http.ResponseWriter, request *http.Request) {
  context.Infof("Made it here URL: %v", request.URL)
  blobs, vals, _ := blobstore.ParseUpload(request)
  file := blobs["screenshot"]
  mechanic := &Mechanic{
    Title: vals.Get("title"),
    Description: vals.Get("description"),
    AddedTime: time.Now(),
    ScreenshotURL: "/serve/?blobKey="+string(file[0].BlobKey),
  }
  context.Infof("Created mechanic: %v", vals.Get("title"))
  currentUser := user.Current(context)
  if currentUser != nil {
    mechanic.AddedBy = currentUser.String()
  }
  mechanic.Approved = false
  key := datastore.NewIncompleteKey(context, "Mechanic", listKey(context, "all_mechanics"))
  newKey, _ := datastore.Put(context, key, mechanic)
  context.Infof("Redirecting: %v", strconv.FormatInt(newKey.IntID(), 10))
  http.Redirect(writer, request, "/mechanics/id?key="+strconv.FormatInt(newKey.IntID(), 10), http.StatusFound)
}

