package controllers

import (
	"github.com/Alexandr-Fox/fox-core/consts/config"
	"github.com/Alexandr-Fox/fox-core/docs"
	"github.com/Alexandr-Fox/fox-core/models"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (app *App) Ping(ctx *fiber.Ctx) error {
	return ctx.JSON(app)
}

func NewApp() *App {
	name, _ := models.GetConfig(config.FieldName, config.ValueName)
	version, _ := models.GetConfig(config.FieldName, config.ValueVersion)

	return &App{
		Name:    name.Value,
		Version: version.Value,
	}
}

// ApiVersions
// При изменении формата входных/выходных параметров, необходимо обновить версию приложения
// и добавить в массив новую строку с версией. Все модели и методы должны поддерживать все версии,
// начиная с той в которой они были созданы. Если есть версия 1.0 и в каком-то контроллере
// меняется структура входных/выходных данных, то мы обновляем версию до 1.1 и все остальные методы
// должны поддерживать версию 1.1 без изменений, а вот измененный метод добавляет в себе проверку на новую версию и в
// зависимости от этого обрабатывает определенный тип входных/выходных данных.
func (app *App) ApiVersions(ctx *fiber.Ctx) error {
	return ctx.JSON([]string{"1.0"})
}

// Docs
// Обязательно передаются Query параметры "version" и "method",
// нужно возвращать описание документации для конкретной версии.
// Пока при переходе между версиями не меняется структура
// входных/выходных параметров метода, можно пропустить обработку "version".
func (app *App) Docs(ctx *fiber.Ctx) error {
	method := ctx.Query("method", app.Version)
	fields := map[string]docs.ControllerDoc{
		"ping": {
			Results: &docs.Result{
				Type: docs.Object,
				Options: &docs.ResultOptions{
					Items: &[]docs.Result{
						{Name: "name", Type: docs.String},
						{Name: "version", Type: docs.String},
					},
				},
			},
		},
		"apiVersions": {
			Results: &docs.Result{
				Type: docs.Array,
				Options: &docs.ResultOptions{
					Item: &docs.Result{Type: docs.String},
				},
			},
		},
	}

	return ctx.JSON(fields[method])
}
