package dataservice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PurplWarrior22/TestingCode/services/madden/swagger"
	"github.com/PurplWarrior22/TestingCode/services/maddendb"
	"github.com/PurplWarrior22/TestingCode/services/models"
	"github.com/PurplWarrior22/TestingCode/services/utilities"
	"gorm.io/gorm"
)

const (
	HISTORIC = "historic"
)

//defines an interface to interact with madden data
type MaddenDataService interface {
	//GetMaddenEntries returns a slice of maddenItem associated with the passed params, it assumes the validity of the params
	GetMaddenEntries(params swagger.GetEntryParams) ([]swagger.MaddenItem, error)
	//GetMaddenById returns a madden entry with the passed id
	GetMaddenById(id int) (swagger.MaddenItem, error)
	//CreateEntry creates a new madden item assuming the validity of the passed item
	CreateEntry(swagger.MaddenItem) (swagger.MaddenItem, error)
	//UpdateEntry updates the passed item, assuming the validity of the item
	UpdateEntry(swagger.MaddenItem) (swagger.MaddenItem, error)
	//DeleteEntry removes the madden item with an id
	DeleteEntry(id int) (error)
	//CreateSummary creates a new summary or returns appropriate error
	CreateSummary(swagger.Summary) (swagger.Summary, error)
	//GetSummary gets the most recent summary or returns appropriate error
	GetSummary() (swagger.Summary, error)
	//CreatePublished creates published state of madden
	CreatePublished(swagger.Published) (swagger.Published, error)
	//GetPublished gets the current published state of madden
	GetPublished() (swagger.Published, error)
}

type pgDataService struct {
	//the database where madden data is stored
	db maddendb.Madden
	//used to build full link to images
	appender utilities.PathBuilder
}

func NewPgDataService(db maddendb.Madden, appender utilities.PathBuilder) MaddenDataService {
	return &pgDataService{db: db, appender: appender}
}

//interface implementation

func (ds *pgDataService) CreateSummary(summary swagger.Summary) (swagger.Summary, error) {
	created, err := ds.db.CreateSummary(maddendb.Summary{Summary: summary.Summary})
	if err != nil {
		return summary, logAndReturnError(err)
	}
	return swagger.Summary{Summary: created.Summary}, nil
}

func (ds *pgDataService) GetSummary() (swagger.Summary, error) {
	summary, err := ds.db.GetSummary()
	if err != nil {
		return swagger.Summary{}, logAndReturnError(err)
	}
	return swagger.Summary{Summary: summary.Summary}, nil
}

func (ds *pgDataService) CreatePublished(published swagger.Published) (swagger.Published, error) {
	created, err := ds.db.CreatePublished(maddendb.Published{Published: published.Published})
	if err != nil {
		return published, logAndReturnError(err)
	}
	return swagger.Published{Published: created.Published}, nil
}

func (ds *pgDataService) GetPublished() (swagger.Published, error) {
	published, err := ds.db.GetPublished()
	if err != nil {
		return swagger.Published{}, logAndReturnError(err)
	}
	return swagger.Published{Published: published.Published}, nil
}

func (ds *pgDataService) GetMaddenEntries(params swagger.GetEntryParams) ([]swagger.MaddenItem, error) {
	items, err := ds.db.GetMaddenItems(*params.PageNumber, *params.PageSize, convertTime(*params.StartDate), convertTime(*params.EndDate), convertToSortField(*params.Sort), historicBool(*params.Historic))
	if err != nil {
		return nil, logAndReturnError(err)
	}
	return ds.convertToSwaggerModels(items), nil
}

func historicBool(param swagger.GetEntryParamsHistoric) bool {
	return param == HISTORIC
}

func (ds *pgDataService) GetMaddenById(id int) (swagger.MaddenItem, error) {
	item, err := ds.db.GetMaddenItemById(uint(id))
	if err != nil {
		return swagger.MaddenItem{}, logAndReturnError(err)
	}
	return ds.convertSingleModel(item), nil
}

