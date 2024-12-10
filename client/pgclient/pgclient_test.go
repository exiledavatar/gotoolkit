package pgclient_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/exiledavatar/gotoolkit/client/pgclient"
)

type StructTest struct {
	IDHash             string    `pg:"_id_hash" pgtype:"text" primarykey:"true"`
	StringField        string    `pg:"string_field" pgtype:"text"`
	StringPointerField *string   `pg:"string_pointer_field" pgtype:"text"`
	BoolField          bool      `pg:"bool_field" `
	IntField           int       `pg:"int_field" `
	FloatField         float64   `pg:"float_field" pgtype:"numeric"`
	TimeField          time.Time `pg:"time_field" pgtype:"timestamp with time zone"`
}

var structTest = StructTest{
	IDHash:      "s0M3Aw350m3T3xt",
	StringField: "Anumber1",
	BoolField:   true,
	IntField:    1138,
	FloatField:  3.14159265,
	TimeField:   time.Now(),
}

func Test(t *testing.T) {
	pgclient.TemplateConfig.Schema = "sample_schema"

	fmt.Println(pgclient.DefaultCreateSchemaText(structTest))
	fmt.Println(pgclient.DefaultDropSchemaText(structTest))
	fmt.Println(pgclient.DefaultCreateTableText(structTest))
	fmt.Println(pgclient.DefaultCreateTempTableText(structTest))
	fmt.Println(pgclient.DefaultDropTableText(structTest))
	fmt.Println(pgclient.DefaultGetText(structTest))
	fmt.Println(pgclient.DefaultGetMostRecentText(structTest))
	fmt.Println(pgclient.DefaultPutText(structTest))
	fmt.Println(pgclient.DefaultPutTempToTableText(structTest))

	// struct slice test
	fmt.Println(pgclient.DefaultCreateTableText([]StructTest{structTest}))

}
