package usecase

import (
	"auth/repository/profile"
	"auth/repository/session"
	"context"
)

type ICore interface {
	CreateSession(ctx context.Context, login string) (string, session.Session, error)
	KillSession(ctx context.Context, sid string) error
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	CreateUserAccount(login string, password string, name string, birthDate string, email string) error
	FindUserAccount(login string, password string) (*profile.UserItem, bool, error)
	FindUserByLogin(login string) (bool, error)
	GetUserName(ctx context.Context, sid string) (string, error)
	GetUserProfile(login string) (*profile.UserItem, error)
	EditProfile(prevLogin string, login string, password string, email string, birthDate string, photo string) error
}