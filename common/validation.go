package common

import (
	"github.com/astaxie/beego/validation"
)

type EnableValidation struct {
	Valid validation.Validation
}

func (this *EnableValidation) GetError(key string) string {
	for _, err := range this.Valid.Errors {
		if err.Key == key {
			return err.Message
		}
	}
	return ""
}
