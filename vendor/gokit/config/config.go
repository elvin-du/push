package config

import (
    "errors"
    "fmt"
    "reflect"
    "strings"
    "io/ioutil"
    "launchpad.net/goyaml"
)

var (
    _data []byte
    _config interface{}
)

func ReadConfig(file string) error {
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return err
    }

    _data = data

    return goyaml.Unmarshal(data, &_config)
}

func Unmarshal(key string, ptr interface{}) error {
    if key == "" {
        return goyaml.Unmarshal(_data, ptr)
    }

    var tmp interface{}
    err := Get(key, &tmp)
    if err != nil {
        return err
    }

    data, err := goyaml.Marshal(tmp)
    if err != nil {
        return err
    }

    return goyaml.Unmarshal(data, ptr)
}

func Get(key string, result interface{}) (err error) {
    defer HandleError(&err)

    keys := strings.Split(key, ":")

    c := _config
    for _, key := range keys {
        if cm, ok := c.(map[interface{}]interface{}); ok {
            if c, ok = cm[key]; ok {
                continue
            }
            return errors.New("key: " + key + " not found")
        }

        return errors.New("not a map")
    }

    v := reflect.ValueOf(result)

    if v.Kind() != reflect.Ptr {
        return errors.New("Need ptr")
    }

    v.Elem().Set(reflect.ValueOf(c))
    return nil
}

func HandleError(err *error) {
    if r := recover(); r != nil {
        *err = errors.New(fmt.Sprint(r))
    }
}
