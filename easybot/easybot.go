package easybot

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"sort"
	"unicode/utf16"
	
	"github.com/telegram-bot-api-bot-api-5.0"
)

/////////////
//VARIABLES//
/////////////

var (
	SEPARATOR = " "
)

//////////
//Errors//
//////////

var (
	ErrorInvalidPattern	= errors.New("easybot: Bad Pattern")
)

//////////////
//Extractors//
//////////////

var ExtractMode func(tgbotapi.Update)(int) = func(tgbotapi.Update)(int) {
	return 0
}

func ExtractParams(update tgbotapi.Update) ([]string) {
	if update.CallbackQuery != nil {
		return strings.Split(update.CallbackQuery.Data, " ")[1:]
	} else {
		return []string{}
	}
}

func ExtractSourceID(update tgbotapi.Update) (int64) {
	if update.Message != nil {
		if update.Message.ForwardFromChat != nil {
			return update.Message.ForwardFromChat.ID
		}
	} 

	return 0
}

func ExtractMediaGroupID(update tgbotapi.Update) (string) {
	if update.Message != nil {
		return update.Message.MediaGroupID
	}
	return ""
}

func ExtractForwardMessageID(update tgbotapi.Update) (int) {
	if update.Message != nil {
		return update.Message.ForwardFromMessageID
	} 
	return 0
}

func ExtractIsForward(update tgbotapi.Update) (bool) {
	return update.Message != nil && (update.Message.ForwardFromChat != nil || update.Message.ForwardFrom != nil)
}

func ExtractChatType(update tgbotapi.Update) (string) {
	if update.ChannelPost != nil {
		if update.ChannelPost.Chat != nil {
			return update.ChannelPost.Chat.Type
		}
	} else if update.CallbackQuery != nil {
		if update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
			return update.CallbackQuery.Message.Chat.Type
		} else {
			return ""
		}
	} else {
		if update.Message != nil && update.Message.Chat != nil {
			return update.Message.Chat.Type
		} else {
			return ""
		}
	}
	return ""
}

func ExtractUserID(update tgbotapi.Update) (int) {
	if update.CallbackQuery != nil {
		return int(update.CallbackQuery.Message.Chat.ID)
	} else {
		return update.Message.From.ID;
	}
}

func ExtractChatID(update tgbotapi.Update) (int64) {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	} else {
		return update.Message.Chat.ID;
	}
}

func ExtractUserName(update tgbotapi.Update) (string) {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.From.UserName
	} else {
		return update.Message.From.UserName
	}
}

func ExtractMessageID(update tgbotapi.Update) (int) {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.MessageID
	} else {
		return update.Message.MessageID
	}
}

func ExtractPhotoID(update tgbotapi.Update) (string) {
	if update.CallbackQuery != nil {
		return ""
	} else {
		if update.Message.Photo == nil {
			return ""
		}

		imgs := update.Message.Photo
		if len(imgs) == 0 {
			return ""
		} else {
			return imgs[0].FileID
		}
	}

}

func ExtractText(update tgbotapi.Update) (string) {
	if update.CallbackQuery != nil {
		return ""
	} else if update.Message.Caption != "" {
		return update.Message.Caption
	} else {
		return update.Message.Text
	}

}

func DefineModeExtractor(f func(tgbotapi.Update)(int)) {
	ExtractMode = f
}

/////////
//Input//
/////////

type InputHandler struct {
	ID 			int
	Mode 		int
	Handler 	func(tgbotapi.Update) (bool)
} 

var (
	InManager 	[]InputHandler
)

func NewInput(mode int, handler func(tgbotapi.Update)(bool)) (int) {
	input := InputHandler{
		ID: len(InManager),
		Mode: mode,
		Handler: handler,
	}

	InManager = append(InManager, input)

	return input.ID
}

///////////
//Command//
///////////

type Command struct {
	ID 			int
	Text 		string
	Handler 	func(tgbotapi.Update, ...string)
}

func NewCommand(text string, handler func(tgbotapi.Update, ...string)) (int) {
	command := Command{
		ID: len(Commands),
		Text: text,
		Handler: handler,
	}

	Commands = append(Commands, command)

	return command.ID
}

//ERASE COMMAND PARAMS//

