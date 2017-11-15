package model

import (
	"errors"
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"

	"gopkg.in/mgo.v2/bson"
)

var (
	E_REGISTRY_NOT_FOUND = errors.New("Registry Not Found")
)

type registry struct{}

var _registry *registry

func RegistryModel() *registry {
	return _registry
}

func (r *registry) Get(id string) (*Registry, error) {
	c := libdb.MainMgoDB().C("registries")

	var reg Registry
	it := c.FindId(id).Iter()
	defer it.Close()
	for it.Next(&reg) {
		return &reg, nil
	}

	err := it.Err()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return nil, E_REGISTRY_NOT_FOUND
}

func (r *registry) Create(appID, devToken, kind string) (*Registry, error) {
	c := libdb.MainMgoDB().C("registries")
	var reg Registry
	reg.AppID = appID
	reg.CreatedAt = util.Timestamp()
	reg.DevToken = devToken
	reg.ID = bson.NewObjectId().Hex()
	reg.Kind = kind
	err := c.Insert(&reg)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &reg, nil
}

func (r *registry) Delete(regID string) error {
	c := libdb.MainMgoDB().C("registries")
	err := c.RemoveId(regID)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}