package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Bono struct {
	Id_persona int     `json: id_persona`
	Prestacion float64 `json: prestacion`
	Tipotra    float64 `json: tipotra`
	Tipoben    float64 `json: tipoben`
	Beneficio  float64 `json: beneficio`
}

var bono Bono
var bonos []Bono
var centroides []Bono
var clusters_centroides []int

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
	response.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(bonos, "", "	")
	io.WriteString(response, string(jsonBytes))
	log.Println("Respuesta del endpoint Data")
}

func resolveGrupos(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	jsonBytes, _ := json.MarshalIndent(clusters_centroides, "", "	")
	io.WriteString(response, string(jsonBytes))
	log.Println("Rpta endpoint Grupos")
}

func manejadorRequest() { //definir endpoints de los servicios
	http.HandleFunc("/", resolveData)
	http.HandleFunc("/grupos", resolveGrupos)

	//establecer el puerto del servicio
	log.Fatal(http.ListenAndServe(":9000", nil)) //nil para no usar el manejador. Fatal imprime si hay alguna excepcion
}

func escogerCentroides(k int) {
	for i := 0; i < k; i++ {
		r := rand.Intn(len(bonos))
		centroides = append(centroides, bonos[r])
	}
	fmt.Println("Centroides iniciales: ")
	for i := range centroides {
		fmt.Println(centroides[i])
	}
}

func distanciaEuclidiana(bono1 Bono, bono2 Bono) float64 {
	suma :=
		math.Pow(bono2.Prestacion-bono1.Prestacion, 2) +
			math.Pow(bono2.Tipotra-bono1.Tipotra, 2) +
			math.Pow(bono2.Tipoben-bono1.Tipoben, 2) +
			math.Pow(bono2.Beneficio-bono1.Beneficio, 2)
	return math.Sqrt(suma)
}

func agrupar() {
	var i int = 0
	for i = 0; i < len(bonos); i++ {
		var auxJ int = 0
		minDist := distanciaEuclidiana(bonos[i], centroides[0])
		for j := range centroides {
			dist := distanciaEuclidiana(bonos[i], centroides[j])
			if dist < minDist {
				minDist = dist
				auxJ = j
			}
		}

		//fmt.Println("agregado valor ", i, " a cluster de ", centroides[auxJ].Id_persona)

		clusters_centroides = append(clusters_centroides, auxJ)

	}
}

func inicializarGrupo() {
	for i := range bonos {
		clusters_centroides[i] = 0
	}
}

func actualizarCentroide(k int) {
	var Bono Bono
	var sumaPrestacion, sumaTipotra, sumaTipoben, sumaBeneficio, cont float64 = 0, 0, 0, 0, 0
	var j int = 0
	for i := 0; i < len(bonos); i++ {
		if clusters_centroides[i] == k {
			sumaPrestacion += bonos[i].Prestacion
			sumaTipotra += bonos[i].Tipotra
			sumaTipoben += bonos[i].Tipoben
			sumaBeneficio += bonos[i].Beneficio
			cont += 1
			j++
		}
	}
	Bono.Id_persona = j
	Bono.Prestacion = sumaPrestacion / cont
	Bono.Tipotra = sumaTipotra / cont
	Bono.Tipoben = sumaTipoben / cont
	Bono.Beneficio = sumaBeneficio / cont
	centroides[k] = Bono
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//raw link:
	url := "https://raw.githubusercontent.com/cjmazzarri/Go-Rest-API/develop/dataset_BC_PC_ago2020.csv?token=AL5F43NOLJ4CVVWZTDNYKYDA2THUI"
	data, err := cargarDatos(url)
	if err != nil {
		panic(err)
	}

	//Poblar los datos obtenidos del archivo CSV accedido por el raw link
	for i, value := range data {
		bono.Id_persona = i
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
	escogerCentroides(3)

	for i := 0; i < 10; i++ {
		agrupar()
		fmt.Println("\nnuevos centroides:")
		for j := range centroides {
			actualizarCentroide(j)
			fmt.Println(centroides[j])
		}
	}

	manejadorRequest()
}
