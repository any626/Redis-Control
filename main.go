package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	// "github.com/therecipe/qt/core"
	"github.com/any626/Redis-Control/services"
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
	// fmt.Println(core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation))
	// connection := &services.RedisConfig{Host: "localhost", Port: 6379}

	// rService := services.NewRedisService(connection)

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationName("Redis-Control")

	// connectionJson, err := json.Marshal(connection)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// _ = os.Mkdir(core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation), 0777)

	// err = ioutil.WriteFile(
	//     core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation) + "/connections.json", connectionJson, 0640)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// dat, err := ioutil.ReadFile(core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation) + "/connections.json")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// fmt.Print(string(dat))

	mainStack := widgets.NewQStackedWidget(nil)

	mainWidget := ui.NewMainWidget()

	connectionWidget := ui.NewConnectionWidget()
	mainStack.AddWidget(connectionWidget)
	mainStack.AddWidget(mainWidget)

	connectionWidget.ConnectButton.ConnectClicked(func(checked bool) {
		connectionConfig := connectionWidget.GetRedisConfig()
		rService := services.NewRedisService(connectionConfig)
		mainWidget.Init(rService)
		mainStack.SetCurrentWidget(mainWidget.QWidget_PTR())
	})

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Redis Control")
	window.SetMinimumSize2(600, 400)

	// qHBoxLayout := widgets.NewQHBoxLayout()

	// databaseListSection := ui.NewDatabaseListSection()
	// err = databaseListSection.AddDatabases(rService)
	// if err != nil {
	//     fmt.Println(err)
	//     return
	// }

	// keyListSection := ui.NewKeyListSection()

	// keyListSection.Layout.KeyListWidget.PopulateKeys(rService)

	// textEdit := widgets.NewQTextEdit(nil)

	// tableWidget := widgets.NewQTableWidget(nil)
	// headerView := tableWidget.HorizontalHeader()
	// headerView.SetStretchLastSection(true)
	// verticalHeaderView := tableWidget.VerticalHeader()
	// verticalHeaderView.SetVisible(false)

	// dataListWidget := widgets.NewQListWidget(nil)
	// dataListWidget.SetAlternatingRowColors(true);

	// blankWidget := widgets.NewQWidget(nil, 0)

	// stackedWidget := widgets.NewQStackedWidget(nil)
	// stackedWidget.AddWidget(textEdit)
	// stackedWidget.AddWidget(tableWidget)
	// stackedWidget.AddWidget(dataListWidget)
	// stackedWidget.AddWidget(blankWidget)
	// stackedWidget.SetCurrentWidget(blankWidget)

	// keyListSection.Layout.KeyListWidget.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
	//     key := item.Text()
	//     model, err := rService.GetModel(key)
	//     if err != nil {
	//         fmt.Println(err)
	//         return
	//     }

	//     switch model.(type) {
	//         case *models.String:
	//             rString := model.(*models.String)
	//             value, err := rString.GetValue()
	//             if err != nil {
	//                 fmt.Println(err)
	//                 stackedWidget.SetCurrentWidget(blankWidget)
	//                 return
	//             }
	//             stackedWidget.SetCurrentWidget(textEdit.QTextEdit_PTR())
	//             textEdit.SetPlainText(value)
	//         case *models.List:
	//             rList := model.(*models.List)
	//             value, err := rList.GetValue()
	//             if err != nil {
	//                 fmt.Println(err)
	//                 stackedWidget.SetCurrentWidget(blankWidget)
	//                 return
	//             }
	//             tableWidget.Clear()
	//             tableWidget.SetColumnCount(2)
	//             tableWidget.SetHorizontalHeaderLabels([]string{"index", "value"})

	//             rowCount := len(value)
	//             tableWidget.SetRowCount(rowCount)
	//             for i, v := range value {
	//                 tableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(i), 0))
	//                 tableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v, 0))
	//             }

	//             stackedWidget.SetCurrentWidget(tableWidget.QTableWidget_PTR())
	//         case *models.Hash:
	//             rHash := model.(*models.Hash)
	//             value, err := rHash.GetValue()
	//             if err != nil {
	//                 fmt.Println(err)
	//                 stackedWidget.SetCurrentWidget(blankWidget)
	//                 return
	//             }
	//             tableWidget.Clear()
	//             stackedWidget.SetCurrentWidget(tableWidget.QTableWidget_PTR())
	//             tableWidget.SetColumnCount(2)
	//             tableWidget.SetHorizontalHeaderLabels([]string{"field", "value"})

	//             rowCount := len(value)
	//             tableWidget.SetRowCount(rowCount)
	//             for i, v := range value {
	//                 tableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(v.Field, 0))
	//                 tableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v.Value, 0))
	//             }
	//         case *models.Set:
	//             rSet := model.(*models.Set)
	//             value, err := rSet.GetValue()
	//             if err != nil {
	//                 fmt.Println(err)
	//                 stackedWidget.SetCurrentWidget(blankWidget)
	//                 return
	//             }

	//             dataListWidget.Clear()
	//             stackedWidget.SetCurrentWidget(dataListWidget.QListWidget_PTR())
	//             for _, v := range value {
	//                 dataListWidget.AddItem(v)
	//             }
	//         case *models.ZSet:
	//             rZSet := model.(*models.ZSet)
	//             value, err := rZSet.GetValue()
	//             if err != nil {
	//                 fmt.Println(err)
	//                 stackedWidget.SetCurrentWidget(blankWidget)
	//                 return
	//             }
	//             tableWidget.Clear()
	//             stackedWidget.SetCurrentWidget(tableWidget.QTableWidget_PTR())
	//             tableWidget.SetColumnCount(3)
	//             tableWidget.SetHorizontalHeaderLabels([]string{"index", "value", "score"})

	//             rowCount := len(value)
	//             tableWidget.SetRowCount(rowCount)
	//             for i, v := range value {
	//                 tableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(i), 0))
	//                 tableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v.Value, 0))
	//                 tableWidget.SetItem(i, 2, widgets.NewQTableWidgetItem2(strconv.FormatFloat(v.Score, 'f', -1, 64), 0))
	//             }
	//         default:
	//             fmt.Println("Unknown Type")
	//             stackedWidget.SetCurrentWidget(blankWidget)
	//             return
	//     }
	// })

	// databaseListSection.List.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
	//     database := item.Text()
	//     dbId, err := strconv.Atoi(database)
	//     if err != nil {
	//         fmt.Println(err)
	//         if (databaseListSection.List.Count() > 0) {
	//             databaseListSection.List.Item(0).SetSelected(true)
	//         }
	//     }
	//     connection.Database = dbId
	//     rService.UpdatePool()
	//     stackedWidget.SetCurrentWidget(blankWidget)

	//     keyListSection.Layout.KeyListWidget.PopulateKeys(rService)
	//     keyListSection.Layout.KeyFilterWidget.Clear()

	// })

	// controlWidget := ui.NewControlWidget()

	// qHBoxLayout.AddWidget(databaseListSection, 1, 0)
	// qHBoxLayout.AddWidget(keyListSection, 2, 0)
	// qHBoxLayout.AddWidget(stackedWidget, 6, 0)
	// qHBoxLayout.AddWidget(tableWidget, 0, 0)
	// qHBoxLayout.AddWidget(controlWidget, 1, 0)

	// centralWidget := widgets.NewQWidget(nil, 0)
	// centralWidget.SetLayout(qHBoxLayout)

	window.SetCentralWidget(mainStack)
	window.Show()

	widgets.QApplication_Exec()
}
