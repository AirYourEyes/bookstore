package store

import (
	mystore "bookstore/store"
	"bookstore/store/factory"
	"errors"
	"sync"
)

type MemStore struct {
	sync.RWMutex
	books map[string]*mystore.Book
}

func init() {
	factory.Register("mem", &MemStore{
		books: make(map[string]*mystore.Book),
	})
}

func (store *MemStore) Create(book *mystore.Book) error {
	store.Lock()
	defer store.Unlock()

	_, ok := store.books[book.Id]
	if ok {
		return errors.New("book already exists")
	}

	store.books[book.Id] = book
	return nil
}

func (store *MemStore) Update(book *mystore.Book) error {
	store.Lock()
	defer store.Unlock()

	bookToUpdate, ok := store.books[book.Id]
	if !ok {
		return errors.New("book not found")
	}

	newName := book.Name
	if newName != bookToUpdate.Name {
		bookToUpdate.Name = newName
	}

	newAuthors := book.Authors
	if newAuthors != nil {
		bookToUpdate.Authors = newAuthors
	}

	newPress := book.Press
	if newPress != "" {
		bookToUpdate.Press = newPress
	}
	return nil
}

func (store *MemStore) Get(id string) (mystore.Book, error) {
	store.RLock()
	defer store.RUnlock()

	book, ok := store.books[id]
	if !ok {
		return mystore.Book{}, errors.New("book not found")
	}
	return *book, nil
}

func (store *MemStore) GetAll() ([]mystore.Book, error) {
	store.RLock()
	defer store.RUnlock()

	resultBooks := make([]mystore.Book, len(store.books))
	i := 0
	for _, book := range store.books {
		resultBooks[i] = *book
		i += 1
	}
	return resultBooks, nil
}

func (store *MemStore) Delete(id string) error {
	store.Lock()
	defer store.Unlock()

	_, ok := store.books[id]
	if !ok {
		return errors.New("book not found")
	}

	delete(store.books, id)
	return nil
}
