package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var transact *FileTransactionLogger

func initializeTransactionLog() error {
	var err error
	transact, err = NewTransactionLogger("/tmp/transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}
	// read the events from the transaction log file
	events, errors := transact.ReadEvents()
	// create a new Event to hold the event from the log file
	e := Event{}
	ok := true
	count := 0

	for ok && err == nil {
		select {
		// condition if we get an error in the Error channel
		case err, ok = <-errors:
		// condition if we get the events in the Events channel
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
				count++
			case EventPut:
				err = Put(e.Key, e.Value)
				count++
			}
		}
	}
	fmt.Printf("%d events have been replayed\n", count)
	transact.Run()

	return err
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	transact.WriteHealthcheck()
	fmt.Fprintln(w, "Healthcheck works...")
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transact.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(string(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transact.WriteGet(key, string(value))
	fmt.Fprintf(w, value)
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(string(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DELETE key=%s\n", key)
	transact.WriteDelete(key)
}
