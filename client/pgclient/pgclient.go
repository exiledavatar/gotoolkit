package pgclient

import (
	"context"
	"reflect"
	"text/template"
	"time"

	"github.com/exiledavatar/gotoolkit/client"
	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/jackc/pgx/v5"
)

type Client struct {
	client.Client[pgx.Conn]
}

func NewConfig(cfg ...client.Config) client.Config {
	cc := DefaultConfig()
	for _, cf := range cfg {
		cc.Connection.Merge(cf.Connection)
		cc.Template.Merge(cf.Template)
	}
	return cc
}

func DefaultConfig() client.Config {
	return client.Config{
		Connection: client.ConnectionConfig{},
		Template:   TemplateConfig,
	}
}

func NewClient(cfg ...client.Config) Client {

	return Client{client.Client[pgx.Conn]{
		Config:    NewConfig(cfg...),
		Templator: PGTemplates,
	}}

}

func (c *Client) Connect(ctx context.Context) error {
	var err error
	c.Conn, err = pgx.Connect(ctx, c.Config.Connection.ConnectionString)
	if err != nil {
		return err
	}
	c.Config.Template.Schema = c.Config.Connection.Schema
	c.Templator.Config.Schema = c.Config.Connection.Schema
	return nil
}

// type Config struct {
// 	Schema           string // explicitly assign schema, will default to postgres
// 	Table            string // explicitly assign table, will attempt to get from TableNameTags or struct type
// 	TableNameTags    []string
// 	FieldNameTags    []string
// 	LastInsertTags   []string // for checking the 'last insert' values in destination
// 	TaggedFieldsOnly bool
// 	DataTypeTag      string
// 	PrimaryKeyTag    string
// }

var TemplateConfig = client.TemplatorConfig{
	Schema:           "public",
	Table:            "",
	TableNameTags:    []string{"table"},
	FieldNameTags:    []string{"pg", "postgres", "db", "sql"},
	LastInsertTags:   []string{"pgli"},
	TaggedFieldsOnly: false, // include all fields by default
	DataTypeTag:      "pgtype",
	PrimaryKeyTag:    "primarykey",
}

var FuncMap = template.FuncMap{
	"pgtype":  GoToPGType,
	"pgtypes": GoToPGTypes,
}

var TemplateData = map[string]any{
	"Config":  TemplateConfig,
	"TypeMap": TypeMap,
}

// type Templator struct {
// 	Config          client.TemplatorConfig
// 	CreateSchema    string
// 	DropSchema      string
// 	CreateTable     string
// 	CreateTempTable string
// 	DropTable       string
// 	Get             string
// 	GetMostRecent   string
// 	Put             string
// 	PutTempToTable  string
// 	FuncMap         template.FuncMap
// 	Data            map[string]any // any additional 'data' passed to templates
// }

