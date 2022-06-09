package maddendb

import (
	"fmt"

	"gorm.io/gorm"
)

//defines and implements madden crud operations

type SortField int

const (
	StartDate SortField = 0
	EndDate   SortField = 1
)

//Madden defines an interface to interact with a madden information data store
type Madden interface {
	//GetSummary returns the most recent madden summary
	GetSummary() (Summary, error)
	//CreateSummary creates a new summary returning an error if something goes wrong
	CreateSummary(summary Summary) (Summary, error)
	//GetPublished returns current state of if madden is published
	GetPublished() (Published, error)
	//CreatePublished creates a new state of madden published
	CreatePublished(published Published) (Published, error)
	//CreateMaintenacneItem creates a new madden item returning an error if anything fails, or if an identical item exists
	CreateMaddenItem(item MaddenItem) (MaddenItem, error)
	//DeleteMaddenItem deletes a madden item given an id, returning an error if one occurs
	DeleteMaddenItem(id uint) error
	//UpdateMaddenItem updates an existing madden item returning an error if anything fails, or if the item did not already exist
	UpdateMaddenItem(item MaddenItem) (MaddenItem, error)
	//GetMaddenItems returns a page of madden items offest by pagenum and size, filtered on start and end date, and sorted by sortField returning an error if anything goes wrong
	//pages are 0 indexed
	GetMaddenItems(pageNum, size int, startDate, endDate int64, sortField SortField, historic bool) ([]MaddenItem, error)
	//GetMaddenItemById returns the madden item with the passed id, or an error if it did not exist or something went wrong
	GetMaddenItemById(id uint) (MaddenItem, error)
	//CreateImage creates a new madden image returning an error if anything fails or an image with the same name exists
	CreateMaddenImage(image MaddenImageFile) (MaddenImageFile, error)
	//UpdateMaddenImage updates an existing madden image, returning an error if anything goes wrong or if the image did not exist
	//the original maddenImageFile entity and the updated entity are returned
	UpdateMaddenImage(image MaddenImageFile) (MaddenImageFile, MaddenImageFile, error)
	//GetMaddenImages returns a page of of Madden images, offset by pagenum and size returning an error if anything goes wrong
	GetMaddenImages(pageNum, size int) ([]MaddenImageFile, error)
	//GetMaddenImagesByName returns a slice of madden images, offset by pagenum and size, with a name similar to, or exactly matching filename
	GetMaddenImagesByName(pageNum, size int, filename string) ([]MaddenImageFile, error)
	//DeleteMaddenImage deletes the image entry with id, returning an error if one occurs
	DeleteMaddenImage(id uint) error
	//SetupDatabase builds all table or migrates existing schemas, this should be the first call any client of this interface makes
	SetupDatabase() error
}

//postgres backed implementation of Madden
type postgresMadden struct {
	db *gorm.DB
}

//postgres constructor
func NewPostgresMaintenace(db *gorm.DB) Madden {
	return &postgresMadden{db: db}
}

//Interface implementation

func (pm *postgresMadden) GetSummary() (Summary, error) {
	summary := Summary{}
	if err := pm.db.Order("created_at desc").Take(&summary).Error; err != nil {
		// There is currently summary, return an empty string
		return summary, nil
	}
	return summary, nil
}

func (pm *postgresMadden) CreateSummary(summary Summary) (Summary, error) {
	created := summary
	if err := pm.db.Create(&created).Error; err != nil {
		fmt.Printf("error creating summary ERROR: %s\n", err.Error())
		return summary, &DbError{Message: "error creating summary", OriginalError: err}
	}
	return created, nil
}

func (pm *postgresMadden) GetPublished() (Published, error) {
	published := Published{}
	if err := pm.db.Order("created_at desc").Take(&published).Error; err != nil {
		// There is currently no published state
		return published, nil
	}
	return published, nil
}

