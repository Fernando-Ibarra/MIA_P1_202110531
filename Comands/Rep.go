package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func DataRep(context []string) {
	name := ""
	pathOut := ""
	id := ""
	// ruta := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			id = tk[1]
		} else if Compare(tk[0], "name") {
			name = tk[1]
		} else if Compare(tk[0], "path") {
			pathOut = tk[1]
		} else if Compare(tk[0], "ruta") {
			// ruta = tk[1]
		}
	}
	if id == "" || pathOut == "" || name == "" {
		Error("REP", "Se necesitan parámetros obligatorios para el comando rep")
		return
	}

	if Compare(name, "mbr") {
		repMBR(id, pathOut)
	} else if Compare(name, "sb") {
		repSuperBlock(id, pathOut)
	} else if Compare(name, "disk") {
		repDisk(id, pathOut)
	} else if Compare(name, "bm_inode") {
		repBM(id, pathOut, "BI")
	} else if Compare(name, "bm_bloc") {
		repBM(id, pathOut, "BB")
	} else if Compare(name, "inode") {
		repInode(id, pathOut)
	} else if Compare(name, "block") {
		repBlock(id, pathOut)
	}
}

func repMBR(id string, pathOut string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}
	letter := id[0]
	currentPath, _ := os.Getwd()
	driveLetter := currentPath + "/MIA/P1/" + string(letter) + ".dsk"

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos")
		return
	}
	pd := aux[0] + ".dot"

	var partitions [4]Structs.Partition
	var logicPartitions []Structs.EBR
	mbr := readDisk(driveLetter)
	partitions[0] = mbr.Mbr_partitions_1
	partitions[1] = mbr.Mbr_partitions_2
	partitions[2] = mbr.Mbr_partitions_3
	partitions[3] = mbr.Mbr_partitions_4

	text := "digraph MBR{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\">REPORTE DE MBR</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">mbr_tamano</td><td bgcolor=\"white\">" + strconv.Itoa(int(mbr.Mbr_tamano)) + "</td></tr>\n"
	fechaC := ""
	for i := 0; i < len(mbr.Mbr_fecha_creacion); i++ {
		if mbr.Mbr_fecha_creacion[i] != 0 {
			fechaC += string(mbr.Mbr_fecha_creacion[i])
		}
	}
	text += "<tr><td bgcolor=\"thistle\">mbr_fecha_creacion</td><td bgcolor=\"thistle\">" + fechaC + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">mbr_dsk_signature</td><td bgcolor=\"white\">" + strconv.Itoa(int(mbr.Mbr_dsk_signature)) + "</td></tr>\n"
	for i := 0; i < len(partitions); i++ {
		if partitions[i].Part_type == 'E' {
			text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\">Particion</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_status</td><td bgcolor=\"white\">" + string(partitions[i].Part_status) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_type</td><td bgcolor=\"thistle\">" + string(partitions[i].Part_type) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_fit</td><td bgcolor=\"white\">" + string(partitions[i].Part_fit) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_start</td><td bgcolor=\"thistle\">" + strconv.Itoa(int(partitions[i].Part_start)) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(partitions[i].Part_s)) + "</td></tr>\n"
			partitionName := ""
			for j := 0; j < len(partitions[i].Part_name); j++ {
				if partitions[i].Part_name[j] != 0 {
					partitionName += string(partitions[i].Part_name[j])
				}
			}
			text += "<tr><td bgcolor=\"thistle\">part_name</td><td bgcolor=\"thistle\">" + partitionName + "</td></tr>\n"

			logicPartitions = GetLogics(partitions[i], driveLetter)
			for k := 0; k < len(logicPartitions); k++ {
				text += "<tr><td colspan=\"2\" bgcolor=\"salmon\">Particion Lógica - EBR</td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_mount</td><td bgcolor=\"white\">" + string(logicPartitions[k].Part_mount) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"lightsalmon\">part_fit</td><td bgcolor=\"lightsalmon\">" + string(logicPartitions[k].Part_fit) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(logicPartitions[k].Part_start)) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"lightsalmon\">part_s</td><td bgcolor=\"lightsalmon\">" + strconv.Itoa(int(logicPartitions[k].Part_s)) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_next</td><td bgcolor=\"white\">" + strconv.Itoa(int(logicPartitions[k].Part_next)) + "</td></tr>\n"
				logicPartitionName := ""

				for m := 0; m < len(logicPartitions[k].Part_name); m++ {
					if logicPartitions[k].Part_name[m] != 0 {
						logicPartitionName += string(logicPartitions[k].Part_name[m])
					}
				}
				text += "<tr><td bgcolor=\"lightsalmon\">part_name</td><td bgcolor=\"lightsalmon\">" + logicPartitionName + "</td></tr>\n"
			}

		} else {
			text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\">Particion</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_status</td><td bgcolor=\"white\">" + string(partitions[i].Part_status) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_type</td><td bgcolor=\"thistle\">" + string(partitions[i].Part_type) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_fit</td><td bgcolor=\"white\">" + string(partitions[i].Part_fit) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_start</td><td bgcolor=\"thistle\">" + strconv.Itoa(int(partitions[i].Part_start)) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(partitions[i].Part_s)) + "</td></tr>\n"
			partitionName := ""
			for j := 0; j < len(partitions[i].Part_name); j++ {
				if partitions[i].Part_name[j] != 0 {
					partitionName += string(partitions[i].Part_name[j])
				}
			}
			text += "<tr><td bgcolor=\"thistle\">part_name</td><td bgcolor=\"thistle\">" + partitionName + "</td></tr>\n"
		}
	}
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	CreateFile(pd)
	WriteFile(text, pd)
	termination := strings.Split(pathOut, ".")
	Execute(pathOut, pd, termination[1])
	Message("REP", "Reporte de MBR se ha generado correctamente en"+pathOut)
}

