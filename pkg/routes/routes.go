package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
     m "github.com/enajera/api-rest/pkg/models"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func Routes() *chi.Mux {

	mux := chi.NewRouter()

	//Se habilita la funcion cross-origin resource sharing (CORS)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
		Debug:            true,
	})

	//globals middlewares
	mux.Use(
		middleware.Logger,    //log de request http
		middleware.Recoverer, //se recupera si hay panico
		middleware.Timeout(60*time.Second),
		middleware.RequestID,
		c.Handler,
	)

	

	//Bienvenida
	mux.Get("/", Welcome)
	//Login
	mux.Post("/login", Login)
	//Obtiene listado de indices
	mux.With(basicAuth).Get("/index", GetIndex)
	//Obtiene listado de correos
	mux.With(basicAuth).Post("/search", SearchHandler)

	

	return mux

}

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "*")

	//Parsea el request
	var req m.Login
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Valida las credenciales
	var resp m.LoginResponse
	resp.Success = checkCredentials(req.User, req.Pass)

	//Devuelve la respuesta
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetIndex trae la lista de indices
func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	//variables en archivos de configuracion
	api := viper.GetString("api")
	user := viper.GetString("user")
	pass := viper.GetString("pass")

	// Realiza la petición HTTP POST al servidor Zincsearch
	req, err := http.NewRequest("GET", api+"index", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error de request: %v ", err), http.StatusBadRequest)
		return

	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, pass)

	//Recibe respuesta de ZincSearch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		http.Error(w, fmt.Sprintf("error al recibir respuesta de ZincSearch: %v ", err),
			http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "error al obtener los indices de ZincSearch", resp.StatusCode)
		return
	}

	// Procesa la respuesta Json y la convierte en estructura Index
	var result m.Index
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, fmt.Sprintf("error al parsear la respuesta Json: %v", err), http.StatusInternalServerError)
		return
	}

	//Devuelve el objeto json del response y lo coloca en el ResponseWriter
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("error al enviar la respuesta Json: %v", err), http.StatusInternalServerError)
		return
	}

}

// SearchHandler obtiene los resultados de la busqueda
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	//Decodifica el body y lo mapea al objeto request
	req, err := ParseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validar el indice
	if req.Index == "" {
		http.Error(w, "no se encontro el campo index", http.StatusBadRequest)
		return
	}

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
func ParseRequest(r *http.Request) (*m.SearchRequest, error) {
	defer r.Body.Close()
	var req m.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("error al parsear el body de la request: %v", err)
	}
	fmt.Println("Search request:", req)
	return &req, nil
}

// Regresa el response de Zincseach
func Buscar(r m.SearchRequest) (m.Response, error) {

	index := r.Index
	var request m.Request
	request.SearchType = r.SearchType
	request.Query = r.Query
	request.From = r.From
	request.MaxResults = r.MaxResults

	//Se construye el cuerpo del request a enviar a ZincSearch
	body, err := json.Marshal(request)
	if err != nil {
		return m.Response{}, fmt.Errorf("error al construir el cuerpo del request a enviar a ZincSearch: %v", err)
	}

	api := viper.GetString("api")
	user := viper.GetString("user")
	pass := viper.GetString("pass")

	// Realiza la petición HTTP POST al servidor Zincsearch
	req, err := http.NewRequest("POST", api+index+"/_search", bytes.NewBuffer(body))
	if err != nil {
		return m.Response{}, fmt.Errorf("error al hacer la peticion POST a ZincSearch:  %v ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, pass)

	//Recibe respuesta de ZincSearch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return m.Response{}, fmt.Errorf("error al procesar respuesta de ZinsSearch: %v ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return m.Response{}, fmt.Errorf("error al obtener respuesta de ZinsSearch: %v ", err)

	}

	// Procesa la respuesta Json y la convierte en estructura Response
	var res m.Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return m.Response{}, fmt.Errorf("error al decodificar estructura del Response: %v ", err)
	}

	return res, nil
}

// Autenticacion basica
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		user, pass, ok := r.BasicAuth()
		if !ok || !checkCredentials(user, pass) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	
		next.ServeHTTP(w, r)
	})
}

// Chequeo de credenciales
func checkCredentials(user, pass string) bool {

	dbUser := viper.GetString("basic_auth_user")
	dbPass := viper.GetString("basic_auth_pass")
	return user == dbUser && pass == dbPass
}

// Bienvenida
func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	banner :=
		`  /$$$$$$            /$$         /$$$$$$$                        /$$    
 /$$__  $$          |__/        | $$__  $$                      | $$    
| $$  \ $$  /$$$$$$  /$$        | $$  \ $$  /$$$$$$   /$$$$$$$ /$$$$$$  
| $$$$$$$$ /$$__  $$| $$ /$$$$$$| $$$$$$$/ /$$__  $$ /$$_____/|_  $$_/  
| $$__  $$| $$  \ $$| $$|______/| $$__  $$| $$$$$$$$|  $$$$$$   | $$    
| $$  | $$| $$  | $$| $$        | $$  \ $$| $$_____/ \____  $$  | $$ /$$
| $$  | $$| $$$$$$$/| $$        | $$  | $$|  $$$$$$$ /$$$$$$$/  |  $$$$/
|__/  |__/| $$____/ |__/        |__/  |__/ \_______/|_______/    \___/ 
          | $$                                                          
          | $$                                                          
          |__/                                                              
	
                            ¡Welcome!		  
                     ZincSearch Api-Rest v1.0          
                     Powered by Elvin Nájera S            
	`
	fmt.Fprint(w, banner)
}
