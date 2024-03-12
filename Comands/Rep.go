package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
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
	}
}

// PENDIENTE CREAR ARCHIVO JGP
func repMBR(id string, pathOut string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido")
		return
	}
	letter := id[0]
	driveLetter := "/home/fernando/Documentos/Universidad/LaboratorioArchivos/Proyectos/Proyecto1/MIA/P1/" + string(letter) + ".dsk"

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
	// fmt.Println(text)
	// CreateFile(pathOut)
	// WriteFile(text, pathOut)
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
	// fmt.Println(text)
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

	aux := strings.Split(path, ".")
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
		tmpLogic += "</tr>\n\n"
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

	fmt.Println(content)
	/*
		b := []byte(content)
		err_ = ioutil.WriteFile(pd, b, 0644)
		if err_ != nil {
			log.Fatal(err_)
		}

		termination := strings.Split(pathOut, ".")
		path2, _ := exec.LookPath("dot")
		cmd, _ := exec.Command(path2, "-T"+termination[1], pd).Output()
		mode := int(0777)
		ioutil.WriteFile(pathOut, cmd, os.FileMode(mode))
		disco := strings.Split(path, "/")
		Message("REP", "Reporte tipo DISK del disco "+disco[len(disco)-1]+", creado correctamente")

	*/
}
