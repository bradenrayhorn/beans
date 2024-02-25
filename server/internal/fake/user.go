package fake

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
)

type userRepository struct{ repository }

var _ beans.UserRepository = (*userRepository)(nil)

func (r *userRepository) Create(ctx context.Context, id beans.ID, username beans.Username, passwordHash beans.PasswordHash) error {
	r.acquire(func() { r.usersMU.Lock() })
	defer r.usersMU.Unlock()

	if _, ok := r.database.users[id]; ok {
		return errors.New("duplicate")
	}

	if len(filter(values(r.database.users), func(it beans.User) bool { return it.Username == username })) > 0 {
		return errors.New("duplicate username")
	}

	r.users[id] = beans.User{ID: id, Username: username, PasswordHash: passwordHash}

	return nil
}

func (r *userRepository) Exists(ctx context.Context, username beans.Username) (bool, error) {
	r.acquire(func() { r.usersMU.RLock() })
	defer r.usersMU.RUnlock()

	res := find(values(r.users), func(it beans.User) bool { return it.Username == username })
	return res != nil, nil
}

func (r *userRepository) Get(ctx context.Context, id beans.ID) (beans.User, error) {
	r.acquire(func() { r.usersMU.RLock() })
	defer r.usersMU.RUnlock()

	if user, ok := r.database.users[id]; ok {
		return user, nil
	}

	return beans.User{}, beans.NewError(beans.ENOTFOUND, "user not found")
}

func (r *userRepository) GetByUsername(ctx context.Context, username beans.Username) (beans.User, error) {
	r.acquire(func() { r.usersMU.RLock() })
	defer r.usersMU.RUnlock()

	res := find(values(r.users), func(it beans.User) bool { return it.Username == username })
	if res == nil {
		return beans.User{}, beans.NewError(beans.ENOTFOUND, "user not found")
	} else {
		return *res, nil
	}
}
