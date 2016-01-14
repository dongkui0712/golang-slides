package main

import (
	"io/ioutil"
	"os"
	"path"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// uid := r.URL.Query().Get("uid")
	uid := "../../../../../../etc/passwd"
	avatar := path.Join("/Users/MC/Workspace/qiniu/web-security/static", uid)

	fd, err := os.Open(avatar)
	if err != nil {
		panic(err.Error())
	}

	b, _ := ioutil.ReadAll(fd)
	println(string(b))
	// })
}
