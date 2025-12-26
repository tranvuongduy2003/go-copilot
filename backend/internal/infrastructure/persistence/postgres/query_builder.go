package postgres

import (
	"fmt"
	"strings"
)

type WhereClause struct {
	conditions []string
	args       []any
	paramIndex int
}

func NewWhereClause() *WhereClause {
	return &WhereClause{
		conditions: make([]string, 0),
		args:       make([]any, 0),
		paramIndex: 1,
	}
}

func NewWhereClauseWithOffset(offset int) *WhereClause {
	return &WhereClause{
		conditions: make([]string, 0),
		args:       make([]any, 0),
		paramIndex: offset,
	}
}

func (w *WhereClause) AddCondition(condition string, args ...any) *WhereClause {
	placeholders := make([]any, len(args))
	for i := range args {
		placeholders[i] = w.paramIndex
		w.paramIndex++
	}
	w.conditions = append(w.conditions, fmt.Sprintf(condition, placeholders...))
	w.args = append(w.args, args...)
	return w
}

func (w *WhereClause) AddConditionIf(condition bool, clause string, args ...any) *WhereClause {
	if condition {
		return w.AddCondition(clause, args...)
	}
	return w
}

func (w *WhereClause) Eq(column string, value any) *WhereClause {
	return w.AddCondition(column+" = $%d", value)
}

func (w *WhereClause) EqIf(condition bool, column string, value any) *WhereClause {
	if condition {
		return w.Eq(column, value)
	}
	return w
}

func (w *WhereClause) Neq(column string, value any) *WhereClause {
	return w.AddCondition(column+" <> $%d", value)
}

func (w *WhereClause) Gt(column string, value any) *WhereClause {
	return w.AddCondition(column+" > $%d", value)
}

func (w *WhereClause) Gte(column string, value any) *WhereClause {
	return w.AddCondition(column+" >= $%d", value)
}

func (w *WhereClause) Lt(column string, value any) *WhereClause {
	return w.AddCondition(column+" < $%d", value)
}

func (w *WhereClause) Lte(column string, value any) *WhereClause {
	return w.AddCondition(column+" <= $%d", value)
}

func (w *WhereClause) Like(column string, value string) *WhereClause {
	return w.AddCondition(column+" LIKE $%d", value)
}

func (w *WhereClause) ILike(column string, value string) *WhereClause {
	return w.AddCondition(column+" ILIKE $%d", value)
}

func (w *WhereClause) In(column string, values ...any) *WhereClause {
	if len(values) == 0 {
		return w
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", w.paramIndex)
		w.paramIndex++
	}
	w.conditions = append(w.conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	w.args = append(w.args, values...)
	return w
}

func (w *WhereClause) IsNull(column string) *WhereClause {
	w.conditions = append(w.conditions, column+" IS NULL")
	return w
}

func (w *WhereClause) IsNotNull(column string) *WhereClause {
	w.conditions = append(w.conditions, column+" IS NOT NULL")
	return w
}

func (w *WhereClause) Between(column string, start, end any) *WhereClause {
	return w.AddCondition(column+" BETWEEN $%d AND $%d", start, end)
}

func (w *WhereClause) Build() (string, []any) {
	if len(w.conditions) == 0 {
		return "", w.args
	}
	return "WHERE " + strings.Join(w.conditions, " AND "), w.args
}

func (w *WhereClause) BuildWithoutKeyword() (string, []any) {
	if len(w.conditions) == 0 {
		return "", w.args
	}
	return strings.Join(w.conditions, " AND "), w.args
}

func (w *WhereClause) Args() []any {
	return w.args
}

func (w *WhereClause) NextParamIndex() int {
	return w.paramIndex
}

type OrderDirection string

const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

type OrderByClause struct {
	orders []string
}

func NewOrderByClause() *OrderByClause {
	return &OrderByClause{
		orders: make([]string, 0),
	}
}

func (o *OrderByClause) Add(column string, direction OrderDirection) *OrderByClause {
	o.orders = append(o.orders, fmt.Sprintf("%s %s", column, direction))
	return o
}

func (o *OrderByClause) AddIf(condition bool, column string, direction OrderDirection) *OrderByClause {
	if condition {
		return o.Add(column, direction)
	}
	return o
}

func (o *OrderByClause) Asc(column string) *OrderByClause {
	return o.Add(column, OrderAsc)
}

func (o *OrderByClause) Desc(column string) *OrderByClause {
	return o.Add(column, OrderDesc)
}

func (o *OrderByClause) Build() string {
	if len(o.orders) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(o.orders, ", ")
}

type PaginationClause struct {
	limit  int
	offset int
}

func NewPaginationClause(page, pageSize int) *PaginationClause {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return &PaginationClause{
		limit:  pageSize,
		offset: (page - 1) * pageSize,
	}
}

func NewPaginationClauseFromOffset(limit, offset int) *PaginationClause {
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return &PaginationClause{
		limit:  limit,
		offset: offset,
	}
}

func (p *PaginationClause) Build() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.limit, p.offset)
}

func (p *PaginationClause) Limit() int {
	return p.limit
}

func (p *PaginationClause) Offset() int {
	return p.offset
}

type ReturningClause struct {
	columns []string
}

func NewReturningClause() *ReturningClause {
	return &ReturningClause{
		columns: make([]string, 0),
	}
}

func (r *ReturningClause) Add(columns ...string) *ReturningClause {
	r.columns = append(r.columns, columns...)
	return r
}

func (r *ReturningClause) All() *ReturningClause {
	r.columns = []string{"*"}
	return r
}

func (r *ReturningClause) Build() string {
	if len(r.columns) == 0 {
		return ""
	}
	return "RETURNING " + strings.Join(r.columns, ", ")
}

type QueryBuilder struct {
	baseQuery    string
	where        *WhereClause
	orderBy      *OrderByClause
	pagination   *PaginationClause
	returning    *ReturningClause
}

func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery: baseQuery,
	}
}

func (qb *QueryBuilder) Where(w *WhereClause) *QueryBuilder {
	qb.where = w
	return qb
}

func (qb *QueryBuilder) OrderBy(o *OrderByClause) *QueryBuilder {
	qb.orderBy = o
	return qb
}

func (qb *QueryBuilder) Paginate(p *PaginationClause) *QueryBuilder {
	qb.pagination = p
	return qb
}

func (qb *QueryBuilder) Returning(r *ReturningClause) *QueryBuilder {
	qb.returning = r
	return qb
}

func (qb *QueryBuilder) Build() (string, []any) {
	parts := []string{qb.baseQuery}
	var args []any

	if qb.where != nil {
		whereClause, whereArgs := qb.where.Build()
		if whereClause != "" {
			parts = append(parts, whereClause)
			args = append(args, whereArgs...)
		}
	}

	if qb.orderBy != nil {
		orderClause := qb.orderBy.Build()
		if orderClause != "" {
			parts = append(parts, orderClause)
		}
	}

	if qb.pagination != nil {
		parts = append(parts, qb.pagination.Build())
	}

	if qb.returning != nil {
		returningClause := qb.returning.Build()
		if returningClause != "" {
			parts = append(parts, returningClause)
		}
	}

	return strings.Join(parts, " "), args
}
