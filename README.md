# http-post-limit
HTTP Post Limit Middleware

# Usage

```go
package main

import (
	"net/http"

	post_limit "github.com/moisespsena-go/http-post-limit"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	
	http.ListenAndServe(":8000", post_limit.New(mux)) 
	// or
	md := post_limit.New(
		mux,
		&post_limit.Opts{
		    MaxPostSize: 1024,
		    ErrorHandler: func(w http.ResponseWriter, r *http.Request) {
		    	http.Error(w, "max post size exceeded", http.StatusBadRequest)
		    },
		},
	)
	http.ListenAndServe(":8000", md)
}
```