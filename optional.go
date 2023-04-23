package qbt_apiv4

import "fmt"

type optional map[string]any

func (opt optional) StringField() map[string]string {
	m := make(map[string]string)
	for k, v := range opt {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}
