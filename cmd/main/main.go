package main

import (
	r "github.com/enajera/api-rest/pkg/routes"
	s "github.com/enajera/api-rest/pkg/server"
)

func main() {
	//Crea los enrutadores
	mux := r.Routes()
	//envia los enrutadores al server
	server := s.NewServer(mux)
	//incia el server y queda escuchando
	server.Run()
}