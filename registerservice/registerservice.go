package registerservice

import (
	"net/http"
	"fmt"
	"log"
	"base-auth/dbservice"
)

var err int

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	repeatPassword := r.FormValue("repeat-password")
	email := r.FormValue("email")
	repeatEmail := r.FormValue("repeat-email")
	phoneNumber := r.FormValue("phone")
	
	//@TODO check password-match
  //@TODO check email-match
  //@TODO validate input types
	
	log.Println(username, password, repeatPassword, email, repeatEmail, phoneNumber)
	exist, err := dbservice.CheckCredentials(username, password)
	log.Println( dbservice.CheckCredentials, "valid")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	if exist {
	 	fmt.Fprint(w, "username exist in database, provide another username!")
		http.Error(w, "username exist in database", http.StatusConflict)
  } else {
 		dbservice.InsertUser(username, password)
	}
}		

func ValidateFormData (){

}