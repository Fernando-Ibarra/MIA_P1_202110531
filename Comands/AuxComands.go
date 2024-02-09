package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

func Compare(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func Error(op string, message string) {
	fmt.Println("\tERROR: " + op + "\n\tTIPO: " + message)
}

func Message(op string, message string) {
	fmt.Println("\tCOMANDO: " + op + "\n\tTIPO: " + message)
}

func Confirm(message string) bool {
	fmt.Println(message + "(y/n")
	var resp string
	fmt.Scanln(&resp)
	if Compare(resp, "y") {
		return true
	}
	return false
}

func ExistedFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func WrittingBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func readBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func readDisk(path string) *Structs.MBR {
	mbr := Structs.MBR{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &mbr)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
	}
	var mDir *Structs.MBR = &mbr
	return mDir
}
