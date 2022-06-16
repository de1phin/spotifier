package main

import (
	spapi "../spotify/api"
	"database/sql"
	"../dbmaster"
	"time"
	"log"
)

var (
	db *sql.DB
)

//time in seconds before next update
const UpdateFrequency int = 1

type ServedUser struct {
	ID				int
	SpotifyID		string
	Token			string
	ExpiresAfter	time.Time
	LastUpdate		int64
}

func ServeUser(user ServedUser) { 

	if time.Now().After(user.ExpiresAfter) {
		refreshToken := (*dbmaster.DatabaseWrapper)(db).GetRefreshToken(user.ID)
		newToken := spapi.RefreshToken(refreshToken)
		expiresAfter := time.Now().Add(time.Second * time.Duration(newToken.ExpiresIn - 60))
		db.Exec("UPDATE Users SET access_token = ?, expires_after = ? WHERE id = ?", newToken.AccessToken, expiresAfter, user.ID)
		log.Println("Refreshed Token for", user.ID)
	}
	
	currentTrack, err := spapi.GetCurrentlyPlayingTrack(user.Token)
	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now().UnixNano()
	_, err = db.Exec("UPDATE Users SET last_update = ? WHERE id = ?", now, user.ID)
	
	timeElapsed := now - user.LastUpdate


	if currentTrack.Track.ID != "" && currentTrack.IsPlaying {
		date := time.Now().Format("2006-01-02")
		(*dbmaster.DatabaseWrapper)(db).AddTimeToTrack(user.ID, currentTrack.Track.ID, timeElapsed, date)
		(*dbmaster.DatabaseWrapper)(db).AddTimeToAlbum(user.ID, currentTrack.Track.Album.ID, timeElapsed, date)
		for _, artist := range currentTrack.Track.Artists {
			(*dbmaster.DatabaseWrapper)(db).AddTimeToArtist(user.ID, artist.ID, timeElapsed, date)
		}
	}

}

func main() {

	var err error
	db, err = dbmaster.DBConnect("Spotifier", "root", "root")
	if err != nil {
		panic(err)
	}

	var user ServedUser

	db.Exec("UPDATE Users SET last_update = ?", time.Now().UnixNano())

	for {

		time.Sleep(time.Second)

		rows, err := db.Query("SELECT id, spotify_id, access_token, expires_after, last_update FROM Users")
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var (
				expiresAfterStr string
			)

			rows.Scan(&user.ID, &user.SpotifyID, &user.Token, &expiresAfterStr, &user.LastUpdate)

			user.ExpiresAfter, _ = time.Parse(dbmaster.TimeStampFormat, expiresAfterStr)

			ServeUser(user)
		}
		rows.Close()
	}
	
}