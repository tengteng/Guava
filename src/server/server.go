package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"

	"controllers/web"
	"system"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	gojiweb "github.com/zenazn/goji/web"
)

func main() {
	filename := flag.String("config", "config.json",
		"Path to configuration file")

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	defer glog.Flush()

	var application = &system.Application{}

	application.Init(filename)
	application.LoadTemplates()
	application.ConnectToDatabase()

	// Setup static files
	static := gojiweb.New()
	static.Get("/assets/*",
		http.StripPrefix("/assets/",
			http.FileServer(
				http.Dir(
					application.Configuration.PublicPath))))

	http.Handle("/assets/", static)

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDatabase)
	goji.Use(application.ApplyAuth)

	controller := &web.Controller{}

	// Routers.

	// Couple of files - in the real world you would use nginx to serve
	// them.
	goji.Get("/robots.txt", http.FileServer(
		http.Dir(application.Configuration.PublicPath)))
	goji.Get("/favicon.ico", http.FileServer(
		http.Dir(application.Configuration.PublicPath+"/images")))

	// Home page
	goji.Get("/", application.Route(controller, "Index"))

	// Sign In routes
	goji.Get("/signin", application.Route(controller, "SignIn"))
	goji.Post("/signin", application.Route(controller, "SignInPost"))

	// Sign Up routes
	goji.Get("/signup", application.Route(controller, "SignUp"))
	goji.Post("/signup", application.Route(controller, "SignUpPost"))

	// KTHXBYE
	goji.Get("/logout", application.Route(controller, "Logout"))

	goji.Get("/create", application.Route(controller, "CreateCourse"))
	goji.Post("/create", application.Route(controller, "CreateCoursePost"))

	goji.Get("/u/:user_id", application.Route(controller, "GetUser"))

	goji.Get("/c/:short_name", application.Route(controller, "GetCourse"))

	goji.Get("/user_courses", application.Route(controller, "GetCreatedCourses"))

	goji.Get("/participated_courses", application.Route(controller, "GetParticipatedCourses"))

	goji.Get("/edit_profile", application.Route(controller, "GetUserProfile"))
	goji.Post("/edit_profile", application.Route(controller, "SaveUserProfile"))

	// TODO(tengteng): Online video page.
	// goji.Get("/v/:video_id", application.Route(controller, "OnlineVideo"))

	goji.Post("/bid", application.Route(controller, "BidCoursePost"))

	goji.Post("/add_review", application.Route(controller, "AppendReviewPost"))

	goji.Post("/participate", application.Route(controller, "ParticipatePost"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}