func (pm *postgresMadden) CreatePublished(published Published) (Published, error) {
	created := published
	if err := pm.db.Create(&created).Error; err != nil {
		fmt.Printf("error creating published state ERROR: %s\n", err.Error())
		return published, &DbError{Message: "error creating published state", OriginalError: err}
	}
	return created, nil
}

func (pm *postgresMadden) DeleteMaddenImage(id uint) error {
	image := MaddenImageFile{Model: gorm.Model{ID: id}}
	if err := pm.db.Delete(&image).Error; err != nil {
		fmt.Printf("error deleting entry with id %d, ERROR: %s\n", id, err.Error())
		return &DbError{Message: fmt.Sprintf("error deleting entry %d", id), OriginalError: err}
	}
	return nil
}

func (pm *postgresMadden) SetupDatabase() error {
	if err := pm.db.AutoMigrate(&MaddenImageFile{}, &MaddenItem{}, &ItemImages{}, &Summary{}, &Published{}); err != nil {
		return &DbError{Message: "Error building database", OriginalError: err}
	}
	return nil
}

func (pm *postgresMadden) CreateMaddenItem(item MaddenItem) (MaddenItem, error) {
	if existed, err := pm.itemExists(item); err != nil || existed {
		return MaddenItem{}, &DbError{Message: "Item Already existed", OriginalError: err}
	}
	//insert the madden item
	insertable := item
	if err := pm.db.Create(&insertable).Error; err != nil {
		return MaddenItem{}, &DbError{Message: "error during Item Creation", OriginalError: err}
	}
	if err := pm.db.Preload("ItemImages").Preload("ItemImages.MaddenImageFile").Find(&insertable).Error; err != nil {
		return insertable, &DbError{Message: "error while retrieving created item", OriginalError: err}
	}
	return insertable, nil
}

func (pm *postgresMadden) UpdateMaddenItem(item MaddenItem) (MaddenItem, error) {
	//error is nil if the item existed
	if err := pm.db.Take(&MaddenItem{}, item.ID).Error; err != nil {
		return MaddenItem{}, &DbError{Message: "Error or item did not exist on update", OriginalError: err}
	}
	insertable := item
	mapped := entryToMap(insertable)
	if err := pm.db.Model(&insertable).Updates(mapped).Error; err != nil {
		return MaddenItem{}, &DbError{Message: "error on update", OriginalError: err}
	}
	if err := pm.db.Preload("ItemImages").Preload("ItemImages.MaddenImageFile").Find(&insertable).Error; err != nil {
		return insertable, &DbError{Message: "error while retrieving updated item", OriginalError: err}
	}

	return insertable, nil
}

func (pm *postgresMadden) DeleteMaddenItem(id uint) error {
	item := MaddenItem{Model: gorm.Model{ID: id}}
	if err := pm.db.Delete(&item).Error; err != nil {
		fmt.Printf("error deleting entry with id %d, ERROR: %s\n", id, err.Error())
		return &DbError{Message: fmt.Sprintf("error deleting entry %d", id), OriginalError: err}
	}
	return nil
}

func (pm *postgresMadden) GetMaddenItems(pageNum, size int, startDate, endDate int64, sortField SortField, historic bool) ([]MaddenItem, error) {
	items := []MaddenItem{}
	if err := pm.db.Offset(pageNum*size).Limit(size).Order(itemOrderString(sortField, startDate, endDate, false)).Order(itemOrderString(sortField, startDate, endDate, true)).Where("begin_date < ? AND end_date > ? AND is_historical = ?", startDate, endDate, historic).Preload("ItemImages").Preload("ItemImages.MaddenImageFile").Find(&items).Error; err != nil {
		fmt.Println(err.Error())
		return nil, &DbError{Message: "error on search", OriginalError: err}
	}
	return items, nil
}

