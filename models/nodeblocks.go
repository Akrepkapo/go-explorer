package models

import (
	"github.com/vmihailenco/msgpack/v5"
)

// NodeBlocks is model
type NodeBlocks struct {
	ID    int64
	Count int64
	Name  string
}

var nodeblocksPrefix = "nodeblocks-"

func (b *NodeBlocks) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(b); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (b *NodeBlocks) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &b); err != nil {
		return err
	}
	return nil
}

func (m *NodeBlocks) SubOneRedis(pos string) error {

	var so NodeBlocks
	f, err := so.Get_rediss(pos)
	if err != nil {
		return err
	}
	if f {
		if so.Count >= 1 {
			so.Count -= 1
			so.Name = nodeblocksPrefix + pos
			val, err := so.Marshal()
			if err != nil {
				return err
			}
			rd := RedisParams{
				Key:   nodeblocksPrefix + pos,
				Value: string(val),
			}
			err = rd.Set()
			if err != nil {
				return err
			}
		}

	}

	return err
}

func (m *NodeBlocks) AddOneFromRedis(pos string) error {

	var so NodeBlocks
	f, err := so.Get_rediss(pos)
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			so.Count = 1
			so.Name = nodeblocksPrefix + pos
			val, err := so.Marshal()
			if err != nil {
				return err
			}
			rd := RedisParams{
				Key:   nodeblocksPrefix + pos,
				Value: string(val),
			}
			err = rd.Set()
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}
	if f {
		so.Count += 1
		so.Name = nodeblocksPrefix + pos
		val, err := so.Marshal()
		if err != nil {
			return err
		}
		rd := RedisParams{
			Key:   nodeblocksPrefix + pos,
			Value: string(val),
		}
		err = rd.Set()
		Key:   nodeblocksPrefix + pos,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}
	return true, err
}

func (m *NodeBlocks) Get_rediss(pos string) (bool, error) {
	rd := RedisParams{
		Key:   nodeblocksPrefix + pos,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}
	return true, err
}
