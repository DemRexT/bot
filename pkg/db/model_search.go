// Code generated by mfd-generator v0.4.5; DO NOT EDIT.

//nolint:all
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const condition = "?.? = ?"

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	appliers []applier
}

func (s *search) apply(query *orm.Query) {
	for _, applier := range s.appliers {
		query.Apply(applier)
	}
}

func (s *search) where(query *orm.Query, table, field string, value interface{}) {
	query.Where(condition, pg.Ident(table), pg.Ident(field), value)
}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query *orm.Query) (*orm.Query, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

type CompanySearch struct {
	search

	ID              *int
	Name            *string
	TgID            *int64
	Scope           *string
	StatusID        *int
	UserName        *string
	Inn             *string
	Phone           *string
	CreatedAt       *time.Time
	NicknameTg      *string
	IDs             []int
	NameILike       *string
	TgIDILike       *int64
	ScopeILike      *string
	UserNameILike   *string
	InnILike        *string
	PhoneILike      *string
	NicknameTgILike *string
}

func (cs *CompanySearch) Apply(query *orm.Query) *orm.Query {
	if cs == nil {
		return query
	}
	if cs.ID != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.ID, cs.ID)
	}
	if cs.Name != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.Name, cs.Name)
	}
	if cs.TgID != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.TgID, cs.TgID)
	}
	if cs.Scope != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.Scope, cs.Scope)
	}
	if cs.StatusID != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.StatusID, cs.StatusID)
	}
	if cs.UserName != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.UserName, cs.UserName)
	}
	if cs.Inn != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.Inn, cs.Inn)
	}
	if cs.Phone != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.Phone, cs.Phone)
	}
	if cs.CreatedAt != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.CreatedAt, cs.CreatedAt)
	}
	if cs.NicknameTg != nil {
		cs.where(query, Tables.Company.Alias, Columns.Company.NicknameTg, cs.NicknameTg)
	}
	if len(cs.IDs) > 0 {
		Filter{Columns.Company.ID, cs.IDs, SearchTypeArray, false}.Apply(query)
	}
	if cs.NameILike != nil {
		Filter{Columns.Company.Name, *cs.NameILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.TgIDILike != nil {
		Filter{Columns.Company.TgID, *cs.TgIDILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.ScopeILike != nil {
		Filter{Columns.Company.Scope, *cs.ScopeILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.UserNameILike != nil {
		Filter{Columns.Company.UserName, *cs.UserNameILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.InnILike != nil {
		Filter{Columns.Company.Inn, *cs.InnILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.PhoneILike != nil {
		Filter{Columns.Company.Phone, *cs.PhoneILike, SearchTypeILike, false}.Apply(query)
	}
	if cs.NicknameTgILike != nil {
		Filter{Columns.Company.NicknameTg, *cs.NicknameTgILike, SearchTypeILike, false}.Apply(query)
	}

	cs.apply(query)

	return query
}

func (cs *CompanySearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if cs == nil {
			return query, nil
		}
		return cs.Apply(query), nil
	}
}

type StudentSearch struct {
	search

	ID              *int
	TgID            *int64
	Name            *string
	City            *string
	Scope           *string
	Email           *string
	StatusID        *int
	Birthday        *string
	CreatedAt       *time.Time
	NicknameTg      *string
	IDs             []int
	TgIDILike       *int64
	NameILike       *string
	CityILike       *string
	ScopeILike      *string
	EmailILike      *string
	BirthdayILike   *string
	NicknameTgILike *string
}

func (ss *StudentSearch) Apply(query *orm.Query) *orm.Query {
	if ss == nil {
		return query
	}
	if ss.ID != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.ID, ss.ID)
	}
	if ss.TgID != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.TgID, ss.TgID)
	}
	if ss.Name != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.Name, ss.Name)
	}
	if ss.City != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.City, ss.City)
	}
	if ss.Scope != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.Scope, ss.Scope)
	}
	if ss.Email != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.Email, ss.Email)
	}
	if ss.StatusID != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.StatusID, ss.StatusID)
	}
	if ss.Birthday != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.Birthday, ss.Birthday)
	}
	if ss.CreatedAt != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.CreatedAt, ss.CreatedAt)
	}
	if ss.NicknameTg != nil {
		ss.where(query, Tables.Student.Alias, Columns.Student.NicknameTg, ss.NicknameTg)
	}
	if len(ss.IDs) > 0 {
		Filter{Columns.Student.ID, ss.IDs, SearchTypeArray, false}.Apply(query)
	}
	if ss.TgIDILike != nil {
		Filter{Columns.Student.TgID, *ss.TgIDILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.NameILike != nil {
		Filter{Columns.Student.Name, *ss.NameILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.CityILike != nil {
		Filter{Columns.Student.City, *ss.CityILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.ScopeILike != nil {
		Filter{Columns.Student.Scope, *ss.ScopeILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.EmailILike != nil {
		Filter{Columns.Student.Email, *ss.EmailILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.BirthdayILike != nil {
		Filter{Columns.Student.Birthday, *ss.BirthdayILike, SearchTypeILike, false}.Apply(query)
	}
	if ss.NicknameTgILike != nil {
		Filter{Columns.Student.NicknameTg, *ss.NicknameTgILike, SearchTypeILike, false}.Apply(query)
	}

	ss.apply(query)

	return query
}

func (ss *StudentSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if ss == nil {
			return query, nil
		}
		return ss.Apply(query), nil
	}
}

type TaskSearch struct {
	search

	ID               *int
	CompanyID        *int
	Scope            *string
	Description      *string
	Link             *string
	Deadline         *time.Time
	ContactSlot      *string
	StatusID         *int
	StudentID        *int
	Budget           *float64
	YougileID        *string
	Name             *string
	Deadline         *string
	Url              *string
	CreatedAt        *time.Time
	IDs              []int
	ScopeILike       *string
	DescriptionILike *string
	LinkILike        *string
	ContactSlotILike *string
	YougileIDILike   *string
	NameILike        *string
	DeadlineILike    *time.Time
	UrlILike         *string
}

func (ts *TaskSearch) Apply(query *orm.Query) *orm.Query {
	if ts == nil {
		return query
	}
	if ts.ID != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.ID, ts.ID)
	}
	if ts.CompanyID != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.CompanyID, ts.CompanyID)
	}
	if ts.Scope != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Scope, ts.Scope)
	}
	if ts.Description != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Description, ts.Description)
	}
	if ts.Link != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Link, ts.Link)
	}
	if ts.Deadline != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Deadline, ts.Deadline)
	}
	if ts.ContactSlot != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.ContactSlot, ts.ContactSlot)
	}
	if ts.StatusID != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.StatusID, ts.StatusID)
	}
	if ts.StudentID != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.StudentID, ts.StudentID)
	}
	if ts.Budget != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Budget, ts.Budget)
	}
	if ts.YougileID != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.YougileID, ts.YougileID)
	}
	if ts.Name != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Name, ts.Name)
	}
	if ts.Deadline != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Deadline, ts.Deadline)
	}
	if ts.Url != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.Url, ts.Url)
	}
	if ts.CreatedAt != nil {
		ts.where(query, Tables.Task.Alias, Columns.Task.CreatedAt, ts.CreatedAt)
	}
	if len(ts.IDs) > 0 {
		Filter{Columns.Task.ID, ts.IDs, SearchTypeArray, false}.Apply(query)
	}
	if ts.ScopeILike != nil {
		Filter{Columns.Task.Scope, *ts.ScopeILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.DescriptionILike != nil {
		Filter{Columns.Task.Description, *ts.DescriptionILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.LinkILike != nil {
		Filter{Columns.Task.Link, *ts.LinkILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.ContactSlotILike != nil {
		Filter{Columns.Task.ContactSlot, *ts.ContactSlotILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.YougileIDILike != nil {
		Filter{Columns.Task.YougileID, *ts.YougileIDILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.NameILike != nil {
		Filter{Columns.Task.Name, *ts.NameILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.DeadlineILike != nil {
		Filter{Columns.Task.Deadline, *ts.DeadlineILike, SearchTypeILike, false}.Apply(query)
	}
	if ts.UrlILike != nil {
		Filter{Columns.Task.Url, *ts.UrlILike, SearchTypeILike, false}.Apply(query)
	}

	ts.apply(query)

	return query
}

func (ts *TaskSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if ts == nil {
			return query, nil
		}
		return ts.Apply(query), nil
	}
}
