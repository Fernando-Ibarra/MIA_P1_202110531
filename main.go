package main

import (
	"Proyecto1/Comands"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var counterDisk = 1

var logued = false

func main() {
	for true {
		fmt.Println("********************* INGRESE UN COMANDO *********************")
		fmt.Println("***** Si desea terminar con la aplicación ingrese \"exit\"")
		fmt.Print("\t")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		election := strings.TrimRight(input, "\r\n")
		if election == "exit" {
			break
		}
		comand := Comand(election)
		election = strings.TrimSpace(election)
		election = strings.TrimLeft(election, comand)
		tokens := Separatorokens(election)
		functions(comand, tokens)
		fmt.Println("\tPresione Enter para continuar...")
		fmt.Scanln()
	}
}

func Comand(text string) string {
	var tkn string
	finished := false
	for i := 0; i < len(text); i++ {
		if finished {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !finished {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				finished = true
			}
		}
	}
	return tkn
}

func Separatorokens(text string) []string {
	var tokens []string
	if text == "" {
		return tokens
	}
	text += " "
	var token string
	state := 0
	for i := 0; i < len(text); i++ {
		c := string(text[i])
		if state == 0 && c == "-" {
			state = 1
		} else if state == 0 && c == "#" {
			continue
		} else if state != 0 {
			if state == 1 {
				if c == "=" {
					state = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(text[i+1]) == " " && string(text[i-1]) == "-" {
					state = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(text[i+1]) == " " && string(text[i-1]) == "-" {
					state = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if state == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					state = 3
					continue
				} else {
					state = 4
				}
			} else if state == 3 {
				if c == "\"" {
					state = 4
					continue
				}
			} else if state == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if state == 4 && c == " " {
				state = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func functions(token string, tks []string) {
	if token != "" {
		if Comands.Compare(token, "EXECUTE") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN EXEC <<<<<<<<<<<<<<<<<<<<")
			ExecFunction(tks)
		} else if Comands.Compare(token, "MKDISK") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MKDISK <<<<<<<<<<<<<<<<<<<<")
			Comands.DataMKDISK(tks, counterDisk)
			counterDisk++
		} else if Comands.Compare(token, "RMDISK") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN RMDISK <<<<<<<<<<<<<<<<<<<<")
			Comands.RMDISK(tks)
		} else if Comands.Compare(token, "FDISK") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN FDISK <<<<<<<<<<<<<<<<<<<<")
			Comands.DataFDISK(tks)
		} else if Comands.Compare(token, "MOUNT") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MOUNT <<<<<<<<<<<<<<<<<<<<")
			Comands.DataMount(tks)
		} else if Comands.Compare(token, "UNMOUNT") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN UNMOUNT <<<<<<<<<<<<<<<<<<<<")
			Comands.DataUnMount(tks)
		} else if Comands.Compare(token, "MKFS") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MKFS <<<<<<<<<<<<<<<<<<<<")
			Comands.DataMkfs(tks)
		} else if Comands.Compare(token, "LOGIN") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN LOGIN <<<<<<<<<<<<<<<<<<<<")
			if logued {
				Comands.Error("LOGIN", "Ya hay un usuario en linea.")
				return
			} else {
				logued = Comands.DataUserLogin(tks)
			}
		} else if Comands.Compare(token, "LOGOUT") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN LOG OUT <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("LOGOUT", "Aún no se ha iniciado sesión")
				return
			} else {
				logued = Comands.LogOut()
			}
		} else if Comands.Compare(token, "MKGRP") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MKGRP <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("MKGRP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataGroup(tks, "MK")
			}
		} else if Comands.Compare(token, "RMGRP") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN RMGRP <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("RMGRP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataGroup(tks, "RM")
			}
		} else if Comands.Compare(token, "MKUSER") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MKUSER <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("MKUSER", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataUser(tks, "MK")
			}
		} else if Comands.Compare(token, "RMUSER") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN RMUSER <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("RMUSER", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataUser(tks, "RM")
			}
		} else if Comands.Compare(token, "MKDIR") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN MKDIR <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("REP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataDir(tks)
			}
		} else if Comands.Compare(token, "REP") {
			fmt.Println(">>>>>>>>>>>>>>>>>>>> FUNCIÓN REP <<<<<<<<<<<<<<<<<<<<")
			if !logued {
				Comands.Error("REP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataRep(tks)
			}
		} else {
			Comands.Error("ANALIZADOR", "NO se reconoce el comando \" "+token+"\" ")
		}
	}
}

func ExecFunction(tokens []string) {
	path := ""
	for i := 0; i < len(tokens); i++ {
		data := strings.Split(tokens[i], "=")
		if Comands.Compare(data[0], "path") {
			path = data[1]
		}
	}
	if path == "" {
		Comands.Error("EXEC", "SE requiere la \"path\" para este comando")
		return
	}
	Exec(path)
}

func Exec(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		text := fileScanner.Text()
		text = strings.TrimSpace(text)
		tk := Comand(text)
		if text != "" {
			if Comands.Compare(tk, "pause") {
				fmt.Println(">>>>>>>>>> FUNCIÓN PAUSA <<<<<<<<<<<<<<<<<<<<")
				var pause string
				Comands.Message("PAUSE", "Presion \"enter\" para continuar...")
				fmt.Scanln(&pause)
				continue
			} else if string(text[0]) == "#" {
				fmt.Println(">>>>>>>>>> FUNCIÓN COMENTARIO <<<<<<<<<<<<<<<<<<<<")
				Comands.Message("COMENTARIO", text)
				continue
			}
			text = strings.TrimLeft(text, tk)
			tokens := Separatorokens(text)
			functions(tk, tokens)
		}
	}
	if err = fileScanner.Err(); err != nil {
		log.Fatalf(
			"Error al leer el archivo: %s",
			err,
		)
	}
	file.Close()
}
