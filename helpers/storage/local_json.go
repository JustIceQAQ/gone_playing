package storage

import (
	"log"
	"os"
)

func SaveToJson(filename string, context []byte) {
	path := "./datas/"
	err := os.WriteFile(path+filename, context, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
