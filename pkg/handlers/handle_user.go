package handlers

import (
	"github.com/julienschmidt/httprouter"
	"go-layouts/models"
	"go-layouts/templates"
	"net/http"
	"strings"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pkg.RenderTemplate(w, r, "users/new", nil)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := models.NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)

	if err != nil {
		if models.IsValidationError(err) {
			pkg.RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}

		panic(err)
	}

	session := models.NewSession(w)
	session.UserID = user.ID

	err = models.GlobalUserStore.Save(user)

	if err != nil {
		panic(err)
		return
	}

	err = models.GlobalSessionStore.Save(session)

	if err != nil {
		panic(err)
		return
	}

	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := models.RequestUser(r)

	pkg.RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
		"isMale": strings.EqualFold(*user.Gender,"male"),
		"isFemale": strings.EqualFold(*user.Gender,"female"),
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := models.RequestUser(r)

	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")
	gender := r.FormValue("gender")

	user, err := models.UpdateUser(currentUser, email, currentPassword, newPassword)
	user, err = models.UpdateExtraInfo(currentUser, gender)

	if err != nil {
		if models.IsValidationError(err) {
			pkg.RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User": user,
			})

			return
		}

		panic(err)
	}

	err = models.GlobalUserStore.Save(user)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}