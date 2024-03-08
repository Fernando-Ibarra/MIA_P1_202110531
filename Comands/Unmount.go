package Comands

import (
	"strconv"
	"strings"
)

func DataUnMount(tokens []string) {
	if len(tokens) > 1 {
		Error("UNMOUNT", "Solo se acepta el párametro id")
		return
	}
	id := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			if id == "" {
				id = tk[1]
			} else {
				Error("UNMOUNT", "Parámetro driveletter repetido en el comando: "+tk[0])
			}
		} else {
			Error("UNMOUNT", "No se esperaba el parámetro "+tk[0])
			error_ = false
			return
		}
	}
	if error_ {
		return
	}
	if id == "" {
		Error("UNMOUNT", "Se require el parámetro id")
		return
	} else {
		unmount(id)
	}
}

func unmount(id string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("UNMOUNT", "El primer identificador no es válido")
		return
	}
	letter := id[0]
	j, _ := strconv.Atoi(string(id[1] - 1))
	if j < 0 {
		Error("UNMOUNT", "El primer identificador no es válido")
		return
	}
	for i := 0; i < 99; i++ {
		if DiskMount[i].Partitions[j].State == 1 {
			if DiskMount[i].Partitions[j].Letter == letter {
				DiskMount[i].Partitions[j].State = 0
				Message("UNMOUNT", "Se ha realizado correctamente el unmount -id="+id)
				return
			} else {
				Error("UNMOUNT", "No se ha podido realizar correctamente el unmount -id="+id)
				return
			}
		}
	}

}
