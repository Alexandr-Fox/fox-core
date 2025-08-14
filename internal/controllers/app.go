package controllers

import (
	config2 "github.com/Alexandr-Fox/fox-core/internal/consts/config"
	docs2 "github.com/Alexandr-Fox/fox-core/internal/docs"
	"github.com/Alexandr-Fox/fox-core/internal/models"
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
	name, _ := models.GetConfig(config2.FieldName, config2.ValueName)
	version, _ := models.GetConfig(config2.FieldName, config2.ValueVersion)

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
	fields := map[string]docs2.ControllerDoc{
		"ping": {
			Results: &docs2.Result{
				Type: docs2.Object,
				Options: &docs2.ResultOptions{
					Items: &[]docs2.Result{
						{Name: "name", Type: docs2.String},
						{Name: "version", Type: docs2.String},
					},
				},
			},
		},
		"apiVersions": {
			Results: &docs2.Result{
				Type: docs2.Array,
				Options: &docs2.ResultOptions{
					Item: &docs2.Result{Type: docs2.String},
				},
			},
		},
	}

	return ctx.JSON(fields[method])
}
