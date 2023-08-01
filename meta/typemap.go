package meta

import (
	"reflect"
	"time"
	// "github.com/ericsgagnon/qgenda/pkg/meta"
	// "github.com/exiledavatar/gotoolkit/meta"
)

// default type maps between go and other systems

// From is for inbound types - it should map the external system's type name as a key
// to go's reflect.Type as a value
type From map[string]reflect.Type

// To is for outbound types - it should map go's reflect.Type as a key to an external type name
type To map[reflect.Type]string

// TypeMap combines To/From
type TypeMap struct {
	From From
	To   To
}

// ToType indirects any value and returns the external type as a string,
// if it exists in the To typemap
func (m TypeMap) ToType(value any) string {
	// switch t, ok := value.(reflect.Type); ok {
	// case true:
	// 	return m.To[t]
	// default:
	// 	t, err := ToValue(value)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return m.To[t.Type()]
	// }
	t, err := ToValue(value)
	if err != nil {
		panic(err)
	}
	return m.To[t.Type()]

}

// FromType takes the external type and returns the reflect.Type,
// if it exists in the From typemap
func (m TypeMap) FromType(externalType string) reflect.Type {
	return m.From[externalType]
}

// Maps are a collection of external To/From mappings,
// they should use external system names as keys.
type TypeMaps map[string]TypeMap

// From returns the go type for the given system/external type,
// if it exists in the typemap
func (m TypeMaps) From(system, externalType string) reflect.Type {
	return m[system].From[externalType]
}

// To takes any go value and returns the exernal type in the given system,
// if it exists in the typemap
func (m TypeMaps) To(system string, value any) string {
	return m[system].To[reflect.TypeOf(value)]
}

var TypeMappings = TypeMaps{
	"postgres": Postgres,
}

var Postgres = TypeMap{
	From: From{
		"text":     reflect.TypeOf("string"),
		"varchar":  reflect.TypeOf("string"),
		"smallint": reflect.TypeOf(int(1)),
		"int":      reflect.TypeOf(int(1)),
		"bigint":   reflect.TypeOf(int(1)),
	},
	To: To{
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

// // PGToGoTypeMap represents the default type mapping
// // we expect to use when retrieving data from postgres
// var PGToGoTypeMap = map[string]string{
// 	"text": "string",
// }

// // GoToPGTypeMap represents the default type mapping
// // we expect to use when sending data to postgres
// var GoToPGTypeMap = map[string]string{
// 	"default":          "bytea[]",
// 	"string":           "text",
// 	"Date":             "date",
// 	"qgenda.Date":      "date",
// 	"Time":             "timestamp with time zone",
// 	"qgenda.Time":      "timestamp with time zone",
// 	"TimeOfDay":        "time without time zone",
// 	"qgenda.TimeOfDay": "time without time zone",
// 	"time.Time":        "timestamp with time zone",
// 	"bool":             "boolean",

// 	"int":     "bigint",
// 	"int8":    "smallint",
// 	"int16":   "smallint",
// 	"int32":   "bigint",
// 	"int64":   "bigint",
// 	"float32": "double precision",
// 	"float64": "double precision",
// }

// func GoToPGType(gotype string) string {
// 	pgtype, ok := GoToPGTypeMap[gotype]
// 	if !ok {
// 		pgtype = "text"
// 	}
// 	return pgtype
// }
