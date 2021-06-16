package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Institucion struct {
	Institucion   string `json: institucion`
	Departamento  string `json: departamento`
	Provincia     string `json: provincia`
	Distrito      string `json: distrito`
	Representante string `json: representante`
	Sector        string `json: sector`
}

func cargarDatos(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)

	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

/* func resuelveListar(response http.ResponseWriter, request *http.Request) {
	//definir tipo de contenido de la respuesta
	response.Header().Set("Content-Type", "application/json")
	//serializar, codificar a json
	jsonBytes, _ := json.MarshalIndent(alumnos, "", "	") //2do obj _ es obj de error
	io.WriteString(response, string(jsonBytes))
	log.Println("Respuesta exitosa")
} */

/* func resuelveBuscarAlumno(response http.ResponseWriter, request *http.Request) {
	//http://localhost:98000/alumno?dni=12345678 etc
	log.Println("Llamada al endpoint /alumno")
	//recuperar parámetros por querystring
	sDni := request.FormValue("dni")
	response.Header().Set("Content-Type", "application/json")

	//logica del endpoint
	//_ pq se ignora el índice del for
	iDni, _ := strconv.Atoi(sDni)
	for _, alumno := range alumnos {
		if alumno.Dni == iDni {
			//codificar
			jsonBytes, _ := json.MarshalIndent(alumno, "", "	")
			io.WriteString(response, string(jsonBytes))
		}
	}

} */

func resuelveCreditos(response http.ResponseWriter, request *http.Request) {
	log.Println("Llamada al endpoint /creditos")
	response.Header().Set("Content-Type", "text/html") //para mandar un HTML
	io.WriteString(response,
		`<doctype html>
	<html>
		<head><title>API</title></head>
		<body>
			<h2>API para el curso de programación concurrente y dist.</h2>
		</body>
	</html>	
	`)
}

func resolveHome(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Home screen")
}

func resolveData(response http.ResponseWriter, request *http.Request) {

}

func manejadorRequest() { //definir endpoints de los servicios
	http.HandleFunc("/home", resolveHome)
	//http.HandleFunc("/data", resolveData)

	//establecer el puerto del servicio
	log.Fatal(http.ListenAndServe(":9000", nil)) //nil para no usar el manejador. Fatal imprime si hay alguna excepcion
}

func main() {

	url := "https://raw.githubusercontent.com/cjmazzarri/Go-Rest-API/develop/IPREDA_Dataset.csv?token=AL5F43MBXVZMOOWSRB47KT3A2NZPA"
	data, err := cargarDatos(url)
	if err != nil {
		panic(err)
	}

	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}
		fmt.Println(row)
	}

	manejadorRequest()

}
