package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Bono struct {
	Id_persona float64 `json: id_persona`
	Prestacion float64 `json: prestacion`
	Tipotra    float64 `json: tipotra`
	Tipoben    float64 `json: tipoben`
	Beneficio  float64 `json: beneficio`
}

var bono Bono
var bonos []Bono

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

func resolveData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(bonos, "", "	")
	io.WriteString(response, string(jsonBytes))
	log.Println("Respuesta del endpoint Data")
}

func manejadorRequest() { //definir endpoints de los servicios
	http.HandleFunc("/", resolveData)

	//establecer el puerto del servicio
	log.Fatal(http.ListenAndServe(":9000", nil)) //nil para no usar el manejador. Fatal imprime si hay alguna excepcion
}

func main() {

	url := "https://raw.githubusercontent.com/cjmazzarri/Go-Rest-API/develop/dataset_BC_PC_ago2020.csv?token=AL5F43NOLJ4CVVWZTDNYKYDA2THUI"
	data, err := cargarDatos(url)
	if err != nil {
		panic(err)
	}

	for _, value := range data {
		bono.Id_persona, _ = strconv.ParseFloat(value[0], 64)
		bono.Prestacion, _ = strconv.ParseFloat(value[1], 64)
		bono.Tipotra, _ = strconv.ParseFloat(value[2], 64)
		bono.Tipoben, _ = strconv.ParseFloat(value[3], 64)
		bono.Beneficio, _ = strconv.ParseFloat(value[4], 64)

		bonos = append(bonos, bono) //Agregar a array bonos
	}

	for idx, row := range bonos {
		// Saltar la primera fila, contiene nombres de tablas
		if idx == 0 {
			continue
		}

		//Pequeña muestra de 10 elementos, porque el dataset tiene más de 11000
		if idx == 10 {
			break
		}
		fmt.Println(row)
	}

	manejadorRequest()

}
