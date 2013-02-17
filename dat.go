package netkit

import (
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
)

// data wrapper
type DataWrapper struct {
	Session 	*mgo.Session
	Database 	*mgo.Database
	C 			*mgo.Collection
}

// return a new data wrapper instance
func NewDataWrapper(host string) *DataWrapper {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	return &DataWrapper{
		Session: session,
	}
}

// set database
func (self *DataWrapper) SetDb(db string) *DataWrapper {
	self.Database = self.Session.DB(db)
	return self
}

// set collection
func (self *DataWrapper) SetC(c string) *DataWrapper {
	self.C = self.Database.C(c)
	return self
}

// insert
func (self *DataWrapper) Insert(v ...interface{}) interface{} {
	err := self.C.Insert(v...)
	if err != nil {
		return err
	}
	return len(v)
}

// update
func (self *DataWrapper) Update(v ...interface{}) interface{} {
	info, err := self.C.UpdateAll(v[0], v[1])
	if err != nil {
		return err
	}
	return info.Updated
}

// return
func (self *DataWrapper) Return(v ...interface{}) interface{} {
	var lmt int
	var sel, set, ret interface{}
	for k, val := range v {
		switch val.(type) {
		case int:
			lmt = val.(int)
			v = append(v[:k], v[k+1:]...) 
		}
	}
	switch len(v) {
	case 1:
		sel, set = bson.M{}, v[0]
	case 2:
		sel, set = v[0], v[1]
	default:
		ret = nil
	}
	switch lmt {
	case 0:
		ret = self.C.Find(sel).All(set)
	case 1:
		//ret = self.C.Find(sel).One(set)
    	ret = self.C.Find(sel).Sort("-_id").One(set)
    default:
    	ret = self.C.Find(sel).Limit(lmt).All(set)
	}
	return ret
}

// delete
func (self *DataWrapper) Delete(v ...interface{}) interface{} {
	info, err := self.C.RemoveAll(v)
	if err != nil {
		return err
	}
	return info.Removed
}