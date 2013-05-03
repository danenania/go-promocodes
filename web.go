package main

import (
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
  "os"
	"time"
  "fmt"
  "math/rand"
)

func main() {
  rand.Seed( time.Now().UTC().UnixNano())
	m := mux.NewRouter()

	//redeem code
	m.HandleFunc("/promocodes/{code}", RedeemCode)

	//create code
	m.HandleFunc("/pc", CreateCode).
		Queries("p", os.Getenv("PROMOPW"))

  m.HandleFunc("/pcall", ListCodes).
    Queries("p", os.Getenv("PROMOPW"))

  port := os.Getenv("PORT")

	http.Handle("/", m)
  fmt.Println("Listening on port " + port + "...")
  log.Fatal(http.ListenAndServe(":"+port, nil))
}

const CodeChars = "abcdefghijklmnopqrstuvwxyz"

func RandString(size int) string {
    buf := make([]byte, size)
    for i := 0; i < size; i++ {
        buf[i] = CodeChars[rand.Intn(len(CodeChars))]
    }
    return string(buf)
}

type PromoCode struct {
	Code     string
	Created  time.Time
	Redeemed time.Time ",omitempty"
}

var (
	mgoSession   *mgo.Session
	databaseName = "progolfscout"
)

func getMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(os.Getenv("MONGOHQ_URL"))
		if err != nil {
			log.Fatal(err)
		}
	}
	return mgoSession.Clone()
}

func getRequestVar(r *http.Request, k string) string{
	vars := mux.Vars(r)
	return vars[k]
}

func CreateCode(w http.ResponseWriter, r *http.Request) {
	session := getMongoSession()
	defer session.Close()
  code := RandString(12)
	coll := session.DB(databaseName).C("promocodes")
	err := coll.Insert(&PromoCode{code, time.Now(), time.Time{}})
  if err != nil {
    fmt.Println(err)
  }
  fmt.Fprint(w, code)
}

func ListCodes(w http.ResponseWriter, r *http.Request) {
  session := getMongoSession()
  defer session.Close()

  coll := session.DB(databaseName).C("promocodes")


  redeemed := coll.Find(bson.M{"redeemed": bson.M{"$ne" : nil}}).Sort("created").Iter()
  valid := coll.Find(bson.M{"redeemed": nil}).Sort("redeemed").Iter()

  resp := "<html><body><h1>Valid</h1>"

  var v PromoCode

  for valid.Next(&v) {
    resp +=  v.Code
    resp += " - <strong>Created</strong> " + v.Created.Format(time.RFC850)
    resp += "<br>"
  }

  resp += "<br><br>"
  resp += "<h1>Redeemed</h1>"

  for redeemed.Next(&v) {
    resp +=  v.Code
    resp += " - <strong>Created</strong> " + v.Created.Format(time.RFC850) + " | <strong>Redeemed</strong> " + v.Redeemed.Format(time.RFC850)
    resp += "<br>"
  }

  resp += "</body></html>"

  fmt.Fprint(w, resp)
}


func RedeemCode(w http.ResponseWriter, r *http.Request) {
	session := getMongoSession()
	defer session.Close()

	code := getRequestVar(r, "code")
  coll := session.DB(databaseName).C("promocodes")
  err := coll.Update(bson.M{"code": code}, bson.M{"$set" : bson.M{"redeemed": time.Now()}})
  fmt.Fprint(w, err == nil)
}