func repSuperBlock(id string, pathOut string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("MKGRP", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco")
		return
	}

	aux := strings.Split(path, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos")
		return
	}
	pd := aux[0] + ".dot"

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo")
		return
	}

	text := "digraph SuperBloque{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"2\" bgcolor=\"palegreen4\">REPORTE DE SuperBloque</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_filesystem_type</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_filesystem_type)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inodes_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inodes_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_blocks_count</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_blocks_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_free_inodes_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_free_inodes_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_free_blocks_count</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_free_blocks_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_mtime</td><td bgcolor=\"palegreen2\">" + string(super.S_mtime[:]) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_umtime</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_umtime)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_mnt_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_mnt_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_magic</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_magic)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inode_s</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inode_s)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_block_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_block_s)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_firts_ino</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_firts_ino)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_firts_blo</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_firts_blo)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_bm_inode_start</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_bm_inode_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_bm_block_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_bm_block_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inode_start</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inode_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_block_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_block_start)) + "</td></tr>\n"
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	file.Close()

	CreateFile(pd)
	WriteFile(text, pd)
	termination := strings.Split(pathOut, ".")
	Execute(pathOut, pd, termination[1])
	Message("REP", "Reporte de Superbloque se ha generado correctamente en"+pathOut)

}

func repDisk(id string, pathOut string) {
	var path string
	GetMount("REP", id, &path)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco")
		return
	}

	var disk Structs.MBR
	file.Seek(0, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}
	file.Close()

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos")
		return
	}
	pd := aux[0] + ".dot"

	folder := ""
	address := strings.Split(pd, "/")

	fileaux, _ := os.Open(strings.ReplaceAll(pd, "\"", ""))
	if fileaux == nil {
		for i := 0; i < len(address); i++ {
			folder += "/" + address[i]
			if _, err_2 := os.Stat(folder); os.IsNotExist(err_2) {
				os.Mkdir(folder, 0777)
			}
		}
		os.Remove(pd)
	} else {
		fileaux.Close()
	}

	partitions := GetPartitions(disk)
	var extended Structs.Partition
	ext := false
	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	content := "digraph Disk{\n"
	content += "rankdir=TB;\n"
	content += "forcelabels=true;\n"
	content += "graph [dpi = \"600\"];\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"
	content += "n1 [ label = <\n"
	content += "<table>\n"
	content += "<tr>\n"

	var positions [5]int64
	var positionsii [5]int64

	positions[0] = disk.Mbr_partitions_1.Part_start - (1 + int64(unsafe.Sizeof(Structs.MBR{})))
	positions[1] = disk.Mbr_partitions_2.Part_start - disk.Mbr_partitions_1.Part_start + disk.Mbr_partitions_1.Part_s
	positions[2] = disk.Mbr_partitions_3.Part_start - disk.Mbr_partitions_2.Part_start + disk.Mbr_partitions_2.Part_s
	positions[3] = disk.Mbr_partitions_4.Part_start - disk.Mbr_partitions_3.Part_start + disk.Mbr_partitions_3.Part_s
	positions[4] = disk.Mbr_tamano + 1 - disk.Mbr_partitions_4.Part_start + disk.Mbr_partitions_4.Part_s

	copy(positionsii[:], positions[:])

	logic := 0
	tmpLogic := ""

	if ext {
		tmpLogic += "<tr>\n"
		auxEBR := Structs.NewEBR()
		file, err = os.Open(strings.ReplaceAll(path, "\"", ""))

		if err != nil {
			Error("REP", "No se ha encontrado el disco")
			return
		}

		file.Seek(extended.Part_start, 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEBR)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}
		file.Close()

		var tamGen int64 = 0
		for auxEBR.Part_next != -1 {
			tamGen += auxEBR.Part_s
			res := float64(auxEBR.Part_s) / float64(disk.Mbr_tamano)
			res = res * 100
			tmpLogic += "<td>\"EBR\"</td>"
			s := fmt.Sprintf("%.2f", res)
			tmpLogic += "<td>\"Lógica \n " + s + "% de la partición extendida</td>\n"

			resta := float64(auxEBR.Part_next) - (float64(auxEBR.Part_start) + float64(auxEBR.Part_s))
			resta = resta / float64(disk.Mbr_tamano)
			resta = resta * 10000.00
			resta = math.Round(resta) / 100.00
			if resta != 0 {
				s = fmt.Sprintf("%f", resta)
				tmpLogic += "<td>\"Lógica\n " + s + "% libre de la partición extendida</td>\n"
				logic++
			}
			logic += 2
			file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
			if err != nil {
				Error("REP", "No se ha encontrado el disco")
				return
			}

			file.Seek(auxEBR.Part_next, 0)
			data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEBR)
			if err_ != nil {
				Error("REP", "Error al leer el archivo")
				return
			}
			file.Close()
		}
		resta := float64(extended.Part_s) - float64(tamGen)
		resta = resta / float64(disk.Mbr_tamano)
		resta = math.Round(resta * 100)
		if resta != 0 {
			s := fmt.Sprintf("%.2f", resta)
			tmpLogic += "<td>\"Libre \n " + s + "% de la partición extendida \"</td>\n"
			logic++
		}
		tmpLogic += "</tr>\n"
		logic += 2
	}
	var tamPrim int64
	for i := 0; i < 4; i++ {
		if partitions[i].Part_type == 'E' {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.3f", res)
			content += "<td COLSPAN='" + strconv.Itoa(logic) + "'>Extendida \n" + s + "% del disco</td>\n"
		} else if partitions[i].Part_start != -1 {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.3f", res)
			content += "<td ROWSPAN='2'>Primaria \n" + s + "% del disco</td>\n"
		}
	}

	if tamPrim != 0 {
		libre := disk.Mbr_tamano - tamPrim
		res := float64(libre) / float64(disk.Mbr_tamano)
		res = math.Round(res * 100)
		s := fmt.Sprintf("%.3f", res)
		content += "<td ROWSPAN='2'>Libre\n" + s + "% del disco</td>"
	}

	content += "</tr>\n"
	content += tmpLogic
	content += "</table>\n"
	content += "> ]\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	termination := strings.Split(pathOut, ".")
	Execute(pathOut, pd, termination[1])
	Message("REP", "Reporte del Disco se ha generado correctamente en"+pathOut)
}

