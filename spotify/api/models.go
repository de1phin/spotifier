package spotify_api

type Image struct {
	URL				string 			`json:"url"`
	Height			int				`json:"height"`
	Width			int				`json:"width"`
}

type Artist struct {
	ID 				string			`json:"id"`
	Name			string			`json:"name"`
	Images			[]Image			`json:"images"`
	TimeListened	int
}

type Track struct {
	ID 				string			`json:"id"`
	Name 			string			`json:"name"`
	Duration		int				`json:"duration_ms"`
	Artists			[]Artist		`json:"artists"`
	Album			Album			`json:"album"`
	TimeListened	int
}

type Item struct {
	Track			Track			`json:"track"` 
}

type Album struct {
	ID				string			`json:"id"`
	Name			string			`json:"name"`
	AlbumType		string			`json:"album_type"`
	Artists 		[]Artist		`json:"artists"`
	Images			[]Image			`json:"images"`
	TimeListened 	int	
}

type PlaybackHistoryResponse struct {
	Items			[]Item			`json:"items"`
}

type CurrentlyPlayingTrackResponse struct {
	Track 			Track			`json:"item"`
	IsPlaying		bool			`json:"is_playing"`
}

type SpotifyError struct {
	Status 			int				`json:"status"`
	Message			string			`json:"message"`
}

type SpotifyErrorItem struct {
	Error			SpotifyError	`json:"error"`
}

type TokenItem struct {
	AccessToken		string			`json:"access_token"`
	TokenType		string			`json:"token_type"`
	RefreshToken	string			`json:"refresh_token"`
	ExpiresIn		int				`json:"expires_in"`
}