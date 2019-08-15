package main

import (
	"github.com/julienschmidt/httprouter"
	"go-layouts/models"
	"go-layouts/pkg"
	"go-layouts/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	router.Handle("GET", "/", handlers.HandleHome)
	router.Handle("GET", "/register", handlers.HandleUserNew)
	router.Handle("POST", "/register", handlers.HandleUserCreate)

	router.Handle("GET", "/login", handlers.HandleSessionNew)
	router.Handle("POST", "/login", handlers.HandleSessionCreate)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("../assets/"),
	)

	router.ServeFiles(
		"/im/*filepath",
		http.Dir("../assets/images"),
	)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", handlers.HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", handlers.HandleUserEdit)
	secureRouter.Handle("POST", "/account", handlers.HandleUserUpdate)

	secureRouter.Handle("GET", "/images/new", handlers.HandleImageNew)
	secureRouter.Handle("POST", "/images/new", handlers.HandleImageCreate)

	middleware := pkg.Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(models.RequireLogin))
	middleware.Add(secureRouter)

	log.Fatal(http.ListenAndServe(":3000", middleware))

	/*unauthenticatedRouter := NewRouter()
	unauthenticatedRouter.GET("/", HandleHome)
	unauthenticatedRouter.GET("/register", HandleUserNew)
	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", HandleImageNew)

	middleware := Middleware{}
	middleware.Add(unauthenticatedRouter)
	middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", middleware)*/
}

type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFound)
	router.NotFound = notFound
	return router
}
