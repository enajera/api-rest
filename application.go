package main

import (
	"io/ioutil"
    "fmt"
	r "github.com/enajera/api-rest/pkg/routes"
	s "github.com/enajera/api-rest/pkg/server"
	"github.com/spf13/viper"
	
)

func main() {

	//Lectura del banner
	banner, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		fmt.Println("Error al abrir el banner.txt")
	}
	
	fmt.Println(string(banner))
	fmt.Println()

    //Viper para leer archivo de configuracion
	viper.AddConfigPath("./pkg/config")
	viper.SetConfigName("config") // nombre del archivo
    viper.SetConfigType("json")   // tipo
    viper.ReadInConfig()
	    
	//Crea los enrutadores
	chi := r.Routes()
	//envia los enrutadores al server
	server := s.NewServer(chi)
	//incia el server y queda escuchando
	server.Run()
}