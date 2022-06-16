package main

import (
	spapi "../spotify/api"
	"../dbmaster"
	"strings"
	"../easybot"
	//"strconv"
	//"time"
	"fmt"
	"log"
	
	tgbotapi "github.com/telegram-bot-api-bot-api-5.0"
)

var (
	HandleTop,
	HandleTopTracks,
	HandleTopAlbums,
	HandleTopArtists, 
	HandlePlaybackHistory func (tgbotapi.Update)
)

var (
	HandleHelp,
	HandleStart,
	HandleCurrent func (tgbotapi.Update, ...string)
)

func InitAllHandlers() {
	HandleStart = func(u tgbotapi.Update, args ...string) {
		userID := easybot.ExtractUserID(u)
		chatID := easybot.ExtractChatID(u)


		rows, err := db.Query("SELECT id FROM Users WHERE id = ?", userID)
		if err == nil {
			defer rows.Close()
			if rows.Next() {
				setMainMenu := MainMenu.Call(chatID, "Ты уже подключен:)")
				bot.Send(setMainMenu)
				return
			}
		}

		params := fmt.Sprintf("%d", userID)
		authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s", spapi.ClientID, spapi.RedirectURI + "%2Fcallback", spapi.Scopes, params)
	
		msg := AuthLinkMenu.Prepare(authURL).Call(chatID, "Для подключения к Spotify, перейди по ссылке:")
		bot.Send(msg)
	}

	HandleHelp = func(u tgbotapi.Update, args ...string) {
		chatID := easybot.ExtractChatID(u)
		message := 
		"Список команд:\n" +
		"/start - Начать работать с ботом (предоставить доступ к данным Spotify)\n" +
		"/current - Текущий трек\n" +
		"/history - История прослушивания\n" +
		"\nДля следующих команд время измеряется с момента подключения к боту (/start)\n" +
		"/toptracks - Рейтинг треков по времени прослушивания\n" +
		"/topartists - Рейтинг исполнителей по времени прослушивания\n" +
		"/topalbums - Рейтинг альбомов по времени прослушивания"

		msg := tgbotapi.NewMessage(chatID, message)
		bot.Send(msg)
	}

	HandlePlaybackHistory = func(u tgbotapi.Update) {
		userID := easybot.ExtractUserID(u)
		chatID := easybot.ExtractChatID(u)
		
		token := (*dbmaster.DatabaseWrapper)(db).GetToken(userID)

		limit := 20
		
		playbackHistory, err := spapi.GetPlaybackHistory(token, limit)

		if err != nil {
			log.Println("GetPlaybackHistoryError:", err)
		}

		var history string

		for i, item := range playbackHistory.Items {
			var artists string
			for j, artist := range item.Track.Artists {
				if j > 0 {
					artists += ", "
				}
				artists += artist.Name
			}
			history += fmt.Sprintf("%d. %s - %s\n", i + 1, artists, item.Track.Name)
		}

		log.Println("HISTORY =", history)
		msg := tgbotapi.NewMessage(chatID, history)
		bot.Send(msg)
	}

	HandleCurrent = func(u tgbotapi.Update, args ...string) {
		userID := easybot.ExtractUserID(u)
		chatID := easybot.ExtractChatID(u)
		
		token := (*dbmaster.DatabaseWrapper)(db).GetToken(userID)
		
		currentlyPlayingTrack, err := spapi.GetCurrentlyPlayingTrack(token)
		
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewMessage(chatID, "ОшибОчька произошла у нас")
			bot.Send(msg)
			return
		}

		var artists string
		for i, artist := range currentlyPlayingTrack.Track.Artists {
			if i > 0 {
				artists += ", "
			}
			artists += artist.Name
		}

		var message string = artists + " - " + currentlyPlayingTrack.Track.Name
		if !currentlyPlayingTrack.IsPlaying {
			message += " (PAUSED)"
		}

		if artists == "" {
			message = "Никакой трек не включен"
		}
		msg := tgbotapi.NewMessage(chatID, message)
		bot.Send(msg)
	}

	HandleTopTracks = func(u tgbotapi.Update) {
		chatID := easybot.ExtractChatID(u)

		limit := 10

		menu := TopMenu.AddDaemonButtons(
			TopChooseDatePeriodButton.ChangeText("все время").MoreData(fmt.Sprintf("tracks-all-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последний месяц").MoreData(fmt.Sprintf("tracks-month-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последняя неделя").MoreData(fmt.Sprintf("tracks-week-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("вчера").MoreData(fmt.Sprintf("tracks-yesterday-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("сегодня").MoreData(fmt.Sprintf("tracks-today-%d", limit)),
		) 

		menu.ChangeLocation(easybot.InlineLocation, 1, 2, 2)
		msg := menu.Call(chatID, "Выберите период, за который хотите увидеть статистику")
		bot.Send(msg)
	}

	HandleTopAlbums = func(u tgbotapi.Update) {
		chatID := easybot.ExtractChatID(u)

		limit := 10

		menu := TopMenu.AddDaemonButtons(
			TopChooseDatePeriodButton.ChangeText("все время").MoreData(fmt.Sprintf("albums-all-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последний месяц").MoreData(fmt.Sprintf("albums-month-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последняя неделя").MoreData(fmt.Sprintf("albums-week-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("вчера").MoreData(fmt.Sprintf("albums-yesterday-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("сегодня").MoreData(fmt.Sprintf("albums-today-%d", limit)),
		) 

		menu.ChangeLocation(easybot.InlineLocation, 1, 2, 2)
		msg := menu.Call(chatID, "Выберите период, за который хотите увидеть статистику")
		bot.Send(msg)
	}

	HandleTopArtists = func(u tgbotapi.Update) {
		chatID := easybot.ExtractChatID(u)

		limit := 5

		menu := TopMenu.AddDaemonButtons(
			TopChooseDatePeriodButton.ChangeText("все время").MoreData(fmt.Sprintf("artists-all-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последний месяц").MoreData(fmt.Sprintf("artists-month-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("последняя неделя").MoreData(fmt.Sprintf("artists-week-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("вчера").MoreData(fmt.Sprintf("artists-yesterday-%d", limit)),
			TopChooseDatePeriodButton.ChangeText("сегодня").MoreData(fmt.Sprintf("artists-today-%d", limit)),
		) 

		menu.ChangeLocation(easybot.InlineLocation, 1, 2, 2)
		msg := menu.Call(chatID, "Выберите период, за который хотите увидеть статистику")
		bot.Send(msg)
	}

	HandleTop = func (u tgbotapi.Update) {
		chatID := easybot.ExtractChatID(u)
		userID := easybot.ExtractUserID(u)
		//msgID := easybot.ExtractMessageID(u)

		data := strings.Split(strings.Split(u.CallbackQuery.Data, " ")[1], "-")
		/*limit, _ := strconv.Atoi(data[2])*/

		log.Println(data)

		//var message string

		switch data[0] {
		case "tracks":
			tracks := (*dbmaster.DatabaseWrapper)(db).GetTopTracks(userID, data[1], 5)

			var imgItems []ImgItem
			for _, track := range tracks {
				var artists string
				for i, artist := range track.Artists {
					if i > 0 {
						artists += ", "
					}
					artists += artist.Name
				}
				timeListened := fmt.Sprintf("%d", track.TimeListened)
				
				minutes := "минут"
				if track.TimeListened % 10 < 5 && track.TimeListened % 10 > 1 {
					minutes += "ы"
				}
				if track.TimeListened % 10 == 1 && track.TimeListened % 100 != 11 {
					minutes = "минута"
				}
				if track.TimeListened % 100 >= 11 && track.TimeListened % 100 <= 19 {
					minutes = "минут"
				}

				item := ImgItem{ImgURL: track.Album.Images[1].URL, Title: track.Name, Subtitle: artists, Additional: timeListened + " " + minutes}
				imgItems = append(imgItems, item)
			}

			GenerateRatingImage(imgItems)
			

		case "albums":
			albums := (*dbmaster.DatabaseWrapper)(db).GetTopAlbums(userID, data[1], 5)

			var imgItems []ImgItem
			for _, album := range albums {

				var artists string
				for j, artist := range album.Artists {
					if j > 0 {
						artists += ", "
					}
					artists += artist.Name
				}

				timeListened := fmt.Sprintf("%d", album.TimeListened)
				
				minutes := "минут"
				if album.TimeListened % 10 < 5 && album.TimeListened % 10 > 1 {
					minutes += "ы"
				}
				if album.TimeListened % 10 == 1 && album.TimeListened % 100 != 11 {
					minutes = "минута"
				}
				if album.TimeListened % 100 >= 11 && album.TimeListened % 100 <= 19 {
					minutes = "минут"
				}

				item := ImgItem{ImgURL: album.Images[1].URL, Title: album.Name, Subtitle: artists, Additional: timeListened + " " + minutes}
				imgItems = append(imgItems, item)
			}

			GenerateRatingImage(imgItems)
		case "artists":
			artists := (*dbmaster.DatabaseWrapper)(db).GetTopArtists(userID, data[1], 5)

			var imgItems []ImgItem
			for _, artist := range artists {

				timeListened := fmt.Sprintf("%d", artist.TimeListened)
				
				minutes := "минут"
				if artist.TimeListened % 10 < 5 && artist.TimeListened % 10 > 1 {
					minutes += "ы"
				}
				if artist.TimeListened % 10 == 1 && artist.TimeListened % 100 != 11 {
					minutes = "минута"
				}
				if artist.TimeListened % 100 >= 11 && artist.TimeListened % 100 <= 19 {
					minutes = "минут"
				}

				item := ImgItem{ImgURL: artist.Images[1].URL, Title: artist.Name, Subtitle: "", Additional: timeListened + " " + minutes}
				imgItems = append(imgItems, item)
			}

			GenerateRatingImage(imgItems)
		}
		
		photoConfig := tgbotapi.NewPhotoUpload(chatID, ImageFolder + "stat.png")
		bot.Send(photoConfig)
		return

	}
}