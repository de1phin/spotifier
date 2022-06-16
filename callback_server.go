package main

import (
	dbmaster "./dbmaster"
	"encoding/json"
	spapi "./spotify/api"
	"net/http"
	"net/url"
	"time"
	"log"

	"github.com/gorilla/mux"
)

func main() {

	db, err := dbmaster.DBConnect("Spotifier", "root", "root")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {

		val, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("ParseQueryError:", err)
			return
		}
		code := val.Get("code") 
		if code == "" {
			log.Println("val.Get Error:", err)
			return
		}
		tgid := val.Get("state")
		
		resp, err :=  spapi.Authorize(code)

		if err != nil {
			log.Println("ERROR:", err)
			return
		}

		userReq, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
		if err != nil {
			log.Println(err)
			return
		}
		userReq.Header.Add("Authorization", "Bearer " + resp.AccessToken)

		userResp, err := http.DefaultClient.Do(userReq)
		if err != nil {
			log.Println("userRESP ERROR:", err)
			return
		}

		var userRes map[string]interface{}
		json.NewDecoder(userResp.Body).Decode(&userRes)
		log.Print(userRes);

		expiresAfter := time.Now().Add(time.Second * time.Duration(resp.ExpiresIn - 60))
		
		_, err = db.Exec("INSERT INTO Users(id, spotify_id, name, access_token, refresh_token, expires_after, last_update) VALUES(?, ?, ?, ?, ?, ?, ?)", tgid, userRes["id"], userRes["display_name"], resp.AccessToken, resp.RefreshToken, expiresAfter, 0)
		if err != nil {
			log.Println("DB Error:", err)
		}

		http.Redirect(w, r, "https://t.me/De1phin_bot?start=", 302)
	})

	
	server := http.Server{
		Handler:router,
		Addr:"localhost:1631",
	}

	log.Fatal(server.ListenAndServe())
}