package crud

import (
	"context"
	"database/sql"

	"github.com/goflower-io/example/mysql/crud/user"
	"github.com/goflower-io/xsql"

	"github.com/goflower-io/xsql/mysql"
)

type Client struct {
	config *xsql.Config
	db     *xsql.DB
	Master *ClientM
	debug  bool
	User   *UserClient
}

type ClientM struct {
	User *UserClient
}

func (c *Client) init() {
	var eqx xsql.ExecQuerier
	var eqxm xsql.ExecQuerier
	eqx = c.db
	eqxm = c.db.Master()
	if c.debug {
		eqx = xsql.Debug(c.db)
		eqxm = xsql.Debug(c.db.Master())
	}
	c.User = &UserClient{eq: eqx, config: c.config}
	c.Master = &ClientM{
		User: &UserClient{eq: eqxm, config: c.config},
	}
}

type Tx struct {
	config *xsql.Config
	tx     *sql.Tx
	User   *UserClient
}

func (tx *Tx) init() {
	tx.User = &UserClient{eq: tx.tx, config: tx.config}
}

func NewClient(config *xsql.Config, debug bool) (*Client, error) {
	db, err := mysql.NewDB(config)
	if err != nil {
		return nil, err
	}
	c := &Client{config: config, db: db, debug: debug}
	c.init()
	return c, nil
}

func NewClientWithDB(db *xsql.DB, debug bool) *Client {
	c := &Client{config: db.Config(), db: db, debug: debug}
	c.init()
	return c
}

func (c *Client) Begin(ctx context.Context) (*Tx, error) {
	return c.BeginTx(ctx, nil)
}

func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := c.db.Master().BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	t := &Tx{tx: tx, config: c.config}
	t.init()
	return t, nil
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

type UserClient struct {
	eq     xsql.ExecQuerier
	config *xsql.Config
}

func (c *UserClient) Find() *xsql.SelectExecutor[*user.User] {
	return user.Find(c.eq).Timeout(c.config.QueryTimeout)
}

func (c *UserClient) Create() *user.Creater {
	a := user.Create(c.eq)
	a.Timeout(c.config.ExecTimeout)
	return a
}

func (c *UserClient) Update() *user.Updater {
	a := user.Update(c.eq)
	a.Timeout(c.config.ExecTimeout)
	return a
}

func (c *UserClient) Delete() *xsql.DeleteExecutor[*user.User] {
	return user.Delete(c.eq).Timeout(c.config.ExecTimeout)
}
