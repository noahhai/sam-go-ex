package utils

import (
	"encoding/json"
	"net/http"
)

func WriteDataResponse(w http.ResponseWriter, data interface{}) {
	resp := respObj{
		Data: data,
	}
	if marshalled, err := json.Marshal(resp); err == nil {
		w.Write(marshalled)
	} else {
		m, _ := json.Marshal(respObj{})
		w.Write(m)
	}
}

type respObj struct {
	Data interface{}
}
