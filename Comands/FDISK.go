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

type Transition struct {
	partition int
	start     int
	end       int
	before    int
	after     int
}

var startValue int

func DataFDISK(tokens []string) {
	size := ""
	unit := "K"
	driveLetter := ""
	tipo := "P"
	fit := "WF"
	name := ""
	// add := ""
	// deleteP := ""
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "size") {
			size = tk[1]
		} else if Compare(tk[0], "unit") {
			unit = tk[1]
		} else if Compare(tk[0], "driveletter") {
			driveLetter = "/home/fernando/Documentos/Universidad/LaboratorioArchivos/Proyectos/Proyecto1/MIA/P1/" + tk[1] + ".dsk"
		} else if Compare(tk[0], "type") {
			tipo = tk[1]
		} else if Compare(tk[0], "fit") {
			fit = tk[1]
		} else if Compare(tk[0], "name") {
			name = tk[1]
		} else if Compare(tk[0], "delete") {
			// deleteP = tk[1]
		} else if Compare(tk[0], "add") {
			// add = tk[1]
		}
	}
	if size == "" || driveLetter == "" || name == "" {
		Error("FDISK", "EL comando FDISK necesita parámetros obligatorios")
		return
	} else {
		generatePartition(size, unit, driveLetter, tipo, fit, name)
	}
}

func generatePartition(s string, u string, d string, t string, f string, n string) {
	startValue = 0
	i, error_ := strconv.Atoi(s)
	if error_ != nil {
		Error("FDISK", "Size debe ser un número entero")
		return
	}
	if i <= 0 {
		Error("FDISK", "Size debe ser mayor que 0")
		return
	}
	if Compare(u, "b") || Compare(u, "k") || Compare(u, "m") {
		if Compare(u, "k") {
			i = i * 1024
		} else if Compare(u, "m") {
			i = i * 1024 * 1024
		}
	} else {
		Error("FDISK", "Unit no contiene los valores esperados")
		return
	}
	if !(Compare(t, "p") || Compare(t, "e") || Compare(t, "l")) {
		Error("FDISK", "Type no contiene los valores esperados")
		return
	}
	if !(Compare(f, "bf") || Compare(f, "ff") || Compare(f, "wf")) {
		Error("FDISK", "Fit no contiene los valores esperados")
	}
	mbr := readDisk(d)
	partitions := GetPartitions(*mbr)
	var between []Transition

	used := 0
	ext := 0
	c := 0

	base := int(unsafe.Sizeof(Structs.MBR{}))
	extended := Structs.NewPartition()

	for j := 0; j < len(partitions); j++ {
		prttn := partitions[j]
		if prttn.Part_status == '1' {
			var trn Transition
			trn.partition = c
			trn.start = int(prttn.Part_start)
			trn.end = int(prttn.Part_start + prttn.Part_s)
			trn.before = trn.start - base
			base = trn.end
			if used != 0 {
				between[used-1].after = trn.start - (between[used-1].end)
			}
			between = append(between, trn)
			used++
			if prttn.Part_type == "e"[0] || prttn.Part_type == "E"[0] {
				ext++
				extended = prttn
			}
		}
		if used == 4 && !Compare(t, "l") {
			Error("FDISK", "Límite de particioens alcanzado")
			return
		} else if ext == 1 && Compare(t, "e") {
			Error("FDISK", "Solo se puede crear una partición extendida")
			return
		}
		c++
	}
	if ext == 0 && Compare(t, "l") {
		Error("FDISK", "Aún no se han creado particiones extendidas, no se puede agregar una lógica")
		return
	}

	if used != 0 {
		between[len(between)-1].after = int(mbr.Mbr_tamano) - between[len(between)-1].end
	}

	comeBack := SearchPartitions(*mbr, n, d)
	if comeBack != nil {
		Error("FDISK", "El nombre "+n+", ya está en uso")
		return
	}
	temporal := Structs.NewPartition()
	temporal.Part_status = '1'
	temporal.Part_s = int64(i)
	temporal.Part_type = strings.ToUpper(t)[0]
	temporal.Part_fit = strings.ToUpper(f)[0]
	copy(temporal.Part_name[:], n)

	if Compare(t, "L") {
		Logic(temporal, extended, d, n)
		return
	}

	mbr = fit(*mbr, temporal, between, partitions, used)
	if mbr == nil {
		return
	}

	file, err := os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
	}
	file.Seek(0, 0)
	var binary2 bytes.Buffer
	binary.Write(&binary2, binary.BigEndian, mbr)
	WrittingBytes(file, binary2.Bytes())
	if Compare(t, "E") {
		ebr := Structs.NewEBR()
		ebr.Part_mount = '0'
		ebr.Part_start = int64(startValue)
		ebr.Part_s = 0
		ebr.Part_next = -1

		file.Seek(int64(startValue), 0)
		var binary3 bytes.Buffer
		binary.Write(&binary3, binary.BigEndian, ebr)
		WrittingBytes(file, binary3.Bytes())
		Message("FDISK", "Partición Extendida: "+n+", creada correctamente")
		return
	}
	file.Close()
	Message("FDISK", "Partición Primaria: "+n+", creada correctamente")
}

