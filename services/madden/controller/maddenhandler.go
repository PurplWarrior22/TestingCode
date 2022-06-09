package controller

import (
	"../services/madden/dataservice"
	"../services/madden/swagger"
	"../services/maddendb"
)

type maddenHandler struct {
	dataservice dataservice.MaddenDataService
}

const (
	MINIMUM_SUMMARY_LENGTH                                = 10
	START_OF_TIME                                         = "1971-01-01T15:00:00.01Z"
	END_OF_TIME                                           = "2230-01-01T15:00:00.01Z"
	PAGE_NUMBER_DEFAULT                                   = 0
	PAGE_SIZE_DEFAULT                                     = 25
	DEFAULT_SORT           swagger.GetEntryParamsSort     = "startDate"
	DEFAULT_HISTORIC       swagger.GetEntryParamsHistoric = "historic"
)

//constructor

func NewMaddenServerHandler(dataservice dataservice.MaddenDataService) swagger.ServerInterface {
	return &maddenHandler{dataservice: dataservice}
}

func (handler *maddenHandler) GetSummary(ctx echo.Context) error {
	summary, err := handler.dataservice.GetSummary()
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, summary)
}

func (handler *maddenHandler) PostSummary(ctx echo.Context) error {
	summary := swagger.Summary{}
	err := ctx.Bind(&summary)
	if err != nil {
		fmt.Printf("Error reading request body ERROR: %s\n", err.Error())
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "unable to read request body",
		})
	}
	if err := validateSummary(summary); err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	created, err := handler.dataservice.CreateSummary(summary)
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, created)
}

func (handler *maddenHandler) GetPublished(ctx echo.Context) error {
	published, err := handler.dataservice.GetPublished()
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, published)
}

func (handler *maddenHandler) PostPublished(ctx echo.Context) error {
	published := swagger.Published{}
	err := ctx.Bind(&published)
	if err != nil {
		fmt.Printf("Error reading request body ERROR: %s\n", err.Error())
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "unable to read request body",
		})
	}
	created, err := handler.dataservice.CreatePublished(published)
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, created)
}

