package ui

import (
	"fmt"
	// "github.com/any626/Redis-Control/services"
    "github.com/any626/Redis-Control/utils"
	"github.com/therecipe/qt/widgets"
	"strconv"
	"github.com/therecipe/qt/core"
)

type ConnectionWidget struct {
	*widgets.QWidget
    Layout                *widgets.QHBoxLayout
    LeftWidget            *widgets.QWidget
    RightWidget           *widgets.QWidget

	LeftLayout            *widgets.QVBoxLayout
	Name                  *widgets.QLineEdit
	Host                  *widgets.QLineEdit
	Port                  *widgets.QLineEdit
	Database              *widgets.QLineEdit
	SSHHost               *widgets.QLineEdit
	SSHPort               *widgets.QLineEdit
	SSHUser               *widgets.QLineEdit
	SSHPassword           *widgets.QLineEdit
    SSHPrivateKey         *widgets.QLineEdit

    RightLayout           *widgets.QVBoxLayout
    Favorites             []*FavoriteWidget

	SaveButton    *widgets.QPushButton
	ConnectButton *widgets.QPushButton
}

func NewConnectionWidget() *ConnectionWidget {

	name := widgets.NewQLineEdit(nil)
    name.SetPlaceholderText("Name")
	name.SetFixedWidth(300)

	host := widgets.NewQLineEdit(nil)
	host.SetPlaceholderText("127.0.0.1")
	host.SetFixedWidth(300)

	port := widgets.NewQLineEdit(nil)
	port.SetPlaceholderText("6379")
	port.SetFixedWidth(300)

	db := widgets.NewQLineEdit(nil)
	db.SetPlaceholderText("0")
	db.SetFixedWidth(300)

	sshHost := widgets.NewQLineEdit(nil)
    sshHost.SetPlaceholderText("SSH Address")
	sshHost.SetFixedWidth(300)

	sshPort := widgets.NewQLineEdit(nil)
    sshPort.SetPlaceholderText("SSH Port")
	sshPort.SetFixedWidth(300)

	sshUser := widgets.NewQLineEdit(nil)
    sshUser.SetPlaceholderText("SSH User")
	sshUser.SetFixedWidth(300)

	sshPassword := widgets.NewQLineEdit(nil)
    sshPassword.SetPlaceholderText("SSH Password")
	sshPassword.SetFixedWidth(300)

    sshPrivateKey := widgets.NewQLineEdit(nil)
    sshPrivateKey.SetPlaceholderText("SSH Private Key")
    sshPrivateKey.SetFixedWidth(300)


	saveButton := widgets.NewQPushButton2("Save", nil)
	saveButton.SetFixedWidth(300)

	connectButton := widgets.NewQPushButton2("Connect", nil)
	connectButton.SetFixedWidth(300)

	leftLayout := widgets.NewQVBoxLayout()
	leftLayout.AddWidget(name, 1, 0)
	leftLayout.AddWidget(host, 1, 0)
	leftLayout.AddWidget(port, 1, 0)
	leftLayout.AddWidget(db, 1, 0)
	leftLayout.AddWidget(sshHost, 1, 0)
	leftLayout.AddWidget(sshPort, 1, 0)
	leftLayout.AddWidget(sshUser, 1, 0)
	leftLayout.AddWidget(sshPassword, 1, 0)
    leftLayout.AddWidget(sshPrivateKey, 1, 0)
	leftLayout.AddWidget(saveButton, 1, 0)
	leftLayout.AddWidget(connectButton, 1, 0)

    rightLayout := widgets.NewQVBoxLayout()
    var favs []*FavoriteWidget

    layout := widgets.NewQHBoxLayout()
    leftWidget := widgets.NewQWidget(nil, 0)
    rightWidget := widgets.NewQWidget(nil, 0)
    layout.AddWidget(leftWidget, 1, 0)
    layout.AddWidget(rightWidget, 1, 0)

    leftWidget.SetLayout(leftLayout)
    rightWidget.SetLayout(rightLayout)


	connWidget := &ConnectionWidget{
        widgets.NewQWidget(nil, 0),
        layout,
        leftWidget,
        rightWidget,
        leftLayout,
        name,
        host,
        port,
        db,
        sshHost,
        sshPort,
        sshUser,
        sshPassword,
        sshPrivateKey,
        rightLayout,
        favs,
        saveButton,
        connectButton,
    }

	connWidget.SetLayout(layout)
	return connWidget
}

func (c *ConnectionWidget) Init() {
    filPath := core.QStandardPaths_WritableLocation(core.QStandardPaths__AppLocalDataLocation) + "/connections.json"
    connections, err := utils.ReadConnections(filPath)
    if err != nil {
        fmt.Println(err)
        return
    }

    var favs []*FavoriteWidget
    for _, v := range connections {
        favWidget := NewFavoriteWidget(v)
        favs = append(favs, favWidget)
        c.RightLayout.AddWidget(favWidget, 0, 1)
    }

    c.Favorites = favs
}

func (c *ConnectionWidget) GetConnection() *utils.Connection {

    name := c.Name.Text()

	host := c.Host.Text()
	if host == "" {
		host = "127.0.0.1"
	}

	port, err := strconv.Atoi(c.Port.Text())
	if err != nil {
		port = 6379
	}

    sshHost := c.SSHHost.Text()

    sshPort, err := strconv.Atoi(c.SSHPort.Text())
    if err != nil {
        sshPort = 0
    }
    sshUser := c.SSHUser.Text()
    sshPassword := c.SSHPassword.Text()
    sshPrivateKey := c.SSHPrivateKey.Text()

    sshEnabled := false
    if sshHost != "" {
        sshEnabled = true
    }

    ssh := utils.Connection{
        Name: name,
        Host: host,
        Port: port,
        SSHEnabled: sshEnabled,
        SSHHost: sshHost,
        SSHPort: sshPort,
        SSHUser: sshUser,
        SSHPassword: sshPassword,
        SSHPrivateKey: sshPrivateKey,
    }

	return &ssh
}

type FavoriteWidget struct {
    *widgets.QWidget
    Layout           *widgets.QHBoxLayout
    Label            *widgets.QLabel
    EditButton       *widgets.QPushButton
    ConnectButton    *widgets.QPushButton
    Connection       *utils.Connection
}

func NewFavoriteWidget(c *utils.Connection) *FavoriteWidget {
    layout := widgets.NewQHBoxLayout()
    label := widgets.NewQLabel2(c.Name, nil, 0)
    editButton := widgets.NewQPushButton2("Edit", nil)
    connectButton := widgets.NewQPushButton2("Connect", nil)

    layout.AddWidget(label, 3, 0)
    layout.AddWidget(editButton, 1, 0)
    layout.AddWidget(connectButton, 1, 0)

    conn := &FavoriteWidget{
        widgets.NewQWidget(nil, 0),
        layout,
        label,
        editButton,
        connectButton,
        c,
    }

    conn.SetLayout(layout)
    return conn
}