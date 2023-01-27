package service

import (
	"chiller/dao"
	"chiller/http_client"
	"log"
)

func BackupGalleryService(token http_client.Bearer) {

	ids, _ := http_client.GetGalleryImageIds(token)

	// call images one by one.
	for _, id := range ids {

		img, err := http_client.GetGalleryImage(id, token)
		if err != nil {
			log.Fatal(err)
		}

		// insert or update the image record
		if err := dao.InsertImage(img); err != nil {
			if err := dao.UpdateImage(img); err != nil {
				log.Fatal(err)
			}
		}

		for _, albumImage := range img.AlbumImages {

			if err := dao.InsertAlbum(albumImage.Album); err != nil {
				if err := dao.UpdateAlbum(albumImage.Album); err != nil {
					log.Fatal(err)
				}
			}

			i := dao.Image{Id: img.Id} // need the image id for xref
			albumImage.Image = i
			xref := dao.XrefRecord[dao.AiXref]{Id: albumImage.Id, Fk_1: albumImage.Album.Id, Fk_2: albumImage.Image.Id}
			dao.InsertGalleryXrefRecord(xref, dao.INSERT_AI) // dump insert errors.
		}
	}

	log.Print("Completed backup activites of gallery-service.")
}
