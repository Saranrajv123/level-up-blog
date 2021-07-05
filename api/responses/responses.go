package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		fmt.Fprintf(writer, "%s", err.Error())
	}
}

func ERROR(writer http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(writer, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	JSON(writer, http.StatusBadRequest, nil)
}