func repBM(id string, pathOut, t string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco")
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	CreateFile(pathOut)
	content := ""
	counter := 1
	ch := '2'
	if t == "BI" {
		file.Seek(super.S_bm_inode_start, 0)
		for i := 0; i < int(super.S_inodes_count); i++ {
			data = readBytes(file, int(unsafe.Sizeof(ch)))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("REP", "Error al leer el archivo")
				return
			}

			element := fromBMtoFile(strconv.Itoa(int(ch)))
			if element == "-1" {
				break
			}
			if counter == 20 {
				content += element + "\n"
				counter = 1
			} else {
				content += element
				counter++
			}
		}
	} else {
		file.Seek(super.S_bm_block_start, 0)
		for i := 0; i < int(super.S_inodes_count); i++ {
			data = readBytes(file, int(unsafe.Sizeof(ch)))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("REP", "Error al leer el archivo")
				return
			}
			element := fromBMtoFile(strconv.Itoa(int(ch)))
			if element == "-1" {
				break
			}
			if counter == 20 {
				content += element + "\n"
				counter = 1
			} else {
				content += element
				counter++
			}
		}
	}

	WriteFile(content, pathOut)
	if t == "BI" {
		Message("REP", "Reporte de los bitmaps de Inodos "+pathOut+", creado correctamente")
	} else {
		Message("REP", "Reporte de los bitmaps de bloques "+pathOut+", creado correctamente")
	}
}

