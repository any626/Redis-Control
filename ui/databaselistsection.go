package ui

import (
	"github.com/any626/Redis-Control/services"
	"github.com/therecipe/qt/widgets"
	"strconv"
)

type DatabaseListSection struct {
	*widgets.QWidget
	Layout *widgets.QHBoxLayout
	List   *widgets.QListWidget
}

func NewDatabaseListSection() *DatabaseListSection {
	dbListSection := &DatabaseListSection{widgets.NewQWidget(nil, 0), widgets.NewQHBoxLayout(), widgets.NewQListWidget(nil)}
	dbListSection.SetLayout(dbListSection.Layout)
	dbListSection.Layout.AddWidget(dbListSection.List, 0, 0)
	return dbListSection
}

func (dbls *DatabaseListSection) AddDatabases(rService *services.RedisService) error {
	dbCount, err := rService.GetDatabaseCount()
	if err != nil {
		return err
	}

	for i := 0; i < dbCount; i++ {
		dbls.List.AddItem(strconv.Itoa(i))
	}
	dbls.List.Item(rService.Config.Database).SetSelected(true)
	return nil
}
