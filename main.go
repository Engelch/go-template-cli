// main.go

package main

import (
	"fmt"
	"os"
	"time"

	ce "github.com/engelch/debugerrorce/v3"
	"github.com/urfave/cli/v2"
)

// ==============================================================================================
// Constants
// ==============================================================================================

const appVersion = "0.0.0" // version of this ap
const appName = "todo"     // name of the application, enclosing directory should have the same name

// ==============================================================================================
// DataTypes
// ==============================================================================================

// AppData contain public data that anyone may know
type appDataType struct {
	AppName    string `json:"AppName" Xml:"AppName" form:"AppName"`          // set by compilation
	AppVersion string `json:"AppVersion" xml:"AppVersion" form:"AppVersion"` // set by compilation
	Started    string `json:"Started" xml:"Started" form:"Started"`          // set when starting the app
}

// AppSecretData contains AppData and app confidential data that should not be logged and not be easily visible
type appSecretDataType struct {
	appDataType
	Debug bool `json:"Debug" xml:"Debug" form:"Debug"`
}

// ==============================================================================================
// Variables
// ==============================================================================================

var appData = appDataType{AppName: appName, AppVersion: appVersion}
var appSecretData = appSecretDataType{}

// ==============================================================================================
// main app implementation
// ==============================================================================================

func startApp() error {
	fmt.Print("This is the app")
	return nil
}

// ==============================================================================================
// Evaluation of CLI options
// ==============================================================================================

func checkOptions(c *cli.Context) error {
	if c == nil {
		ce.ErrorExit(9, "Context structure is empty.")
	}
	// ListeningIpPort() is not checked, effort vs result
	return nil
}

// commandLineOptions just separates the definition of command line options ==> creating a shorter main
func commandLineOptions() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			Value:       false,
			Usage:       "enable debug mode (optional)",
			EnvVars:     []string{"SVS_FDA_DB_WRAPPER_DEBUG"},
			Destination: &appSecretData.Debug, // not directly required, but if debug enabled should be visible from the status call
		},
	}
}

// main entry point
// EXIT 9 cli options or environment vars with invalid data
// EXIT 10 connection to DB
// EXIT 99 app returned, not expected
func main() {
	var err error
	app := &cli.App{}
	app.Flags = commandLineOptions()
	app.Name = appName
	app.Version = appVersion
	app.Usage = app.Name + ":" + app.Version + ":CRUD REST-API wrapper for psql for OAuth2 implementation"

	app.Action = func(c *cli.Context) error {
		started := ":svc started at UTC:" + time.Now().UTC().Format(time.RFC3339)
		appData.Started = started
		appSecretData.AppName = app.Name
		appSecretData.AppVersion = app.Version
		appSecretData.Started = appData.Started

		if c.Bool("debug") {
			ce.CondDebugSet(true)
			ce.CondDebugln("Debug enabled.")
		}
		err = checkOptions(c) // exits the app in case of error, wrong arguments => exit code 99
		if err != nil {
			ce.ErrorExit(9, err.Error())
		}

		return startApp()
	}

	err = app.Run(os.Args)
	if err != nil {
		ce.ErrorExit(99, "Error from app:"+err.Error())
	}
}

// EOF
