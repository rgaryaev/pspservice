package lstnr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rgaryaev/pspservice/config"
	"github.com/rgaryaev/pspservice/storage"
)

const defaultTimeoutInSec = 3

type passportStatus struct {
	Series string `json:"series"`
	Number string `json:"number"`
	Status string `json:"status"`
}

type baseHandler struct {
	mu                 sync.RWMutex
	requestCount       uint64
	storage            *storage.Storage
	passportPerRequest uint
}

// serveAsGET processing request with single passport number in GET paprametres
// /passport?series=xxxx&number=yyyyyy
func (h *baseHandler) serveAsGET(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
		fmt.Fprintf(w, "number of recieved requests = %d\n", h.requestCount)
		return
	case "/passport":
		h.mu.Lock()
		h.requestCount++
		h.mu.Unlock()
		series := r.URL.Query().Get("series")
		number := r.URL.Query().Get("number")

		if len(series) == 0 || len(number) == 0 {
			badRequestMsg(&w, "incorrect paramenters,  expected as /passport?series=xxxx&number=yyyyyy")
			return
		}

		passportIs := func(status string) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s %s : %s\n", series, number, status)
		}

		inList, err := (*(h.storage)).IsPassportInList(series, number)
		if err != nil {
			// pure error - something wrong
			msg := "Internal error"
			http.Error(w, msg, http.StatusBadGateway)
			return
		}
		// do response as json or as html
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			passport := passportStatus{series, number, ""}
			if inList {
				passport.Status = "non-valid"
			} else {
				passport.Status = "valid"
			}
			json := json.NewEncoder(w)
			json.SetIndent("", "\t")
			err := json.Encode(&passport)
			if err != nil {
				http.Error(w, "json was not encoded for response", http.StatusBadRequest)
				return
			}

		} else {
			//as html
			if inList {
				//  incorect format is considered as a non-valid passport
				passportIs("non-valid")
			} else {
				passportIs("valid")
			}
		}
		return
	}
	//
	badRequestMsg(&w, "")
}

// serveAsPOST processing request with json in the Body
func (h *baseHandler) serveAsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var passportList []passportStatus

	switch r.URL.Path {
	case "/passport":
		{
			h.mu.Lock()
			h.requestCount++
			h.mu.Unlock()

			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()
			jsonNotParsed := func(err error) {
				http.Error(w, "input json was not decoded", http.StatusBadRequest)

			}
			err := dec.Decode(&passportList)
			if err != nil {
				jsonNotParsed(err)
				return
			}
			if len(passportList) == 0 {
				jsonNotParsed(err)
				return
			}
			// check all passport in the list
			for index, passport := range passportList {
				inList, err := (*(h.storage)).IsPassportInList(passport.Series, passport.Number)
				if err != nil {
					// pure error - something wrong
					msg := "Internal error"
					http.Error(w, msg, http.StatusBadGateway)
					return
				}
				if inList {
					//  incorect format is considered as a non-valid passport
					passportList[index].Status = "non-valid"
				} else {
					passportList[index].Status = "valid"
				}
			}
			json := json.NewEncoder(w)
			json.SetIndent("", "\t")
			err = json.Encode(&passportList)

			if err != nil {
				http.Error(w, "json is not encoded in response", http.StatusBadRequest)
				return
			}

			return
		}
	}
	badRequestMsg(&w, "unknown url path, expected /passport, recieved :"+r.URL.Path)
}

func (h *baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	///fmt.Println("Method : " + http.MethodPost + "; Content-Type is: " + r.Header.Get("Content-Type"))

	if r.Method == http.MethodPost {
		h.serveAsPOST(w, r)
	} else if r.Method == http.MethodGet {
		h.serveAsGET(w, r)
	} else {
		badRequestMsg(&w, "unexpected method, only GET or POST are allowed")
	}
}

func badRequestMsg(w *http.ResponseWriter, msg string) {
	http.Error(*w, "Unsupported request: "+msg, http.StatusBadRequest)
}

// StartListener - start http listner
func StartListener(cfg *config.Configuration, storage *storage.Storage) error {

	var handler *baseHandler = new(baseHandler)
	handler.passportPerRequest = cfg.Listener.MaxPassportPerRequest
	//
	if storage == nil {
		return errors.New("http listener:  passport data storage is not initialized")
	}
	handler.storage = storage

	m := http.NewServeMux()
	srv := &http.Server{
		Addr:    cfg.Listener.Address + ":" + cfg.Listener.Port,
		Handler: m,
		//ReadTimeout: defaultTimeoutInSec * time.Second,
		//WriteTimeout: defaultTimeoutInSec * time.Second,
	}
	m.Handle("/", handler)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
