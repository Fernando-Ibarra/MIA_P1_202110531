package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type DiskMounted struct {
	Path       [150]byte
	State      byte
	Partitions [26]PartitionMounted
}

type PartitionMounted struct {
	Letter byte
	State  byte
	Name   [20]byte
}

var DiskMount [99]DiskMounted

var alphabet = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func DataMount(tokens []string) {
	driveLetter := ""
	name := ""
	for i := 0; i < len(tokens); i++ {
		current := tokens[i]
		command := strings.Split(current, "=")
		if Compare(command[0], "name") {
			name = command[1]
		} else if Compare(command[0], "driveletter") {
			driveLetter = "/home/fernando/Documentos/Universidad/LaboratorioArchivos/Proyectos/Proyecto1/MIA/P1/" + command[1] + ".dsk"
		}
	}
	if driveLetter == "" || name == "" {
		Error("MOUNT", "El comando MOUNT requiere parámetros obligatorios")
		return
	}
	mount(driveLetter, name)
	listMount()
}

func mount(d string, n string) {
	file, error_ := os.Open(strings.ReplaceAll(d, "\"", ""))
	if error_ != nil {
		Error("MOUNT", "No se ha podido abrir el archivo")
		return
	}
	disk := Structs.NewMBR()
	file.Seek(0, 0)

	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("MOUNT", "Error al leer el archivo")
		return
	}
	err := file.Close()
	if err != nil {
		return
	}

	partition := SearchPartitions(disk, n, d)
	if partition == nil {
		Error("MOUNT", "ERROR")
		return
	}
	if partition.Part_type == 'E' || partition.Part_type == 'L' {
		var name [16]byte
		copy(name[:], n)
		if partition.Part_name == name && partition.Part_type == 'E' {
			Error("MOUNT", "No se puede montar una partición extendida")
			return
		} else {
			ebrs := GetLogics(*partition, d)
			founded := false
			if len(ebrs) != 0 {
				for i := 0; i < len(ebrs); i++ {
					ebr := ebrs[i]
					nameEbr := ""
					for j := 0; j < len(ebr.Part_name); j++ {
						if ebr.Part_name[j] != 0 {
							nameEbr += string(ebr.Part_name[j])
						}
					}
					if Compare(nameEbr, n) && ebr.Part_mount == '1' {
						founded = true
						n = nameEbr
						break
					} else if nameEbr == n && ebr.Part_mount == '0' {
						Error("MOUNT", "No se puede montar una partición lógica eliminada")
						return
					}
				}
				if !founded {
					Error("MOUNT", "No se encontró la partición lógica")
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		var path [150]byte
		copy(path[:], d)
		if DiskMount[i].Path == path {
			for j := 0; j < 26; j++ {
				var name [20]byte
				copy(name[:], n)
				if DiskMount[i].Partitions[j].Name == name {
					Error("MOUNT", "Ya se ha montado la partición "+n)
					return
				}
				if DiskMount[i].Partitions[j].State == 0 {
					DiskMount[i].Partitions[j].State = 1
					DiskMount[i].Partitions[j].Letter = alphabet[j]
					copy(DiskMount[i].Partitions[j].Name[:], n)
					res := strconv.Itoa(i+1) + string(alphabet[j])
					Message("MOUNT", "Se ha realizado correctamente el mount -id="+res+"31")
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		if DiskMount[i].State == 0 {
			DiskMount[i].State = 1
			copy(DiskMount[i].Path[:], d)
			for j := 0; j < 26; j++ {
				if DiskMount[i].Partitions[j].State == 0 {
					DiskMount[i].Partitions[j].State = 1
					DiskMount[i].Partitions[j].Letter = alphabet[j]
					copy(DiskMount[i].Partitions[j].Name[:], n)
					res := strconv.Itoa(i+1) + string(alphabet[j])
					Message("MOUNT", "Se ha realizado correctamente el mount -id="+res+"31")
					return
				}
			}
		}
	}
}

func listMount() {
	fmt.Println("\n------------------------------ LISTADO DE MOUNTS ------------------------------")
	for i := 0; i < 99; i++ {
		for j := 0; j < 26; j++ {
			if DiskMount[i].Partitions[j].State == 1 {
				name := ""
				for k := 0; k < len(DiskMount[i].Partitions[j].Name); k++ {
					if DiskMount[i].Partitions[j].Name[k] != 0 {
						name += string(DiskMount[i].Partitions[j].Name[k])
					}
				}
				fmt.Println("\t id=" + strconv.Itoa(i+1) + string(alphabet[j]) + "31, Nombre: " + name)
			}
		}
	}
}
