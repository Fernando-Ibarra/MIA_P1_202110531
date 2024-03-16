package Comands

import (
	"Proyecto1/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func DataFile(context []string, partition Structs.Partition, pth string) {
	rBoolean := false
	size := ""
	path := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "r") {
			rBoolean = true
		} else if Compare(tk[0], "size") {
			size = tk[1]
		}
	}
	size = "0"
	if path == "" {
		Error("MKFILE", "Se necesitan parametros obligatorios para crear un directorio")
		return
	}
	tmp := GetPath(path)
	mkfile(tmp, rBoolean, partition, pth, size)
}

func mkfile(path []string, r bool, partition Structs.Partition, pth string, s string) {
	copyPath := path
	size, err := strconv.Atoi(s)
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	content := ""
	for i := 0; i < size; i++ {
		content += strconv.Itoa(i)
	}

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco")
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo")
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo")
		return
	}

	var newf string
	if len(path) == 0 {
		Error("MKDIR", "No se ha brindado un path valido")
		return
	}

	var past int64
	var bi int64
	var bb int64
	fnd := false
	inodetmp := Structs.NewInodos()
	foldertmp := Structs.NewDirectoriesBlocks()

	inodeFile := Structs.NewInodos()
	fileBlock := Structs.FilesBlocks{}

	newf = path[len(path)-1]
	var father int64

	var aux []string
	for i := 0; i < len(path); i++ {
		aux = append(aux, path[i])
	}
	path = aux
	var stack string
	fileWritten := false
	fatherSpace := false

	for v := 0; v < len(path)-1; v++ {
		fnd = false
		for i := 0; i < 16; i++ {
			if i < 16 {
				if inode.I_block[i] != -1 {
					folder = Structs.NewDirectoriesBlocks()
					file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)

					data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &folder)
					if err_ != nil {
						Error("MKDIR", "Error al leer el archivo")
						return
					}

					for j := 0; j < 4; j++ {
						nameFolder := ""
						for nam := 0; nam < len(folder.B_content[j].B_name); nam++ {
							if folder.B_content[j].B_name[nam] == 0 {
								continue
							}
							nameFolder += string(folder.B_content[j].B_name[nam])
						}
						if Compare(nameFolder, path[v]) {
							stack += "/" + path[v]
							fnd = true
							father = folder.B_content[j].B_inodo
							inode = Structs.NewInodos()
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*folder.B_content[j].B_inodo, 0)

							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inode)
							if err_ != nil {
								Error("MKDIR", "Error al leer el archivo")
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("MKDIR", "No tiene permisos para crear carpetas en este directorio")
								return
							}

							break
						}
					}

				} else {
					break
				}
			}
		}
		if !fnd {
			if r {
				stack += "/" + path[v]
				mkdir(GetPath(stack), false, partition, pth)
				file.Seek(spr.S_inode_start, 0)

				data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &inode)
				if err_ != nil {
					Error("MKDIR", "Error al leer el archivo")
					return
				}

				if v == len(path)-2 {
					stack += "/" + path[v+1]
					mkdir(GetPath(stack), false, partition, pth)
					file.Seek(spr.S_inode_start, 0)
					data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &inode)
					if err_ != nil {
						Error("MKDIR", "Error al leer el archivo")
						return
					}
					return
				}
			} else {
				address := ""
				for i := 0; i < len(path); i++ {
					address += "/" + path[i]
				}
				Error("MKDIR", "No se pudo crear el directorio "+address+", no existen directorios")
				return
			}
		}
	}
	/*
		Por si el padre tiene una carpeta donde hay un espacio libre para
	*/
	fnd = false
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			if i < 16 {
				folderAux := folder
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("MKDIR", "Error al leer el archivo")
					return
				}
				nameAux1 := ""
				for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
					if folder.B_content[2].B_name[nam] == 0 {
						continue
					}
					nameAux1 += string(folder.B_content[2].B_name[nam])
				}

				nameAux2 := ""
				for nam := 0; nam < len(folderAux.B_content[2].B_name); nam++ {
					if folderAux.B_content[2].B_name[nam] == 0 {
						continue
					}
					nameAux2 += string(folderAux.B_content[2].B_name[nam])
				}
				padre := ""
				for k := 0; k < len(path); k++ {
					if k >= 1 {
						padre = path[k-1]
					}
				}
				if padre == nameAux1 {
					continue
				}
				for j := 0; j < 4; j++ {
					if folder.B_content[j].B_inodo == -1 {
						past = inode.I_block[i]
						bi = GetFree(spr, pth, "BI")
						if bi == -1 {
							Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad")
							return
						}
						bb = GetFree(spr, pth, "BB")
						if bb == -1 {
							Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad")
							return
						}

						if strings.Contains(newf, ".") {
							inodetmp.I_uid = int64(Logged.Uid)
							inodetmp.I_gid = int64(Logged.Gid)
							inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

							dateNow := time.Now().String()
							copy(inodetmp.I_atime[:], spr.S_mtime[:])
							copy(inodetmp.I_ctime[:], dateNow)
							copy(inodetmp.I_mtime[:], dateNow)
							inodetmp.I_type = 1
							inodetmp.I_perm = 664
							inodetmp.I_block[0] = bb

							folder.B_content[j].B_inodo = bi
							copy(folder.B_content[j].B_name[:], newf)

							inodeFile.I_uid = int64(Logged.Uid)
							inodeFile.I_gid = int64(Logged.Gid)
							inodeFile.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

							dateNow2 := time.Now().String()
							copy(inodeFile.I_atime[:], spr.S_mtime[:])
							copy(inodeFile.I_ctime[:], dateNow2)
							copy(inodeFile.I_mtime[:], dateNow2)
							inodeFile.I_type = 0
							inodeFile.I_perm = 664
							inodeFile.I_block[0] = bb + 2

							fileBlock = Structs.FilesBlocks{}

							copy(fileBlock.B_content[:], content)

							fileWritten = true
							fatherSpace = true
							fnd = true
							i = 20
							break
						} else {
							inodetmp.I_uid = int64(Logged.Uid)
							inodetmp.I_gid = int64(Logged.Gid)
							inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

							dateNow := time.Now().String()
							copy(inodetmp.I_atime[:], spr.S_mtime[:])
							copy(inodetmp.I_ctime[:], dateNow)
							copy(inodetmp.I_mtime[:], dateNow)
							inodetmp.I_type = 0
							inodetmp.I_perm = 664
							inodetmp.I_block[0] = bb

							copy(foldertmp.B_content[0].B_name[:], ".")
							foldertmp.B_content[0].B_inodo = bi
							copy(foldertmp.B_content[1].B_name[:], "..")
							foldertmp.B_content[1].B_inodo = father
							copy(foldertmp.B_content[2].B_name[:], "-")
							copy(foldertmp.B_content[3].B_name[:], "-")

							folder.B_content[j].B_inodo = bi
							copy(folder.B_content[j].B_name[:], newf)
							fnd = true
							i = 20
							break
						}
					}
				}
			}
		} else {
			break
		}
	}
	/*
	 Encontrando un espacio donde se puede escribir el bloque/carpeta si hay espacio libre en un carpeta padre. Se usa un nuevo inodo
	*/
	if !fnd {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] == -1 {
				if i < 16 {
					bi = GetFree(spr, pth, "BI")
					if bi == -1 {
						Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad")
						return
					}
					past = GetFree(spr, pth, "BB")
					if past == -1 {
						Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad")
						return
					}

					bb = GetFree(spr, pth, "BB")

					if strings.Contains(newf, ".") {
						fileBlock = Structs.FilesBlocks{}

						tam := len(content)
						var contents []string
						if tam > 10 {
							for tam > 10 {
								auxFile := ""
								for j := 0; j < 10; j++ {
									auxFile += string(content[i])
								}
								contents = append(contents, auxFile)
								content = strings.ReplaceAll(content, auxFile, "")
								tam = len(content)
							}
							if tam < 10 && tam != 0 {
								contents = append(contents, content)
							}
						} else {
							contents = append(contents, content)
						}

						if len(contents) > 16 {
							Error("MKFILE", "Se ha llenado la cantidad de archivo posibles y no es posible generar más")
							return
						}

						inodeFile.I_uid = int64(Logged.Uid)
						inodeFile.I_gid = int64(Logged.Gid)
						inodeFile.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

						dateNow := time.Now().String()
						copy(inodeFile.I_atime[:], spr.S_mtime[:])
						copy(inodeFile.I_ctime[:], dateNow)
						copy(inodeFile.I_mtime[:], dateNow)
						inodeFile.I_type = 0
						inodeFile.I_perm = 664
						inodeFile.I_block[0] = bb

						if len(contents) > 0 {
							//
						}
						for k := 0; k < len(content); k++ {
							fileBlock.B_content[k] = content[k]
						}
						fileWritten = true
						file.Close()
						inode.I_block[i] = past
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("MKDIR", "No se ha encontrado el disco")
							return
						}
						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*father, 0)
						var binInodo bytes.Buffer
						binary.Write(&binInodo, binary.BigEndian, inode)
						WrittingBytes(file, binInodo.Bytes())
						file.Close()
						break
					} else {
						folder = Structs.NewDirectoriesBlocks()
						copy(folder.B_content[0].B_name[:], ".")
						folder.B_content[0].B_inodo = bi
						copy(folder.B_content[1].B_name[:], "..")
						folder.B_content[1].B_inodo = father
						folder.B_content[2].B_inodo = bi
						copy(folder.B_content[2].B_name[:], newf)
						copy(folder.B_content[3].B_name[:], "-")

						inodetmp.I_uid = int64(Logged.Uid)
						inodetmp.I_gid = int64(Logged.Gid)
						inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

						dateNow := time.Now().String()
						copy(inodetmp.I_atime[:], spr.S_mtime[:])
						copy(inodetmp.I_ctime[:], dateNow)
						copy(inodetmp.I_mtime[:], dateNow)
						inodetmp.I_type = 0
						inodetmp.I_perm = 664
						inodetmp.I_block[0] = bb

						copy(foldertmp.B_content[0].B_name[:], ".")
						foldertmp.B_content[0].B_inodo = bi
						copy(foldertmp.B_content[1].B_name[:], ".")
						foldertmp.B_content[1].B_inodo = father
						copy(foldertmp.B_content[2].B_name[:], "-")
						copy(foldertmp.B_content[3].B_name[:], "-")
						file.Close()

						copy(folder.B_content[2].B_name[:], newf)
						inode.I_block[i] = past
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("MKDIR", "No se ha encontrado el disco")
							return
						}
						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*father, 0)
						var binInodo bytes.Buffer
						binary.Write(&binInodo, binary.BigEndian, inode)
						WrittingBytes(file, binInodo.Bytes())
						file.Close()
						break
					}
				}
			}
		}
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco")
		return
	}

	if fileWritten {
		if fatherSpace {
			file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
			var binInodeTmp bytes.Buffer
			binary.Write(&binInodeTmp, binary.BigEndian, inodetmp)
			WrittingBytes(file, binInodeTmp.Bytes())

			file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*(bi+1), 0)
			var binInodeFile bytes.Buffer
			binary.Write(&binInodeFile, binary.BigEndian, inodetmp)
			WrittingBytes(file, binInodeFile.Bytes())

			file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*bb+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*bb, 0)
			var binFolderTmp bytes.Buffer
			binary.Write(&binFolderTmp, binary.BigEndian, fileBlock)
			WrittingBytes(file, binFolderTmp.Bytes())

			file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
			var binFile bytes.Buffer
			binary.Write(&binFile, binary.BigEndian, folder)
			WrittingBytes(file, binFile.Bytes())
		} else {

		}

		updateBm(spr, pth, "BI")
		updateBm(spr, pth, "BB")
		ruta := ""
		for i := 0; i < len(copyPath); i++ {
			ruta += "/" + copyPath[i]
		}
		Message("MKFILE", "Se ha creado el archivo en la ruta "+ruta)
	} else {
		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
		var binInodeTmp bytes.Buffer
		binary.Write(&binInodeTmp, binary.BigEndian, inodetmp)
		WrittingBytes(file, binInodeTmp.Bytes())

		file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*bb+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*bb, 0)
		var binFolderTmp bytes.Buffer
		binary.Write(&binFolderTmp, binary.BigEndian, foldertmp)
		WrittingBytes(file, binFolderTmp.Bytes())

		file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
		var binFolder bytes.Buffer
		binary.Write(&binFolder, binary.BigEndian, folder)
		WrittingBytes(file, binFolder.Bytes())

		updateBm(spr, pth, "BI")
		updateBm(spr, pth, "BB")

		ruta := ""
		for i := 0; i < len(copyPath); i++ {
			ruta += "/" + copyPath[i]
		}
		Message("MKFILE", "Se ha creado el directorio "+ruta)
		file.Close()
	}
}
