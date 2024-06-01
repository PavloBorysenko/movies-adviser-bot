package telegram

const msgHelp = `I can save and keep you Movies. Also I can offer you them to watch.

In order to save the movie, just send me al link to it.
/saves [serie title or link] to save series
/savef [film title or link] to save film
/films get all films
/series get all series
/search [search text] to find movie
/deletef [search text] to remove film
/deletes [search text] to remove serie
In order to get a random movie from your list, send me command /rnd.`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedMovies   = "You have no saved Movies 🙊"
	msgNoFoundMovies   = "Unfortunately I couldn't find this movie 🙊"
	msgSaved          = " \U0001F3A5 Saved!  \U0001F4BE"
	msgAlreadyExists  = "You have already have this movie in your list 🤗"
	msgWasDeleted = " \U0001F4E4 Was deleted"
	msgToDelete = "To delete this  \U0001F6AB"
)
