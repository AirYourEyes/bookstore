package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type BookStoreServer struct {
	store          store.Store
	internalServer *http.Server
}

func NewBookStoreServer(addr string, store store.Store) *BookStoreServer {
	bookStoreServer := &BookStoreServer{
		store: store,
		internalServer: &http.Server{
			Addr: addr,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/book", bookStoreServer.createBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", bookStoreServer.updateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", bookStoreServer.getBookHandler).Methods("GET")
	router.HandleFunc("/book", bookStoreServer.getAllBookHandler).Methods("GET")
	router.HandleFunc("/book/{id}", bookStoreServer.deleteBookHandler).Methods("DELETE")

	bookStoreServer.internalServer.Handler = middleware.Logging(middleware.Validating(router))
	return bookStoreServer
}

func (bookStoreServer *BookStoreServer) createBookHandler(writer http.ResponseWriter, request *http.Request) {
	jsonDecoder := json.NewDecoder(request.Body)

	var book store.Book
	if err := jsonDecoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bookStoreServer.store.Create(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bookStoreServer *BookStoreServer) updateBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}

	jsonDecoder := json.NewDecoder(request.Body)

	var book store.Book
	if err := jsonDecoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	book.Id = id
	err := bookStoreServer.store.Update(&book)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bookStoreServer *BookStoreServer) getBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
	}

	book, err := bookStoreServer.store.Get(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response(writer, book)
}

func response(writer http.ResponseWriter, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (bookStoreServer *BookStoreServer) getAllBookHandler(writer http.ResponseWriter, request *http.Request) {
	allBooks, err := bookStoreServer.store.GetAll()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	response(writer, allBooks)
}

func (bookStoreServer *BookStoreServer) deleteBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}

	err := bookStoreServer.store.Delete(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bookStoreServer *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bookStoreServer.internalServer.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bookStoreServer *BookStoreServer) Shutdown(ctx context.Context) error {
	err := bookStoreServer.internalServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
