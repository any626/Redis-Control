package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"github.com/therecipe/qt/core"
	"github.com/any626/Redis-Control/services"
    "github.com/any626/Redis-Control/utils"
	// "github.com/any626/Redis-Control/models"
	// "github.com/garyburd/redigo/redis"
	// "fmt"
	// "log"
	// "encoding/json"
	// "io/ioutil"
	// "strconv"
	// "errors"
	// "reflect"
	// "strings"
	// "github.com/therecipe/qt/core"
	// "github.com/therecipe/qt/gui"
	"github.com/any626/Redis-Control/ui"
)

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationName("Redis-Control")

	mainStack := widgets.NewQStackedWidget(nil)

	mainWidget := ui.NewMainWidget()

	connectionWidget := ui.NewConnectionWidget()
	mainStack.AddWidget(connectionWidget)
	mainStack.AddWidget(mainWidget)

    connectionWidget.Init()

	connectionWidget.ConnectButton.ConnectClicked(func(checked bool) {
		connConfig := connectionWidget.GetConnection()
		rService := services.NewRedisService(connConfig)
		mainWidget.Init(rService)
		mainStack.SetCurrentWidget(mainWidget.QWidget_PTR())
	})

    for _, v := range connectionWidget.Favorites {
        v.ConnectButton.ConnectClicked(func(checked bool) {
            rService := services.NewRedisService(v.Connection)
            mainWidget.Init(rService)
            mainStack.SetCurrentWidget(mainWidget.QWidget_PTR())
        })
    }
    connectionWidget.SaveButton.ConnectClicked(func(checked bool) {
        filPath := core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation) + "/connections.json"
        connection := connectionWidget.GetConnection()
        utils.SaveConnection(connection, filPath)
    })

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Redis Control")
	window.SetMinimumSize2(600, 400)

	window.SetCentralWidget(mainStack)
	window.Show()

	widgets.QApplication_Exec()
}
