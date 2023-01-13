# ZincSearch - Api-rest :postbox:

Api-rest en GO que se conecta a Zincsearch y trae información de correos electrónicos.

## Documentación

El código presentado define un paquete llamado "routes" que contiene funciones para manejar rutas en un servidor web escrito en Go. El paquete importa varios paquetes de Go, así como un paquete personalizado llamado "models" (m.)

La función principal es "Routes()", que devuelve una nueva instancia de "chi.Mux", que es un router HTTP para Go. Dentro de esta función se habilitan varios middlewares globales, como "Logger", "Recoverer", "Timeout", "RequestID" y "CORS" (Cross-Origin Resource Sharing).

Se definen varias rutas, cada una con su propia función handler:

#### "/": "Welcome"
#### "POST /login": "Login"
#### "GET /index" (con autenticación básica): "GetIndex"
#### "POST /search" (con autenticación básica): "SearchHandler"

Mensaje de bienvenida

### Función "Login" 
Maneja una solicitud de inicio de sesión, donde se extraen las credenciales enviadas en el cuerpo de la solicitud y se valida mediante la función "checkCredentials". Se devuelve una respuesta JSON indicando si las credenciales son válidas o no.

### Función "GetIndex"

Maneja una solicitud para obtener una lista de índices, haciendo una solicitud HTTP GET al servidor Zincsearch utilizando las credenciales almacenadas en un archivo de configuración. Se devuelve la respuesta recibida del servidor Zincsearch en formato JSON.

### Función "SearchHandler"
Maneja una solicitud para buscar correos, que recibe en el cuerpo de la solicitud un objeto JSON con los parámetros de búsqueda. La función hace una solicitud HTTP POST al servidor Zincsearch con los parámetros de búsqueda y las credenciales almacenadas en un archivo de configuración. Se devuelve la respuesta recibida del servidor Zincsearch en formato JSON.

Request:
{
    "index": "enronmail",
    "search_type": "match",
    "query": {
        "term": "elvin"
    },
    "from": 0,
    "max_results": 20
}
