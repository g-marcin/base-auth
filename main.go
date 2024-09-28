package main

import (
  "fmt"
  "net/http"
  "base-auth/dbservice"
  "base-auth/registerservice"
  "log"
)

func rootHandler(w http.ResponseWriter, r *http.Request){
  fmt.Fprint(w, "root handler")
}

func htmxActionHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello from HTMX!")
}

func serveHtmx (w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/index.html")
}

func handlePrivateRoute (w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "private route!")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		middlewareHandler(w, r, next)
	})
}

func middlewareHandler(w http.ResponseWriter, r *http.Request, next http.Handler){
	cookie, err := r.Cookie("auth-token")
	if err != nil || cookie.Value != "secret-token" {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
  }
	next.ServeHTTP(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	valid, err := dbservice.CheckCredentials(username, password)
	log.Println(valid, "valid")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	if valid {
	  http.SetCookie(w, &http.Cookie{
        Name:  "auth-token",
        Value: "secret-token",
        Path:  "/",
				HttpOnly: true,
				Secure: true,
				SameSite: http.SameSiteStrictMode,
				
    })
		fmt.Fprint(w, "Login sucessful, access to /private granted")
  } else {
  	fmt.Fprint(w, "please provide valid credentials!")
		http.Error(w, "invalid credetials", http.StatusUnauthorized)
	}
}		

func main(){
	 fs := http.FileServer(http.Dir("./static"))
   http.Handle("/", fs)

	
	http.HandleFunc("/htmx-action", htmxActionHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", registerservice.HandleRegister)
	
	privateRoute := http.HandlerFunc(handlePrivateRoute)	
	http.Handle("/private", AuthMiddleware(privateRoute))

  fmt.Println("Starting server on port :8080")
  if err := http.ListenAndServe(":8080", nil); err!= nil {
    fmt.Println("Error starting server:", err)
  }
}
