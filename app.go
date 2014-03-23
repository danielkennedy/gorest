package main

import (
	"code.google.com/p/log4go"
  "github.com/ant0ine/go-json-rest"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
	"fmt"
	"net/http"
	"os"
)

type (
  Tweets []Tweet
  Tweet struct {
    Id        bson.ObjectId `bson:"_id,omitempty"`
    Text      string
    Author    string
    Timestamp string
    Followers int64
    Retweets  int64
    Sentiment string
    Score     float64
  }
)

var (
  session *mgo.Session
  err error
  collection *mgo.Collection
	log = make(log4go.Logger)
)

func GetTweets(w *rest.ResponseWriter, r *rest.Request) {
  var (
    tweets Tweets
    err   error
  )
  if err = collection.Find(bson.M{}).All(&tweets); err != nil {
    log.Debug("%v", err)
    http.Error(w, "500 Internal Server Error", 500)
    return
  }
  w.WriteJson(&tweets)
}

func GetTweet(w *rest.ResponseWriter, r *rest.Request) {
  var tweet Tweet
  if err = collection.FindId(bson.ObjectIdHex(r.PathParam("id"))).One(&tweet); err != nil {
    rest.NotFound(w, r)
    return
  }
  w.WriteJson(&tweet)
}

func PutTweet(w *rest.ResponseWriter, r *rest.Request) {
  var tweet Tweet
  err = r.DecodeJsonPayload(&tweet)
  if err != nil {
    rest.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = collection.Insert(&tweet)
  if err != nil {
    rest.Error(w, err.Error(), http.StatusInternalServerError)
  }
  // FIXME: Fill in the ID assigned by mongo before returning to user
  w.WriteJson(&tweet)
}

func PostTweet(w *rest.ResponseWriter, r *rest.Request) {
  var (
    tweet Tweet
    changeInfo *mgo.ChangeInfo
    oid = bson.ObjectIdHex(r.PathParam("id"))
  )
  err = r.DecodeJsonPayload(&tweet)
  if err != nil {
    rest.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if changeInfo, err = collection.UpsertId(oid, &tweet); err != nil {
    rest.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Println(changeInfo)
  tweet.Id = oid
  w.WriteJson(&tweet)
}

func DeleteTweet(w *rest.ResponseWriter, r *rest.Request) {
  if err = collection.RemoveId(bson.ObjectIdHex(r.PathParam("id"))); err != nil {
    rest.NotFound(w, r)
    return
  }
  w.WriteHeader(http.StatusOK)
}

func main() {
  var (
    port string
    url string
  )

  if port = os.Getenv("VCAP_APP_PORT"); port == "" {
    port = "8080"
  }

  if url = os.Getenv("MONGO_URL"); url == "" {
    url = "mongodb://localhost:27017"
  }

	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

  session, err = mgo.Dial(url)
  if (err != nil) {
    panic(err)
  }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  collection = session.DB("tweets").C("tweets")

  handler := rest.ResourceHandler{}
  handler.SetRoutes(
    rest.Route{"GET", "/tweets", GetTweets},
    rest.Route{"PUT", "/tweets", PutTweet},
    rest.Route{"GET", "/tweets/:id", GetTweet},
    rest.Route{"POST", "/tweets/:id", PostTweet},
    rest.Route{"DELETE", "/tweets/:id", DeleteTweet},
  )
  log.Debug("Listening at port %v\n", port)
  if err := http.ListenAndServe(":"+port, &handler); err != nil {
    panic(err)
  }

/*
  err = c.Insert(&Tweet{"golang is just ok", "I Said It", 12345, 9, 11, "neutral", 0})
  if err != nil {
    panic(err)
  }

  result := Person{}
  err = c.Find(bson.M{"name": "Ale"}).One(&result)
  if err != nil {
    panic(err)
  }

  fmt.Println("Phone:", result.Phone)
*/
}

func handler(res http.ResponseWriter, req *http.Request) {
	// Dump ENV
	fmt.Fprint(res, "ENV:\n")
	env := os.Environ()
	for _, e := range env {
		fmt.Fprintln(res, e)
	}
  fmt.Fprintf(res, "Serving request for %s", req.URL.Path[1:])
}

// FIXME: Initializer to handle connection
// FIXME: Destructor to handle clean up


