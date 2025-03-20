package service

// Package mysql is a osin storage implementation for mysql.

import (
	"database/sql"
	"fmt"
	"oauth2/infrastructure/svc"
	"strings"
	"time"

	"github.com/openshift/osin"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	// driver for mysql db
	_ "github.com/go-sql-driver/mysql"
)

var schemas = []string{`CREATE TABLE IF NOT EXISTS {prefix}client (
	id           varchar(255) NOT NULL PRIMARY KEY,
	secret       varchar(255) NOT NULL,
	extra        text,
	redirect_uri varchar(255) NOT NULL,
	created_at   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)`, `CREATE TABLE IF NOT EXISTS {prefix}token (
	id            varchar(255) NOT NULL PRIMARY KEY,
	client_id     varchar(255) NOT NULL,
	type          varchar(20) NOT NULL,    -- 'authorize' 或 'access'
	access_token  varchar(255),            -- 访问令牌
	refresh_token varchar(255),            -- 刷新令牌
	code          varchar(255),            -- 授权码
	expires_in    int NOT NULL,
	scope         varchar(255),
	redirect_uri  varchar(255) NOT NULL,
	state         varchar(255),
	extra         text,
	created_at    timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at    timestamp NULL,
	INDEX idx_refresh (refresh_token),
	INDEX idx_expires (expires_at),
	INDEX idx_access_token (access_token),
	INDEX idx_code (code),
	FOREIGN KEY (client_id) REFERENCES {prefix}client(id) ON DELETE CASCADE
)`}

// Storage implements interface "github.com/RangelReale/osin".Storage and interface "github.com/felipeweb/osin-mysql/storage".Storage
type Storage struct {
	db          sqlx.SqlConn
	tablePrefix string
}

// New returns a new mysql storage instance.
func NewStorage(svcCtx *svc.ServiceContext, tablePrefix string) *Storage {
	return &Storage{
		db:          svcCtx.DB,
		tablePrefix: tablePrefix,
	}
}

// CreateSchemas creates the schemata, if they do not exist yet in the database. Returns an error if something went wrong.
func (s *Storage) CreateSchemas() error {
	for _, schema := range schemas {
		schema = strings.Replace(schema, "{prefix}", s.tablePrefix, -1)
		_, err := s.db.Exec(schema)
		if err != nil {
			return fmt.Errorf("创建表失败: %v", err)
		}
	}
	return nil
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *Storage) Close() {
}

// GetClient loads the client by id
func (s *Storage) GetClient(id string) (osin.Client, error) {
	var result struct {
		Id          string         `db:"id"`
		Secret      string         `db:"secret"`
		RedirectUri string         `db:"redirect_uri"`
		Extra       sql.NullString `db:"extra"`
	}

	query := fmt.Sprintf("SELECT id, secret, redirect_uri, extra FROM %sclient WHERE id = ?", s.tablePrefix)
	err := s.db.QueryRowPartial(&result, query, id)

	if err == sql.ErrNoRows {
		return nil, osin.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("获取客户端失败: %v", err)
	}

	client := &osin.DefaultClient{
		Id:          result.Id,
		Secret:      result.Secret,
		RedirectUri: result.RedirectUri,
	}

	if result.Extra.Valid {
		client.UserData = result.Extra.String
	}

	return client, nil
}

// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
func (s *Storage) UpdateClient(c osin.Client) error {
	query := fmt.Sprintf("UPDATE %sclient SET secret=?, redirect_uri=?, extra=? WHERE id=?", s.tablePrefix)
	_, err := s.db.Exec(query,
		c.GetSecret(),
		c.GetRedirectUri(),
		toString(c.GetUserData()),
		c.GetId(),
	)
	if err != nil {
		return fmt.Errorf("更新客户端失败: %v", err)
	}
	return nil
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (s *Storage) CreateClient(c osin.Client) error {
	data := toString(c.GetUserData())

	if _, err := s.db.Exec(fmt.Sprintf("INSERT INTO %sclient (id, secret, redirect_uri, extra) VALUES (?, ?, ?, ?)", s.tablePrefix), c.GetId(), c.GetSecret(), c.GetRedirectUri(), data); err != nil {
		return err
	}
	return nil
}

// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
func (s *Storage) RemoveClient(id string) (err error) {
	if _, err = s.db.Exec(fmt.Sprintf("DELETE FROM %sclient WHERE id=?", s.tablePrefix), id); err != nil {
		return err
	}
	return nil
}

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {
	query := fmt.Sprintf(`INSERT INTO %stoken (
		id, client_id, type, code, expires_in, scope, 
		redirect_uri, state, extra, expires_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, s.tablePrefix)

	_, err := s.db.Exec(query,
		data.Code,
		data.Client.GetId(),
		"authorize",
		data.Code,
		data.ExpiresIn,
		data.Scope,
		data.RedirectUri,
		data.State,
		toString(data.UserData),
		data.ExpireAt(),
	)

	if err != nil {
		return fmt.Errorf("保存授权数据失败: %v", err)
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var result struct {
		ClientId    string         `db:"client_id"`
		Code        string         `db:"code"`
		ExpiresIn   int32          `db:"expires_in"`
		Scope       string         `db:"scope"`
		RedirectUri string         `db:"redirect_uri"`
		State       string         `db:"state"`
		Extra       sql.NullString `db:"extra"`
		CreatedAt   time.Time      `db:"created_at"`
		ExpiresAt   time.Time      `db:"expires_at"`
	}

	query := fmt.Sprintf(`SELECT client_id, code, expires_in, scope, 
		redirect_uri, state, extra, created_at, expires_at 
		FROM %stoken WHERE code = ? AND type = 'authorize'`, s.tablePrefix)

	err := s.db.QueryRowPartial(&result, query, code)

	if err == sqlx.ErrNotFound {
		return nil, osin.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("加载授权数据失败: %v", err)
	}

	// 检查是否过期
	if result.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("授权码已过期")
	}

	// 加载客户端信息
	client, err := s.GetClient(result.ClientId)
	if err != nil {
		return nil, err
	}

	data := &osin.AuthorizeData{
		Client:      client,
		Code:        result.Code,
		ExpiresIn:   result.ExpiresIn,
		Scope:       result.Scope,
		RedirectUri: result.RedirectUri,
		State:       result.State,
		CreatedAt:   result.CreatedAt,
	}

	if result.Extra.Valid {
		data.UserData = result.Extra.String
	}

	return data, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) error {
	query := fmt.Sprintf("DELETE FROM %stoken WHERE code = ? AND type = 'authorize'", s.tablePrefix)
	_, err := s.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("删除授权码失败: %v", err)
	}
	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(data *osin.AccessData) error {
	query := fmt.Sprintf(`INSERT INTO %stoken (
		id, client_id, type, access_token, refresh_token,
		expires_in, scope, redirect_uri, extra, expires_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, s.tablePrefix)

	err := s.db.Transact(func(session sqlx.Session) error {
		_, err := session.Exec(query,
			data.AccessToken,
			data.Client.GetId(),
			"access",
			data.AccessToken,
			data.RefreshToken,
			data.ExpiresIn,
			data.Scope,
			data.RedirectUri,
			toString(data.UserData),
			time.Now().Add(time.Duration(data.ExpiresIn)*time.Second),
		)
		return err
	})

	if err != nil {
		return fmt.Errorf("保存访问令牌失败: %v", err)
	}
	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(token string) (*osin.AccessData, error) {
	var result struct {
		ClientId     string         `db:"client_id"`
		AccessToken  string         `db:"access_token"`
		RefreshToken string         `db:"refresh_token"`
		ExpiresIn    int32          `db:"expires_in"`
		Scope        string         `db:"scope"`
		RedirectUri  string         `db:"redirect_uri"`
		Extra        sql.NullString `db:"extra"`
		CreatedAt    time.Time      `db:"created_at"`
		ExpiresAt    time.Time      `db:"expires_at"`
	}

	query := fmt.Sprintf(`SELECT client_id, access_token, refresh_token,
		expires_in, scope, redirect_uri, extra, created_at, expires_at
		FROM %stoken WHERE access_token = ? AND type = 'access'`, s.tablePrefix)

	err := s.db.QueryRowPartial(&result, query, token)

	if err == sqlx.ErrNotFound {
		return nil, osin.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("加载访问令牌失败: %v", err)
	}

	// 检查令牌是否过期
	if result.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("访问令牌已过期")
	}

	client, err := s.GetClient(result.ClientId)
	if err != nil {
		return nil, err
	}

	data := &osin.AccessData{
		Client:       client,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		Scope:        result.Scope,
		RedirectUri:  result.RedirectUri,
		CreatedAt:    result.CreatedAt,
	}

	if result.Extra.Valid {
		data.UserData = result.Extra.String
	}

	return data, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(token string) error {
	query := fmt.Sprintf("DELETE FROM %stoken WHERE access_token = ?", s.tablePrefix)
	err := s.db.Transact(func(session sqlx.Session) error {
		_, err := session.Exec(query, token)
		return err
	})

	if err != nil {
		return fmt.Errorf("删除访问令牌失败: %v", err)
	}
	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(token string) (*osin.AccessData, error) {
	var accessToken string
	query := fmt.Sprintf(`SELECT access_token FROM %stoken 
		WHERE refresh_token = ? AND type = 'access'`, s.tablePrefix)

	err := s.db.QueryRowPartial(&struct {
		AccessToken string `db:"access_token"`
	}{}, query, token)

	if err == sqlx.ErrNotFound {
		return nil, osin.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("加载刷新令牌失败: %v", err)
	}

	return s.LoadAccess(accessToken)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(token string) error {
	query := fmt.Sprintf("DELETE FROM %stoken WHERE refresh_token = ?", s.tablePrefix)
	_, err := s.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("删除刷新令牌失败: %v", err)
	}
	return nil
}

// CreateClientWithInformation Makes easy to create a osin.DefaultClient
func (s *Storage) CreateClientWithInformation(id string, secret string, redirectURI string, userData interface{}) osin.Client {
	return &osin.DefaultClient{
		Id:          id,
		Secret:      secret,
		RedirectUri: redirectURI,
		UserData:    userData,
	}
}

// Convert any type to string.
func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
