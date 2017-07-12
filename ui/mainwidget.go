package ui

import (
	"fmt"
	"github.com/any626/Redis-Control/models"
	"github.com/any626/Redis-Control/services"
	"github.com/therecipe/qt/widgets"
	"strconv"
)

type MainWidget struct {
	*widgets.QWidget
	DatabaseList *DatabaseListSection
	KeyList      *KeyListSection

	StackedWidget  *widgets.QStackedWidget
	TextEdit       *widgets.QTextEdit
	TableWidget    *widgets.QTableWidget
	DataListWidget *widgets.QListWidget
	BlankWidget    *widgets.QWidget

	Layout *widgets.QHBoxLayout
}

func NewMainWidget() *MainWidget {
	mainWidget := widgets.NewQWidget(nil, 0)
	dbList := NewDatabaseListSection()
	keyList := NewKeyListSection()

	textEdit := widgets.NewQTextEdit(nil)
	tableWidget := widgets.NewQTableWidget(nil)
	dataListWidget := widgets.NewQListWidget(nil)
	blankWidget := widgets.NewQWidget(nil, 0)

	stackedWidget := widgets.NewQStackedWidget(nil)
	stackedWidget.AddWidget(textEdit)
	stackedWidget.AddWidget(tableWidget)
	stackedWidget.AddWidget(dataListWidget)
	stackedWidget.AddWidget(blankWidget)
	stackedWidget.SetCurrentWidget(blankWidget)

	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(dbList, 1, 0)
	layout.AddWidget(keyList, 3, 0)
	layout.AddWidget(stackedWidget, 6, 0)

	mainWidget.SetLayout(layout)

	return &MainWidget{mainWidget, dbList, keyList, stackedWidget, textEdit, tableWidget, dataListWidget, blankWidget, layout}
}

func (m *MainWidget) Init(rService *services.RedisService) {
	err := m.DatabaseList.AddDatabases(rService)
	if err != nil {
		fmt.Println(err)
		return
	}

	m.KeyList.Layout.KeyListWidget.PopulateKeys(rService)

	headerView := m.TableWidget.HorizontalHeader()
	headerView.SetStretchLastSection(true)
	verticalHeaderView := m.TableWidget.VerticalHeader()
	verticalHeaderView.SetVisible(false)

	m.DataListWidget.SetAlternatingRowColors(true)

	m.KeyList.Layout.KeyListWidget.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		key := item.Text()
		model, err := rService.GetModel(key)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch model.(type) {
		case *models.String:
			rString := model.(*models.String)
			value, err := rString.GetValue()
			if err != nil {
				fmt.Println(err)
				m.StackedWidget.SetCurrentWidget(m.BlankWidget)
				return
			}
			m.StackedWidget.SetCurrentWidget(m.TextEdit.QTextEdit_PTR())
			m.TextEdit.SetPlainText(value)
		case *models.List:
			rList := model.(*models.List)
			value, err := rList.GetValue()
			if err != nil {
				fmt.Println(err)
				m.StackedWidget.SetCurrentWidget(m.BlankWidget)
				return
			}
			m.TableWidget.Clear()
			m.TableWidget.SetColumnCount(2)
			m.TableWidget.SetHorizontalHeaderLabels([]string{"index", "value"})

			rowCount := len(value)
			m.TableWidget.SetRowCount(rowCount)
			for i, v := range value {
				m.TableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(i), 0))
				m.TableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v, 0))
			}

			m.StackedWidget.SetCurrentWidget(m.TableWidget.QTableWidget_PTR())
		case *models.Hash:
			rHash := model.(*models.Hash)
			value, err := rHash.GetValue()
			if err != nil {
				fmt.Println(err)
				m.StackedWidget.SetCurrentWidget(m.BlankWidget)
				return
			}
			m.TableWidget.Clear()
			m.StackedWidget.SetCurrentWidget(m.TableWidget.QTableWidget_PTR())
			m.TableWidget.SetColumnCount(2)
			m.TableWidget.SetHorizontalHeaderLabels([]string{"field", "value"})

			rowCount := len(value)
			m.TableWidget.SetRowCount(rowCount)
			for i, v := range value {
				m.TableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(v.Field, 0))
				m.TableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v.Value, 0))
			}
		case *models.Set:
			rSet := model.(*models.Set)
			value, err := rSet.GetValue()
			if err != nil {
				fmt.Println(err)
				m.StackedWidget.SetCurrentWidget(m.BlankWidget)
				return
			}

			m.DataListWidget.Clear()
			m.StackedWidget.SetCurrentWidget(m.DataListWidget.QListWidget_PTR())
			for _, v := range value {
				m.DataListWidget.AddItem(v)
			}
		case *models.ZSet:
			rZSet := model.(*models.ZSet)
			value, err := rZSet.GetValue()
			if err != nil {
				fmt.Println(err)
				m.StackedWidget.SetCurrentWidget(m.BlankWidget)
				return
			}
			m.TableWidget.Clear()
			m.StackedWidget.SetCurrentWidget(m.TableWidget.QTableWidget_PTR())
			m.TableWidget.SetColumnCount(3)
			m.TableWidget.SetHorizontalHeaderLabels([]string{"index", "value", "score"})

			rowCount := len(value)
			m.TableWidget.SetRowCount(rowCount)
			for i, v := range value {
				m.TableWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(i), 0))
				m.TableWidget.SetItem(i, 1, widgets.NewQTableWidgetItem2(v.Value, 0))
				m.TableWidget.SetItem(i, 2, widgets.NewQTableWidgetItem2(strconv.FormatFloat(v.Score, 'f', -1, 64), 0))
			}
		default:
			fmt.Println("Unknown Type")
			m.StackedWidget.SetCurrentWidget(m.BlankWidget)
			return
		}
	})

	m.DatabaseList.List.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		database := item.Text()
		dbId, err := strconv.Atoi(database)
		if err != nil {
			fmt.Println(err)
			if m.DatabaseList.List.Count() > 0 {
				m.DatabaseList.List.Item(0).SetSelected(true)
			}
		}
		rService.Config.Database = dbId
		rService.UpdatePool()
		m.StackedWidget.SetCurrentWidget(m.BlankWidget)

		m.KeyList.Layout.KeyListWidget.PopulateKeys(rService)
		m.KeyList.Layout.KeyFilterWidget.Clear()
	})

}
