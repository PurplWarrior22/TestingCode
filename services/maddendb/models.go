package maddendb

import (
	"gorm.io/gorm"
)

type MaddenImageFile struct {
	//gorm model for auditing data
	gorm.Model
	//name of the file, serves as a file identifier, applications are responsible for pathing to/from/reaching out to the correct service
	FileName string `gorm:"size:500;unique;not null"`
	//name of the thumbnail
	Thumbnail string `gorm:"size:500;unique;not null"`
}

type MaddenItem struct {
	//gorm model for auditing data
	gorm.Model
	//time the madden item will start
	BeginDate int64
	//time the madden item will end
	EndDate int64
	//a summary of the madden item
	Summary string
	//additional details about the madden item
	Details string
	//historical flag, currently no use
	IsHistorical bool
	//Join table reference
	ItemImages []ItemImages
}

//an image entity to be associated with a madden item
type ItemImages struct {
	gorm.Model
	//status one of FMC PMC NMC
	Status string `gorm:"not null"`
	//the image file associated with this item image
	MaddenImageFile MaddenImageFile
	//id associated with madden entity
	MaddenItemId uint
	//id associated with the image file
	MaddenImageFileId uint
}

type Summary struct {
	gorm.Model
	Summary string `gorm:"not null"`
}

type Published struct {
	gorm.Model
	Published bool
}

//madden item hooks

//need to delete existing itemimage associations prior to updating
func (item *MaddenItem) BeforeUpdate(tx *gorm.DB) (err error) {

	if err := tx.Unscoped().Where("madden_item_id=?", item.ID).Delete(&ItemImages{}).Error; err != nil {
		return err
	}
	return nil
}