func (pm *postgresMadden) GetMaddenItemById(id uint) (MaddenItem, error) {
	item := MaddenItem{}
	if err := pm.db.Preload("ItemImages").Preload("ItemImages.MaddenImageFile").First(&item, id).Error; err != nil {
		//some error other than the record didn't exist
		if !(err == gorm.ErrRecordNotFound) {
			return MaddenItem{}, &DbError{Message: err.Error(), OriginalError: err}
		}
		return MaddenItem{}, &DbError{Message: fmt.Sprintf("item with ID: %d did not exist", item.ID)}
	}
	return item, nil
}

func (pm *postgresMadden) CreateMaddenImage(image MaddenImageFile) (MaddenImageFile, error) {
	inserted := image
	if err := pm.db.Where("file_name=?", image.FileName).Take(&MaddenImageFile{}).Error; err != nil {
		if !(err == gorm.ErrRecordNotFound) {
			return inserted, &DbError{Message: "error during check for existing image", OriginalError: err}
		}
	} else {
		return inserted, &DbError{Message: fmt.Sprintf("image with filename %s already exists", image.FileName)}
	}
	if err := pm.db.Create(&inserted).Error; err != nil {
		return inserted, &DbError{Message: "error while inserting image into database", OriginalError: err}
	}
	return inserted, nil
}

func (pm *postgresMadden) UpdateMaddenImage(image MaddenImageFile) (MaddenImageFile, MaddenImageFile, error) {
	original := MaddenImageFile{}
	updateable := image
	if err := pm.db.Take(&original, image.ID).Error; err != nil {
		if !(err == gorm.ErrRecordNotFound) {
			return original, updateable, &DbError{Message: "error during check for existing image", OriginalError: err}
		}
		return original, updateable, &DbError{Message: fmt.Sprintf("image with id %d did not exist", image.ID)}
	}
	if err := pm.db.Preload("ItemImages").Preload("ItemImages.MaddenImageFile").Updates(updateable).Error; err != nil {
		return original, updateable, &DbError{Message: fmt.Sprintf("error while updating item with id %d", image.ID), OriginalError: err}
	}
	return original, updateable, nil
}

func (pm *postgresMadden) GetMaddenImages(pageNum, size int) ([]MaddenImageFile, error) {
	images := []MaddenImageFile{}
	if err := pm.db.Offset(pageNum * size).Limit(size).Order(imageOrderString()).Find(&images).Error; err != nil {
		return images, &DbError{Message: "error while searching for images", OriginalError: err}
	}
	return images, nil
}

func (pm *postgresMadden) GetMaddenImagesByName(pageNum, size int, filename string) ([]MaddenImageFile, error) {
	images := []MaddenImageFile{}
	if err := pm.db.Offset(pageNum*size).Limit(size).Where("file_name ~ ?", filename).Order(imageOrderString()).Find(&images).Error; err != nil {
		return images, &DbError{Message: "error while searching for images", OriginalError: err}
	}

	return images, nil
}

//Implementation helpers

//itemExists checks if a madden item with identical fields exists already, returns true if the item already existed, false if it did not
func (pm *postgresMadden) itemExists(item MaddenItem) (bool, error) {
	if err := pm.db.Where("end_date=? AND begin_date=? AND summary=? AND details=?", item.EndDate, item.BeginDate, item.Summary, item.Details).Take(&MaddenItem{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, &DbError{Message: err.Error(), OriginalError: err}
	}
	return true, nil
}

func itemOrderString(sortField SortField, startDate, endDate int64, reverse bool) string {
	if sortField == StartDate {
		if reverse {
			return "end_date asc"
		} else {
			return "begin_date asc"
		}
	}
	if reverse {
		return "begin_date asc"
	} else {
		return "end_date asc"
	}
}

//no required image order so order by creation time
func imageOrderString() string {
	return "created_at asc"
}

//false is a 0 value so need to convert the entire entry or gorm wont update isHistorical
func entryToMap(entry MaddenItem) map[string]interface{} {
	return map[string]interface{}{
		"id":            entry.ID,
		"begin_date":    entry.BeginDate,
		"end_date":      entry.EndDate,
		"summary":       entry.Summary,
		"details":       entry.Details,
		"is_historical": entry.IsHistorical,
	}
}
