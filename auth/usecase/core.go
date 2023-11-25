package usecase

import (
	"auth/repository/profile"
	"auth/repository/session"
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2023_2_Vkladyshi/configs"
	"github.com/go-park-mail-ru/2023_2_Vkladyshi/errors"
)

type Core struct {
	sessions   session.SessionRepo
	mutex      sync.RWMutex
	lg         *slog.Logger
	users      profile.IUserRepo
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetCore(cfg_sql configs.DbDsnCfg, cfg_csrf configs.DbRedisCfg, cfg_sessions configs.DbRedisCfg, lg *slog.Logger) (*Core, error) {
	session, err := session.GetSessionRepo(cfg_sessions, lg)

	if err != nil {
		lg.Error("Session repository is not responding")
		return nil, err
	}

	users, err := profile.GetUserRepo(cfg_sql, lg)
	if err != nil {
		lg.Error("cant create repo")
		return nil, err
	}

	core := Core{
		sessions:   *session,
		lg:         lg.With("module", "core"),
		users:      users,
	}
	return &core, nil
}

func (core *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	core.mutex.RLock()
	login, err := core.sessions.GetUserLogin(ctx, sid, core.lg)
	core.mutex.RUnlock()

	if err != nil {
		return "", err
	}

	return login, nil
}

func (core *Core) CreateSession(ctx context.Context, login string) (string, session.Session, error) {
	sid := RandStringRunes(32)

	newSession := session.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	core.mutex.Lock()
	sessionAdded, err := core.sessions.AddSession(ctx, newSession, core.lg)
	core.mutex.Unlock()

	if !sessionAdded && err != nil {
		return "", session.Session{}, err
	}

	if !sessionAdded {
		return "", session.Session{}, nil
	}

	return sid, newSession, nil
}

func (core *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	core.mutex.RLock()
	found, err := core.sessions.CheckActiveSession(ctx, sid, core.lg)
	core.mutex.RUnlock()

	if err != nil {
		return false, err
	}

	return found, nil
}

func (core *Core) KillSession(ctx context.Context, sid string) error {
	core.mutex.Lock()
	_, err := core.sessions.DeleteSession(ctx, sid, core.lg)
	core.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (core *Core) CreateUserAccount(login string, password string, name string, birthDate string, email string) error {
	if matched, _ := regexp.MatchString(`@`, email); !matched {
		return errors.InvalideEmail
	}
	err := core.users.CreateUser(login, password, name, birthDate, email)
	if err != nil {
		core.lg.Error("create user error", "err", err.Error())
		return fmt.Errorf("CreateUserAccount err: %w", err)
	}

	return nil
}

func (core *Core) FindUserAccount(login string, password string) (*profile.UserItem, bool, error) {
	user, found, err := core.users.GetUser(login, password)
	if err != nil {
		core.lg.Error("find user error", "err", err.Error())
		return nil, false, fmt.Errorf("FindUserAccount err: %w", err)
	}
	return user, found, nil
}

func (core *Core) FindUserByLogin(login string) (bool, error) {
	found, err := core.users.FindUser(login)
	if err != nil {
		core.lg.Error("find user error", "err", err.Error())
		return false, fmt.Errorf("FindUserByLogin err: %w", err)
	}

	return found, nil
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}

func (core *Core) GetUserProfile(login string) (*profile.UserItem, error) {
	profile, err := core.users.GetUserProfile(login)
	if err != nil {
		core.lg.Error("GetUserProfile error", "err", err.Error())
		return nil, fmt.Errorf("GetUserProfile err: %w", err)
	}

	return profile, nil
}

func (core *Core) EditProfile(prevLogin string, login string, password string, email string, birthDate string, photo string) error {
	err := core.users.EditProfile(prevLogin, login, password, email, birthDate, photo)
	if err != nil {
		core.lg.Error("Edit profile error", "err", err.Error())
		return fmt.Errorf("Edit profile error: %w", err)
	}

	return nil
}