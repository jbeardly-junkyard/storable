package storable

import (
	"gopkg.in/mgo.v2/bson"
)

type DocumentBase interface {
	GetId() bson.ObjectId
	SetId(bson.ObjectId)
	IsNew() bool
	SetIsNew(isNew bool)
}

type Document struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`

	//Tracks if the document has been saved or recovered from the db or not
	isNew bool
}

func (d *Document) SetId(id bson.ObjectId) {
	d.Id = id
}

func (d *Document) GetId() bson.ObjectId {
	return d.Id
}

func (d *Document) SetIsNew(isNew bool) {
	d.isNew = isNew
}

func (d *Document) IsNew() bool {
	return d.isNew
}
