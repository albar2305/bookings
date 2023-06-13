package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/albar2305/bookings/pkg/config"
	"github.com/albar2305/bookings/pkg/handlers"
	"github.com/albar2305/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when production

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Sessions = session

	tc, err := render.CreteTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplates(&app)
	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	fmt.Println(fmt.Sprintf("starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
