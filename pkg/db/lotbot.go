package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type LotbotRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewLotbotRepo returns new repository
func NewLotbotRepo(db orm.DB) LotbotRepo {
	return LotbotRepo{
		db: db,
		filters: map[string][]Filter{
			Tables.Company.Name: {StatusFilter},
			Tables.Student.Name: {StatusFilter},
			Tables.Task.Name:    {StatusFilter},
		},
		sort: map[string][]SortField{
			Tables.Company.Name: {{Column: Columns.Company.CreatedAt, Direction: SortDesc}},
			Tables.Student.Name: {{Column: Columns.Student.CreatedAt, Direction: SortDesc}},
			Tables.Task.Name:    {{Column: Columns.Task.CreatedAt, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Company.Name: {TableColumns},
			Tables.Student.Name: {TableColumns},
			Tables.Task.Name:    {TableColumns, Columns.Task.Company, Columns.Task.Student},
		},
	}
}

// WithTransaction is a function that wraps LotbotRepo with pg.Tx transaction.
func (lr LotbotRepo) WithTransaction(tx *pg.Tx) LotbotRepo {
	lr.db = tx
	return lr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (lr LotbotRepo) WithEnabledOnly() LotbotRepo {
	f := make(map[string][]Filter, len(lr.filters))
	for i := range lr.filters {
		f[i] = make([]Filter, len(lr.filters[i]))
		copy(f[i], lr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	lr.filters = f

	return lr
}

/*** Company ***/

// FullCompany returns full joins with all columns
func (lr LotbotRepo) FullCompany() OpFunc {
	return WithColumns(lr.join[Tables.Company.Name]...)
}

// DefaultCompanySort returns default sort.
func (lr LotbotRepo) DefaultCompanySort() OpFunc {
	return WithSort(lr.sort[Tables.Company.Name]...)
}

// CompanyByID is a function that returns Company by ID(s) or nil.
func (lr LotbotRepo) CompanyByID(ctx context.Context, id int, ops ...OpFunc) (*Company, error) {
	return lr.OneCompany(ctx, &CompanySearch{ID: &id}, ops...)
}

// OneCompany is a function that returns one Company by filters. It could return pg.ErrMultiRows.
func (lr LotbotRepo) OneCompany(ctx context.Context, search *CompanySearch, ops ...OpFunc) (*Company, error) {
	obj := &Company{}
	err := buildQuery(ctx, lr.db, obj, search, lr.filters[Tables.Company.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// CompaniesByFilters returns Company list.
func (lr LotbotRepo) CompaniesByFilters(ctx context.Context, search *CompanySearch, pager Pager, ops ...OpFunc) (companies []Company, err error) {
	err = buildQuery(ctx, lr.db, &companies, search, lr.filters[Tables.Company.Name], pager, ops...).Select()
	return
}

// CountCompanies returns count
func (lr LotbotRepo) CountCompanies(ctx context.Context, search *CompanySearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, lr.db, &Company{}, search, lr.filters[Tables.Company.Name], PagerOne, ops...).Count()
}

// AddCompany adds Company to DB.
func (lr LotbotRepo) AddCompany(ctx context.Context, company *Company, ops ...OpFunc) (*Company, error) {
	q := lr.db.ModelContext(ctx, company)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Company.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return company, err
}

// UpdateCompany updates Company in DB.
func (lr LotbotRepo) UpdateCompany(ctx context.Context, company *Company, ops ...OpFunc) (bool, error) {
	q := lr.db.ModelContext(ctx, company).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Company.ID, Columns.Company.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteCompany set statusId to deleted in DB.
func (lr LotbotRepo) DeleteCompany(ctx context.Context, id int) (deleted bool, err error) {
	company := &Company{ID: id, StatusID: StatusDeleted}

	return lr.UpdateCompany(ctx, company, WithColumns(Columns.Company.StatusID))
}

/*** Student ***/

// FullStudent returns full joins with all columns
func (lr LotbotRepo) FullStudent() OpFunc {
	return WithColumns(lr.join[Tables.Student.Name]...)
}

// DefaultStudentSort returns default sort.
func (lr LotbotRepo) DefaultStudentSort() OpFunc {
	return WithSort(lr.sort[Tables.Student.Name]...)
}

// StudentByID is a function that returns Student by ID(s) or nil.
func (lr LotbotRepo) StudentByID(ctx context.Context, id int, ops ...OpFunc) (*Student, error) {
	return lr.OneStudent(ctx, &StudentSearch{ID: &id}, ops...)
}

// OneStudent is a function that returns one Student by filters. It could return pg.ErrMultiRows.
func (lr LotbotRepo) OneStudent(ctx context.Context, search *StudentSearch, ops ...OpFunc) (*Student, error) {
	obj := &Student{}
	err := buildQuery(ctx, lr.db, obj, search, lr.filters[Tables.Student.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// StudentsByFilters returns Student list.
func (lr LotbotRepo) StudentsByFilters(ctx context.Context, search *StudentSearch, pager Pager, ops ...OpFunc) (students []Student, err error) {
	err = buildQuery(ctx, lr.db, &students, search, lr.filters[Tables.Student.Name], pager, ops...).Select()
	return
}

// CountStudents returns count
func (lr LotbotRepo) CountStudents(ctx context.Context, search *StudentSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, lr.db, &Student{}, search, lr.filters[Tables.Student.Name], PagerOne, ops...).Count()
}

// AddStudent adds Student to DB.
func (lr LotbotRepo) AddStudent(ctx context.Context, student *Student, ops ...OpFunc) (*Student, error) {
	q := lr.db.ModelContext(ctx, student)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Student.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return student, err
}

// UpdateStudent updates Student in DB.
func (lr LotbotRepo) UpdateStudent(ctx context.Context, student *Student, ops ...OpFunc) (bool, error) {
	q := lr.db.ModelContext(ctx, student).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Student.ID, Columns.Student.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteStudent set statusId to deleted in DB.
func (lr LotbotRepo) DeleteStudent(ctx context.Context, id int) (deleted bool, err error) {
	student := &Student{ID: id, StatusID: StatusDeleted}

	return lr.UpdateStudent(ctx, student, WithColumns(Columns.Student.StatusID))
}

/*** Task ***/

// FullTask returns full joins with all columns
func (lr LotbotRepo) FullTask() OpFunc {
	return WithColumns(lr.join[Tables.Task.Name]...)
}

// DefaultTaskSort returns default sort.
func (lr LotbotRepo) DefaultTaskSort() OpFunc {
	return WithSort(lr.sort[Tables.Task.Name]...)
}

// TaskByID is a function that returns Task by ID(s) or nil.
func (lr LotbotRepo) TaskByID(ctx context.Context, id int, ops ...OpFunc) (*Task, error) {
	return lr.OneTask(ctx, &TaskSearch{ID: &id}, ops...)
}

// OneTask is a function that returns one Task by filters. It could return pg.ErrMultiRows.
func (lr LotbotRepo) OneTask(ctx context.Context, search *TaskSearch, ops ...OpFunc) (*Task, error) {
	obj := &Task{}
	err := buildQuery(ctx, lr.db, obj, search, lr.filters[Tables.Task.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// TasksByFilters returns Task list.
func (lr LotbotRepo) TasksByFilters(ctx context.Context, search *TaskSearch, pager Pager, ops ...OpFunc) (tasks []Task, err error) {
	err = buildQuery(ctx, lr.db, &tasks, search, lr.filters[Tables.Task.Name], pager, ops...).Select()
	return
}

// CountTasks returns count
func (lr LotbotRepo) CountTasks(ctx context.Context, search *TaskSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, lr.db, &Task{}, search, lr.filters[Tables.Task.Name], PagerOne, ops...).Count()
}

// AddTask adds Task to DB.
func (lr LotbotRepo) AddTask(ctx context.Context, task *Task, ops ...OpFunc) (*Task, error) {
	q := lr.db.ModelContext(ctx, task)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Task.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return task, err
}

// UpdateTask updates Task in DB.
func (lr LotbotRepo) UpdateTask(ctx context.Context, task *Task, ops ...OpFunc) (bool, error) {
	q := lr.db.ModelContext(ctx, task).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Task.ID, Columns.Task.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteTask set statusId to deleted in DB.
func (lr LotbotRepo) DeleteTask(ctx context.Context, id int) (deleted bool, err error) {
	task := &Task{ID: id, StatusID: StatusDeleted}

	return lr.UpdateTask(ctx, task, WithColumns(Columns.Task.StatusID))
}
