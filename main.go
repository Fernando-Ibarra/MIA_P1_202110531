package main

import (
	"Proyecto1/Comands"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var CounterDisk = 1

var logued = false

func main() {
	for true {
		fmt.Println("..................... INGRESE UN COMANDO .....................")
		fmt.Println("-----> Si desea terminar con la aplicación ingrese \"exit\"")
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
			fmt.Println("")
			fmt.Println("-------------------> COMANDO EXEC <-------------------")
			ExecFunction(tks)
		} else if Comands.Compare(token, "MKDISK") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKDISK <-------------------")
			Comands.DataMKDISK(tks, CounterDisk, &CounterDisk)
			CounterDisk++
		} else if Comands.Compare(token, "RMDISK") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO RMDISK <-------------------")
			Comands.RMDISK(tks)
		} else if Comands.Compare(token, "FDISK") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO FDISK <-------------------")
			Comands.DataFDISK(tks)
		} else if Comands.Compare(token, "MOUNT") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MOUNT <-------------------")
			Comands.DataMount(tks)
		} else if Comands.Compare(token, "UNMOUNT") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO UNMOUNT <-------------------")
			Comands.DataUnMount(tks)
		} else if Comands.Compare(token, "MKFS") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKFS <-------------------")
			Comands.DataMkfs(tks)
		} else if Comands.Compare(token, "LOGIN") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO LOGIN <-------------------")
			if logued {
				Comands.Error("LOGIN", "Ya hay un usuario en linea.")
				return
			} else {
				logued = Comands.DataUserLogin(tks)
			}
		} else if Comands.Compare(token, "LOGOUT") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO LOG OUT <-------------------")
			if !logued {
				Comands.Error("LOGOUT", "Aún no se ha iniciado sesión")
				return
			} else {
				logued = Comands.LogOut()
			}
		} else if Comands.Compare(token, "MKGRP") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKGRP <-------------------")
			if !logued {
				Comands.Error("MKGRP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataGroup(tks, "MK")
			}
		} else if Comands.Compare(token, "RMGRP") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO RMGRP <-------------------")
			if !logued {
				Comands.Error("RMGRP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataGroup(tks, "RM")
			}
		} else if Comands.Compare(token, "CHGRP") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO CHGRP <-------------------")
			if !logued {
				Comands.Error("CHGRP", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataUser(tks, "CH")
			}
		} else if Comands.Compare(token, "MKUSR") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKUSR <-------------------")
			if !logued {
				Comands.Error("MKUSER", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataUser(tks, "MK")
			}
		} else if Comands.Compare(token, "RMUSR") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO RMUSR <-------------------")
			if !logued {
				Comands.Error("RMUSER", "Aún no se ha iniciado sesión")
				return
			} else {
				Comands.DataUser(tks, "RM")
			}
		} else if Comands.Compare(token, "MKDIR") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKDIR <-------------------")
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p)
				Comands.DataDir(tks, partition, p)
			}
		} else if Comands.Compare(token, "MKFILE") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MKFILE <-------------------")
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p)
				Comands.DataFile(tks, partition, p)
			}
		} else if Comands.Compare(token, "CAT") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO CAT <-------------------")
			if !logued {
				Comands.Error("CAT", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("CAT", Comands.Logged.Id, &p)
				Comands.DataCat(tks, partition, p)
			}
		} else if Comands.Compare(token, "CHMOD") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO CHMOD <-------------------")
			if !logued {
				Comands.Error("CHMOD", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("CHMOD", Comands.Logged.Id, &p)
				Comands.DataChmod(tks, partition, p)
			}
		} else if Comands.Compare(token, "CHOWN") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO CHOWN <-------------------")
			if !logued {
				Comands.Error("CHOWN", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p)
				Comands.DataChown(tks, partition, p)
			}
		} else if Comands.Compare(token, "RENAME") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO RENAME <-------------------")
			if !logued {
				Comands.Error("RENAME", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p)
				Comands.DataRename(tks, partition, p)
			}
		} else if Comands.Compare(token, "MOVE") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO MOVE <-------------------")
			if !logued {
				Comands.Error("MOVE", "Aún no se ha iniciado sesión")
				return
			} else {
				var p string
				partition := Comands.GetMount("MOVE", Comands.Logged.Id, &p)
				Comands.DataMove(tks, partition, p)
			}
		} else if Comands.Compare(token, "REP") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO REP <-------------------")
			Comands.DataRep(tks)
		} else if Comands.Compare(token, "JSON") {
			fmt.Println("")
			fmt.Println("-------------------> COMANDO JSON <-------------------")
			response := Comands.MakeJson()
			fmt.Println(response)
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
