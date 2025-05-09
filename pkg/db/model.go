// Code generated by mfd-generator v0.4.5; DO NOT EDIT.

//nolint:all
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"
)

var Columns = struct {
	Company struct {
		ID, Name, TgID, Inn, Scope, UserName, Phone, StatusID string
	}
	Student struct {
		ID, TgID, Name, Birthday, City, Scope, Email, StatusID string
	}
	Task struct {
		ID, CompanyID, Scope, Description, Images, Link, Deadline, ContactSlottext, StatusID, StudentID, ContactSlot string

		Company, Student string
	}
}{
	Company: struct {
		ID, Name, TgID, Inn, Scope, UserName, Phone, StatusID string
	}{
		ID:       "companyId",
		Name:     "name",
		TgID:     "tgId",
		Inn:      "inn",
		Scope:    "scope",
		UserName: "userName",
		Phone:    "phone",
		StatusID: "statusId",
	},
	Student: struct {
		ID, TgID, Name, Birthday, City, Scope, Email, StatusID string
	}{
		ID:       "studentId",
		TgID:     "tgId",
		Name:     "name",
		Birthday: "birthday",
		City:     "city",
		Scope:    "scope",
		Email:    "email",
		StatusID: "statusId",
	},
	Task: struct {
		ID, CompanyID, Scope, Description, Images, Link, Deadline, ContactSlottext, StatusID, StudentID, ContactSlot string

		Company, Student string
	}{
		ID:              "taskId",
		CompanyID:       "companyId",
		Scope:           "scope",
		Description:     "description",
		Images:          "images",
		Link:            "link",
		Deadline:        "deadline",
		ContactSlottext: "contactSlot text",
		StatusID:        "statusId",
		StudentID:       "studentId",
		ContactSlot:     "contactSlot",

		Company: "Company",
		Student: "Student",
	},
}

var Tables = struct {
	Company struct {
		Name, Alias string
	}
	Student struct {
		Name, Alias string
	}
	Task struct {
		Name, Alias string
	}
}{
	Company: struct {
		Name, Alias string
	}{
		Name:  "companies",
		Alias: "t",
	},
	Student: struct {
		Name, Alias string
	}{
		Name:  "students",
		Alias: "t",
	},
	Task: struct {
		Name, Alias string
	}{
		Name:  "tasks",
		Alias: "t",
	},
}

type Company struct {
	tableName struct{} `pg:"companies,alias:t,discard_unknown_columns"`

	ID       int         `pg:"companyId,pk"`
	Name     string      `pg:"name,use_zero"`
	TgID     int64       `pg:"tgId,use_zero"`
	Inn      int         `pg:"inn,use_zero"`
	Scope    string      `pg:"scope,use_zero"`
	UserName interface{} `pg:"-"` // unsupported
	Phone    int         `pg:"phone,use_zero"`
	StatusID int         `pg:"statusId,use_zero"`
}

type Student struct {
	tableName struct{} `pg:"students,alias:t,discard_unknown_columns"`

	ID       int       `pg:"studentId,pk"`
	TgID     int64     `pg:"tgId,use_zero"`
	Name     string    `pg:"name,use_zero"`
	Birthday time.Time `pg:"birthday,use_zero"`
	City     string    `pg:"city,use_zero"`
	Scope    string    `pg:"scope,use_zero"`
	Email    string    `pg:"email,use_zero"`
	StatusID int       `pg:"statusId,use_zero"`
}

type Task struct {
	tableName struct{} `pg:"tasks,alias:t,discard_unknown_columns"`

	ID              int       `pg:"taskId,pk"`
	CompanyID       int       `pg:"companyId,use_zero"`
	Scope           string    `pg:"scope,use_zero"`
	Description     string    `pg:"description,use_zero"`
	Images          []string  `pg:"images,array,use_zero"`
	Link            string    `pg:"link,use_zero"`
	Deadline        time.Time `pg:"deadline,use_zero"`
	ContactSlottext string    `pg:"contactSlot text,use_zero"`
	StatusID        int       `pg:"statusId,use_zero"`
	StudentID       *int      `pg:"studentId"`
	ContactSlot     string    `pg:"contactSlot,use_zero"`

	Company *Company `pg:"fk:companyId,rel:has-one"`
	Student *Student `pg:"fk:studentId,rel:has-one"`
}