// {{- $fields := .Struct.Fields.WithTagTrue .Config.FieldNameTags -}}
// these templates assume you use meta.Struct and overwrite members as needed
var PGTemplates = client.Templator{
	Config:       TemplateConfig,
	FuncMap:      FuncMap,
	Data:         TemplateData,
	CreateSchema: `create schema if not exists {{ .Config.Schema | tolower }}`,
	DropSchema:   `drop schema if exists {{ .Config.Schema | tolower }}`,
	CreateTable: `{{- "\n" -}}
	CREATE TABLE IF NOT EXISTS {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }} (
		{{- $fields := .Struct.Fields -}}
		{{- if .Config.TaggedFieldsOnly -}}{{- $fields = .Struct.Fields.WithTagTrue .Config.FieldNameTags -}}{{- end -}}
		{{- $names := $fields.TagNames .Config.FieldNameTags | tolowerslices -}}
		{{- $tagtypes := $fields.NonEmptyTagValues .Config.DataTypeTag -}}
		{{- $defaulttypes := pgtypes $fields.TypeNames -}}
		{{- $types := coalesce $tagtypes $defaulttypes "text" -}}
		{{- $columnDefs := joinslices "\t" ",\n\t" $names $types -}}
		{{- print "\n\t" $columnDefs -}}
		{{- $primarykeyfields := $fields.WithTagTrue .Config.PrimaryKeyTag -}}
		{{- $primarykey := $primarykeyfields.TagNames .Config.FieldNameTags | join ", " -}}
		{{- if ne $primarykey "" -}}{{- printf ",\n\tPRIMARY KEY ( %s )" $primarykey -}}{{- end -}}
		{{- "\n)" -}}
		
		`,
	CreateTempTable: `
		create temp table _tmp_{{ .Struct.TagName .Config.TableNameTags | tolower }} (
		like {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }}
		excluding constraints ) 
		`,
	DropTable: `drop table if exists {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }}`,
	Put: `{{- "\n" -}}
		insert into {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }} ( 
			{{- $fields := .Struct.Fields -}}
			{{- if .Config.TaggedFieldsOnly -}}{{- $fields = .Struct.Fields.WithTagTrue .Config.FieldNameTags -}}{{- end -}}
			{{- $names := $fields.TagNames .Config.FieldNameTags | tolowerslices -}}
			{{- "\n\t" -}}{{- $names | join ",\n\t" }}{{- "\n" -}}
			) values (
				{{- "\n\t" -}}:{{- $names | join ",\n\t:" -}}
				{{- "\n) on conflict (" -}}
				{{- $primarykeyfields := $fields.WithTagTrue .Config.PrimaryKeyTag -}}
				{{- $primarykeyfields.TagNames .Config.FieldNameTags | tolowerslices | join ", " -}}
				) do nothing
				`,
	PutTempToTable: `{{- "\n" -}}
		{{- $primarykeyfields := .Struct.Fields.WithTagTrue .Config.PrimaryKeyTag -}}
		{{- $primarykey := $primarykeyfields.TagNames .Config.FieldNameTags | join ", " -}}

		insert into {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }} (
		select distinct 
		{{- if $primarykeyfields }}
		{{- $pkey := $primarykeyfields.TagNames .Config.FieldNameTags | join ", tmp." -}}
		{{ "" }} on ( tmp.{{ $pkey }} ) 
		{{- end }}		
		tmp.*
		from _tmp_{{ .Struct.TagName .Config.TableNameTags | tolower }} tmp 	
		) {{ if $primarykey }} on conflict ( {{ $primarykey }} ) do nothing{{ end }}
		`,
	Get: `{{- "\n" -}}
		select
			{{- $fields := .Struct.Fields -}}
			{{- if .Config.TaggedFieldsOnly -}}{{- $fields = .Struct.Fields.WithTagTrue .Config.FieldNameTags -}}{{- end -}}
			{{- $names := $fields.TagNames .Config.FieldNameTags | tolowerslices -}}
			{{- "\n\t" -}}{{- $names | join ",\n\t" }}{{- "\n" -}}
		from
			{{ .Struct.TagIdentifier .Config.TableNameTags | tolower }}
		{{ if .rowlimit -}}limit {{ .rowlimit }}{{- end }}
		{{- "\n)" -}}
		`,
	GetMostRecent: `{{- "\n" -}}
		select
			{{- $fields := .Struct.Fields.WithTagTrue .Config.LastInsertTags -}}
			{{- names := $fields.NonEmptyTagValues .Config.LastInsertTags -}}
			{{- "\n\t" -}}{{- $names | join ",\n\t" }}{{- "\n" -}}
		from
			{{ .Struct.TagIdentifier .Config.TableNameTags | tolower }}
		{{- "\n)" -}}
	`,
}

var XTemp string = `{{- "\n" -}}
	insert into {{ .Struct.TagIdentifier "table" | tolower }} (
		select distinct on (tmp._id_hash) tmp.*
		from _tmp_{{- .Struct.TagName "table" }} tmp 
		{{- if ne .Struct.Parent nil }}
		inner join _tmp_{{ .Struct.Parent.TagName "table" }} ptmp
		{{- $parentprimarykey := (index ( .Struct.Parent.Fields.WithTagTrue "primarykey" ) 0 ).TagName "db" -}}
		{{- $parentpkey := ( index ( .Struct.Fields.WithTagTrue "parentprimarykey" ) 0  ).TagName .Config.FieldNameTags | tolower }} 
		on tmp.{{ $parentpkey }} = ptmp.{{ $parentprimarykey -}}
		{{ end }}
		where not exists (
			select 1
			from {{ .Struct.TagIdentifier "table" | tolower }} dst
			{{- $pkey := index ( ( .Struct.Fields.WithTagTrue "primarykey" ).TagNames .Config.FieldNameTags ) 0 | tolower }}
			where dst.{{ $pkey }} = tmp.{{ $pkey }}
		) 
	)
`

// this is just copy pasted from meta for now - idk when to break it out into gotoolkit/typemap
var TypeMaps = meta.TypeMaps{
	"postgres": TypeMap,
}