func GetPartitions(disk Structs.MBR) []Structs.Partition {
	var v []Structs.Partition
	v = append(v, disk.Mbr_partitions_1)
	v = append(v, disk.Mbr_partitions_2)
	v = append(v, disk.Mbr_partitions_3)
	v = append(v, disk.Mbr_partitions_4)
	return v
}

func SearchPartitions(mbr Structs.MBR, name string, path string) *Structs.Partition {
	var partitions [4]Structs.Partition
	partitions[0] = mbr.Mbr_partitions_1
	partitions[1] = mbr.Mbr_partitions_2
	partitions[2] = mbr.Mbr_partitions_3
	partitions[3] = mbr.Mbr_partitions_4

	ext := false
	extended := Structs.NewPartition()
	for i := 0; i < len(partitions); i++ {
		partition := partitions[i]
		if partition.Part_status == "1"[0] {
			nameP := ""
			for j := 0; j < len(partition.Part_name); j++ {
				if partition.Part_name[j] != 0 {
					nameP += string(partition.Part_name[j])
				}
			}
			if Compare(nameP, name) {
				return &partition
			} else if partition.Part_type == "E"[0] || partition.Part_type == "e"[0] {
				ext = true
				extended = partition
			}
		}
	}

	if ext {
		ebrs := GetLogics(extended, path)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_mount == '1' {
				nameE := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					nameE += string(ebr.Part_name[j])
				}
				if Compare(nameE, name) {
					tmp := Structs.NewPartition()
					tmp.Part_status = '1'
					tmp.Part_type = 'L'
					tmp.Part_fit = ebr.Part_fit
					tmp.Part_start = ebr.Part_start
					tmp.Part_s = ebr.Part_s
					copy(tmp.Part_name[:], ebr.Part_name[:])
					return &tmp
				}
			}
		}
	}
	return nil
}

func GetLogics(partition Structs.Partition, path string) []Structs.EBR {
	var ebrs []Structs.EBR
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	tmp := Structs.NewEBR()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return nil
	}
	for {
		if int(tmp.Part_next) != -1 && int(tmp.Part_mount) != 0 {
			ebrs = append(ebrs, tmp)
			file.Seek(tmp.Part_next, 0)
			data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &tmp)
			if err_ != nil {
				Error("FDISK", "Error al leer el archivo")
				return nil
			}
		} else {
			file.Close()
			break
		}
	}
	return ebrs
}

func Logic(p Structs.Partition, e Structs.Partition, d string, n string) {
	logic := Structs.NewEBR()
	logic.Part_mount = '1'
	logic.Part_fit = p.Part_fit
	logic.Part_s = p.Part_s
	logic.Part_next = -1
	copy(logic.Part_name[:], p.Part_name[:])

	file, err := os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
	}

	file.Seek(0, 0)
	tmp := Structs.NewEBR()
	file.Seek(e.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return
	}
	size := int64(0)
	for {
		size += int64(unsafe.Sizeof(tmp)) + tmp.Part_s
		if tmp.Part_mount == '0' && tmp.Part_next == -1 {
			logic.Part_start = tmp.Part_start
			logic.Part_next = logic.Part_start + logic.Part_s + int64(unsafe.Sizeof(tmp))
			if e.Part_s-size <= logic.Part_s {
				Error("FDISK", "No hay espacio para más particiones lógicas")
			}
			file.Seek(logic.Part_start, 0)
			var binary2 bytes.Buffer
			binary.Write(&binary2, binary.BigEndian, logic)
			WrittingBytes(file, binary2.Bytes())
			addLogic := Structs.NewEBR()
			addLogic.Part_mount = '0'
			addLogic.Part_next = -1
			addLogic.Part_start = logic.Part_next
			file.Seek(addLogic.Part_start, 0)
			var binary3 bytes.Buffer
			binary.Write(&binary3, binary.BigEndian, addLogic)
			WrittingBytes(file, binary3.Bytes())
			Message("FDISK", "Partición Lógica: "+n+", creada correctamente")
			return
		}
	}
	file.Close()
}

func fit(mbr Structs.MBR, t Structs.Partition, b Transition, p []Structs.Partition, u int) *Structs.MBR {
	if u == 0 {

	}
}
