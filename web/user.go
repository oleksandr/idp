package web

import (
	"net/http"

	"github.com/oleksandr/idp/usecases"
)

//
// UserWebHandler is a collection of CRUD methods for Users API
//
type UserWebHandler struct {
	UserInteractor usecases.UserInteractor
}

func (handler *UserWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	/*
	   defer r.Body.Close()
	   var a User
	   err := json.NewDecoder(r.Body).Decode(&a)
	   if err != nil {
	       log.Println(err.Error())
	       w.WriteHeader(http.StatusBadRequest)
	       return
	   }

	   a.ID = "abcd"
	   a.CreatedOn.Time = time.Now().UTC()
	   a.UpdatedOn = a.CreatedOn
	*/

	//w.Header().Add("Location", fmt.Sprintf("/accounts/%v", a.ID))
	//w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(map[string]User{"account": a})
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	//params := context.Get(r, ctxParamsKey).(httprouter.Params)
	//log.Printf("%#v", params)
	//json.NewEncoder(w).Encode(params.ByName("id"))
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Modify(w http.ResponseWriter, r *http.Request) {
	//params := context.Get(r, ctxParamsKey).(httprouter.Params)
	//log.Printf("%#v", params)
	//json.NewEncoder(w).Encode(params.ByName("id"))
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//params := context.Get(r, ctxParamsKey).(httprouter.Params)
	//log.Printf("%#v", params.ByName("id"))
	//w.WriteHeader(http.StatusNoContent)
	w.WriteHeader(http.StatusNotImplemented)
}
