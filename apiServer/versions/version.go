package versions

import (
	"OSS/apiServer/es"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
