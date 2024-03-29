package handlers

import (
	"go-layouts/models"
	"go-layouts/templates"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func HandleImageNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pkg.RenderTemplate(w,r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.FormValue("url")  != ""  {
		HandleImageCreateFromURL(w, r)
		return
	}

	HandleImageCreateFromFile(w, r)
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request)  {
	user := models.RequestUser(r)

	image := models.NewImage(user)

	image.Description = r.FormValue("description")

	err := image.CreateFromURL(r.FormValue("url"))

	if err != nil {
		if models.IsValidationError(err) {
			pkg.RenderTemplate(w, r, "images/new", map[string]interface{} {
				"Error": err,
				"ImageURL" : r.FormValue("url"),
				"Image": image,
			})

			return
		}

		panic(err)
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := models.RequestUser(r)
	image := models.NewImage(user)

	image.Description = r.FormValue("description")

	file, headers, err := r.FormFile("file")

	if file == nil {
		pkg.RenderTemplate(w, r, "images/new", map[string]interface{} {
			"Error": models.ErrNoImage,
			"Image": image,
		})

		return
	}

	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = image.CreateFromFile(file, headers)

	if err != nil {
		pkg.RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": err,
			"Image": image,
		})

		return
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)

}