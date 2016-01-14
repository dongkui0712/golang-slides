package main

import (
	"os/exec"
	"path"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// uid := r.URL.Query().Get("uid")
	uid := "../../../../../../usr/local/var/"
	avatar := path.Join("/Users/MC/Workspace/qiniu/web-security/static", uid)

	result, err := exec.Command("ls", "-l", avatar).Output()
	if err != nil {
		panic(err.Error())
	}
	println(string(result))
	// })
}
