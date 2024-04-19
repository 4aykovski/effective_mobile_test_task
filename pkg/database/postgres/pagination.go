package postgres

import "fmt"

func AddPaginationToStmt(stmt string, args []interface{}, limit, offset int) (string, []interface{}) {
	if limit > 0 {
		args = append(args, limit)
		stmt += fmt.Sprintf(" LIMIT $%d", len(args))
	}
	if offset > 0 {
		args = append(args, offset)
		stmt += fmt.Sprintf(" OFFSET $%d", len(args))
	}
	return stmt, args
}