func (handler *maddenHandler) GetEntry(ctx echo.Context, params swagger.GetEntryParams) error {
	filledParams := fillParamDefaults(params)
	if !paramsValid(filledParams) {
		return ctx.JSON(http.StatusBadRequest, swagger.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters",
		})
	}
	items := []swagger.MaddenItem{}
	var err error
	if params.Id != nil {
		items, err = handler.getSingleItem(params)
	} else {
		items, err = handler.dataservice.GetMaddenEntries(filledParams)
	}

	if err != nil {
		ctx.JSON(utilities.StatusCodeError(err), swagger.Error{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, swagger.MaddenItems{
		Entries: items,
	})
}

func (handler *maddenHandler) PostEntry(ctx echo.Context) error {
	itemBody := swagger.MaddenItem{}
	if err := ctx.Bind(&itemBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	if err := createEntryValid(itemBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	created, err := handler.dataservice.CreateEntry(itemBody)
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, created)
}

func (handler *maddenHandler) PutEntryMaddenId(ctx echo.Context, maddenId int) error {
	itemBody := swagger.MaddenItem{}
	if err := ctx.Bind(&itemBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	if err := updateEntryValid(itemBody, maddenId); err != nil {
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	updated, err := handler.dataservice.UpdateEntry(itemBody)
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, updated)
}

func (handler *maddenHandler) DeleteEntryMaddenId(ctx echo.Context, maddenId int) error {
	if err := deleteEntryValid(maddenId); err != nil {
		return ctx.JSON(http.StatusBadRequest, swagger.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	err := handler.dataservice.DeleteEntry(maddenId)
	if err != nil {
		return ctx.JSON(utilities.StatusCodeError(err), swagger.ErrorResponse{
			Code:    utilities.StatusCodeError(err),
			Message: err.Error(),
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

//implementation helpers

//getSingleItem retrieves a single item and returns it as the single item in a slice
func (handler *maddenHandler) getSingleItem(params swagger.GetEntryParams) ([]swagger.MaddenItem, error) {
	item, err := handler.dataservice.GetMaddenById(*params.Id)
	if err != nil {
		return nil, err
	}
	return []swagger.MaddenItem{item}, nil
}

//deleteEntryValid ensures that a given id is valid
func deleteEntryValid(id int) error {
	if id < 0 {
		return fmt.Errorf("item id was invalid must be positive integer")
	}
	return nil
}

//updateEntryValid ensures additional fields in an update entry are present and valid
func updateEntryValid(item swagger.MaddenItem, id int) error {
	if err := createEntryValid(item); err != nil {
		return err
	}
	if item.Id == nil || *item.Id < 0 {
		return fmt.Errorf("item id was invalid must be positive integer")
	}
	if *item.Id != id {
		return fmt.Errorf("path id must match id of item in body")
	}
	return nil
}

//createEntryValid returns an error with detailed message if the passed MaddenItem is invalid
func createEntryValid(item swagger.MaddenItem) error {
	if item.Images == nil || len(item.Images) < 1 || len(item.Images) > 2 {
		return fmt.Errorf("must supply at least one and no more than 2 images")
	}
	if len(item.Summary) < MINIMUM_SUMMARY_LENGTH {
		return fmt.Errorf("item summary must be at least %d in length", MINIMUM_SUMMARY_LENGTH)
	}
	if _, err := time.Parse(time.RFC3339, item.EndDate); err != nil {
		return fmt.Errorf("time format of end date was not valid, expect RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, item.StartDate); err != nil {
		return fmt.Errorf("time format of start date was not valid, expect RFC3339")
	}
	for _, image := range item.Images {
		if err := validateImage(image); err != nil {
			return err
		}
	}
	return nil
}

//validateImage returns an error with detailed message if the passed MaddenImage is invalid
func validateImage(image swagger.MaddenImage) error {
	if image.Id < 0 {
		return fmt.Errorf("image id was invalid")
	}
	if image.Status != "FMC" && image.Status != "PMC" && image.Status != "NMC" {
		return fmt.Errorf("status was not a valid status must be one of FMC PMC NMC")
	}
	return nil
}

//paramsValid returns true if the passed parameters are valid and safe, false otherwise
func paramsValid(params swagger.GetEntryParams) bool {
	return *params.Id >= 0 && *params.PageNumber >= 0 && *params.PageSize >= 0 && validDate(*params.EndDate) && validDate(*params.StartDate)
}

//fillParamDefaults replaces any nil pointers with their default values
func fillParamDefaults(params swagger.GetEntryParams) swagger.GetEntryParams {
	if params.EndDate == nil {
		params.EndDate = utilities.StrPtr(START_OF_TIME)
	}
	if params.StartDate == nil {
		params.StartDate = utilities.StrPtr(END_OF_TIME)
	}
	if params.PageNumber == nil {
		params.PageNumber = utilities.IntPtr(PAGE_NUMBER_DEFAULT)
	}
	if params.PageSize == nil {
		params.PageSize = utilities.IntPtr(PAGE_SIZE_DEFAULT)
	}
	if params.Sort == nil {
		params.Sort = entryParamSortPtr(DEFAULT_SORT)
	}
	if params.Id == nil {
		params.Id = utilities.IntPtr(0)
	}
	if params.Historic == nil {
		params.Historic = entryParamHistoricPtr("non-historic")
	}
	return params
}

//validateSummary confirms a text summary field is present and of a minimum length
func validateSummary(summary swagger.Summary) error {
	if len(summary.Summary) < 10 {
		return models.NewDataServiceError("Summary did not meet minimum length of 10", http.StatusBadRequest)
	}
	return nil
}

//validDate returns true if the passed string conforms to RFC3339
func validDate(dateString string) bool {
	_, err := time.Parse(time.RFC3339, dateString)
	return err == nil
}

func entryParamHistoricPtr(param swagger.GetEntryParamsHistoric) *swagger.GetEntryParamsHistoric {
	return &param
}

func entryParamSortPtr(entryParam swagger.GetEntryParamsSort) *swagger.GetEntryParamsSort {
	return &entryParam
}
