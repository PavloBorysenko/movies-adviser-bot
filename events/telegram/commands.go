package telegram

import (
	"context"
	"errors"
	"log"
	"strings"

	"movies-adviser-bot/lib/e"
	"movies-adviser-bot/storage"
)
const (
	SaveCmd = "/save"
	SaveSeriesCmd = "/saves"
	SaveFilmCmd = "/savef"
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	AllFilmsCmd = "/films"
	AllSeriesCmd = "/series"
	FindCmd = "/search"
	DeleteFilmCmd = "/deletef"
	DeleteSeriesCmd = "/deletes"
	TypeFilm = "Film"
	TypeSeries = "Series"
)

var commands = map[string]string{
	SaveSeriesCmd: TypeSeries,
	SaveFilmCmd: TypeFilm,
	RndCmd: "",
	HelpCmd: "",
	StartCmd: "",
	AllFilmsCmd: TypeFilm,
	AllSeriesCmd: TypeSeries,
	FindCmd: "",
	DeleteFilmCmd: TypeFilm,
	DeleteSeriesCmd: TypeSeries,
}

var types = [2]string{TypeFilm, TypeSeries}


func (p *Processor) doCmd(text string, chatID int, username string) error {

	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	message, cmd := getCmd(text);
	ftype := getMovieType(cmd)

	switch cmd {	
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case AllFilmsCmd:
		return p.sendAll(chatID, username, ftype)
	case AllSeriesCmd:
		return p.sendAll(chatID, username, ftype)
	case SaveSeriesCmd:
		return p.savePage(chatID, message, username, ftype )
	case SaveFilmCmd:	
		return p.savePage(chatID, message, username, ftype )
	case FindCmd:
		return p.searchAllType(chatID, username, message)
	case DeleteFilmCmd:
		return p.deleteOne(chatID, username, messages, ftype)	
	case DeleteSeriesCmd:
		return p.deleteOne(chatID, username, message, ftype)	
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string, ftype string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	movie := &storage.Movie{
		Title:      pageURL,
		UserName: username,
		Type:    ftype,
	}

	isExists, err := p.storage.IsExists(context.Background(), movie)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), movie); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}
func (p *Processor) sendAll(chatID int, username string, ftype string ) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't get all records", err) }()

	movies, err := p.storage.GetAll(context.Background(), username, ftype)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, msgNoSavedMovies)
	}
	for _, movie := range movies {
		if err := p.tg.SendMessage(chatID, prepareMessage(movie) ); err != nil {
			return err
		}
	}


	return nil
}
func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	movie, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedMovies) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedMovies) {
		return p.tg.SendMessage(chatID, msgNoSavedMovies)
	}

	if err := p.tg.SendMessage(chatID, prepareMessage(movie)  ); err != nil {
		return err
	}

	return nil
}
func (p *Processor) searchAllType(chatID int, username string, searchtext string) error {
    for _, value := range types {

		err := p.sendOne(chatID, username, searchtext, value)
		if err != nil {
			return p.tg.SendMessage(chatID, msgNoFoundMovies + ":" + value)
		}
	}
	return nil
}
func (p *Processor) sendOne(chatID int, username string, searchtext string, ftype string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't sendOne", err) }()
	
	movie, err := p.storage.FindOne(context.Background(), username, searchtext, ftype)
	
	if err != nil && !errors.Is(err, storage.ErrNofoundMovies) {
		return err
	}
	if errors.Is(err, storage.ErrNofoundMovies) {
		return p.tg.SendMessage(chatID, msgNoFoundMovies)
	}

	if err := p.tg.SendMessage(chatID, prepareMessage(movie) ); err != nil {
		return err
	}

	return nil
}
func (p *Processor) deleteOne(chatID int, username string, searchtext string, ftype string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't deleteOne", err) }()
  
	movie, err := p.storage.FindOne(context.Background(), username, searchtext, ftype)
	log.Print(err)
	if err != nil && !errors.Is(err, storage.ErrNofoundMovies) {
		return err
	}
	if errors.Is(err, storage.ErrNofoundMovies) {
		return p.tg.SendMessage(chatID, msgNoFoundMovies)
	}

	if err := p.tg.SendMessage(chatID,  msgWasDeleted + ": " + movie.Title ); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), movie)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
func prepareMessage (film *storage.Movie) string{

	return " " + film.Type + ": " + film.Title + getDeleteCommand(film)
}
func getDeleteCommand(film *storage.Movie) string{

	switch film.Type {	
	case TypeFilm:
		return  " \n - " + msgToDelete + ": " +  DeleteFilmCmd + " " + film.Title
	case TypeSeries:
		return  " \n - " + msgToDelete + ": " + DeleteSeriesCmd + " " + film.Title
	}

	return ""
}


func getCmd(text string) (string, string){

	for key, _ := range commands {
		i := strings.Index(text, key)

		if i == 0 {
			return strings.TrimSpace(strings.ReplaceAll(text, key, "")), key
		}
	}		
	return text, ""
}
func getMovieType(key string) string {
	ftype := TypeSeries
	if t, ok :=commands[key]; ok {
		return t
	}
	return ftype
}
