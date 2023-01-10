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
		log.Printf("Updated image record %d in the backup image database.", image.Id)
	}
	return err
}
