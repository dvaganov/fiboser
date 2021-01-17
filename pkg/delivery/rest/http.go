package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type (
	requestHandelr func(r *http.Request) (interface{}, error)
)

func NewHttpRequestHandler(handler requestHandelr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		res, err := handler(r)
		if err, ok := err.(httpError); ok {
			res, _ := json.Marshal(err)
			w.WriteHeader(err.Status)
			w.Write(res)
			return
		}

		if err != nil {
			err := NewInternalError(err)
			log.Println(err.Unwrap())

			res, _ := json.Marshal(err)
			w.WriteHeader(err.Status)
			w.Write(res)
			return
		}

		jsonRes, _ := json.Marshal(res)
		w.Write(jsonRes)
	}
}
