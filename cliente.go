package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

var nickname string

func cliente() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := "0|null|" + nickname
	err = gob.NewEncoder(c).Encode(msg)

	for {
		var msge string
		err = gob.NewDecoder(c).Decode(&msge)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			mens := strings.Split(msge, "|")

			opcion := mens[0]
			mensaje := mens[1]
			//Bienvaenida
			if opcion == "0" {
				fmt.Println(mensaje)
			}
			//Imprime mensajes resividos
			if opcion == "1" {
				fmt.Println(mensaje)
			}
		}
	}
	c.Close()
}

func menuPrincipal() {
	var op int64
	for {
		fmt.Println("1) Enviar Mensaje")
		fmt.Println("2) Enviar Archivo")
		fmt.Println("3) Mostrar Chat")
		fmt.Println("0) Salir")
		fmt.Scanln(&op)

		switch op {
		case 1:
			c, err := net.Dial("tcp", ":9999")
			if err != nil {
				fmt.Println(err)
				return
			}
			var mensaje string

			consoleReader := bufio.NewReader(os.Stdin)
			fmt.Println("Mensaje: ")
			input, _ := consoleReader.ReadString('\n')

			mensF := strings.Split(input, "\n")
			mensaje = mensF[0]
			msg := "1|" + mensaje + "|" + nickname
			err = gob.NewEncoder(c).Encode(msg)
			c.Close()
		case 2:
			c, err := net.Dial("tcp", ":9999")
			if err != nil {
				fmt.Println(err)
				return
			}
			var mensaje string
			fmt.Println("Ruta del archivo: ")
			fmt.Scanln(&mensaje)
			msg := "2|" + mensaje + "|" + nickname
			err = gob.NewEncoder(c).Encode(msg)
			c.Close()
		case 3:
			c, err := net.Dial("tcp", ":9999")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("-------------------")
			fmt.Println("Ventana de Chat: \n")
			msg := "3|null|" + nickname
			err = gob.NewEncoder(c).Encode(msg)
			c.Close()
			var input string
			fmt.Scanln(&input)
		case 0:
			return
		}
	}
}

func clienteEND() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := "4|null|" + nickname
	err = gob.NewEncoder(c).Encode(msg)
	c.Close()
}

func main() {
	fmt.Println("Nickname: ")
	fmt.Scanln(&nickname)

	go cliente()

	menuPrincipal()

	clienteEND()
	fmt.Println("Desconectado")
}
