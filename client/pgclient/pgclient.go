package pgclient

import (
	"database/sql"
	"errors"

	"github.com/exiledavatar/gotoolkit/meta"
)

type Config struct {
	Schema        string // explicitly assign schema, will default to postgres
	Table         string // explicitly assign table, will attempt to get from TableNameTags or struct type
	TableNameTags []string
	FieldNameTags []string
	DataTypeTag   string
	PrimaryKeyTag string
}

var TemplateConfig = Config{
	Schema:        "postgres",
	TableNameTags: []string{"table"},
	FieldNameTags: []string{"pg", "postgres", "db", "sql"},
	DataTypeTag:   "pgtype",
	PrimaryKeyTag: "primarykey",
}

type PGClient struct {
}

func (c *PGClient) CreateSchema(schema string) (sql.Result, error) {

	return nil, errors.New("TODO")
}

type Templates struct {
	Config        Config
	CreateSchema  string
	DropSchema    string
	CreateTable   string
	DropTable     string
	Get           string
	GetMostRecent string
	Put           string
}

func Test() {
	s, err := meta.ToStruct(PGTemplates)
	if err != nil {
		panic(err)
	}
	s.Tags.Append("table", "test")
}

// these templates assume you use meta.Struct and overwrite members as needed
var PGTemplates = Templates{
	CreateSchema: `create schema if not exists {{ .Struct.LastNameSpace | tolower }}`,
	DropSchema:   `drop schema if exists {{ .Struct.LastNameSpace | tolower }}`,
	CreateTable: `{{- "\n" -}}
	CREATE TABLE IF NOT EXISTS {{ .Struct.TagIdentifier .Config.TableNameTags | tolower }} (
		{{- $fields := .Struct.Fields.WithAnyTagTrueSlice .Config.FieldNameTags -}}
		{{- $names := $fields.TagNames .Config.FieldNameTags -}}
		{{- $types := $fields.NonEmptyTagValues .Config.DataTypeTag -}}
		{{- $columnDefs := joinslices "\t" ",\n\t" $names $types -}}
		{{- print "\n\t" $columnDefs -}}
		{{- $primarykeyfields := .Struct.Fields.WithTagTrue .Config.PrimaryKeyTag -}}
		{{- $primarykey := $primarykeyfields.TagNames .Config.FieldNameTags | join ", " -}}
		{{- if ne $primarykey "" -}}{{- printf ",\n\tPRIMARY KEY ( %s )" $primarykey -}}{{- end -}}
		{{- "\n)" -}}
	
	`,
	DropTable: `drop table if exists {{ .Struct.TagIdentifier .Config.TableNameTags }}`,
	Put:       ``,
}

var XTemp string = `{{- "\n" -}}
	insert into {{ .Struct.TagIdentifier "table" | tolower }} (
		select distinct on (tmp._id_hash) tmp.*
		from _tmp_{{- .Struct.TagName "table" }} tmp 
		{{- if ne .Struct.Parent nil }}
		inner join _tmp_{{ .Struct.Parent.TagName "table" }} ptmp
		{{- $parentprimarykey := (index ( .Struct.Parent.Fields.WithTagTrue "primarykey" ) 0 ).TagName "db" -}}
		{{- $parentpkey := ( index ( .Struct.Fields.WithTagTrue "parentprimarykey" ) 0  ).TagName "db" }} 
		on tmp.{{ $parentpkey }} = ptmp.{{ $parentprimarykey -}}
		{{ end }}
		where not exists (
			select 1
			from {{ .Struct.TagIdentifier "table" }} dst
			{{- $pkey := index ( ( .Struct.Fields.WithTagTrue "primarykey" ).TagNames "db" ) 0 }}
			where dst.{{ $pkey }} = tmp.{{ $pkey }}
		) 
	)
`
