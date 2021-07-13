package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/julienschmidt/httprouter"
)

/* Media Feed */

func (a *App) listFeedSource() {

}

func (a *App) subscribeFeedSource() {

}

func (a *App) unsubscribeFeedSource() {

}

func (a *App) collectFeed() {

}

/* Media */

func (a *App) linkMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	var input struct {
		Link        string `json:"link"`
	}

	var output struct {
		Name      string `json:"name"`
		MediaType string `json:"media_type"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, err)
		return
	}

	switch {
	default:
		res, err := http.Get(input.Link)
		if err != nil {
			a.badRequestResponse(w, err)
			return
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			a.badRequestResponse(w, err)
			return
		}

		output.Name = doc.Find("title").Text()
		output.MediaType = "articles"
	}

	err = a.writeJSON(w, http.StatusOK, output)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *App) listMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := a.Models.Media.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, 200, media)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *App) createMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	var input struct {
		Description string  `json:"description"`
		MediaType   string  `json:"media_type"`
		RelatedLink *string `json:"related_link"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, err)
		return
	}

	var media Media
	media.Description = input.Description
	media.MediaType = input.MediaType
	media.RelatedLink = input.RelatedLink

	err = a.Models.Media.Insert(&media)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusCreated, media)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *App) updateMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	idString := ps.ByName("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		a.notFoundResponse(w)
		return
	}

	media, err := a.Models.Media.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			a.notFoundResponse(w)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input Media

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, err)
		return
	}

	media.ID = id
	media.Description = input.Description
	media.Rating = input.Rating
	media.Notes = input.Notes
	media.RelatedLink = input.RelatedLink
	media.Cart = input.Cart
	media.Completed = input.Completed
	media.Progress = input.Progress
	media.CompletedDate = input.CompletedDate

	err = a.Models.Media.Update(media)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *App) deleteMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)

	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	idString := ps.ByName("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		a.notFoundResponse(w)
		return
	}

	err = a.Models.Media.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			a.notFoundResponse(w)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* Todos */

func (a *App) listTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if err != nil {
		a.unauthorizedResponse(w)
		return
	}

	todos, err := a.Models.Todo.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	if !authenticated {
		filterd_todos := []*Todo{}
		for i := range todos {
			if !todos[i].Private {
				filterd_todos = append(filterd_todos, todos[i])
			}
		}
		todos = filterd_todos
	}

	err = a.writeJSON(w, 200, todos)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *App) createTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	var input struct {
		Description string     `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		Private     bool       `json:"private"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, err)
		return
	}

	var todo Todo
	todo.Description = input.Description
	todo.DueDate = input.DueDate
	todo.Private = input.Private
	err = a.Models.Todo.Insert(&todo)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusCreated, todo)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *App) updateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	idString := ps.ByName("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		a.notFoundResponse(w)
		return
	}

	todo, err := a.Models.Todo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			a.notFoundResponse(w)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input Todo

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, err)
		return
	}

	todo.ID = id
	todo.Description = input.Description
	todo.Cart = input.Cart
	todo.Completed = input.Completed
	todo.Progress = input.Progress
	todo.DueDate = input.DueDate
	todo.CompletedDate = input.CompletedDate
	todo.Private = input.Private

	err = a.Models.Todo.Update(todo)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *App) deleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}

	idString := ps.ByName("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		a.notFoundResponse(w)
		return
	}

	err = a.Models.Todo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			a.notFoundResponse(w)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* Daily */

func (a *App) listDaily() {

}

func (a *App) createDaily() {

}

func (a *App) updateDaily() {

}

func (a *App) deleteDaily() {

}

/* Reminder */

func (a *App) listReminder() {

}

func (a *App) createReminder() {

}

func (a *App) updateReminder() {

}

func (a *App) deleteReminder() {

}

/* Misc */

func (a *App) updateSticky() {

}

func (a *App) activityLog() {

}

func (a *App) confirmToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticated, err := a.IsAuthenticated(r)
	if !authenticated || err != nil {
		a.unauthorizedResponse(w)
		return
	}
	res := struct{ response string }{"token is valid"}
	err = a.writeJSON(w, http.StatusOK, &res)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

var badAuthError = errors.New("Invalid authentication password")

func (a *App) IsAuthenticated(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false, nil
	}

	getHashed := func(s string) string {
		b := sha256.Sum256([]byte(s))
		converted := b[:]
		return hex.EncodeToString(converted)
	}

	if cookie.Value == getHashed(a.Password) {
		return true, nil
	}

	return false, badAuthError
}
