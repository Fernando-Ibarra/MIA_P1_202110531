package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strings"
	"unsafe"
)

func DataDir(context []string) {
	rBoolean := false
	path := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "r") {
			rBoolean = true
		}
	}
	if path == "" {
		Error("MKDIR", "Se necesitan parametros obligatorios para crear una carpeta")
		return
	} else if path != "" {
		mkdir(path, rBoolean)
	} else {
		Error("MKDIR", "No se reconoce este comando")
		return
	}
}

func mkdir(path string, rBoolean bool) {
	folders := strings.Split(path, "/")
	folders = folders[1:]
	id := Logged.Id
	var driveletter string
	partition := GetMount("MKDIR", id, &driveletter)
	if string(partition.Part_status) == "0" {
		Error("MKDIR", "No se encontro la partición montada con el id: "+id)
		return
	}
	file, err := os.Open(strings.ReplaceAll(driveletter, "\"", ""))
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco")
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo")
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo")
		return
	}

	currentFolder := Structs.NewDirectoriesBlocks()
	space := false
	size := int(super.S_inodes_count)
	for j := 0; j < size; j++ {
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+(int64(j)*int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))), 0)
		if j > 1 {
			newFolder := Structs.NewDirectoriesBlocks()
			copy(newFolder.B_content[0].B_name[:], ".")
			newFolder.B_content[0].B_inodo = 0
			copy(newFolder.B_content[1].B_name[:], "..")
			newFolder.B_content[1].B_inodo = 0
			copy(newFolder.B_content[2].B_name[:], folders[0])
			newFolder.B_content[2].B_inodo = 1
			data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &currentFolder)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo")
				return
			}
			for i := 0; i < len(currentFolder.B_content); i++ {
				if currentFolder.B_content[2].B_inodo != -1 {
					space = true
					folders = folders[1:]
					break
				}
			}
			if space {
				var binFolder bytes.Buffer
				binary.Write(&binFolder, binary.BigEndian, newFolder)
				WrittingBytes(file, binFolder.Bytes())
				folderName := ""
				for m := 0; m < len(newFolder.B_content[2].B_name); m++ {
					if newFolder.B_content[2].B_name[m] != 0 {
						folderName += string(newFolder.B_content[2].B_name[m])
					}
				}
				Message("MKDIR", "Carpeta "+folderName+", se ha creado correctamente")
				break
			}
		}

	}

}
