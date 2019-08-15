package handlers

import (
	"github.com/julienschmidt/httprouter"
	"go-layouts/models"
	"go-layouts/templates"
	"net/http"
)

func HandleSessionNew(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
	next := r.URL.Query().Get("next")

	pkg.RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Next": next,
	})
}

func HandleSessionCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")

	user, err := models.FindUser(username, password)

	if err != nil {
		if models.IsValidationError(err) {
			pkg.RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"Error" : err,
				"Next": next,
			})
			return
		}
		panic(err)
	}

	session := models.FindOrCreateSession(w, r)
	session.UserID = user.ID
	err = models.GlobalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	if next == "" {
		next = "/"
	}

	http.Redirect(w, r, next+"?flash=Signed+in", http.StatusFound)
}

func HandleSessionDestroy(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
	session := models.RequestSession(r)
	if session != nil {
		err := models.GlobalSessionStore.Delete(session)

		if err != nil {
			panic(err)
		}
	}

	pkg.RenderTemplate(w, r, "sessions/destroy", nil)
}

