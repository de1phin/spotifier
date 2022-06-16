package dbmaster

import (
	"database/sql"
	"log"
)

type DatabaseWrapper sql.DB

func (dbwrapper *DatabaseWrapper) GetToken(id int) (token string) {
	db := (*sql.DB)(dbwrapper)

	rows, err := db.Query("SELECT access_token FROM Users WHERE id = ?", id)
	if err != nil {
		log.Println("GetTokenError(Query):", err)
		return ""
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("GetTokenError(Next)")
		return ""
	}
	err = rows.Scan(&token)
	if err != nil {
		log.Println("GetTokenError(Scan):", err)
		return ""
	}
	return token
}

func (dbwrapper *DatabaseWrapper) GetRefreshToken(ID int) (refreshToken string) {
	db := (*sql.DB)(dbwrapper)

	rows, err := db.Query("SELECT refresh_token FROM Users WHERE id = ?", ID)
	if err != nil {
		log.Println("GetRefreshTokenError:", err)
		return ""
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("GetRefreshTokenError(Next):", err)
		return ""
	}

	err = rows.Scan(&refreshToken)
	if err != nil {
		log.Println("GetRefreshTokenError(Scan):", err)
		return ""
	}

	return refreshToken
}

func (dbwrapper *DatabaseWrapper) UpdateToken(token, newToken string) {
	db := (*sql.DB)(dbwrapper)

	_, err := db.Query("UPDATE Users SET access_token = ? WHERE access_token = ?", newToken, token)
	if err != nil {
		log.Println("UpdateTokensError:", err)
	}
}

func (dbwrapper *DatabaseWrapper) AddTimeToTrack(userID int, trackID string, seconds int64, date string) {
	db := (*sql.DB)(dbwrapper)
	
	rows, err := db.Query("SELECT time_listened FROM TrackStatistics WHERE user_id = ? AND track_id = ? AND statistic_date = ?", userID, trackID, date)
	if err != nil {
		log.Println("AddTimeToTrackError(Select):", err)
		return 
	}
	defer rows.Close()
	if !rows.Next() {
		_, err := db.Exec("INSERT INTO TrackStatistics(user_id, track_id, time_listened, statistic_date) VALUES(?, ?, ?, ?)", userID, trackID, seconds, date)
		if err != nil {
			log.Println("AddTimeToTrackError(Insert):", err)
		}
	} else {
		var timeListened int64
		rows.Scan(&timeListened)
		_, err := db.Exec("UPDATE TrackStatistics SET time_listened = ? WHERE user_id = ? AND track_id = ? AND statistic_date = ?", timeListened + seconds, userID, trackID, date)
		if err != nil {
			log.Println("AddTimeToTrackError(Update):", err)
		}
	}
}

func (dbwrapper *DatabaseWrapper) AddTimeToArtist(userID int, artistID string, seconds int64, date string) {
	db := (*sql.DB)(dbwrapper)

	rows, err := db.Query("SELECT time_listened FROM ArtistStatistics WHERE user_id = ? AND artist_id = ? AND statistic_date = ?", userID, artistID, date)
	if err != nil {
		log.Println("AddTimeToArtistError(Select):", err)
		return 
	}
	defer rows.Close()
	if !rows.Next() {
		_, err := db.Exec("INSERT INTO ArtistStatistics(user_id, artist_id, time_listened, statistic_date) VALUES(?, ?, ?, ?)", userID, artistID, seconds, date)
		if err != nil {
			log.Println("AddTimeToArtistError(Insert):", err)
		}
	} else {
		var timeListened int64
		rows.Scan(&timeListened)
		_, err := db.Exec("UPDATE ArtistStatistics SET time_listened = ? WHERE user_id = ? AND artist_id = ? AND statistic_date = ?", timeListened + seconds, userID, artistID, date)
		if err != nil {
			log.Println("AddTimeToArtistError(Update):", err)
		}
	}
}

func (dbwrapper *DatabaseWrapper) AddTimeToAlbum(userID int, albumID string, seconds int64, date string) {
	db := (*sql.DB)(dbwrapper)

	rows, err := db.Query("SELECT time_listened FROM AlbumStatistics WHERE user_id = ? AND album_id = ? AND statistic_date = ?", userID, albumID, date)
	if err != nil {
		log.Println("AddTimeToAlbumError(Select):", err)
		return 
	}
	defer rows.Close()
	if !rows.Next() {
		_, err := db.Exec("INSERT INTO AlbumStatistics(user_id, album_id, time_listened, statistic_date) VALUES(?, ?, ?, ?)", userID, albumID, seconds, date)
		if err != nil {
			log.Println("AddTimeToAlbumError(Insert):", err)
		}
	} else {
		var timeListened int64
		rows.Scan(&timeListened)
		_, err := db.Exec("UPDATE AlbumStatistics SET time_listened = ? WHERE user_id = ? AND album_id = ? AND statistic_date = ?", timeListened + seconds, userID, albumID, date)
		if err != nil {
			log.Println("AddTimeToAlbumError(Update):", err)
		}
	}
}