package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"strings"
	"testingGorillaw/api"
	"testingGorillaw/databaseFuncs"
	"testingGorillaw/env"
)

var (
	Key   = env.Key
	Store = env.Store
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	res := fmt.Sprint("hello world")
	fmt.Fprint(w, res)
}

func testVarsHandler(w http.ResponseWriter, r *http.Request) {
	res := fmt.Sprint("test vars")
	log.Print(mux.Vars(r))
	fmt.Fprint(w, res)
}

func productCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["idProduct"]
	nameCategory := vars["nameCategory"]
	res := fmt.Sprintf("Product id = %v\nCategory name = %s", id, nameCategory)
	fmt.Fprint(w, res)
}

func MiddlewareAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[MIDDLEWARE] ", r.URL)
		log.Println(strings.Contains(r.URL.String(), "login"))
		session, _ := env.Store.Get(r, "session.id")
		if strings.Contains(r.URL.String(), "login") || strings.Contains(r.URL.String(), "reg") {
			authMatch := session.Values["authenticated"]
			if authMatch != nil && (authMatch.(bool)) {
				log.Println("1234")
				http.Redirect(w, r, "/catalog", http.StatusSeeOther)
			}
		} else if strings.Contains(r.URL.String(), "catalog") {
			authMatch := session.Values["authenticated"]
			log.Println(authMatch)
			if authMatch == nil || !(authMatch.(bool)) {
				log.Println("12345")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	render := render.New()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db := databaseFuncs.CreateDBWithDefaultConfig()
	monitors, err := databaseFuncs.GetAllMonitors(db)
	if err != nil {
		log.Println(err)
	}
	log.Println(monitors)
	//log.Println(databaseFuncs.CheckUserByLoginInDB(db, "123333"))
	//log.Println(databaseFuncs.InsertDisplay(db, 10.3, "1234", "OLED", false))
	//log.Println(databaseFuncs.InsertMonitor(db, 10, false, false, 10))
	//log.Println(databaseFuncs.CheckSignInData(db, "wf", "445"))
	//log.Print(db)
	router := mux.NewRouter()

	router.Use(MiddlewareAuth)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, http.StatusOK, "login", nil)
	}).Methods("GET")
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, http.StatusOK, "login", nil)
	}).Methods("GET")
	router.HandleFunc("/reg", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, http.StatusOK, "reg", nil)
	}).Methods("GET")
	router.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, http.StatusOK, "catalog", nil)
	}).Methods("GET")

	router.HandleFunc("/test/{nameCategory}/{idProduct}", productCategoryHandler).Methods("GET")
	router.HandleFunc("/api/login", api.SignInPost).Methods("POST")
	router.HandleFunc("/api/healthcheck", api.HealthCheck).Methods("GET")
	router.HandleFunc("/api/logout", api.LogoutHandler).Methods("GET")
	router.HandleFunc("/api/reg", api.CreateUserPost).Methods("POST")
	router.HandleFunc("/api/addMonitorWithMatrixData", api.InsertMonitorWithDataDisplay).Methods("POST")
	router.HandleFunc("/api/addMatrix", api.InsertMatrix).Methods("POST")
	router.HandleFunc("/api/addMonitor", api.InsertMonitor).Methods("POST")
	router.HandleFunc("/api/getAllMonitors", api.GetAllMonitors).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)

}
