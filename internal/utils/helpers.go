package queryHelpers

import (
	"database/sql"
	"fmt"
)

type QueryOrderDirectionType string

var QueryDirection = struct {
	DESC QueryOrderDirectionType
	ASC  QueryOrderDirectionType
}{
	DESC: "DESC",
	ASC:  "ASC",
}

type ColumnOperatorType string

var ColumnOperator = struct {
	EQUAL    ColumnOperatorType
	CONTAINS ColumnOperatorType
}{
	EQUAL:    "EQUAL",
	CONTAINS: "CONTAINS",
}

type QueryConfig struct {
	FromTable         string                  `json:"from_table"`
	WhereColumnName   string                  `json:"where_column_name"`
	Operator          ColumnOperatorType      `json:"operator"`
	SearchValue       string                  `json:"search_value"`
	OrderByColumnName string                  `json:"order_by_column_name"`
	Direction         QueryOrderDirectionType `json:"direction"`
	Skip              int                     `json:"skip"`
	Limit             int                     `json:"limit"`
}

type QueryBuilder struct {
	DB    *sql.DB
	query QueryConfig
}

func (builder *QueryBuilder) FromTable(tableName string) *QueryBuilder {
	builder.query.SearchValue = "*"
	builder.query.Operator = ColumnOperator.CONTAINS
	builder.query.OrderByColumnName = ""
	builder.query.Direction = QueryDirection.ASC
	builder.query.Limit = -1
	builder.query.Skip = -1

	builder.query.FromTable = tableName
	return builder
}

func (builder *QueryBuilder) WhereColumn(colName string) *QueryBuilder {
	builder.query.WhereColumnName = colName
	return builder
}

func (builder *QueryBuilder) Equal(searchValue string) *QueryBuilder {
	builder.query.Operator = ColumnOperator.EQUAL
	builder.query.SearchValue = searchValue
	return builder
}

func (builder *QueryBuilder) Search(operator ColumnOperatorType, searchValue string) *QueryBuilder {
	builder.query.Operator = operator
	builder.query.SearchValue = searchValue
	return builder
}

func (builder *QueryBuilder) OrderBy(columnName string, direction QueryOrderDirectionType) *QueryBuilder {
	builder.query.OrderByColumnName = columnName
	builder.query.Direction = direction
	return builder
}

func (builder *QueryBuilder) Contains(searchValue string) *QueryBuilder {
	builder.query.Operator = ColumnOperator.CONTAINS
	builder.query.SearchValue = searchValue
	return builder
}

func (builder *QueryBuilder) Skip(skip int) *QueryBuilder {
	builder.query.Skip = skip
	return builder
}

func (builder *QueryBuilder) Limit(limit int) *QueryBuilder {
	builder.query.Limit = limit
	return builder
}

func (builder *QueryBuilder) LeftJoin(joinTableName string, table1JoinColumn string, table2JoinColumn string) {

}

func (builder *QueryBuilder) GetQuery() string {
	columnName := builder.query.WhereColumnName
	var columnOperator string
	var columnValue string
	whereQuery := ""
	if builder.query.SearchValue != "*" {
		switch builder.query.Operator {
		case ColumnOperator.CONTAINS:
			columnOperator = "ILIKE"
			columnValue = "'%" + builder.query.SearchValue + "%'"
		case ColumnOperator.EQUAL:
			columnOperator = "ILIKE"
			columnValue = "'" + builder.query.SearchValue + "'"
		}
		whereQuery = fmt.Sprintf("WHERE %v %v %v", columnName, columnOperator, columnValue)
	}
	queryStr := fmt.Sprintf("SELECT * FROM %v %v ", builder.query.FromTable, whereQuery)
	if builder.query.OrderByColumnName != "" {
		queryStr += fmt.Sprintf(" ORDER BY %v %v", builder.query.OrderByColumnName, builder.query.Direction)
	}
	if builder.query.Skip != -1 {
		queryStr += fmt.Sprintf(" OFFSET %v", builder.query.Skip)
	}
	if builder.query.Limit != -1 {
		queryStr += fmt.Sprintf(" LIMIT %v", builder.query.Limit)
	}
	return queryStr
}
