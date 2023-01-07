package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		middleware.Timeout(60 * time.Second),
		middleware.RequestID,
	)

	mux.Get("/index", GetIndex)
	mux.Post("/search", SearchHandler)
	
	return mux
}

// GetIndex trae la lista de indices
func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Realiza la petición HTTP POST al servidor Zincsearch
	req, err := http.NewRequest("GET", "http://localhost:4080/api/index",nil)
	if err != nil {
		 fmt.Errorf("Error al hacer la peticion GET a ZincSearch ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "Complexpass#123")

	//Recibe respuesta de ZincSearch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Procesa la respuesta Json y la convierte en estructura Index
	var result m.Index
	_ = json.NewDecoder(resp.Body).Decode(&result)

	//Devuelve el objeto json del response y lo coloca en el ResponseWriter
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
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

	results := m.Results{
		Total: res.Hits.Total.Value,
		Data:  m.EmailFields(res),
	}
	
	//Devuelve el objeto json del response de ZincSearch y lo coloca en el ResponseWriter
	err = json.NewEncoder(w).Encode(results)
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
		return m.Response{}, fmt.Errorf("Error al construir el cuerpo del request a enviar a ZincSearch: %v", err)
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
