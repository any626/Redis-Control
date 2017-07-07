package ui

import (
    "github.com/therecipe/qt/widgets"
    "github.com/any626/GoRedis/services"
    "fmt"
    "strings"
)

type KeyListSection struct {
    *widgets.QWidget
    Layout *KeyListLayout
}

func NewKeyListSection() *KeyListSection {
    keyListSection := &KeyListSection{widgets.NewQWidget(nil, 0), NewKeyListLayout()}
    keyListSection.SetLayout(keyListSection.Layout)
    return keyListSection
}

type KeyListLayout struct {
    *widgets.QVBoxLayout
    KeyFilterWidget *widgets.QLineEdit
    KeyListWidget *KeyListWidget
}

func NewKeyListLayout() *KeyListLayout {
    keyListLayout := &KeyListLayout{widgets.NewQVBoxLayout(), widgets.NewQLineEdit(nil), NewKeyListWidget()}
    keyListLayout.AddWidget(keyListLayout.KeyFilterWidget, 0, 0)
    keyListLayout.AddWidget(keyListLayout.KeyListWidget, 0, 0)


    keyListLayout.KeyFilterWidget.ConnectTextChanged(func (text string) {
        keyListCount := keyListLayout.KeyListWidget.Count()

        for i := 0; i < keyListCount; i++ {
            item := keyListLayout.KeyListWidget.Item(i)
            if !strings.HasPrefix(item.Text(), text) {
                item.SetHidden(true)
            } else if item.IsHidden() {
                item.SetHidden(false)
            }
        }
    })

    return keyListLayout
}

type KeyListWidget struct {
    *widgets.QListWidget
}

func NewKeyListWidget() *KeyListWidget {
    return &KeyListWidget{widgets.NewQListWidget(nil)}
}

func (keyList *KeyListWidget) PopulateKeys(rService *services.RedisService) {
    keysCh := make(chan string)
    keysErrorCh := make(chan error)
    keyList.Clear()
    go rService.GetKeys(keysCh, keysErrorCh)
    go func() {
        Loop:
            for {
                select {
                case k := <-keysCh:
                    keyList.AddItem(k)
                    break
                case keyErr := <-keysErrorCh:
                    if keyErr != nil {
                        fmt.Println(keyErr)
                    }
                    break Loop
                }
                if keysCh == nil && keysErrorCh == nil {
                    break Loop
                }
            }
    }()
}