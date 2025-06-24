package client

import "html/template"

type TemplatorConfig struct {
	Schema           string // explicitly assign schema, most systems default to schema used in connection
	Table            string // explicitly assign table, will attempt to get from TableNameTags or struct type
	TableNameTags    []string
	FieldNameTags    []string
	LastInsertTags   []string // for checking the 'last insert' values in destination
	TaggedFieldsOnly bool
	DataTypeTag      string
	PrimaryKeyTag    string
}

func (tc *TemplatorConfig) Merge(cfg ...TemplatorConfig) *TemplatorConfig {
	for _, cf := range cfg {
		if cf.Schema != "" {
			tc.Schema = cf.Schema
		}
		if cf.Table != "" {
			tc.Table = cf.Table
		}
		if cf.TableNameTags != nil {
			tc.TableNameTags = cf.TableNameTags
		}
		if cf.FieldNameTags != nil {
			tc.FieldNameTags = cf.FieldNameTags
		}
		if cf.LastInsertTags != nil {
			tc.LastInsertTags = cf.LastInsertTags
		}
		// there's no way to test validity of a bool
		tc.TaggedFieldsOnly = cf.TaggedFieldsOnly
		if cf.DataTypeTag != "" {
			tc.DataTypeTag = cf.DataTypeTag
		}
		if cf.PrimaryKeyTag != "" {
			tc.PrimaryKeyTag = cf.PrimaryKeyTag
		}
	}
	return tc
}

// Templator is a collection of common templates for our client
type Templator struct {
	Config          TemplatorConfig
	CreateSchema    string
	DropSchema      string
	CreateTable     string
	CreateTempTable string
	DropTable       string
	Get             string
	GetMostRecent   string
	Put             string
	PutTempToTable  string
	FuncMap         template.FuncMap
	Data            map[string]any // any additional 'data' passed to templates

}

func NewTemplatorConfig(cfg ...TemplatorConfig) TemplatorConfig {
	tc := &TemplatorConfig{}
	tc.Merge(cfg...)
	return *tc
}

func (tc TemplatorConfig) ToTemplator() Templator {
	return Templator{
		Config: tc,
	}
}