func EraseParameters(update *tgbotapi.Update) {
	splt := strings.Split(update.Message.Text, " ")
	if len(splt) == 0 {
		return
	}
	update.Message.Text = splt[0]
}

//Get Command Name//

func GetCommand(update tgbotapi.Update) string {
	if update.Message == nil {
		return ""
	}
	
	splt := strings.Split(update.Message.Text, " ")
	if len(splt) == 0 || splt[0][0] != '/' {
		return ""
	}

	return splt[0]
}

//////////
//Button//
//////////

//Class of button types
type ButtonType int

const (
	//Button which send signal about click
	DataButton ButtonType = iota

	//Button which send signal about click with params
	DynamicDataButton

	//Button which redirect to url
	UrlButton

	//Button which can change own url
	DynamicUrlButton

	//Button which send message with button text
	TextButton

	//Button which can switch own status
	ToggleButton
)

//Button Class
type Button struct {
	Type 		ButtonType
	Text 		string
	Data 		string
	Url 		string
	On 			string
	Off			string
	Handler 	func(tgbotapi.Update)
}

//Create new Text Button
func NewTextButton(text string, handler func(tgbotapi.Update)) (Button) {
	return Button{
		Type: TextButton,
		Text: text,
		Handler: handler,
	}
}

//Create new Data Button
func NewDataButton(text string, data string, handler func(tgbotapi.Update)) (Button) {
	return Button{
		Type: DataButton,
		Text: text,
		Data: data,
		Handler: handler,
	}
}

//Create new Data Button
func NewDynamicDataButton(text string, data string, handler func(tgbotapi.Update)) (Button) {
	return Button{
		Type: DynamicDataButton,
		Text: text,
		Data: data,
		Handler: handler,
	}
}

//Create new Url Button
func NewUrlButton(text string, url string) (Button) {
	return Button{
		Type: UrlButton,
		Text: text,
		Url: url,
	}
}

//Create new Dynamic Url Button
func NewDynamicUrlButton(text string) (Button) {
	return Button{
		Type: DynamicUrlButton,
		Text: text,
		Url: "%s",
	}
}

//Create new Toggle Button
func NewToggleButton(text string, off string, data string, handler func(tgbotapi.Update)) (Button) {
	return Button{
		Type: ToggleButton,
		Text: text,
		Data: data,
		Off: off,
		On: text,
		Handler: handler,
	}
}


////////
//Menu//
////////

//Class of location types
type LocationType int

const (
	//Without buttons
	NoKeyboard LocationType = iota

	//Inline button placement
	InlineLocation 

	//Simple button placement
	SimpleLocation 
)

//Class which describe how buttons locate
type ButtonsLocation struct {
	Type 		LocationType
	Pattern 	[]int
}


//Menu Class
type Menu struct {
	ID			int
	Text 		string
	Buttons 	[]Button
	Location 	ButtonsLocation
}

var (
	Lobby		[]Menu
	AdminLobby 	[]Menu
	Commands	[]Command
)

func LobbyPop() {
	if len(Lobby) == 0 {
		return
	}
	Lobby = Lobby[:len(Lobby) - 1]
}

func NewMenu(text string, buttons ...Button) (int) {
	location := ButtonsLocation{
		Type: NoKeyboard,
		Pattern: []int{len(buttons)},
	}	

	menu := Menu{
		ID: len(Lobby),
		Text: text,
		Buttons: buttons,
		Location: location,
	}
	Lobby = append(Lobby, menu)

	return menu.ID
}

func NewAdminMenu(text string, buttons ...Button) (int) {
	location := ButtonsLocation{
		Type: NoKeyboard,
		Pattern: []int{len(buttons)},
	}	

	menu := Menu{
		ID: len(AdminLobby),
		Text: text,
		Buttons: buttons,
		Location: location,
	}
	AdminLobby = append(AdminLobby, menu)

	return menu.ID
}

func (m *Menu) AddButtons(buttons ...Button) {
	for i := 0;i < len(buttons);i++ {
		m.Buttons = append(m.Buttons, buttons[i])
	}

	m.Location.Pattern = []int{len(m.Buttons)}
}

func (m *Menu) ChangeLocation(loc_type LocationType, pattern ...int) (error) {
	sum := 0
	for i := 0;i < len(pattern);i++ {
		sum += pattern[i]
	}
	if sum != len(m.Buttons) {
		log.Printf("Wait %d\nGot %d", len(m.Buttons), sum)
		return ErrorInvalidPattern
	}

	m.Location = ButtonsLocation{
		Type: loc_type,
		Pattern: pattern,
	}

	return nil
}

