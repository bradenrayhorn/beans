package service

import "github.com/bradenrayhorn/beans/server/beans"

type service struct {
	ds                beans.DataSource
	sessionRepository beans.SessionRepository
}

type All struct {
	MonthCategory beans.MonthCategoryService
	User          beans.UserService
}

func NewServices(datasource beans.DataSource, sessionRepository beans.SessionRepository) *All {
	service := service{datasource, sessionRepository}

	return &All{
		MonthCategory: &monthCategoryService{service},
		User:          &userService{service},
	}
}
