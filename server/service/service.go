package service

import "github.com/bradenrayhorn/beans/server/beans"

type service struct {
	ds                beans.DataSource
	sessionRepository beans.SessionRepository
}

type All struct {
	User beans.UserService
}

func NewServices(datasource beans.DataSource, sessionRepository beans.SessionRepository) *All {
	service := service{datasource, sessionRepository}

	return &All{
		User: &userService{service},
	}
}