func (m Menu) Prepare(attr ...string) (Menu) {
	pos := 0
	for i := 0;i < len(m.Buttons);i++ {
		if pos == len(attr) {
			return m
		}

		if m.Buttons[i].Type == DynamicUrlButton {
			m.Buttons[i].Url = attr[pos]
			pos++
		}
	}
	return m
}

func (m Menu) Mark(attr ...bool) (Menu) {
	pos := 0
	for i := 0;i < len(m.Buttons);i++ {
		if pos == len(attr) {
			return m
		}

		if m.Buttons[i].Type == ToggleButton {
			log.Printf("%s == %s",m.Buttons[i].Text, m.Buttons[i].Off)

			if !attr[pos] {
				m.Buttons[i].Text = m.Buttons[i].Off
			} else {
				m.Buttons[i].Text = m.Buttons[i].On
			}

			pos++
		}
	}
	return m
}

func (m Menu) Call(chat_id int64, attr ...interface{}) (tgbotapi.MessageConfig) {
	text := fmt.Sprintf(m.Text, attr...)
	msg := tgbotapi.NewMessage(chat_id, text)

	if m.Location.Type == InlineLocation {
		var buttons [][]tgbotapi.InlineKeyboardButton
		pos := 0
		for i := 0;i < len(m.Location.Pattern);i++ {
			var row []tgbotapi.InlineKeyboardButton
			for j := 0;j < m.Location.Pattern[i];j++ {
				var button tgbotapi.InlineKeyboardButton
				if m.Buttons[pos].Type == UrlButton || m.Buttons[pos].Type == DynamicUrlButton {
					button = tgbotapi.NewInlineKeyboardButtonURL(m.Buttons[pos].Text, m.Buttons[pos].Url)
				} else if m.Buttons[pos].Type == DataButton || m.Buttons[pos].Type == DynamicDataButton || m.Buttons[pos].Type == ToggleButton {
					button = tgbotapi.NewInlineKeyboardButtonData(m.Buttons[pos].Text, m.Buttons[pos].Data)
				}

				row = append(row, button)
				pos++
			}

			buttons = append(buttons, row)
		}
		markup := tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		}

		msg.ReplyMarkup = markup
	} else if m.Location.Type == SimpleLocation {
		var buttons [][]tgbotapi.KeyboardButton
		pos := 0
		for i := 0;i < len(m.Location.Pattern);i++ {
			var row []tgbotapi.KeyboardButton
			for j := 0;j < m.Location.Pattern[i];j++ {
				row = append(row, tgbotapi.NewKeyboardButton(m.Buttons[pos].Text))
				pos++
			}

			buttons = append(buttons, row)
		}

		markup := tgbotapi.ReplyKeyboardMarkup{
			Keyboard: buttons,
			ResizeKeyboard: true,
		}

		msg.ReplyMarkup = markup
	}

	return msg
}

func (m Menu) Recall(chat_id int64, message_id int, attr ...interface{}) (tgbotapi.EditMessageTextConfig, tgbotapi.EditMessageReplyMarkupConfig) {
	msg := m.Call(chat_id, attr...)
	erm := tgbotapi.NewEditMessageReplyMarkup(chat_id, message_id, msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup))
	et := tgbotapi.NewEditMessageText(chat_id, message_id, msg.Text)

	return et, erm
}

func (m Menu) FastRecall(chat_id int64, message_id int, attr ...interface{}) (tgbotapi.EditMessageTextConfig) {
	msg := m.Call(chat_id, attr...)
	// erm := tgbotapi.NewEditMessageReplyMarkup(chat_id, message_id, )
	et := tgbotapi.NewEditMessageText(chat_id, message_id, msg.Text)

	if m.Location.Type == NoKeyboard {
		return et
	}

	rm := msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	et.ReplyMarkup = &rm

	return et//, erm
}

func Clear(chat_id int64, message_id int) (tgbotapi.DeleteMessageConfig) {
	return tgbotapi.NewDeleteMessage(chat_id, message_id)
}



//////////////////
//Daemon buttons//
//////////////////

