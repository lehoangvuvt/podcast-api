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
	SelectColumns     []string                `json:"select_columns"`
}

type QueryBuilder struct {
	DB    *sql.DB
	query QueryConfig
}

func (builder *QueryBuilder) Select(columns ...string) *QueryBuilder {
	builder.query.SelectColumns = []string{}
	builder.query.SelectColumns = append(builder.query.SelectColumns, columns...)
	return builder
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
			columnOperator = "="
			columnValue = "'" + builder.query.SearchValue + "'"
		}
		whereQuery = fmt.Sprintf("WHERE %v %v %v", columnName, columnOperator, columnValue)
	}
	var selectedColumns string
	for i, column := range builder.query.SelectColumns {
		if i < len(builder.query.SelectColumns)-1 {
			selectedColumns += column + ", "
		} else {
			selectedColumns += column
		}
	}
	queryStr := fmt.Sprintf("SELECT %v FROM %v %v ", selectedColumns, builder.query.FromTable, whereQuery)
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

func (builder *QueryBuilder) GetOne() *sql.Row {
	builder.query.Limit = 1
	builder.query.Skip = 0
	queryStr := builder.GetQuery()
	row := builder.DB.QueryRow(queryStr)
	return row
}

func (builder *QueryBuilder) GetMany() (*sql.Rows, error) {
	queryStr := builder.GetQuery()
	rows, err := builder.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