var TypeMap = meta.TypeMap{
	From: meta.From{
		"text":     reflect.TypeOf("string"),
		"varchar":  reflect.TypeOf("string"),
		"smallint": reflect.TypeOf(int(1)),
		"int":      reflect.TypeOf(int(1)),
		"bigint":   reflect.TypeOf(int(1)),
	},
	To: meta.To{
		reflect.TypeOf(string("string")): "text",
		reflect.TypeOf(bool(true)):       "boolean",
		reflect.TypeOf(int(1)):           "bigint",
		reflect.TypeOf(int8(1)):          "smallint",
		reflect.TypeOf(int16(1)):         "smallint",
		reflect.TypeOf(int32(1)):         "bigint",
		reflect.TypeOf(int64(1)):         "bigint",
		reflect.TypeOf(float32(1.0)):     "float4",
		reflect.TypeOf(float64(1.0)):     "float8",
		reflect.TypeOf(time.Time{}):      "timestamp with time zone",
		reflect.TypeOf([]byte{}):         "bytea[]",
		nil:                              "bytea[]", // serves as a default
	},
}

// PGToGoTypeMap represents the default type mapping
// we expect to use when retrieving data from postgres
var PGToGoTypeMap = map[string]string{
	"text": "string",
}

// GoToPGTypeMap represents the default type mapping
// we expect to use when sending data to postgres
var GoToPGTypeMap = map[string]string{
	"default":          "bytea[]",
	"string":           "text",
	"Date":             "date",
	"qgenda.Date":      "date",
	"Time":             "timestamp with time zone",
	"qgenda.Time":      "timestamp with time zone",
	"TimeOfDay":        "time without time zone",
	"qgenda.TimeOfDay": "time without time zone",
	"time.Time":        "timestamp with time zone",
	"bool":             "boolean",

	"int":     "bigint",
	"int8":    "smallint",
	"int16":   "smallint",
	"int32":   "bigint",
	"int64":   "bigint",
	"float32": "double precision",
	"float64": "double precision",
}

func GoToPGType(gotype string) string {
	pgtype, ok := GoToPGTypeMap[gotype]
	if !ok {
		pgtype = "text"
	}
	return pgtype
}

func GoToPGTypes(gotypes []string) []string {
	out := []string{}
	for _, gt := range gotypes {
		out = append(out, GoToPGType(gt))
	}
	return out
}

// These are all a work in progress. Defaults act the same and rely on the package's FuncMap, TemplateData, and TemplateConfig.
// You can customize by modifying these.
func DefaultCreateSchemaText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.CreateSchema, &TemplateConfig, FuncMap, nil)
}

func DefaultDropSchemaText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.DropSchema, &TemplateConfig, FuncMap, nil)
}

func DefaultCreateTableText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.CreateTable, &TemplateConfig, FuncMap, nil)
}

func DefaultCreateTempTableText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.CreateTempTable, &TemplateConfig, FuncMap, nil)
}

func DefaultDropTableText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.DropTable, &TemplateConfig, FuncMap, nil)
}

func DefaultGetText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.Get, &TemplateConfig, FuncMap, nil)
}
func DefaultGetMostRecentText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.GetMostRecent, &TemplateConfig, FuncMap, nil)
}
func DefaultPutText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.Put, &TemplateConfig, FuncMap, nil)
}
func DefaultPutTempToTableText(value any) (string, error) {
	return TemplateToText(value, PGTemplates.PutTempToTable, &TemplateConfig, FuncMap, nil)
}

// TemplateToText is the base function for all XXXToText functions. You can expand Funcmap and TemplateData or use overrides to remove/replace them.
func TemplateToText(value any, tpl string, cfg *client.TemplatorConfig, funcMap template.FuncMap, data map[string]any) (string, error) {
	str, err := meta.ToStruct(value)
	if err != nil {
		return "", err
	}

	for k, v := range funcMap {
		FuncMap[k] = v
	}

	for k, v := range data {
		TemplateData[k] = v
	}
	tcfg := TemplateConfig
	if cfg != nil {
		tcfg = *cfg
	}
	TemplateData["Config"] = tcfg
	// update str with any relevant config items
	if tcfg.Schema != "" {
		str.NameSpace = []string{tcfg.Schema}
	}
	if tcfg.Table != "" {
		str.Name = tcfg.Table
	}

	TemplateData["TypeMap"] = TypeMap

	return str.ExecuteTemplate(
		tpl,
		FuncMap,
		TemplateData,
	)

}

func (c Client) TemplateToText(value any, tpl string) (string, error) {
	str, err := meta.ToStruct(value)
	if err != nil {
		return "", err
	}

	return str.ExecuteTemplate(
		tpl,
		c.Templator.FuncMap,
		c.Templator.Data,
	)
}
