package utils

import (
    "os"
    "strings"
    "io/ioutil"
    "encoding/json"
)

type Connection struct {
    Name string `json:"name"`
    Host string `json:"host"`
    Port int `json:"port"`
    Database int `json:"database"`
    SSHEnabled bool `json:"sshEnabled"`
    SSHHost string `json:"sshHost"`
    SSHPort int `json:"sshPort"`
    SSHUser string `json:"sshUser"`
    SSHPassword string `json:sshPassword"`
    SSHPrivateKey string `json:"sshPrivateKey"`
}

func SaveConnection(conn *Connection, filePath string) error {
    var connections []*Connection
    if _, err := os.Stat(filePath); !os.IsNotExist(err) {
        // file exists
        raw, err := ioutil.ReadFile(filePath)
        if err != nil {
            return err
        }
        err = json.Unmarshal(raw, &connections)
        if err != nil {
            return err
        }
    } else {
        // make directory just in case
        _ = os.Mkdir(strings.TrimRight(filePath, "connections.json"), 0777)
    }
    connections = append(connections, conn)

    connectionsJson, err := json.Marshal(connections)

    err = ioutil.WriteFile(filePath, connectionsJson, 0640)
    if err != nil {
        return err
    }
    return nil
}

func UpdateConnection(connections []*Connection, filePath string) error {
    connectionsJson, err := json.Marshal(connections)
    if err != nil {
        return nil
    }
    err = ioutil.WriteFile(filePath, connectionsJson, 0640)
    if err != nil {
        return err
    }
    return nil
}

func ReadConnections(filePath string) ([]*Connection, error) {
    raw, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    var connections []*Connection
    err = json.Unmarshal(raw, &connections)
    if err != nil {
        return nil, err
    }
    return connections, nil
}