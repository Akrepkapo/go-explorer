
import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"strconv"

	"github.com/IBAX-io/go-explorer/conf"
)

// StateParameter is model
type StateParameter struct {
	ecosystem  int64
	ID         int64  `gorm:"primary_key;not null"`
	Name       string `gorm:"not null;size:100"`
	Value      string `gorm:"not null"`
	Conditions string `gorm:"not null"`
}

// TableName returns name of table
func (sp *StateParameter) TableName() string {
	if sp.ecosystem == 0 {
		sp.ecosystem = 1
	}
	return `1_parameters`
}

// SetTablePrefix is setting table prefix
func (sp *StateParameter) SetTablePrefix(prefix string) {
	pre, _ := strconv.ParseInt(prefix, 10, 64)
	sp.ecosystem = pre
}

// Get is retrieving model from database
func (sp *StateParameter) Get(name string) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("ecosystem = ? and name = ?", sp.ecosystem, name).First(sp))
}

// Get is retrieving model from database
func (sp *StateParameter) GetMintAmount() (string, error) {
	var sp1, sp2 StateParameter
	f, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "mint_balance").First(sp))
	if err != nil {
		return "0", err
	}

	f1, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "foundation_balance").First(&sp1))
	if err != nil {
		return "0", err
	}

	f2, err := isFound(GetDB(nil).Where("ecosystem = ? and name = ?", sp.ecosystem, "assign_rule").First(&sp2))
	if err != nil {
		return "0", err
	}
	if !f || !f1 || !f2 {
		return "0", nil
	}

	ret := make(map[int64]AssignRules, 1)
	err = json.Unmarshal([]byte(sp2.Value), &ret)
	if err != nil {
		return "0", err
	}

	as3, ok3 := ret[3]
	as6, ok6 := ret[6]
	if !ok3 || !ok6 {
		return "0", nil
	}

	tfa, err := decimal.NewFromString(as3.TotalAmount)
	if err != nil {
		return "0", err
	}
	tma, err := decimal.NewFromString(as6.TotalAmount)
	if err != nil {
		return "0", err
	}

	ma, err := decimal.NewFromString(sp.Value)
	if err != nil {
		return "0", err
	}
	fa, err := decimal.NewFromString(sp1.Value)
	if err != nil {
		return "0", err
	}
	if fa.LessThanOrEqual(tfa) && ma.LessThanOrEqual(tma) {
		mb := tma.Sub(ma)
		fb := tfa.Sub(fa)
		tt := mb.Add(fb)
		return tt.String(), nil
	} else {
		return "0", errors.New("assign rules err")
	}

	return "0", nil
}
