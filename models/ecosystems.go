package models

import "github.com/IBAX-io/go-explorer/conf"

	TypeWithdraw   int64
}

func (sys *Ecosystem) TableName() string {
	return "1_ecosystems"
}

func (sys *Ecosystem) Get(id int64) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return isFound(conf.GetDbConn().Conn().First(sys, "id = ?", id))
}
func GetAllSystemStatesIDs() ([]int64, []string, error) {
	if !IsTable("1_ecosystems") {
		//return nil, fmt.Errorf("%s does not exists", ecosysTable)
		return nil, nil, nil
	}

	ecosystems := new([]Ecosystem)
	if err := conf.GetDbConn().Conn().Order("id asc").Find(&ecosystems).Error; err != nil {
		return nil, nil, err
	}

	ids := make([]int64, len(*ecosystems))
	names := make([]string, len(*ecosystems))
	for i, s := range *ecosystems {
		ids[i] = s.ID
		names[i] = s.Name
	}

	return ids, names, nil
}
