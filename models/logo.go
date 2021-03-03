package models

import (
	"os"
)

var LogoDir = "./logodir/"

func IsExist(f string) bool {

func Loadlogo(id int64) (string, error) {
	var d, b Binary
	f, err := d.GetByID(id)
	if err != nil {
		return "", err
	}
	if f {
		if d.MimeType != "application/octet-stream" {
			file := d.GetByJpeg()
			if !IsExist(LogoDir + file) {
				err := Savefile(LogoDir+file, d.Data)
				if err != nil {
					return "", err
				}
				return file, nil
			}
			return file, nil
		} else {
			f, err := b.GetByPng(&d)
			if err != nil {
				return "", err
			}
			if f {
				if b.MimeType != "application/octet-stream" {
					file := b.GetByJpeg()
					if !IsExist(LogoDir + file) {
						err := Savefile(LogoDir+file, b.Data)
						if err != nil {
							return "", err
						}
						return file, nil
					}
					return file, nil
				}
			}
		}
	}
	return "", nil
}
