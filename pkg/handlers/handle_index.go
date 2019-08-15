package handlers

import (
	"github.com/julienschmidt/httprouter"
	"go-layouts/models"
	"go-layouts/templates"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Home Page

	images, err := models.GlobalImageStore.FindAll(0)

	if err != nil {
		panic(err)
	}

	pkg.RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
