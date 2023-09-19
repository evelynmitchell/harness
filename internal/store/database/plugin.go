// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.PluginStore = (*pluginStore)(nil)

const (
	pluginColumns = `
	plugin_uid
	,plugin_description
	,plugin_type
	,plugin_version
	,plugin_logo
	,plugin_spec
	`
)

// NewPluginStore returns a new PluginStore.
func NewPluginStore(db *sqlx.DB) *pluginStore {
	return &pluginStore{
		db: db,
	}
}

type pluginStore struct {
	db *sqlx.DB
}

// Create creates a new entry in the plugin datastore.
func (s *pluginStore) Create(ctx context.Context, plugin *types.Plugin) error {
	const pluginInsertStmt = `
	INSERT INTO plugins (
		plugin_uid
		,plugin_description
		,plugin_type
		,plugin_version
		,plugin_logo
		,plugin_spec
	) VALUES (
		:plugin_uid
		,:plugin_description
		,:plugin_type
		,:plugin_version
		,:plugin_logo
		,:plugin_spec
	) RETURNING plugin_uid`

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(pluginInsertStmt, plugin)
	if err != nil {
		return database.ProcessSQLErrorf(err, "Failed to bind plugin object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&plugin.UID); err != nil {
		return database.ProcessSQLErrorf(err, "plugin query failed")
	}

	return nil
}

// Find finds a version of a plugin
func (s *pluginStore) Find(ctx context.Context, name, version string) (*types.Plugin, error) {
	const pluginFindStmt = `
	SELECT` + pluginColumns +
		`FROM plugins
	WHERE plugin_uid = $1 AND plugin_version = $2
	`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.Plugin)
	if err := db.GetContext(ctx, dst, pluginFindStmt, name, version); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to find pipeline")
	}

	return dst, nil
}

// List returns back the list of plugins along with their associated schemas.
func (s *pluginStore) List(
	ctx context.Context,
	filter types.ListQueryFilter,
) ([]*types.Plugin, error) {
	stmt := database.Builder.
		Select(pluginColumns).
		From("plugins")

	if filter.Query != "" {
		stmt = stmt.Where("LOWER(plugin_uid) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query)))
	}

	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []*types.Plugin{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed executing custom list query")
	}

	return dst, nil
}

// ListAll returns back the full list of plugins in the database.
func (s *pluginStore) ListAll(
	ctx context.Context,
) ([]*types.Plugin, error) {
	stmt := database.Builder.
		Select(pluginColumns).
		From("plugins")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []*types.Plugin{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed executing custom list query")
	}

	return dst, nil
}

// Count of plugins matching the filter criteria.
func (s *pluginStore) Count(ctx context.Context, filter types.ListQueryFilter) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("plugins")

	if filter.Query != "" {
		stmt = stmt.Where("LOWER(plugin_uid) LIKE ?", fmt.Sprintf("%%%s%%", filter.Query))
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(err, "Failed executing count query")
	}
	return count, nil
}

// Update updates a plugin row.
func (s *pluginStore) Update(ctx context.Context, p *types.Plugin) error {
	const pluginUpdateStmt = `
	UPDATE plugins
	SET
		plugin_description = :plugin_description
		,plugin_type = :plugin_type
		,plugin_version = :plugin_version
		,plugin_logo = :plugin_logo
		,plugin_spec = :plugin_spec
	WHERE plugin_uid = :plugin_uid`
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(pluginUpdateStmt, p)
	if err != nil {
		return database.ProcessSQLErrorf(err, "Failed to bind plugin object")
	}

	_, err = db.ExecContext(ctx, query, arg...)
	if err != nil {
		return database.ProcessSQLErrorf(err, "Failed to update plugin")
	}

	return nil
}
