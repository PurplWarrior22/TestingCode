package models

import (
	"strings"
)

//holds shared models

//OthEntry holds the standardized fields from the OTH record stored in elasticsearch
type OthEntry struct {
	Title                     string   `json:"title" validate:"required,max=40"`
	ProductType               string   `json:"productType"  validate:"required,max=40"`
	Filename                  string   `json:"filename"  validate:"required,max=40,alphanum"`
	ProductDateDayOfMonth     string   `json:"productDateDayOfMonth"  validate:"required,number,max=3"`
	ProductDateMonth          string   `json:"productDateMonth"  validate:"required,max=6,alpha"`
	ProductDateYear           string   `json:"productDateYear"  validate:"required,max=4,alphanum"`
	ProductStartTime          string   `json:"productStartTime"  validate:"required,number,max=4"`
	ProductEndTime            string   `json:"productEndTime"  validate:"required,number,max=4"`
	ProductFormat             string   `json:"productFormat"  validate:"required,alpha,max=20"`
	ProductStartTimestamp     string   `json:"productStartTimestamp,omitempty"  validate:"required,max=40"`
	ProductEndTimestamp       string   `json:"productEndTimestamp,omitempty"  validate:"required,max=40"`
	PublicationDateDayOfMonth string   `json:"publicationDateDayOfMonth"  validate:"required,number,max=4"`
	PublicationDateMonth      string   `json:"publicationDateMonth"  validate:"required,max=6,alphanum"`
	PublicationDateYear       string   `json:"publicationDateYear"  validate:"required,number,max=4"`
	PublicationTimestamp      string   `json:"publicationTimestamp,omitempty" validate:"max=40"`
	KmlPresent                bool     `json:"kmlPresent" xml:"kmlPresent"`
	ImageFormat               string   `json:"imageFormat" xml:"imageFormat" validate:"max=20"`
	ReportSummary             string   `json:"reportSummary"  validate:"required"`
	ListOfCountries           []string `json:"listOfCountries"  validate:"required"`
	ListOfLocations           []string `json:"listOfLocations"  validate:"required"`
	ListOfRequirementNumbers  []string `json:"listOfRequirementNumbers"  validate:"required"`
	ListOfCustomers           []string `json:"listOfCustomers"  validate:"required"`
	ListOfTargetTypes         []string `json:"listOfTargetTypes"  validate:"required"`
	ListOfSensors             []string `json:"listOfSensors"  validate:"required"`
	SystemStatus              string   `json:"systemStatus" xml:"systemStatus" validate:"max=50"`
	RelatedReports            []string `json:"relatedReports,omitempty" xml:"relatedReports"`
	Highlighted               string   `json:"highlighted" xml:"highlighted" validate:"max=20"`
	ImageLink                 string   `json:"imageLink,omitempty"  validate:"required,url"`
	ThumbnailLink             string   `json:"thumbnailLink,omitempty"  validate:"required,url"`
	ReportLink                string   `json:"reportLink,omitempty"  validate:"required,url"`
}

//TrimFields trims all string fields removing spaces and endlines returning a new trimmed othentry
func (entry OthEntry) TrimFields() OthEntry {
	return OthEntry{
		Title:                     TrimString(entry.Title),
		ProductType:               TrimString(entry.ProductType),
		Filename:                  TrimString(entry.Filename),
		ProductDateDayOfMonth:     TrimString(entry.ProductDateDayOfMonth),
		ProductDateMonth:          TrimString(entry.ProductDateMonth),
		ProductDateYear:           TrimString(entry.ProductDateYear),
		ProductStartTime:          TrimString(entry.ProductStartTime),
		ProductEndTime:            TrimString(entry.ProductEndTime),
		ProductFormat:             TrimString(entry.ProductFormat),
		ProductStartTimestamp:     TrimString(entry.ProductStartTimestamp),
		ProductEndTimestamp:       TrimString(entry.ProductEndTimestamp),
		PublicationDateDayOfMonth: TrimString(entry.PublicationDateDayOfMonth),
		PublicationDateMonth:      TrimString(entry.PublicationDateMonth),
		PublicationDateYear:       TrimString(entry.PublicationDateYear),
		PublicationTimestamp:      TrimString(entry.PublicationTimestamp),
		KmlPresent:                entry.KmlPresent,
		ImageFormat:               TrimString(entry.ImageFormat),
		ReportSummary:             entry.ReportSummary,
		ListOfCountries:           TrimFromList(entry.ListOfCountries),
		ListOfLocations:           TrimFromList(entry.ListOfLocations),
		ListOfRequirementNumbers:  TrimFromList(entry.ListOfRequirementNumbers),
		ListOfCustomers:           TrimFromList(entry.ListOfCustomers),
		ListOfTargetTypes:         TrimFromList(entry.ListOfTargetTypes),
		ListOfSensors:             TrimFromList(entry.ListOfSensors),
		SystemStatus:              TrimString(entry.SystemStatus),
		RelatedReports:            TrimFromList(entry.RelatedReports),
		Highlighted:               TrimString(entry.Highlighted),
		ImageLink:                 TrimString(entry.ImageLink),
		ThumbnailLink:             TrimString(entry.ThumbnailLink),
		ReportLink:                TrimString(entry.ReportLink),
	}
}

