package ragserver

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

func requestToJSON(req *http.Request, v any) error {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		return fmt.Errorf("unsupported content type: %s", mediaType)
	}
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(v)
}

func writeJSONResponse(w http.ResponseWriter, v any) error {
	js, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	return err
}
