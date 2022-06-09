package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"../services/maddendb"
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
)

var (
	//object under test
	postgresMaint maddendb.Madden
)

func init() {
	var err error
	postgresMaint, err = maddendb.BuildPostgresMaddenFromEnvironment()
	if err != nil {
		fmt.Printf("failure during db setup ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	err = buildTestDbHook()
	if err != nil {
		fmt.Printf("error building test database connection, ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	postgresMaint.SetupDatabase()
}

//Setup

func setup(t *testing.T) func(t *testing.T) {
	return tearDown
}

func tearDown(t *testing.T) {
	if err := db.Unscoped().Where("1=1").Delete(&maddendb.ItemImages{}).Error; err != nil {
		t.Errorf("error cleaning up tables ERROR: %s\n", err.Error())
	}
	if err := db.Unscoped().Where("1=1").Delete(&maddendb.MaddenItem{}).Error; err != nil {
		t.Errorf("error cleaning up tables ERROR: %s\n", err.Error())
	}
	if err := db.Unscoped().Where("1=1").Delete(&maddendb.MaddenImageFile{}).Error; err != nil {
		t.Errorf("error cleaning up tables ERROR: %s\n", err.Error())
	}
	if err := db.Unscoped().Where("1=1").Delete(&maddendb.Summary{}).Error; err != nil {
		t.Errorf("error cleaning up tables ERROR: %s\n", err.Error())
	}
}

//Tests

func TestDelete(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	err := postgresMaint.DeleteMaddenImage(1)
	if err != nil {
		t.Errorf("expected nil error on delete got ERROR: %s\n", err.Error())
	}
	items, err := postgresMaint.GetMaddenImages(0, 25)
	if err != nil {
		t.Errorf("got error while confirming madden items inserted ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 4, len(items))
	for _, item := range items {
		assert.NotEqual(t, 1, item.ID)
	}
}

func TestMaddenEntryDelete(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultItems(t)
	var removeId = uint(1)
	_, err := postgresMaint.GetMaddenItemById(removeId)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	err = postgresMaint.DeleteMaddenItem(removeId)
	if err != nil {
		t.Errorf("expected nil error on delete got ERROR: %s\n", err.Error())
	}
	_, err = postgresMaint.GetMaddenItemById(removeId)
	if err == nil {
		t.Errorf("error on item delete: %s\n", err.Error())
		t.FailNow()
	}
}

func TestCreate(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)
	item, err := postgresMaint.CreateMaddenItem(createDefaultItem())
	if err != nil {
		t.Errorf("expected non error but got error: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, "Im a summary", item.Summary)
}

func TestCreateOnItemExists(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)
	t1 := time.Now().UTC().Unix()
	t2 := time.Now().UTC().Unix()
	item := maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details"}
	duplicate := maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details"}
	created, err := postgresMaint.CreateMaddenItem(item)
	if err != nil {
		t.Errorf("expected non error but got error: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, created.BeginDate, duplicate.BeginDate)
	assert.Equal(t, duplicate.Summary, item.Summary)
	assert.Equal(t, duplicate.EndDate, item.EndDate)
	assert.Equal(t, duplicate.Details, item.Details)
	assert.Equal(t, duplicate.BeginDate, item.BeginDate)
	_, err = postgresMaint.CreateMaddenItem(item)
	if err == nil {
		t.Errorf("expected error on duplicate insert, but got no error\n")
		t.FailNow()
	}
}

func TestValidUpdate(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := createDefaultItem()
	item.IsHistorical = true
	inserted, err := postgresMaint.CreateMaddenItem(item)
	fmt.Println("INSERTED ID")
	fmt.Println(inserted.ID)
	if err != nil {
		t.Errorf("expected non error but got error: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.Details, item.Details)
	assert.Equal(t, inserted.BeginDate, item.BeginDate)
	assert.Equal(t, inserted.EndDate, item.EndDate)
	assert.Equal(t, inserted.Summary, item.Summary)
	assert.Equal(t, inserted.IsHistorical, item.IsHistorical)
	inserted.Summary = "whoops i needed to update the summary"
	inserted.IsHistorical = false
	updated, err := postgresMaint.UpdateMaddenItem(inserted)
	if err != nil {
		t.Errorf("expected non error but got error: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.Details, updated.Details)
	assert.Equal(t, inserted.BeginDate, updated.BeginDate)
	assert.Equal(t, inserted.EndDate, updated.EndDate)
	assert.Equal(t, inserted.Summary, updated.Summary)
	assert.Equal(t, inserted.IsHistorical, updated.IsHistorical)
}

func TestInvalidUpdate(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := createDefaultItem()
	inserted, err := postgresMaint.CreateMaddenItem(item)
	if err != nil {
		t.Errorf("expected non error but got error: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.Details, item.Details)
	assert.Equal(t, inserted.BeginDate, item.BeginDate)
	assert.Equal(t, inserted.EndDate, item.EndDate)
	assert.Equal(t, inserted.Summary, item.Summary)
	inserted.ID = 42
	_, err = postgresMaint.UpdateMaddenItem(item)
	if err == nil {
		t.Errorf("expected error but got none")
		t.FailNow()
	}
}

func TestSearch(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultItems(t)
	items, err := postgresMaint.GetMaddenItems(0, 10, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), maddendb.StartDate, false)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 4, len(items))
	for i, item := range items {
		index := len(items) - i
		assert.Equal(t, fmt.Sprintf("Item%d", index), item.Summary)
		assert.Equal(t, fmt.Sprintf("Details%d", index), item.Details)
		if index-1 < 3 {
			assert.Equal(t, 1, len(item.ItemImages))
		} else {
			assert.Equal(t, 2, len(item.ItemImages))
		}
		for j, image := range item.ItemImages {
			fmt.Println(image.MaddenImageFile)
			assert.Equal(t, fmt.Sprintf("f%d", j+index), image.MaddenImageFile.FileName)
		}
	}
}

func TestSort(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertSortTestItems(t)
	items, err := postgresMaint.GetMaddenItems(0, 10, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), maddendb.StartDate, false)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 4, len(items))
	assert.Equal(t, "Item3", items[0].Summary)
	assert.Equal(t, "Item2", items[1].Summary)
	assert.Equal(t, "Item4", items[2].Summary)
	assert.Equal(t, "Item1", items[3].Summary)
}

func TestSummaryCreate(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	s := maddendb.Summary{Summary: "hello i am a summary"}
	saved, err := postgresMaint.CreateSummary(s)
	if err != nil {
		t.Errorf("got error on create summary expected none, ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, saved.Summary, s.Summary)
}

func TestSummaryReturnsMostRecent(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	s := maddendb.Summary{Summary: "hello i am a summary"}
	saved, err := postgresMaint.CreateSummary(s)
	if err != nil {
		t.Errorf("got error on create summary expected none, ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, saved.Summary, s.Summary)
	s = maddendb.Summary{Summary: "I now have new text"}
	_, err = postgresMaint.CreateSummary(s)
	if err != nil {
		t.Errorf("got error on create summary expected none, ERROR: %s\n", err.Error())
		t.FailNow()
	}
	shouldBeS, err := postgresMaint.GetSummary()
	if err != nil {
		t.Errorf("got error on summar search, expected none ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, shouldBeS.Summary, s.Summary)
}


func TestSingleItemSearch(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultItems(t)
	items, err := postgresMaint.GetMaddenItems(3, 1, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), maddendb.StartDate, false)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "Details1", items[0].Details)
}

func TestEmptyItems(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultItems(t)
	items, err := postgresMaint.GetMaddenItems(0, 20, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).Unix(), maddendb.StartDate, false)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 0, len(items))
}


func TestSearchById(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	t1 := time.Now().UTC().Unix()
	t2 := time.Now().UTC().Unix()
	item := maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details"}
	inserted, err := postgresMaint.CreateMaddenItem(item)
	if err != nil {
		t.Errorf("error on item insert ERROR: %s\n", err.Error())
		t.FailNow()
	}
	foundById, err := postgresMaint.GetMaddenItemById(inserted.ID)
	if err != nil {
		t.Errorf("error on item search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, foundById.Details, inserted.Details)
	assert.Equal(t, foundById.ID, inserted.ID)
}

func TestErrorOnSearchNoItem(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	t1 := time.Now().UTC().Unix()
	t2 := time.Now().UTC().Unix()
	item := maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details"}
	inserted, err := postgresMaint.CreateMaddenItem(item)
	if err != nil {
		t.Errorf("error on item insert ERROR: %s\n", err.Error())
		t.FailNow()
	}
	_, err = postgresMaint.GetMaddenItemById(inserted.ID + 1)
	if err == nil {
		t.Errorf("expected error on item search but got none")
		t.FailNow()
	}
}

func TestShouldntInsertWithoutImage(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	t1 := time.Now().UTC().Unix()
	t2 := time.Now().UTC().Unix()
	item := maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details", ItemImages: []maddendb.ItemImages{{MaddenImageFileId: 5435}}}
	_, err := postgresMaint.CreateMaddenItem(item)
	if err == nil {
		t.Errorf("Expected failure on insert where image did not exist but got none")
	}
}

func TestImageInsert(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := maddendb.MaddenImageFile{FileName: "file1", Thumbnail: "thumb1"}
	inserted, err := postgresMaint.CreateMaddenImage(item)
	if err != nil {
		t.Errorf("error while inserting item ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.FileName, item.FileName)
}

func TestInvalidImageInsert(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := maddendb.MaddenImageFile{FileName: "file1", Thumbnail: "thumb1"}
	inserted, err := postgresMaint.CreateMaddenImage(item)
	if err != nil {
		t.Errorf("error while inserting item ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.FileName, item.FileName)
	_, err = postgresMaint.CreateMaddenImage(inserted)
	if err == nil {
		t.Errorf("expected error on duplicate filename insert, got none")
	}
}

func TestValidImageUpdate(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := maddendb.MaddenImageFile{FileName: "file1", Thumbnail: "thumb1"}
	inserted, err := postgresMaint.CreateMaddenImage(item)
	if err != nil {
		t.Errorf("error while inserting item ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.FileName, item.FileName)
	inserted.FileName = "file2"
	original, updated, err := postgresMaint.UpdateMaddenImage(inserted)
	assert.Equal(t, "file1", original.FileName)
	if err != nil {
		t.Errorf("error while updating item ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.FileName, updated.FileName)
	assert.Equal(t, inserted.ID, updated.ID)
}

func TestInvalidImageUpdate(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	item := maddendb.MaddenImageFile{FileName: "file1", Thumbnail: "thumb1"}
	inserted, err := postgresMaint.CreateMaddenImage(item)
	if err != nil {
		t.Errorf("error while inserting item ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, inserted.FileName, item.FileName)
	inserted.ID++
	_, _, err = postgresMaint.UpdateMaddenImage(inserted)
	if err == nil {
		t.Errorf("error while updating item expected, but got none\n")
	}
}

func TestImageSearch(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	items, err := postgresMaint.GetMaddenImages(0, 10)
	if err != nil {
		t.Errorf("error on image search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 5, len(items))
	for i, item := range items {
		assert.Equal(t, fmt.Sprintf("f%d", i+1), item.FileName)
	}
}

func TestPagedImageSearch(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	items, err := postgresMaint.GetMaddenImages(0, 1)
	if err != nil {
		t.Errorf("error on image search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 1, len(items))
	for i, item := range items {
		assert.Equal(t, fmt.Sprintf("f%d", i+1), item.FileName)
	}
}

func TestNonOnePagedImageSearch(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	items, err := postgresMaint.GetMaddenImages(1, 1)
	if err != nil {
		t.Errorf("error on image search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "f2", items[0].FileName)
}

func TestImageSearchByName(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	items, err := postgresMaint.GetMaddenImagesByName(0, 10, "f1")
	if err != nil {
		t.Errorf("error on image search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 1, len(items))
	for i, item := range items {
		assert.Equal(t, fmt.Sprintf("f%d", i+1), item.FileName)
	}
}

func TestImageSearchByNameAll(t *testing.T) {
	tearDown := setup(t)
	defer tearDown(t)
	insertDefaultImages(t)
	items, err := postgresMaint.GetMaddenImagesByName(0, 10, "f")
	if err != nil {
		t.Errorf("error on image search ERROR: %s\n", err.Error())
		t.FailNow()
	}
	assert.Equal(t, 5, len(items))
	for i, item := range items {
		assert.Equal(t, fmt.Sprintf("f%d", i+1), item.FileName)
	}
}

//Test helpers
func createDefaultItem() maddendb.MaddenItem {
	t1 := time.Now().UTC().Unix()
	t2 := time.Now().UTC().Unix()
	return maddendb.MaddenItem{BeginDate: t1, EndDate: t2, Summary: "Im a summary", Details: "these are details"}
}

func insertDefaultImages(t *testing.T) {
	images := []maddendb.MaddenImageFile{
		{
			FileName:  "f1",
			Model:     gorm.Model{ID: 1},
			Thumbnail: "t1",
		},
		{
			FileName:  "f2",
			Model:     gorm.Model{ID: 2},
			Thumbnail: "t2",
		},
		{
			FileName:  "f3",
			Model:     gorm.Model{ID: 3},
			Thumbnail: "t3",
		},
		{
			FileName:  "f4",
			Model:     gorm.Model{ID: 4},
			Thumbnail: "t4",
		},
		{
			FileName:  "f5",
			Model:     gorm.Model{ID: 5},
			Thumbnail: "t5",
		},
	}
	for _, image := range images {
		if _, err := postgresMaint.CreateMaddenImage(image); err != nil {
			t.Errorf("error on image insert ERROR: %s\n", err.Error())
			t.FailNow()
		}
	}
}

func insertDefaultItems(t *testing.T) {
	insertDefaultImages(t)
	items := []maddendb.MaddenItem{
		{
			Summary:   "Item1",
			Details:   "Details1",
			BeginDate: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2022, 2, 2, 2, 2, 2, 2, time.UTC).Unix(),

			Model: gorm.Model{ID: 1},
			ItemImages: []maddendb.ItemImages{
				{
					Status:               "FMC",
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 1}},
				},
			},
		},
		{
			Summary:   "Item2",
			Details:   "Details2",
			BeginDate: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2021, 2, 2, 2, 2, 2, 2, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 2}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 2},
		},
		{
			Summary:   "Item3",
			Details:   "Details3",
			BeginDate: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2020, 2, 2, 2, 2, 2, 2, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 3}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 3},
		},

		{
			Summary:   "Item4",
			Details:   "Details4",
			BeginDate: time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2019, 2, 2, 2, 2, 2, 2, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 4}},
					Status:               "FMC",
				},
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 5}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 4},
		},
	}

	for _, item := range items {
		if _, err := postgresMaint.CreateMaddenItem(item); err != nil {
			t.Errorf("error on default item insert ERROR: %s\n", err.Error())
			t.FailNow()
		}
	}
}

// This is a slightly different set of items used to test if it's sorting correctly
func insertSortTestItems(t *testing.T) {
	insertDefaultImages(t)
	items := []maddendb.MaddenItem{
		{
			Summary:   "Item1",
			Details:   "Details1",
			BeginDate: time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC).Unix(),

			Model: gorm.Model{ID: 1},
			ItemImages: []maddendb.ItemImages{
				{
					Status:               "FMC",
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 1}},
				},
			},
		},
		{
			Summary:   "Item2",
			Details:   "Details2",
			BeginDate: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2022, 1, 3, 1, 1, 1, 1, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 2}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 2},
		},
		{
			Summary:   "Item3",
			Details:   "Details3",
			BeginDate: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2022, 1, 2, 1, 1, 1, 1, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 3}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 3},
		},

		{
			Summary:   "Item4",
			Details:   "Details4",
			BeginDate: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC).Unix(),
			EndDate:   time.Date(2022, 2, 2, 2, 2, 2, 2, time.UTC).Unix(),
			ItemImages: []maddendb.ItemImages{
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 4}},
					Status:               "FMC",
				},
				{
					MaddenImageFile: maddendb.MaddenImageFile{Model: gorm.Model{ID: 5}},
					Status:               "FMC",
				},
			},
			Model: gorm.Model{ID: 4},
		},
	}

	for _, item := range items {
		if _, err := postgresMaint.CreateMaddenItem(item); err != nil {
			t.Errorf("error on default item insert ERROR: %s\n", err.Error())
			t.FailNow()
		}
	}
}
