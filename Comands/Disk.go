package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

func DataMKDISK(tokens []string, counterDisks int, counterDisksR *int) {
	*counterDisksR = counterDisks
	size := ""
	fit := ""
	unit := ""
	currentPath, _ := os.Getwd()
	path := currentPath + "/MIA/P1/" + string(getNameDisk(counterDisks)) + ".dsk"
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "fit") {
			if fit == "" {
				fit = tk[1]
			} else {
				Error("MKDISK", "Parametro fit repetido en el comando"+tk[0])
				counterDisks--
				*counterDisksR = counterDisks
				return
			}
		} else if Compare(tk[0], "size") {
			if size == "" {
				size = tk[1]
			} else {
				Error("MKDISK", "Parametro sizse repetido en el comendo"+tk[0])
				counterDisks--
				*counterDisksR = counterDisks
				return
			}
		} else if Compare(tk[0], "unit") {
			if unit == "" {
				unit = tk[1]
			} else {
				Error("MKDISK", "Parametro unit repetido en el comendo"+tk[0])
				counterDisks--
				*counterDisksR = counterDisks
				return
			}
		} else {
			Error("MKDISK", "No se esperaba el parametro "+tk[0])
			counterDisks--
			*counterDisksR = counterDisks
			error_ = true
			return
		}
	}
	if fit == "" {
		fit = "FF"
	}

	if unit == "" {
		unit = "M"
	}

	if error_ {
		return
	}

	if size == "" {
		Error("MKDISK", "Se requiere párametro Size para este comando de forma obligatoria")
		return
	} else if !Compare(fit, "BF") && !Compare(fit, "FF") && !Compare(fit, "WF") {
		Error("MKDISK", "Se obtuvo un valor de fit no esperado")
		return
	} else if !Compare(unit, "k") && !Compare(unit, "m") {
		Error("MKDISK", "Se obtuvo un valor de unit no esperado")
		return
	} else {
		makeFile(size, fit, unit, path)
	}
}

func makeFile(s string, f string, u string, path string) {
	var disk = Structs.NewMBR()
	size, err := strconv.Atoi(s)
	if err != nil {
		Error("MKDISK", "Size debe ser un número entero")
		return
	}
	if size <= 0 {
		Error("MKDISK", "Size debe ser mayor a 0")
		return
	}
	if Compare(u, "M") {
		size = 1024 * 1024 * size
	} else if Compare(u, "K") {
		size = 1024 * size
	}
	f = string(f[0])
	disk.Mbr_tamano = int64(size)
	fecha := time.Now().String()
	copy(disk.Mbr_fecha_creacion[:], fecha)
	aleatorio, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	entero, _ := strconv.Atoi(aleatorio.String())
	disk.Mbr_dsk_signature = int64(entero)
	copy(disk.Dsk_fit[:], string(f[0]))
	disk.Mbr_partitions_1 = Structs.NewPartition()
	disk.Mbr_partitions_2 = Structs.NewPartition()
	disk.Mbr_partitions_3 = Structs.NewPartition()
	disk.Mbr_partitions_4 = Structs.NewPartition()

	if ExistedFile(path) {
		_ = os.Remove(path)
	}

	if !strings.HasSuffix(path, "dsk") {
		Error("MKDISK", "Extensión de archivos no válida")
		return
	}

	folder := ""
	address := strings.Split(path, "/")
	for i := 0; i < len(address)-1; i++ {
		folder += "/" + address[i]
		if _, err_ := os.Stat(folder); os.IsNotExist(err_) {
			os.Mkdir(folder, 0777)
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		Error("MKDISK", "No se pudo crear el disco")
		return
	}
	var empty int8 = 0
	s1 := &empty
	var num int64 = 0
	num = int64(size)
	num = num - 1
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s1)
	WritingBytes(file, binario.Bytes())

	file.Seek(num, 0)

	var binar2 bytes.Buffer
	binary.Write(&binar2, binary.BigEndian, s1)
	WritingBytes(file, binar2.Bytes())

	file.Seek(0, 0)
	disk.Mbr_tamano = num + 1

	var binar3 bytes.Buffer
	binary.Write(&binar3, binary.BigEndian, disk)
	WritingBytes(file, binar3.Bytes())
	file.Close()
	nameDisk := strings.Split(path, "/")
	Message("MKDISK", "¡DISCO "+nameDisk[len(nameDisk)-1]+" CREADO EXITOSAMENTE!")
}

func getNameDisk(number int) string {
	if number <= 26 {
		return string(rune('A' - 1 + number))
	}
	firstLetter := 'A' + (number-1)/26 - 1
	secondLetter := 'A' + (number-1)%26
	if secondLetter == 'A'-1 {
		secondLetter = 'Z'
		firstLetter++
	}
	return string(rune(firstLetter)) + string(rune(secondLetter))
}

func RMDISK(tokens []string) {
	if len(tokens) > 1 {
		Error("RMDISK", "Solo se acepta el párametro driveletter")
		return
	}
	driveLetter := ""
	nameDisk := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "driveletter") {
			if driveLetter == "" {
				currentPath, _ := os.Getwd()
				driveLetter = currentPath + "/MIA/P1/" + tk[1] + ".dsk"
				nameDisk = tk[1]
			} else {
				Error("RMDISK", "Parámetro driveletter repetido en el comando: "+tk[0])
			}
		} else {
			Error("RMDISK", "No se esperaba el parámetro "+tk[0])
			error_ = false
			return
		}
	}
	if error_ {
		return
	}
	if driveLetter == "" {
		Error("RMDISK", "Se requiere el parámetro path")
		return
	} else {
		if !ExistedFile(driveLetter) {
			Error("RMDISK", "No se encontró el disco en la ruta indicada")
			return
		}
		if !strings.HasSuffix(driveLetter, "dsk") {
			Error("RMDISK", "Extensión de archivo no válida")
			return
		}
		if Confirm("¿Desea eliminar el disco " + nameDisk + ".dsk ?") {
			err := os.Remove(driveLetter)
			if err != nil {
				Error("RMDISK", "Error al intentar eliminar el archivo.")
				return
			}
			Message("RMDISK", "Disco ubicado en "+driveLetter+", ha sido eliminado correctamente")
			return
		} else {
			Message("RMDISK", "Eliminación del disco "+driveLetter+", cancelada exitosamente.")
		}
	}
}
