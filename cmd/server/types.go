package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/script"
	"github.com/project-safari/zebra/model"
)

func handleTypes() httprouter.Handle {
	allTypes := model.Factory()

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()

		typeReq := &struct {
			Types []string `json:"types"`
		}{Types: []string{}}

		if err := script.ReadJSON(ctx, req, typeReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		typeRes := &struct {
			Types []zebra.Type `json:"types"`
		}{Types: []zebra.Type{}}

		if len(typeReq.Types) == 0 {
			// return all types
			typeRes.Types = allTypes.Types()
		} else {
			for _, t := range typeReq.Types {
				if aType, ok := allTypes.Type(t); ok {
					typeRes.Types = append(typeRes.Types, aType)
				}
			}
		}

		script.WriteJSON(ctx, res, typeRes)
	}
}