type DaemonButton struct {
	ID 			int
	Button
}

var DaemonStorage []DaemonButton

func NewDaemonDataButton(text string, data string, handler func(tgbotapi.Update)) (int) {
	button := DaemonButton{
		ID: len(DaemonStorage),
		Button: Button {
			Type: DataButton,
			Text: text,
			Data: data,
			Handler: handler,
		},
	}

	DaemonStorage = append(DaemonStorage, button)

	return button.ID
}

func NewDaemonDynamicDataButton(text string, data string, handler func(tgbotapi.Update)) (int) {
	button := DaemonButton{
		ID: len(DaemonStorage),
		Button: Button {
			Type: DynamicDataButton,
			Text: text,
			Data: data,
			Handler: handler,
		},
	}

	DaemonStorage = append(DaemonStorage, button)

	return button.ID
}

func NewDaemonDynamicURLButton(text string) (int) {
	button := DaemonButton{
		ID: len(DaemonStorage),
		Button: Button {
			Type: DynamicUrlButton,
			Text: text,
		},
	}

	DaemonStorage = append(DaemonStorage, button)

	return button.ID
}

func (db DaemonButton) ChangeText(newText string) (DaemonButton) {
	db.Text = newText
	return db
}

func CastDaemonButtons(buttons ...DaemonButton) ([]tgbotapi.InlineKeyboardButton) {
	var typed_buttons []tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		if btn.Type == DataButton || btn.Type == DynamicDataButton {
			typed_button := tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data)
			typed_buttons = append(typed_buttons, typed_button)
		} else if btn.Type == DynamicUrlButton {
			typed_button := tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.Url)
			typed_buttons = append(typed_buttons, typed_button)
		}
	}

	return typed_buttons
}

func (db DaemonButton) MoreData(param string) (DaemonButton) {
	if db.Type == DynamicDataButton {
		db.Data += SEPARATOR + param
		return db
	} else {
		return db
	}
}

func (db DaemonButton) Prepare(url string) (DaemonButton) {
	db.Url = url
	return db
}

func (m Menu) AddDaemonButtons(buttons ...DaemonButton) (Menu) {
	for _, button := range buttons {
		m.Buttons = append(m.Buttons, button.Button)
	}
	m.Location.Pattern = append(m.Location.Pattern, len(buttons))
	return m
}

////////////////////////
//Photo Config Wrapper//
////////////////////////

type PhotoConfigWrapper tgbotapi.PhotoConfig

func (pc *PhotoConfigWrapper) AddDaemonButtons(buttons ...DaemonButton) {
	typed_buttons := CastDaemonButtons(buttons...)
	if pc.ReplyMarkup != nil {
		keyboard := pc.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, typed_buttons)
		pc.ReplyMarkup = keyboard
	} else {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(typed_buttons,)
		pc.ReplyMarkup = keyboard
	}
}


