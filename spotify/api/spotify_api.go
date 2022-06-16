package spotify_api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"log"
)

const (
	RedirectURI string = "http%3A%2F%2Fc38c-95-106-167-101.ngrok.io"
	ClientSecret string = "9bd26f301d4048f7a5e50477038a1ecd"
	ClientID string = "57de80e0b28144bc8e20c9b71a8b5553"
	Scopes string = "user-read-private%20user-read-email%20user-library-read%20playlist-read-private%20user-read-currently-playing%20user-top-read%20user-read-playback-state%20user-read-recently-played"
)

const (
	AccessTokenExpiredError int = 401
)

var (
	AuthString string = base64.StdEncoding.EncodeToString([]byte(ClientID + ":" + ClientSecret))
)

func RefreshToken(refreshToken string) (token TokenItem) {
	
	data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", refreshToken)
	rawData := []byte(data)
	body := bytes.NewReader(rawData)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic " + AuthString)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("RefreshTokenError:", err)
		return token
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	newToken := TokenItem{}
	json.Unmarshal(bytes, &newToken)

	return newToken
}

func Authorize(code string) (respData TokenItem, err error) {
	
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s/callback", code, RedirectURI)
	rawData := []byte(data)
	body := bytes.NewReader(rawData)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	
	if err != nil {
		return respData, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic " + AuthString)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return respData, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &respData)

	return respData, nil
}

func GetPlaybackHistory(token string, limit int) (playbackHistory PlaybackHistoryResponse, err error){
	URL := fmt.Sprintf("https://api.spotify.com/v1/me/player/recently-played?limit=%d", limit)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return playbackHistory, err
	} 

	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	
	if err != nil {
		return playbackHistory, err
	}
	
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("GetPlaybackHistoryReadAllError:", err)
	}
	json.Unmarshal(bytes, &playbackHistory)

	return playbackHistory, nil
}

func GetCurrentlyPlayingTrack(token string) (currentTrack CurrentlyPlayingTrackResponse, err error) {
	URL := fmt.Sprintf("https://api.spotify.com/v1/me/player/currently-playing")
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("NewRequestError:", err)
		return currentTrack, err
	}

	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	
	if err != nil {
		log.Println("Send:", err)
		return currentTrack, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return currentTrack, err
	}
	json.Unmarshal(bytes, &currentTrack)

	return currentTrack, nil
}

func GetTrackById(trackID, token string) (track Track, err error) {
	URL := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("GetTrackByIdError:", err)
		return track, err
	}

	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("GetTrackByIdError:", err)
		return track, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &track)

	return track, nil
}

func GetAlbumById(albumID, token string) (album Album, err error) {
	URL := fmt.Sprintf("https://api.spotify.com/v1/albums/%s", albumID)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("GetAlbumByIdError:", err)
		return album, err
	}

	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("GetAlbumByIdError:", err)
		return album, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &album)

	return album, nil
}

func GetArtistById(artistID, token string) (artist Artist, err error) {
	URL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", artistID)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("GetArtistByIdError:", err)
		return artist, err
	}

	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("GetArtistByIdError:", err)
		return artist, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &artist)

	return artist, nil
}