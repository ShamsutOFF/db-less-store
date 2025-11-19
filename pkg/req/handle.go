package req

import (
	"db-less-store/pkg/res"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.SendJsonResponse(*w, err.Error(), 400)
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		res.SendJsonResponse(*w, err.Error(), 400)
		return nil, err
	}
	return &body, err
}
