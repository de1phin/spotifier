package main

import (
	"../easybot"
)


var (
	TopChooseDatePeriodButton *easybot.DaemonButton
)

func InitDaemonDataButtons() {
	var button int
	
	button = easybot.NewDaemonDynamicDataButton("сегодня", "toptracks", HandleTop)
	TopChooseDatePeriodButton = &easybot.DaemonStorage[button]
}

var (
	MainMenu,
	AuthLinkMenu,
	TopMenu *easybot.Menu	
)

func InitMenu() {

	InitDaemonDataButtons()

	var menu int
	menu = easybot.NewMenu("%s", 
		easybot.NewDynamicUrlButton("ссылка"),
	)
	AuthLinkMenu = &easybot.Lobby[menu]
	AuthLinkMenu.ChangeLocation(easybot.InlineLocation, 1)

	menu = easybot.NewMenu("%s")
	TopMenu = &easybot.Lobby[menu]

	menu = easybot.NewMenu("%s",
		easybot.NewTextButton("История прослушиваний", HandlePlaybackHistory),
		easybot.NewTextButton("Рейтинг треков", HandleTopTracks),
		easybot.NewTextButton("Рейтинг исполнителей", HandleTopArtists),
		easybot.NewTextButton("Рейтинг альбомов", HandleTopAlbums),
	)
	MainMenu = &easybot.Lobby[menu]
	MainMenu.ChangeLocation(easybot.SimpleLocation, 2, 2)

}

func InitAllInstances() {
	easybot.NewCommand("/help", HandleHelp)
	easybot.NewCommand("/start", HandleStart)
	//easybot.NewCommand("/history", HandlePlaybackHistory)
	easybot.NewCommand("/current", HandleCurrent)
	//easybot.NewCommand("/toptracks", HandleTopTracks)
	//easybot.NewCommand("/topsongs", HandleTopTracks)
	//easybot.NewCommand("/topartists", HandleTopArtists)
	//easybot.NewCommand("/topalbums", HandleTopAlbums)

	InitMenu()
}