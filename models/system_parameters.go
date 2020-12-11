package models

/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

import (
	"encoding/json"

	"github.com/IBAX-io/go-explorer/conf"
)

// SystemParameter is model
type SystemParameter struct {
	ID         int64  `gorm:"primary_key;not null;"`
	Name       string `gorm:"not null;size:255"`
	Value      string `gorm:"not null"`
	Conditions string `gorm:"not null"`
}

type SystemParameterResult struct {
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Rets  []SystemParameter `json:"rets"`
}

// TableName returns name of table
func (sp SystemParameter) TableName() string {
	return "1_system_parameters"
}

// Get is retrieving model from database
func (sp *SystemParameter) Get(name string) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("name = ?", name).First(sp))
}

// GetJSONField returns fields as json
func (sp *SystemParameter) GetJSONField(jsonField string, name string) (string, error) {
	var result string
	err := conf.GetDbConn().Conn().Table("1_system_parameters").Where("name = ?", name).Select(jsonField).Row().Scan(&result)
	return result, err
}

// GetValueParameterByName returns value parameter by name
func (sp *SystemParameter) GetValueParameterByName(name, value string) (*string, error) {
	var result *string
	err := conf.GetDbConn().Conn().Raw(`SELECT value->'`+value+`' FROM "1_system_parameters" WHERE name = ?`, name).Row().Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ToMap is converting SystemParameter to map
func (sp *SystemParameter) ToMap() map[string]string {
	result := make(map[string]string, 0)
	result["name"] = sp.Name
	result["value"] = sp.Value
	result["conditions"] = sp.Conditions
	return result
}

// Update is update model
func (sp SystemParameter) Update(value string) error {
	return conf.GetDbConn().Conn().Model(sp).Where("name = ?", sp.Name).Update(`value`, value).Error
}

// SaveArray is saving array
func (sp *SystemParameter) SaveArray(list [][]string) error {
	ret, err := json.Marshal(list)
	if err != nil {
		return err
	}
	ns := "%" + name + "%"
	if err := conf.GetDbConn().Conn().Table(sp.TableName()).Where("name like ?", ns).Count(&num).Error; err != nil {
		return num, rets, err
	}

	err = conf.GetDbConn().Conn().Table(sp.TableName()).Where("name like ?", ns).
		Order(order).Offset((page - 1) * size).Limit(size).Find(&rets).Error

	return num, rets, err
}
