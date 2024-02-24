package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bradenrayhorn/beans/server/beans"
)

type userService struct{ service }

var _ beans.UserService = (*userService)(nil)

func (s *userService) GetAuth(ctx context.Context, sessionID beans.SessionID) (*beans.AuthContext, error) {
	session, err := s.sessionRepository.Get(sessionID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.ErrorUnauthorized
		}

		return nil, fmt.Errorf("GetAuth find session: %w", err)
	}

	return beans.NewAuthContext(session.UserID, session.ID), nil
}