func (ds *pgDataService) CreateEntry(item swagger.MaddenItem) (swagger.MaddenItem, error) {
	created, err := ds.db.CreateMaddenItem(swaggerToEntry(item, 0))
	if err != nil {
		return swagger.MaddenItem{}, logAndReturnError(err)
	}
	return ds.convertSingleModel(created), nil
}

func (ds *pgDataService) UpdateEntry(item swagger.MaddenItem) (swagger.MaddenItem, error) {
	updated, err := ds.db.UpdateMaddenItem(swaggerToEntry(item, uint(*item.Id)))
	if err != nil {
		return swagger.MaddenItem{}, logAndReturnError(err)
	}
	return ds.convertSingleModel(updated), nil
}

func (ds *pgDataService) DeleteEntry(id int) (error) {
	deleted := ds.db.DeleteMaddenItem(uint(id))
	if deleted != nil {
		return logAndReturnError(deleted)
	}

	return nil
}

//helpers
func logAndReturnError(err error) error {
	fmt.Printf("Error during database action ERROR: %s\n", err.Error())
	switch converted := err.(type) {
	case *maddendb.DbError:
		fmt.Println(converted.OriginalError)
		return models.NewDataServiceError(err.Error(), http.StatusInternalServerError)
	default:
		return err
	}
}

func (ds *pgDataService) convertToSwaggerModels(items []maddendb.MaddenItem) []swagger.MaddenItem {
	swaggerItems := []swagger.MaddenItem{}
	for _, item := range items {
		swaggerItems = append(swaggerItems, ds.convertSingleModel(item))
	}
	return swaggerItems
}

func (ds *pgDataService) convertSingleModel(item maddendb.MaddenItem) swagger.MaddenItem {
	return swagger.MaddenItem{
		Details:    item.Details,
		Summary:    item.Summary,
		EndDate:    time.Unix(item.EndDate, 0).UTC().Format(time.RFC3339),
		StartDate:  time.Unix(item.BeginDate, 0).UTC().Format(time.RFC3339),
		Historical: &item.IsHistorical,
		Id:         uintPtr(int(item.ID)),
		Images:     ds.convertToSwaggerImages(item.ItemImages),
	}
}

func (ds *pgDataService) convertToSwaggerImages(images []maddendb.ItemImages) []swagger.MaddenImage {
	converted := []swagger.MaddenImage{}
	for _, image := range images {
		converted = append(converted, swagger.MaddenImage{
			Id:            int(image.MaddenImageFileId),
			ImageLink:     utilities.StrPtr(ds.appender.BuildFullPath(image.MaddenImageFile.FileName)),
			ThumbnailLink: utilities.StrPtr(ds.appender.BuildFullPath(image.MaddenImageFile.Thumbnail)),
			Status:        swagger.MaddenImageStatus(image.Status),
		})
	}
	return converted
}

func convertTime(dateString string) int64 {
	//safe to assume time validity
	timeObj, _ := time.Parse(time.RFC3339, dateString)
	return timeObj.Unix()
}

func convertToSortField(sortField swagger.GetEntryParamsSort) maddendb.SortField {
	switch sortField {
	case "startDate":
		return maddendb.StartDate
	default:
		return maddendb.EndDate
	}
}

func swaggerToEntry(item swagger.MaddenItem, id uint) maddendb.MaddenItem {
	return maddendb.MaddenItem{
		Model:        gorm.Model{ID: id},
		Details:      item.Details,
		Summary:      item.Summary,
		IsHistorical: nullSafeHistoric(item.Historical),
		BeginDate:    convertTime(item.StartDate),
		EndDate:      convertTime(item.EndDate),
		ItemImages:   swaggerToImages(item.Images, id),
	}
}

func swaggerToImages(images []swagger.MaddenImage, itemId uint) []maddendb.ItemImages {
	itemImages := []maddendb.ItemImages{}
	for _, image := range images {
		itemImages = append(itemImages, maddendb.ItemImages{
			MaddenImageFileId: uint(image.Id),
			Status:                 string(image.Status),
		})
	}
	return itemImages
}

func nullSafeHistoric(item *bool) bool {
	if item == nil {
		return false
	}
	return *item
}

func uintPtr(i int) *int {
	return &i
}
