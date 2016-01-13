package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/qiniu/log"
	"labix.org/v2/mgo"
)

const startIdx = 1452652841

type Config struct {
	Addr    string `json:"addr"`
	MgoAddr string `json:"mongo_addr"`
	MgoDB   string `json:"mongo_db"`
	MgoColl string `json:"mongo_coll"`
}

func Run(conf Config) error {
	session, err := mgo.Dial(conf.MgoAddr)
	if err != nil {
		return err
	}
	coll := session.DB(conf.MgoDB).C(conf.MgoColl)

	ist, err := NewInstance(coll)
	if err != nil {
		return err
	}
	return ist.Run(conf.Addr)
}

//-----------------------------------------------------------------------------

func NewInstance(coll *mgo.Collection) (*Instance, error) {
	ist := &Instance{coll, 0}
	if err := ist.pickIdx(); err != nil {
		return nil, err
	}
	return ist, nil
}

type Entry struct {
	URL string `bson:"url"`
}

type M map[string]interface{}

type Instance struct {
	coll *mgo.Collection
	idx  int64
}

func (p *Instance) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/short", p.HandleShort)
	mux.HandleFunc("/", p.HandleRedirect)
	log.Info("running at:", addr)
	return http.ListenAndServe(addr, mux)
}

func (p *Instance) HandleShort(resp http.ResponseWriter, req *http.Request) {

	URL := req.FormValue("url")

	if URL == "" {
		replyErr(resp, "need url arg")
		return
	}

	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = "http://" + URL
	}

	short, err := p.Short(URL)
	if err != nil {
		replyErr(resp, err.Error())
		return
	}

	resp.Write([]byte(short))
	return
}

func (p *Instance) HandleRedirect(resp http.ResponseWriter, req *http.Request) {

	short := strings.TrimPrefix(req.URL.Path, "/")
	URL, err := p.Redirect(short)
	if err != nil {
		replyErr(resp, err.Error())
		return
	}
	resp.Header().Add("Location", URL)
	resp.WriteHeader(301)
	return
}

func replyErr(resp http.ResponseWriter, err string) {
	msg := fmt.Sprintf("err: %s\n", err)
	resp.WriteHeader(500)
	resp.Write([]byte(msg))
}

func (p *Instance) pickIdx() (err error) {

	p.idx = startIdx

	var doc map[string]interface{}
	if err = p.coll.Find(nil).Sort("-_id").Limit(1).One(&doc); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}

	n, err := strconv.ParseInt(doc["_id"].(string), 36, 64)
	if err != nil {
		return
	}
	p.idx = n
	return
}

func (p *Instance) Short(URL string) (string, error) {

	URL = strings.ToLower(URL)
	id := atomic.AddInt64(&p.idx, 1)
	idstr := strconv.FormatInt(id, 36)
	doc := &Entry{
		URL: URL,
	}

	if _, err := p.coll.UpsertId(idstr, doc); err != nil {
		return "", err
	}

	return idstr, nil
}

func (p *Instance) Redirect(short string) (URL string, err error) {

	var entry Entry
	if err = p.coll.FindId(short).One(&entry); err != nil {
		return
	}
	return entry.URL, nil
}
