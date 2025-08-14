package boot

import (
	"errors"
	"fmt"
	"github.com/Alexandr-Fox/fox-core/models"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App AppConfig `yaml:"app"`
}

type AppConfig struct {
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
	Port    uint   `yaml:"port"`
}

func defaultConfig() Config {
	return Config{
		App: AppConfig{
			Version: "0.1",
			Name:    "fox-core",
			Port:    3000,
		},
	}
}

// FlattenConfig преобразует структуру в map[string]string с ключами вида "section.subsection.key"
func FlattenConfig(config interface{}) map[string]string {
	result := make(map[string]string)
	v := reflect.ValueOf(config)
	flatten("", v, result)
	return result
}

func flatten(prefix string, v reflect.Value, result map[string]string) {
	// Разыменовываем указатели
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	// Получаем тип значения
	t := v.Type()

	// Обходим поля структуры
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Пропускаем неэкспортируемые поля
		if !fieldValue.CanInterface() {
			continue
		}

		// Получаем тег yaml, если есть
		yamlTag := fieldType.Tag.Get("yaml")
		if yamlTag == "" {
			yamlTag = strings.ToLower(fieldType.Name)
		} else {
			// Убираем модификаторы вроде `,omitempty`
			yamlTag = strings.Split(yamlTag, ",")[0]
		}

		// Формируем ключ
		key := prefix
		if key != "" {
			key += "." + yamlTag
		} else {
			key = yamlTag
		}

		// Обработка вложенных структур
		if fieldValue.Kind() == reflect.Struct {
			flatten(key, fieldValue, result)
		} else {
			// Простые типы — конвертируем в строку
			result[key] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}
}

func LoadConfig(filename string) (Config, error) {
	cfg := defaultConfig()

	// Проверяем, существует ли файл
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует — используем defaults
			fmt.Printf("Config file %s not found, using defaults\n", filename)

			prepareCfg := FlattenConfig(cfg)
			for key := range prepareCfg {
				c := models.Config{
					Name:  key,
					Value: prepareCfg[key],
				}
				c.Create()
			}

			return cfg, nil
		}
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	// Парсим YAML в map, чтобы проверить, какие ключи есть
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return cfg, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	// Теперь парсим в структуру — только те поля, что есть в YAML, перезапишут defaults
	// Мы создаём временную структуру и заполняем её из YAML
	tempCfg := Config{}
	if err := yaml.Unmarshal(data, &tempCfg); err != nil {
		return cfg, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Применяем только те поля, которые присутствуют в YAML
	if version, ok := raw["app"].(map[string]interface{})["version"]; ok {
		if v, valid := version.(string); valid {
			cfg.App.Version = v
		}
	}
	if name, ok := raw["app"].(map[string]interface{})["name"]; ok {
		if v, valid := name.(string); valid {
			cfg.App.Name = v
		}
	}
	if port, ok := raw["app"].(map[string]interface{})["port"]; ok {
		switch v := port.(type) {
		case int:
		case int64:
		case float64:
			{
				cfg.App.Port = uint(v)
				break
			}
		default:
			return cfg, errors.New("port error")
		}
	}

	prepareCfg := FlattenConfig(cfg)
	for key := range prepareCfg {
		c := models.Config{
			Name:  key,
			Value: prepareCfg[key],
		}
		c.Create()
	}

	return cfg, nil
}
