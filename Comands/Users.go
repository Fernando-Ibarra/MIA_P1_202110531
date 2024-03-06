package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func DataUser(context []string, action string) {
	user := ""
	pass := ""
	grp := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "user") {
			user = tk[1]
		} else if Compare(tk[0], "pass") {
			pass = tk[1]
		} else if Compare(tk[0], "grp") {
			grp = tk[1]
		}
	}
	if Compare(action, "MK") {
		if user == "" || pass == "" || grp == "" {
			Error(action+"USER", "Se necesitan parámetros obligatorios para crear un usuario")
			return
		}
		if len(user) > 10 || len(pass) > 10 || len(grp) > 10 {
			Error(action+"USER", "La cantidad maxima de caracteres que se pueden usar son 10")
			return
		}
		mkuser(user, pass, grp)
	} else if Compare(action, "RM") {
		if user == "" {
			Error(action+"USER", "Se necesitan parametros obligatorios para eliminar un usuario")
			return
		}
		rmuser(user)
	} else {
		Error(action+"USER", "No se reconoce este comando")
		return
	}
}

func mkuser(user string, pass string, grp string) {
	if !Compare(Logged.User, "root") {
		Error("MKUSER", "Solo el usuario \"root\" puede acceder a estos comandos")
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("MKUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKUSER", "No se ha encontrado el disco")
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKUSER", "Error al leer el archivo")
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKUSER", "Error al leer el archivo")
		return
	}

	var fb Structs.FilesBlocks
	txt := ""
	for block := 1; block < 16; block++ {
		if inode.I_block[block-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(block-1), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)
		if err_ != nil {
			Error("MKUSER", "Error al leer el archivo")
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	exists := false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'G' || line[2] == 'g') && line[0] != '0' {
			in := strings.Split(line, ",")
			if in[2] == grp {
				exists = true
				break
			}
		}
	}

	if !exists {
		Error("MKUSER", "No se encontro el grupo \""+grp+"\".")
		return
	}

	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if line[2] == 'U' || line[2] == 'u' {
			c++
			in := strings.Split(line, ",")
			if in[3] == user {
				if line[0] != '0' {
					Error("MKUSER", "El nombre "+user+", ya esta en uso")
					return
				}
			}
		}
	}

	txt += strconv.Itoa(c+1) + ",U," + grp + "," + user + "," + pass + "\n"
	tam := len(txt)
	var cadenaS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenaS = append(cadenaS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenaS = append(cadenaS, txt)
		}
	} else {
		cadenaS = append(cadenaS, txt)
	}

	if len(cadenaS) > 16 {
		Error("MKUSER", "Se ha llenado la cantidad de archivos posibles y no se puede generar más")
		return
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKUSER", "No se ha encontrado el disco")
		return
	}

	for i := 0; i < len(cadenaS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WrittingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenaS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WrittingBytes(file, bin1.Bytes())
	}
	for i := 0; i < len(cadenaS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WrittingBytes(file, inodos.Bytes())

	Message("MKUSER", "Usuario "+user+", creado correctamente")
	file.Close()
}

func rmuser(user string) {
	if !Compare(Logged.User, "root") {
		Error("RMUSER", "Solo el usuario \"root\" puede acceder a estos comandos")
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("RMUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMUSER", "No se ha encontrado el disco")
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("RMUSER", "Error al leer el archivo")
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RMUSER", "Error al leer el archivo")
		return
	}

	var fb Structs.FilesBlocks
	txt := ""
	for block := 1; block < 16; block++ {
		if inode.I_block[block-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(block-1), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)
		if err_ != nil {
			Error("RMUSER", "Error al leer el archivo")
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	exists := false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'G' || line[2] == 'g') && line[0] != '0' {
			in := strings.Split(line, ",")
			if in[3] == user {
				exists = true
				aux += strconv.Itoa(0) + ",U," + in[2] + "," + in[3] + "," + in[4] + "\n"
				continue
			}
		}
		aux += line + "\n"
	}

	if !exists {
		Error("MKUSER", "No se encontro el usuario \""+user+"\".")
		return
	}

	txt = aux
	tam := len(txt)
	var cadenaS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenaS = append(cadenaS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenaS = append(cadenaS, txt)
		}
	} else {
		cadenaS = append(cadenaS, txt)
	}

	if len(cadenaS) > 16 {
		Error("RMUSER", "Se ha llenado la cantidad de archivos posibles y no se puede generar más")
		return
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("RMUSER", "No se ha encontrado el disco")
		return
	}

	for i := 0; i < len(cadenaS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WrittingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenaS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WrittingBytes(file, bin1.Bytes())
	}
	for i := 0; i < len(cadenaS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WrittingBytes(file, inodos.Bytes())

	Message("RMUSER", "Usuario "+user+", eliminado correctamente")
	file.Close()
}