func fromBMtoFile(ch string) string {
	if ch == "48" {
		return "0"
	} else if ch == "49" {
		return "1"
	} else {
		return "-1"
	}
	return "-1"
}

func repInode(id string, pathOut string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco")
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos")
		return
	}
	pd := aux[0] + ".dot"

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	content := "digraph Inodos{\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"

	var inodes []Structs.Inodos
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	for i := 0; i < len(inodes); i++ {
		content += "A" + strconv.Itoa(i)
		content += "[label= <"
		content += "<table border=\"1\" cellborder=\"0\">\n"
		content += "<tr><td colspan=\"2\" >Inodo " + strconv.Itoa(i) + "</td></tr>\n"
		content += "<tr><td>I_uid</td><td>" + strconv.Itoa(int(inodes[i].I_uid)) + "</td></tr>\n"
		content += "<tr><td>I_gid</td><td>" + strconv.Itoa(int(inodes[i].I_gid)) + "</td></tr>\n"
		content += "<tr><td>I_s</td><td>" + strconv.Itoa(int(inodes[i].I_s)) + "</td></tr>\n"
		atime := ""
		for k := 0; k < len(inodes[i].I_atime); k++ {
			if inodes[i].I_atime[k] != 0 {
				atime += string(inodes[i].I_atime[k])
			}
		}
		content += "<tr><td>I_atime</td><td>" + atime + "</td></tr>\n"
		ctime := ""
		for k := 0; k < len(inodes[i].I_ctime); k++ {
			if inodes[i].I_ctime[k] != 0 {
				ctime += string(inodes[i].I_ctime[k])
			}
		}
		content += "<tr><td>I_ctime</td><td>" + ctime + "</td></tr>\n"
		mtime := ""
		for k := 0; k < len(inodes[i].I_mtime); k++ {
			if inodes[i].I_mtime[k] != 0 {
				mtime += string(inodes[i].I_mtime[k])
			}
		}
		content += "<tr><td>I_mtime</td><td>" + mtime + "</td></tr>\n"
		for j := 0; j < len(inodes[i].I_block); j++ {
			content += "<tr><td>I_block " + strconv.Itoa(i+1) + " </td><td>" + strconv.Itoa(int(inodes[i].I_block[j])) + "</td></tr>\n"
		}
		content += "<tr><td>I_type</td><td>" + strconv.Itoa(int(inodes[i].I_type)) + "</td></tr>\n"
		content += "<tr><td>I_perm</td><td>" + strconv.Itoa(int(inodes[i].I_perm)) + "</td></tr>\n"
		content += "</table>\n"
		content += ">]"
	}

	content += "\n"

	for i := 0; i < len(inodes); i++ {
		if i == 0 {
			content += "A" + strconv.Itoa(i)
		} else {
			content += " -> " + "A" + strconv.Itoa(i)
		}
	}

	content += "\n"
	content += "{ rank=same "
	for i := 0; i < len(inodes); i++ {
		content += "A" + strconv.Itoa(i) + " "
	}
	content += "}"

	content += "\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	termination := strings.Split(pathOut, ".")
	Execute(pathOut, pd, termination[1])
	Message("REP", "Reporte de Inodos se ha generado correctamente en"+pathOut)
}

