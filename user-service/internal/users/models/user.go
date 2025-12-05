package models

import (
	"github.com/google/uuid"
)

type User struct {
	UID          uuid.UUID `db:"uid"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	FirstName    string    `db:"first_name"`
	SecondName   string    `db:"second_name"`
	Age          int32     `db:"age"`
	Salary       float64   `db:"salary"`
	WorkSphereID int64     `db:"work_sphere"`
}

type CreateUserDTO struct {
	Username     string
	Password     string
	FirstName    string
	SecondName   string
	Age          int32
	Salary       float64
	WorkSphereID int64
}

type UpdateUserDTO struct {
	UID          uuid.UUID
	Username     string
	FirstName    string
	SecondName   string
	Age          int32
	Salary       float64
	WorkSphereID int64
}
