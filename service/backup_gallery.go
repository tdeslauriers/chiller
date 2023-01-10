package service

import (
	"chiller/http_client"
	"log"
)

func backupGalleryService() {

	ids, err := http_client.GetGalleryImageIds()
	if err != nil {
		log.Fatal(err)
	}

	// call images one by one.
	for _, v := range ids {

	}
}
