package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	m "github.com/enajera/api-rest/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes() *chi.Mux {
	mux := chi.NewRouter()

	//globals middlewares
	mux.Use(
		middleware.Logger,    //log every http request
		middleware.Recoverer, //recover if a panic occurs
	)

	mux.Get("/", InitHandler)
	mux.Post("/search", SearchHandler)
	mux.Post("/search/{param}", SearchParamHandler)

	return mux
}

// InitHandler metodo de inicio
func InitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	res := map[string]interface{}{"message": "No results"}
	_ = json.NewEncoder(w).Encode(res)
}

// SearchParamHandler busqueda de un parametro
func SearchParamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	param := chi.URLParam(r, "param")
	res := map[string]interface{}{"message": param}
	_ = json.NewEncoder(w).Encode(res)
}

// SearchHandler obtiene los resultados de la busqueda
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//Decodifica el body y lo mapea al objeto request
	req := ParseRequest(r)

	//Obtengo la respuesta 
	res, err := Buscar(*req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := m.EmailFields(res)
	//Devuelve el objeto json del response de ZincSearch y lo coloca en el ResponseWriter
	err = json.NewEncoder(w).Encode(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Parsea el Body a un objeto Request
func ParseRequest(r *http.Request) *m.Request {
	body := r.Body
	defer body.Close()
	var req m.Request
	_ = json.NewDecoder(body).Decode(&req)

	return &req
}

// Regresa el response de Zincseach
func Buscar(request m.Request) (m.Response, error) {

	//Se construye el cuerpo del request a enviar a ZincSearch
	body, err := json.Marshal(request)
	if err != nil {
		return m.Response{}, fmt.Errorf("Error al construir el cuerpo del request a enviar a ZincSearch ", err)
	}

	// Realiza la petición HTTP POST al servidor Zincsearch
	req, err := http.NewRequest("POST", "http://localhost:4080/api/enronmail/_search", bytes.NewBuffer(body))
	if err != nil {
		return m.Response{}, fmt.Errorf("Error al hacer la peticion POST a ZincSearch ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "Complexpass#123")

	//Recibe respuesta de ZincSearch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return m.Response{}, fmt.Errorf("Error al obtener respuesta de ZinsSearch ", err)
	}
	defer resp.Body.Close()

	// Procesa la respuesta Json y la convierte en estructura Response
	var res m.Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return m.Response{}, fmt.Errorf("Error al decodificar estructura del Response ", err)
	}

	return res, nil
}