func repBlock(id string, pathOut string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}

	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco")
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos")
		return
	}
	pd := aux[0] + ".dot"

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	var inodes []Structs.Inodos
	inode = Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	counter := 0
	content := "digraph Bloques{\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"

	file.Seek(super.S_inode_start, 0)
	for v := 0; v < len(inodes); v++ {
		file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(v), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("MKDIR", "Error al leer el archivo")
			return
		}
		if inode.I_type == 0 {
			for i := 0; i < 16; i++ {
				if i < 16 {
					if inode.I_block[i] != -1 {
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)

						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("MKDIR", "Error al leer el archivo")
							return
						}

						content += "\nA" + strconv.Itoa(counter)
						content += "[label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td colspan=\"2\" >Bloque Carpeta " + strconv.Itoa(counter) + "</td></tr>\n"
						content += "<tr><td>B_name</td><td>B_Inodo</td></tr>\n"
						for j := 0; j < 4; j++ {
							name := ""
							for nam := 0; nam < len(folder.B_content[j].B_name); nam++ {
								if folder.B_content[j].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[j].B_name[nam])
							}
							content += "<tr><td>" + name + "</td><td>" + strconv.Itoa(int(folder.B_content[j].B_inodo)) + "</td></tr>\n"
						}
						content += "</table>\n"
						content += ">]"
						counter++
					}
				}
			}
		} else if inode.I_type == 1 {
			for i := 0; i < 16; i++ {
				if i < 16 {
					if inode.I_block[i] != -1 {
						var folderAux Structs.FilesBlocks
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i+1), 0)

						data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folderAux)
						if err_ != nil {
							Error("MKDIR", "Error al leer el archivo")
							return
						}

						content += "\nA" + strconv.Itoa(counter)
						content += "[label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td> Bloque Archivo " + strconv.Itoa(counter) + "</td></tr>\n"
						folderContent := ""
						for k := 0; k < len(folderAux.B_content); k++ {
							if folderAux.B_content[k] == 0 {
								continue
							}
							regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúüñ,]+$`)
							if !regex.MatchString(string(folderAux.B_content[k])) {
								continue
							} else {
								folderContent += string(folderAux.B_content[k])
							}
						}
						content += "<tr><td>" + folderContent + "</td></tr>\n"
						content += "</table>\n"
						content += ">]"
						counter++
					}
				}
			}
		} else {
			continue
		}

	}

	content += "\n"

	for i := 0; i < counter; i++ {
		if i == 0 {
			content += "A" + strconv.Itoa(i)
		} else {
			content += " -> " + "A" + strconv.Itoa(i)
		}
	}

	content += "\n"
	content += "{ rank=same "
	for i := 0; i < counter; i++ {
		content += "A" + strconv.Itoa(i) + " "
	}
	content += "}"
	content += "\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	termination := strings.Split(pathOut, ".")
	Execute(pathOut, pd, termination[1])
	Message("REP", "Reporte de Bloques se ha generado correctamente en"+pathOut)
}