func (pc *PhotoConfigWrapper) ReshapeMenu(shape ...int) (error) {
	markup := pc.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	keyboard := markup.InlineKeyboard
	var buttons []tgbotapi.InlineKeyboardButton
	var size int

	for _, i := range shape {
		size += i
	}

	for i := 0;i < len(keyboard);i++ {
		buttons = append(buttons, keyboard[i]...)
	}

	if size != len(buttons) {
		return ErrorInvalidPattern
	}

	var n_row int
	var row []tgbotapi.InlineKeyboardButton
	var new_keyboard [][]tgbotapi.InlineKeyboardButton

	for i := 0;i < size;i++ {
		row = append(row, buttons[i])

		if len(row) == shape[n_row] {
			new_keyboard = append(new_keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
			n_row++
		}
	}

	pc.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(new_keyboard...)

	return nil
}
/*
func (pc *PhotoConfigWrapper) FastRecall(chatID int64, mes int) (tgbotapi.EditMessageCaptionConfig) {

	ec := tgbotapi.NewEditMessageCaption(chatID, mes, pc.Caption)
	ec.File = pc.File
	ec.FileID = pc.FileID
	ec.FileSize = pc.FileSize
	
	if pc.ReplyMarkup == nil {
		return ec
	}
	rm := pc.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	ec.ReplyMarkup = &rm
	
	return ec
}*/

//////////////
//Middleware//
//////////////

var (
	middleware func(u *tgbotapi.Update)(bool) = func(u *tgbotapi.Update)(bool) {	return	true	}
)

func DefineMiddleware(mdw func(u *tgbotapi.Update)(bool)) {
	middleware = mdw
}

func Middleware(u *tgbotapi.Update)(bool) {
	return middleware(u)
}

var (
	dropMode func(u tgbotapi.Update) = func(u tgbotapi.Update) {return}
)

func DefineDropMode(drm func(u tgbotapi.Update)) {
	dropMode = drm
}

func DropMode(u tgbotapi.Update) {
	dropMode(u)
}

//////////////
//Bot runner//
//////////////


func Run(updates *tgbotapi.UpdatesChannel) {

	NextUpdate:
	for update := range *updates {
		if !Middleware(&update) {
			continue
		}

		if update.Message == nil && update.InlineQuery != nil {
			fmt.Printf("Recieve inline mode\n")
			
		}
		
		if update.CallbackQuery != nil {
			fmt.Printf("Recieve callback\n")
			MergedLobby := append(Lobby, AdminLobby...)
			for _, menu := range MergedLobby {
				if menu.Location.Type != InlineLocation {
					continue
				}

				for _, button := range menu.Buttons {
					if (button.Type == DataButton || button.Type == ToggleButton) && button.Data == update.CallbackQuery.Data {
						DropMode(update)
						button.Handler(update)
						continue NextUpdate
					}
				}
			}

			for _, button := range DaemonStorage {
				if (button.Type == DataButton || button.Type == ToggleButton) && button.Data == update.CallbackQuery.Data {
					DropMode(update)
					button.Handler(update)
					continue NextUpdate
				} else if button.Type == DynamicDataButton && button.Data == strings.Split(update.CallbackQuery.Data, " ")[0] {
					DropMode(update)
					button.Handler(update)
					continue NextUpdate
				}
			}
		}

		if update.Message != nil && update.Message.Text != "" {
			fmt.Printf("Recieve message\n")
			if update.Message.Text[0] == '/' {
				splt := strings.Split(update.Message.Text, " ")
				command_name := splt[0]
				params := splt[1:]
				for _, command := range Commands {
					if command.Text == command_name {
						DropMode(update)
						command.Handler(update, params...)
						continue NextUpdate
					}
				}
			} else {
				MergedLobby := append(Lobby, AdminLobby...)
				for _, menu := range MergedLobby {
					if menu.Location.Type != SimpleLocation {
						continue
					}

					for _, button := range menu.Buttons {
						if button.Type == TextButton && button.Text == update.Message.Text {
							DropMode(update)
							button.Handler(update)
							continue NextUpdate
						}
					}
				}
			}
		}

		mode := ExtractMode(update)
		for _, input := range InManager {
			if input.Mode == mode {
				res := input.Handler(update)
				if res {
					DropMode(update)
				}
				continue NextUpdate
			}
		}
	}

}


///////

type Cell struct {
	Offset 		int
	Liter 		[]uint16
}

func InsertEntities(text string, entities []tgbotapi.MessageEntity) string {
	var cells []Cell	
	for _, ent := range entities {
		if ent.IsBold() {
			cells = append(cells, Cell{
				Offset: ent.Offset,
				Liter: utf16.Encode([]rune("*")),
			})

			cells = append(cells, Cell{
				Offset: (ent.Offset + ent.Length),
				Liter: utf16.Encode([]rune("*")),
			})
		} else if ent.IsTextLink() {
			cells = append(cells, Cell{
				Offset: ent.Offset,
				Liter: utf16.Encode([]rune("[")),
			})

			cells = append(cells, Cell{
				Offset: (ent.Offset + ent.Length),
				Liter: utf16.Encode([]rune(fmt.Sprintf("](%s)", ent.URL))),
			})
		}
	}
	runes := []rune(text)
	encr := utf16.Encode(runes)
	sort.SliceStable(cells, func(i, j int) bool {
		return cells[i].Offset > cells[j].Offset
	})
	log.Printf("%+v", cells)
	
	for _, cell := range cells {
		log.Printf("DO %s", string(utf16.Decode(encr)))
		encr = append(encr[:cell.Offset], append(cell.Liter, encr[cell.Offset:]...)...)
		log.Printf("POSLE %s", string(utf16.Decode(encr)))
	}
	return string(utf16.Decode(encr))
}