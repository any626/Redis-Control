package ui

import (
	// "fmt"
	"github.com/any626/Redis-Control/services"
	"github.com/therecipe/qt/widgets"
	"strconv"
	// "github.com/therecipe/qt/core"
)

type ConnectionWidget struct {
	*widgets.QWidget
	Layout                *widgets.QVBoxLayout
	Label                 *widgets.QLineEdit
	Host                  *widgets.QLineEdit
	Port                  *widgets.QLineEdit
	Database              *widgets.QLineEdit
	SSHAddress            *widgets.QLineEdit
	SSHPort               *widgets.QLineEdit
	SSHUser               *widgets.QLineEdit
	SSHPrivateKey         *widgets.QLineEdit
	SSHPrivateKeyPassword *widgets.QLineEdit

	SaveButton    *widgets.QPushButton
	ConnectButton *widgets.QPushButton
}

func NewConnectionWidget() *ConnectionWidget {

	label := widgets.NewQLineEdit(nil)
	label.SetFixedWidth(300)

	host := widgets.NewQLineEdit(nil)
	host.SetPlaceholderText("127.0.0.1")
	host.SetFixedWidth(300)

	port := widgets.NewQLineEdit(nil)
	port.SetPlaceholderText("6379")
	port.SetFixedWidth(300)

	db := widgets.NewQLineEdit(nil)
	db.SetPlaceholderText("0")
	db.SetFixedWidth(300)

	sshAddress := widgets.NewQLineEdit(nil)
	sshAddress.SetFixedWidth(300)

	sshPort := widgets.NewQLineEdit(nil)
	sshPort.SetFixedWidth(300)

	sshUser := widgets.NewQLineEdit(nil)
	sshUser.SetFixedWidth(300)

	sshPrivateKey := widgets.NewQLineEdit(nil)
	sshPrivateKey.SetFixedWidth(300)

	sshPrivateKeyPassword := widgets.NewQLineEdit(nil)
	sshPrivateKeyPassword.SetFixedWidth(300)

	saveButton := widgets.NewQPushButton2("Save", nil)
	saveButton.SetFixedWidth(300)

	connectButton := widgets.NewQPushButton2("Connect", nil)
	connectButton.SetFixedWidth(300)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(label, 1, 0)
	layout.AddWidget(host, 1, 0)
	layout.AddWidget(port, 1, 0)
	layout.AddWidget(db, 1, 0)
	layout.AddWidget(sshAddress, 1, 0)
	layout.AddWidget(sshPort, 1, 0)
	layout.AddWidget(sshUser, 1, 0)
	layout.AddWidget(sshPrivateKey, 1, 0)
	layout.AddWidget(sshPrivateKeyPassword, 1, 0)
	layout.AddWidget(saveButton, 1, 0)
	layout.AddWidget(connectButton, 1, 0)

	connWidget := &ConnectionWidget{widgets.NewQWidget(nil, 0), layout, label, host, port, db, sshAddress, sshPort, sshUser, sshPrivateKey, sshPrivateKeyPassword, saveButton, connectButton}

	connWidget.SetLayout(layout)
	return connWidget
}

func (c *ConnectionWidget) GetRedisConfig() *services.RedisConfig {
	host := c.Host.Text()
	if host == "" {
		host = "127.0.0.1"
	}

	port, err := strconv.Atoi(c.Port.Text())
	if err != nil {
		port = 6379
	}

	db, err := strconv.Atoi(c.Database.Text())
	if err != nil {
		db = 0
	}

	return &services.RedisConfig{Host: host, Port: port, Database: db}
}
