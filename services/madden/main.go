package main

import (

	"../services/madden/controller"
	"../services/madden/dataservice"
	"../services/madden/swagger"
	"../services/maddendb"
)

const (
	IMAGE_BASE_PATH_ENV = "IMAGE_PATH"
	SERVER_PORT_ENV     = "SERVER_PORT"
)

var (
	pathBuilder utilities.PathBuilder
	maintData   dataservice.MaddenDataService
	serverPort  = "8080"
)

func init() {
	imageBase, err := utilities.GetEnvAndLogOrError(IMAGE_BASE_PATH_ENV)
	if err != nil {
		os.Exit(1)
	}
	serverPort = utilities.GetEnvDefaultAndLog(SERVER_PORT_ENV, serverPort)
	pathBuilder = utilities.NewSimpleAppender(imageBase)
	db, err := maddendb.BuildPostgresMaddenFromEnvironment()
	if err != nil {
		fmt.Printf("unable to build pg database connection due to ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	if err := db.SetupDatabase(); err != nil {
		fmt.Printf("Error while building database connection")
		convertedErr, _ := err.(*maddendb.DbError)
		fmt.Println(err.Error())
		fmt.Println(convertedErr.OriginalError.Error())
		os.Exit(1)
	}
	maddenData = dataservice.NewPgDataService(db, pathBuilder)
}

//build and run the madden db server
func main() {
	handler := controller.NewMaddenServerHandler(maddenData)
	e := echo.New()
	echopprof.Wrap(e)
	e.Use(middleware.Logger())
	swagger.RegisterHandlers(e, handler)
	log.Fatal(e.Start(fmt.Sprintf(":%s", serverPort)))
}