//OthRecord holds the contents of a csv line and a Done chan to indicate complete processing
type OthRecord struct {
	CsvContents map[string]string
	DoneChan    chan bool
}

//JsonIngestEntry used to pass json file contents with a channel to indicate ingestion is complete
type JsonIngestEntry struct {
	Entry    OthEntry
	DoneChan chan bool
}

//Identical to above but we need a separate struct to deserialize lists in xml
type OthXmlEntry struct {
	Title                     string             `json:"title" xml:"title" validate:"required"`
	ProductType               string             `json:"productType" xml:"productType" validate:"required"`
	Filename                  string             `json:"filename" xml:"filename" validate:"required"`
	ProductDateDayOfMonth     string             `json:"productDateDayOfMonth" xml:"productDateDayOfMonth" validate:"required"`
	ProductDateMonth          string             `json:"productDateMonth" xml:"productDateMonth" validate:"required"`
	ProductDateYear           string             `json:"productDateYear" xml:"productDateYear" validate:"required"`
	ProductStartTime          string             `json:"productStartTime" xml:"productStartTime" validate:"required"`
	ProductEndTime            string             `json:"productEndTime" xml:"productEndTime" validate:"required"`
	ProductFormat             string             `json:"productFormat" xml:"productFormat" validate:"required"`
	ProductStartTimestamp     string             `json:"productStartTimestamp,omitempty" xml:"productStartTimestamp" validate:"required"`
	ProductEndTimestamp       string             `json:"productEndTimestamp,omitempty" xml:"productEndTimestamp" validate:"required"`
	PublicationDateDayOfMonth string             `json:"publicationDateDayOfMonth" xml:"publicationDateDayOfMonth" validate:"required"`
	PublicationDateMonth      string             `json:"publicationDateMonth" xml:"publicationDateMonth" validate:"required"`
	PublicationDateYear       string             `json:"publicationDateYear" xml:"publicationDateYear" validate:"required"`
	PublicationTimestamp      string             `json:"publicationTimestamp,omitempty" xml:"publicationTimestamp" `
	KmlPresent                bool               `json:"kmlPresent" xml:"kmlPresent"`
	ImageFormat               string             `json:"imageFormat" xml:"imageFormat"`
	ReportSummary             string             `json:"reportSummary" xml:"reportSummary" validate:"required"`
	ListOfCountries           Countries          `json:"listOfCountries" xml:"listOfCountries" validate:"required"`
	ListOfLocations           Locations          `json:"listOfLocations" xml:"listOfLocations" validate:"required"`
	ListOfRequirementNumbers  RequirementNumbers `json:"listOfRequirementNumbers" xml:"listOfRequirementNumbers" validate:"required"`
	ListOfCustomers           Customers          `json:"listOfCustomers" xml:"listOfCustomers" validate:"required"`
	ListOfTargetTypes         Targets            `json:"listOfTargetTypes" xml:"listOfTargetTypes" validate:"required"`
	ListOfSensors             Sensors            `json:"listOfSensors" xml:"listOfSensors" validate:"required"`
	SystemStatus              string             `json:"systemStatus" xml:"systemStatus"`
	RelatedReports            Reports            `json:"relatedReports,omitempty" xml:"relatedReports"`
	Highlighted               string             `json:"highlighted" xml:"highlighted"`
	ImageLink                 string             `json:"imageLink" xml:"imageLink" validate:"required,url"`
	ThumbnailLink             string             `json:"thumbnailLink" xml:"thumbnailLink" validate:"required,url"`
	ReportLink                string             `json:"reportLink" xml:"reportLink" validate:"required,url"`
}

//xml list deserializers
type Locations struct {
	Elements []string `xml:"location"`
}

type RequirementNumbers struct {
	Elements []string `xml:"element"`
}

type Sensors struct {
	Elements []string `xml:"sensor"`
}

type Customers struct {
	Elements []string `xml:"customer"`
}

type Targets struct {
	Elements []string `xml:"target"`
}

type Reports struct {
	Elements []string `xml:"report"`
}

type Countries struct {
	Elements []string `xml:"country"`
}

//we need a convenient way to get to these lists without writing branching code in non model libraries
//XmlList defines an interface for any xml list to return its string elements
type XmlList interface {
	GetElements() []string
}

func (c Reports) GetElements() []string {
	return c.Elements
}

func (c Countries) GetElements() []string {
	return c.Elements
}

func (c Targets) GetElements() []string {
	return c.Elements
}

func (c Customers) GetElements() []string {
	return c.Elements
}

func (c Sensors) GetElements() []string {
	return c.Elements
}

func (c RequirementNumbers) GetElements() []string {
	return c.Elements
}

func (c Locations) GetElements() []string {
	return c.Elements
}

//helpful string utilities

func TrimFromList(items []string) []string {
	retItems := make([]string, len(items))
	for i, _ := range items {
		retItems[i] = TrimString(items[i])
	}
	return retItems
}

func TrimString(s string) string {
	return strings.Trim(s, " \n")
}
