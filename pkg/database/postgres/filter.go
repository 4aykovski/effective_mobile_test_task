package postgres

import (
	"fmt"
	"reflect"

	"github.com/4aykovski/effective_mobile_test_task/pkg/api/filter"
	"github.com/4aykovski/effective_mobile_test_task/pkg/tag"
)

func AddFilterToStmt(stmt string, args []interface{}, filterOptions filter.Options, model interface{}) (string, []interface{}) {
	stmt += fmt.Sprintf(" WHERE 1=1 ")

	tags := map[string]string{}
	carType := reflect.TypeOf(model)
	for i := 0; i < carType.NumField(); i++ {
		jsonTag := tag.ParseJsonTag(carType.Field(i).Tag.Get("json"))
		dbTag := carType.Field(i).Tag.Get("db")
		tags[jsonTag] = dbTag
	}

	for _, field := range filterOptions.Fields() {
		if dbTag, ok := tags[field.Name]; ok {
			op, err := filter.ParseOperator(field.Op)
			if err != nil {
				continue
			}
			stmt += fmt.Sprintf(" AND %s %s $%d", dbTag, op, len(args)+1)
			args = append(args, field.Value)
		}
	}

	return stmt, args
}
