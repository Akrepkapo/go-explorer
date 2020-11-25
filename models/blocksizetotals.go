package models

import (
	"github.com/vmihailenco/msgpack/v5"
)

type BlockSizeTotal struct {
	ID    int64
	Count int64
	Name  string
}

var blocksizetotalPrefix = "blocksizetotal-"

func (b *BlockSizeTotal) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(b); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (b *BlockSizeTotal) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &b); err != nil {
		return err
	}
	return nil
}

func (m *BlockSizeTotal) ReSetRedis(blockid, pos int64) error {

	var so BlockSizeTotal
	f, err := so.Get_Redis()
	if err != nil {
		return err
	}
	if f {
		so.ID = blockid
		so.Count += pos
		if so.Count < 0 {
			so.Count = 0
		}
		so.Name = blocksizetotalPrefix + "lastet"
		val, err := so.Marshal()
		if err != nil {
			return err
		}
		rd := RedisParams{
			Key:   blocksizetotalPrefix + "lastet",
			Value: string(val),
		}
		err = rd.Set()
		if err != nil {
		}
		so.Name = blocksizetotalPrefix + "lastet"
		val, err := so.Marshal()
		if err != nil {
			return err
		}
		rd := RedisParams{
			Key:   blocksizetotalPrefix + "lastet",
			Value: string(val),
		}
		err = rd.Set()
		if err != nil {
			return err
		}
	}
	return err
}

func (m *BlockSizeTotal) Get_Redis() (bool, error) {
	rd := RedisParams{
		Key:   blocksizetotalPrefix + "lastet",
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
