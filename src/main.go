package main

import (
	"backend/src/db"
	"backend/src/handlers"
	"backend/src/utils"

	"fmt"
	"net/http"
	"os"
)

func init() {
	utils.InitEnvVariables()
}

func main() {

	// TODO:
	db.InitDB()
	defer db.DB.Close()

	mux := http.NewServeMux()
	noCorsMux := utils.EnableCORS(mux)

	// public routes
	mux.HandleFunc("/", handlers.HandleRoot)
	mux.HandleFunc("/api/register", handlers.HandleRegister)
	mux.HandleFunc("/api/login", handlers.HandleLogin)

	// private routes
	mux.Handle("/api/logout", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleLogout)))
	// mux.Handle("/api/graph", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleGraph)))
	mux.Handle("/api/getuser", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleGetUser)))
	// mux.Handle("/api/run", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleRun)))

	mux.Handle("GET /api/projects", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleGETProjects)))
	mux.Handle("POST /api/projects", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandlePOSTProject)))
	mux.Handle("DELETE /api/projects/{id}", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleDELETEProject)))

	mux.Handle("POST /api/projects/{id}/run", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleRun)))

	mux.Handle("GET /api/projects/{id}", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleGETProjectData)))
	mux.Handle("PUT /api/projects/{id}", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandlePUTProjectData)))
	mux.Handle("GET /api/projecttreeids/{id}", utils.AuthenticateRequest(http.HandlerFunc(handlers.HandleGETProjectTreeIDs)))

	PORT := os.Getenv("PORT")
	PORT = ":" + PORT
	fmt.Printf("server listening at port %s \n\n\n", PORT)
	http.ListenAndServe(PORT, noCorsMux)
}
