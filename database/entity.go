package database

import (
	"errors"
	"fmt"
	"github.com/Alexandr-Fox/fox-core/consts"
	"gorm.io/gorm"
	"strings"
	"time"
)

type LoadQuery struct {
	Preload []string
	Filter  interface{}
}

type Entity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (e *Entity) Create() error {
	return GetCoreDB().Create(e).Error
}

func (e *Entity) Delete() error {
	return GetCoreDB().Delete(&e, fmt.Sprintf(consts.QueryEqOne, consts.FieldID), e.ID).Error
}

func (e *Entity) Save() error {
	return GetCoreDB().Save(&e).Error
}

func (e *Entity) First(name string, value interface{}) error {
	if err := GetCoreDB().Where(fmt.Sprintf(consts.QueryEqOne, name), value).First(&e).Error; err != nil {
		return err
	}

	return nil
}
func (e *Entity) Find(name string, value interface{}) (entities []Entity, err error) {
	if err = GetCoreDB().Where(fmt.Sprintf(consts.QueryEqOne, name), value).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (e *Entity) Load(query ...LoadQuery) (entities []Entity, err error) {
	if len(query) > 1 {
		return []Entity{}, errors.New(consts.ErrorArgumentsCount)
	}

	qb := GetCoreDB().Model(&e)

	if len(query) > 0 {
		for _, v := range query[0].Preload {
			qb = qb.Preload(v)
		}

		query, values := ParseQuery(query[0].Filter)
		qb = qb.Where(query, values...)
	}

	err = qb.Find(&entities).Error

	return entities, err
}

func ParseQuery(data interface{}) (string, []interface{}) {
	if mapped, ok := data.(map[string]interface{}); ok {
		if name, ok := mapped["name"]; ok {
			value := mapped["value"]
			op := "="

			if v, ok := mapped["op"]; ok {
				op = v.(string)
			}

			if value == nil {
				return fmt.Sprintf("(`%s` %s NULL)", name, op), []interface{}{}
			}

			return fmt.Sprintf("(`%s` %s ?)", name, op), []interface{}{value}
		}

		if items, ok := mapped["items"]; ok {
			op := "and"

			if v, ok := mapped["op"]; ok {
				op = v.(string)
			}

			var wheres []string
			var values []interface{}

			for _, item := range items.([]interface{}) {
				query, value := ParseQuery(item)
				wheres = append(wheres, query)
				values = append(values, value...)
			}

			if len(wheres) == 1 {
				return wheres[0], values
			}

			return fmt.Sprintf("(%s)", strings.Join(wheres[:], fmt.Sprintf(" %s ", op))), values
		}
	} else if items, ok := data.([]interface{}); ok {
		op := "and"

		if v, ok := mapped["op"]; ok {
			op = v.(string)
		}

		var wheres []string
		var values []interface{}

		for _, item := range items {
			query, value := ParseQuery(item)
			wheres = append(wheres, query)
			values = append(values, value...)
		}

		if len(wheres) == 1 {
			return wheres[0], values
		}

		return fmt.Sprintf("(%s)", strings.Join(wheres[:], fmt.Sprintf(" %s ", op))), values
	}

	return consts.DefaultEmpty, []interface{}{}
}
