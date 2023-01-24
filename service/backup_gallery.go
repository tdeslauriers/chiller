package service

import (
	"chiller/dao"
	"chiller/http_client"
	"log"
	"sync"
)

func BackupGalleryService(token http_client.Bearer) {

	ids, _ := http_client.GetGalleryImageIds(token)

	// capture xref and album records
	ais := make(chan dao.AlbumImages, len(ids)*6) // most imgs should have <= 6 xrefs

	var wgTables sync.WaitGroup
	wgTables.Add(len(ids))

	// call images one by one.
	for _, v := range ids {
		go func(id int64) {
			defer wgTables.Done()

			img, err := http_client.GetGalleryImage(id, token)
			if err != nil {
				log.Printf("Failed to get img Id: %d - err: %v", id, err)
			}

			// insert or update the image record
			if err := dao.InsertImage(img); err != nil {
				if err := dao.UpdateImage(img); err != nil {
					log.Fatal(err)
				}
			}

			// send to channel
			for _, albumImage := range img.AlbumImages {
				i := dao.Image{Id: img.Id} // need the image id for xref
				albumImage.Image = i
				ais <- albumImage
			}

		}(v)
	}
	wgTables.Wait()
	close(ais)

	var wgXref sync.WaitGroup
	wgXref.Add(len(ais))

	for ai := range ais {

		go func(ai dao.AlbumImages) {
			defer wgXref.Done()

			// album must go first so xref doesnt error
			if err := dao.InsertAlbum(ai.Album); err != nil {
				if err := dao.UpdateAlbum(ai.Album); err != nil {
					log.Fatal(err)
				}
			}

			xref := dao.XrefRecord[dao.AiXref]{Id: ai.Id, Fk_1: ai.Album.Id, Fk_2: ai.Image.Id}
			if err := dao.InsertGalleryXrefRecord(xref, dao.INSERT_AI); err != nil {
				log.Print(err)
			}
		}(ai)
	}

	wgXref.Wait()
	log.Print("Completed backup activites of gallery-service.")
}
