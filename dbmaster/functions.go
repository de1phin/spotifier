package dbmaster

import (
	spapi "../spotify/api"
	"database/sql"
	"time"
	"log"
	//"fmt"
)

func (db *DatabaseWrapper) GetTopTracks(userID int, dateKeyword string, limit int) (tracks []spapi.Track) {

	userToken := db.GetToken(userID)

	var (
		dateAfter,
		dateBefore string
	)

	switch dateKeyword {
	case "today": 
		dateAfter = time.Now().Format("2006-01-02")
		dateBefore = dateAfter
	case "yesterday": 
		dateAfter = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		dateBefore = dateAfter
	case "week":
		dateAfter = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		dateBefore  = time.Now().Format("2006-01-02")
	case "month":
		dateAfter = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		dateBefore = time.Now().Format("2006-01-02")
	default:
		dateAfter = "2021-01-01"
		dateBefore = time.Now().Format("2006-01-02")
	}

	rows, err := (*sql.DB)(db).Query("SELECT track_id, SUM(time_listened) AS summa FROM TrackStatistics WHERE user_id = ? AND statistic_date >= ? AND statistic_date <= ? GROUP BY track_id, user_id ORDER BY summa DESC LIMIT ?", userID, dateAfter, dateBefore, limit)
	if err != nil {
		log.Println("TopTracksError:", err)
		return tracks
	}
	defer rows.Close()

	for rows.Next() {
		var (
			trackID string
			timeListened int64
		)

		rows.Scan(&trackID, &timeListened)
		timeListened = ((timeListened / 1e9) + 30) / 60

		track, _ := spapi.GetTrackById(trackID, userToken)
		track.TimeListened = int(timeListened)
		tracks = append(tracks, track)

	}

	return tracks
}

func (db *DatabaseWrapper) GetTopArtists(userID int, dateKeyword string, limit int) (artists []spapi.Artist) {

	userToken := db.GetToken(userID)

	var (
		dateAfter,
		dateBefore string
	)

	switch dateKeyword {
	case "today": 
		dateAfter = time.Now().Format("2006-01-02")
		dateBefore = dateAfter
	case "yesterday": 
		dateAfter = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		dateBefore = dateAfter
	case "week":
		dateAfter = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		dateBefore  = time.Now().Format("2006-01-02")
	case "month":
		dateAfter = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		dateBefore = time.Now().Format("2006-01-02")
	default:
		dateAfter = "2021-01-01"
		dateBefore = time.Now().Format("2006-01-02")
	}

	rows, err := (*sql.DB)(db).Query("SELECT artist_id, SUM(time_listened) AS summa FROM ArtistStatistics WHERE user_id = ? AND statistic_date >= ? AND statistic_date <= ? GROUP BY artist_id, user_id ORDER BY summa DESC LIMIT ?", userID, dateAfter, dateBefore, limit)
	if err != nil {
		log.Println("GetTopArtistsError:", err)
		return artists
	}
	defer rows.Close()

	for rows.Next() {
		var (
			artistID string
			timeListened int64
		)

		rows.Scan(&artistID, &timeListened)
		timeListened = ((timeListened / 1e9) + 30) / 60

		artist, _ := spapi.GetArtistById(artistID, userToken)
		artist.TimeListened = int(timeListened)
		artists = append(artists, artist)

	}

	return artists
}

func (db *DatabaseWrapper) GetTopAlbums(userID int, dateKeyword string, limit int) (albums []spapi.Album) {

	userToken := db.GetToken(userID)

	var (
		dateAfter,
		dateBefore string
	)

	switch dateKeyword {
	case "today": 
		dateAfter = time.Now().Format("2006-01-02")
		dateBefore = dateAfter
	case "yesterday": 
		dateAfter = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		dateBefore = dateAfter
	case "week":
		dateAfter = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		dateBefore  = time.Now().Format("2006-01-02")
	case "month":
		dateAfter = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		dateBefore = time.Now().Format("2006-01-02")
	default:
		dateAfter = "2021-01-01"
		dateBefore = time.Now().Format("2006-01-02")
	}

	rows, err := (*sql.DB)(db).Query("SELECT album_id, SUM(time_listened) AS summa FROM AlbumStatistics WHERE user_id = ? AND statistic_date >= ? AND statistic_date <= ? GROUP BY album_id, user_id ORDER BY summa DESC LIMIT ?", userID, dateAfter, dateBefore, limit)
	if err != nil {
		log.Println("GetTopAlbumsError:", err)
		return albums
	}
	defer rows.Close()

	for rows.Next() {
		var (
			albumID string
			timeListened int64
		)

		rows.Scan(&albumID, &timeListened)
		timeListened = ((timeListened / 1e9) + 30) / 60

		album, _ := spapi.GetAlbumById(albumID, userToken)
		album.TimeListened = int(timeListened)
		albums = append(albums, album)

	}

	return albums
}