package postgres

func AddPaginationToStmt(stmt string, args []interface{}, limit, offset int) (string, []interface{}) {
	if limit > 0 {
		stmt += " LIMIT $1"
		args = append(args, limit)
	}
	if offset > 0 {
		stmt += " OFFSET $2"
		args = append(args, offset)
	}
	return stmt, args
}
