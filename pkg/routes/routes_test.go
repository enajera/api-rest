package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/enajera/api-rest/pkg/models"
	"github.com/spf13/viper"
)

func TestLogin(t *testing.T) {
	
	// setear configuracion necesaria
   	viper.Set("basic_auth_user", "admin")
	viper.Set("basic_auth_pass", "root")

    //Crea una nueva petición de prueba
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"user":"admin", "pass":"root"}`)))
    if err != nil {
        t.Fatal(err)
    }

    // Crea un ResponseRecorder para analizar el resultado de la llamada al servidor
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Login)

    // Llamar al servidor con la petición de prueba
    handler.ServeHTTP(rr, req)

    // Comprueba el estado de la respuesta
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("bad status code: got %v want %v",
            status, http.StatusOK)
    }

    // Comprueba el contenido de la respuesta
    var resp m.LoginResponse
    json.Unmarshal(rr.Body.Bytes(), &resp)
    if !resp.Success {
        t.Errorf("respuesta inesperada: se obtuvo %v -> se esperaba %v",
            resp.Success, true)
    }
}

func TestGetIndex(t *testing.T) {

	// setear configuracion necesaria
	viper.Set("api", "https://playground.dev.zincsearch.com/api/")
	viper.Set("user", "admin")
	viper.Set("pass", "Complexpass#123")

	//Se crea un objeto tipo r *http.Request
	r := httptest.NewRequest("GET", "/index", nil)
	// Se crea un objeto tipo w http.ResponseWriter
	w := httptest.NewRecorder()

	// se crea la instancia de GetIndex
	GetIndex(w, r)

	// Chequea el status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("bad status code: Se obtuvo: %v -> Se esperaba: %v",
			status, http.StatusOK)
	}

	//Testeo del body
	got := m.Index{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Errorf("no se pudo procesar el json: %v", err)
	}

	//Estructura que al menos contiene un elemento
	want := &m.Index{
		List: []struct {
			Name string `json:"name"`
		}{
			{Name: "enronmail"},
		},
	}

	//Se evalua que al menos contenga un valor igual en la estructura:
	if !atLeastContains(&got, want) {
		t.Errorf("el resultado no contiene el valor esperado")
	}

}

func atLeastContains(got *m.Index, want *m.Index) bool {
	for _, v := range got.List {
		if v.Name == want.List[0].Name {
			return true
		}
	}
	return false
}

func TestSearchHandler(t *testing.T) {

	// setear configuracion necesaria
	viper.Set("api", "http://localhost:4080/api/")
	viper.Set("user", "admin")
	viper.Set("pass", "Complexpass#123")

	//Se crea un objeto tipo *http.Request
	sreq := bytes.NewBufferString(`{
        "index":"enronmail",
        "search_type": "match",
        "query":
        {
            "term": "elvin"
           
        },
        "from": 0,
        "max_results": 20
             
    }`)
	r := httptest.NewRequest("POST", "/search", sreq)
	// Se crea un objeto tipo w http.ResponseWriter
	w := httptest.NewRecorder()

	// se crea la instancia de SearchHandler
	SearchHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status inesperado: Se obtuvo %v-> Se esperaba %v Body -> %v", w.Code, http.StatusOK, w.Body)
	}

}
