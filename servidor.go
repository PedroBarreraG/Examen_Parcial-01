package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

var clientes []net.Conn
var clientesNombres []string
var activosClientes []int64

var respaldoMensajes []string

func servidor() {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		if msg != "" {
			mens := strings.Split(msg, "|")

			opcion := mens[0]
			mensaje := mens[1]
			usuario := mens[2]

			//Registra Cliente (Conecta)
			if opcion == "0" {
				clientes = append(clientes, c)
				clientesNombres = append(clientesNombres, usuario)
				activosClientes = append(activosClientes, 1)

				respuesta := "0|Bienvenido " + usuario
				err2 := gob.NewEncoder(c).Encode(respuesta)
				if err2 != nil {
					fmt.Println(err)
				}
				fmt.Println("Se conecto: ", usuario)
			}

			//Cliente envio Mensaje
			if opcion == "1" {
				enviarMensajes(usuario, mensaje)
			}

			//Cliente envio Archivo
			if opcion == "2" {
				enviarArchivo(usuario, mensaje)
			}

			//Cliente quiere ver mensajes
			if opcion == "3" {
				clienteVerMensajes(usuario)
			}

			//Desconectar Cliente
			if opcion == "4" {
				i := sacarIndiceCliente(usuario)
				activosClientes[i] = 0
				fmt.Println("Se desconecto: ", usuario)
			}
		}
	}

}

func sacarNomArchivo(archivo string) string {
	var nombre string
	ruta := strings.Split(archivo, "\\")
	for _, v := range ruta {
		nombre = v
	}
	return nombre
}

func clienteVerMensajes(userAct string) {
	j := -1

	for i, v := range clientesNombres {
		if v == userAct {
			j = i
		}
	}

	if j == -1 {
		return
	}
	fmt.Println("----------------------------------")
	for _, v := range respaldoMensajes {
		mens := strings.Split(v, ":")
		if mens[0] == userAct {
			mensajeEnv := "1|Tú:"
			fmt.Println(v)
			for i, m := range mens {
				if i != 0 {
					mensajeEnv = mensajeEnv + m
				}
			}
			err := gob.NewEncoder(clientes[j]).Encode(mensajeEnv)
			if err != nil {
				continue
			}
		} else {
			mensajeEnv := "1|" + v
			fmt.Println(v)
			err := gob.NewEncoder(clientes[j]).Encode(mensajeEnv)
			if err != nil {
				continue
			}
		}
	}
	fmt.Println("----------------------------------")
}

func enviarArchivo(usuario string, archivo string) {
	mensaje := sacarNomArchivo(archivo)
	respald := usuario + ": " + mensaje
	respaldoMensajes = append(respaldoMensajes, respald)
	fmt.Println(respald)

	mensajeDes := "1|" + usuario + ": " + mensaje
	mensajeOri := "1|" + "Tú: " + mensaje
	for i, v := range clientesNombres {
		if v == usuario {
			err := gob.NewEncoder(clientes[i]).Encode(mensajeOri)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			err := gob.NewEncoder(clientes[i]).Encode(mensajeDes)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func enviarMensajes(usuario string, mensaje string) {
	respald := usuario + ": " + mensaje
	respaldoMensajes = append(respaldoMensajes, respald)
	fmt.Println(respald)

	mensajeDes := "1|" + usuario + ": " + mensaje
	mensajeOri := "1|" + "Tú: " + mensaje
	for i, v := range clientesNombres {
		if v == usuario {
			err := gob.NewEncoder(clientes[i]).Encode(mensajeOri)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			err := gob.NewEncoder(clientes[i]).Encode(mensajeDes)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func sacarIndiceCliente(nombre string) int64 {
	var i int64
	i = 0
	for _, v := range clientesNombres {
		if v == nombre {
			return i
		}
		i = i + 1
	}
	return -1
}

func serverVerMensajes() {
	fmt.Println("-------------------")
	for _, v := range respaldoMensajes {
		fmt.Println(v)
	}
	fmt.Println("-------------------")
}

func respaldarMensajes() {
	file, err := os.Create("respaldo_mensajes.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, v := range respaldoMensajes {
		file.WriteString(v)
		file.WriteString("\n")
	}
	fmt.Println("Los mensajes han sido repaldados")
}

func menuPrincipalServer() {
	var op int64
	for {
		fmt.Println("1) Mostrar los mensajes")
		fmt.Println("2) Respaldar los mensajes")
		fmt.Println("0) Terminar servidor")
		fmt.Scanln(&op)

		switch op {
		case 1:
			serverVerMensajes()
		case 2:
			respaldarMensajes()
		case 0:
			return
		}
	}
}

func main() {
	go servidor()
	menuPrincipalServer()
	//var input string
	//fmt.Scanln(&input)
}
