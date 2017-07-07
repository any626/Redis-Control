package ui

import (
    "github.com/therecipe/qt/widgets"
)

type ControlWidget struct {
    *widgets.QWidget
    Layout *ControlLayout
}

func NewControlWidget() *ControlWidget {
    controlWidget := &ControlWidget{widgets.NewQWidget(nil, 0), NewControlLayout()}
    // controlWidget.Layout = NewControlLayout()

    controlWidget.SetLayout(controlWidget.Layout)
    return controlWidget
}

type ControlLayout struct {
    *widgets.QVBoxLayout
    DeleteButton *widgets.QPushButton
    SaveButton *widgets.QPushButton
}

func NewControlLayout() *ControlLayout {
    controlLayout := &ControlLayout{widgets.NewQVBoxLayout(), widgets.NewQPushButton2("Delete", nil), widgets.NewQPushButton2("Save", nil)}
    controlLayout.AddWidget(controlLayout.DeleteButton, 0, 0)
    controlLayout.AddWidget(controlLayout.SaveButton, 0, 0)

    return controlLayout
}