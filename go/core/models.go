package core

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"sort"
)

type PropertyKind string

const (
	KindString PropertyKind = "String"
	KindEnum   PropertyKind = "Enum"
	KindBool   PropertyKind = "Bool"
	KindUser   PropertyKind = "User"
	KindTag    PropertyKind = "Tag"
)

type PropertyDef struct {
	Name    string       `json:"name" yaml:"name"`
	Kind    PropertyKind `json:"kind" yaml:"kind"`
	Values  []string     `json:"values" yaml:"values"`
	Default string       `json:"default" yaml:"default"`
}

type Model struct {
	Name       string        `json:"name" yaml:"name"`
	Properties []PropertyDef `json:"properties" yaml:"properties"`
	Template   []byte        `json:"template" yaml:"template"`
}

func ReadModels(path string) ([]Model, error) {
	modelsMap := make(map[string]Model)

	path = filepath.Join(path, ProjectModelsFolder)
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		fileName := info.Name()
		filePath := filepath.Join(path, fileName)
		ext := filepath.Ext(fileName)
		name := fileName[0 : len(fileName)-len(ext)]

		switch ext {
		case ".md":
			if template, err := ioutil.ReadFile(filePath); err == nil {
				if model, found := modelsMap[name]; found {
					model.Template = template
				} else {
					modelsMap[name] = Model{
						Name:       name,
						Template:   template,
						Properties: make([]PropertyDef, 0),
					}
				}
				logrus.Debugf("Loaded template %s: %s", filePath, template)
			} else {
				logrus.Warnf("File %s contains an invalid template", filePath)
			}
		case ".yaml":
			{
				var properties []PropertyDef
				if err := ReadYaml(filePath, &properties); err == nil {
					if model, found := modelsMap[name]; found {
						model.Properties = properties
						modelsMap[name] = model
					} else {
						modelsMap[name] = Model{Name: name, Properties: properties}
					}
					logrus.Debugf("Loaded model file %s: %v", filePath, properties)
				} else {
					logrus.Warnf("File %s contains an invalid model", filePath)
				}
			}
		}

	}
	models := make([]Model, 0, len(modelsMap))
	for _, model := range modelsMap {
		models = append(models, model)
		logrus.Debugf("Added model %s: \nProperties: %v\nTemplate: `%s`", model.Name,
			model.Properties, string(model.Template))
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})
	return models, nil
}
