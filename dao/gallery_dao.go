package dao

import "log"

const GALLERY_BACKUP_DB = "backup_gallery"

// image crud
func InsertImage(image Image) (err error) {

	db := dbConn(GALLERY_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO image (id, filename, title, description, date, published, thumbnail, presentation, image) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		image.Id,
		image.Filename,
		image.Title,
		image.Description,
		image.Date,
		image.Published,
		image.Thumbnail,
		image.Presentation,
		image.Image)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Image record %d inserted into backup gallery database.", id)
	return err
}

func InsertAlbum(a Album) (err error) {

	db := dbConn(GALLERY_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO album (id, album) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(a.Id, a.Album)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Album record %d inserted into backup gallery database.", id)
	return err
}

func UpdateImage(image Image) (err error) {

	db := dbConn(GALLERY_BACKUP_DB)
	defer db.Close()

	// other fields will not change
	query := "UPDATE image SET title = ?, description = ?, published =? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(image.Title, image.Description, image.Published, image.Id)

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()
	if count > 0 {
		log.Printf("Updated image record %d in the backup gallery database.", image.Id)
	}
	return err
}

func UpdateAlbum(a Album) (err error) {

	db := dbConn(GALLERY_BACKUP_DB)
	defer db.Close()

	query := "UPDATE album SET album = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(a.Album, a.Id)

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()
	if count > 0 {
		log.Printf("Updated album record %d in the backup gallery database.", a.Id)
	}
	return err
}

func InsertGalleryXrefRecord[T row](r XrefRecord[T], query string) (err error) {

	db := dbConn(GALLERY_BACKUP_DB)
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(
		r.Id,
		r.Fk_1,
		r.Fk_2)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("%T xref record %d inserted into backup gallery database.", r, id)
	return err
}
