package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bueti/yeoman/internal/datasource"
)

type DatasourceRow struct {
	ID   string
	Name sql.NullString
	URL  sql.NullString
}

func convertDatasourceRowToDatasource(c DatasourceRow) datasource.Datasource {
	return datasource.Datasource{
		ID:   c.ID,
		Name: c.Name.String,
		URL:  c.URL.String,
	}
}

func (d *Database) GetDatasource(ctx context.Context, uuid string) (datasource.Datasource, error) {
	var dsRow DatasourceRow
	row := d.Client.QueryRowContext(ctx, `SELECT id, name, url FROM datasources WHERE id = $1`, uuid)
	err := row.Scan(&dsRow.ID, &dsRow.Name, &dsRow.URL)
	if err != nil {
		return datasource.Datasource{}, fmt.Errorf("error fetching the datasource by uuid: %w", err)
	}

	return convertDatasourceRowToDatasource(dsRow), nil
}

func (d *Database) PostDatasource(ctx context.Context, ds datasource.Datasource) (datasource.Datasource, error) {
	dsRow := DatasourceRow{
		Name: sql.NullString{String: ds.Name, Valid: true},
		URL:  sql.NullString{String: ds.URL, Valid: true},
	}
	rows, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO datasources
		(name, url)
		VALUES
		(:name, :url)`,
		dsRow,
	)
	if err != nil {
		return datasource.Datasource{}, fmt.Errorf("failed to insert datasource: %w", err)
	}
	if err := rows.Close(); err != nil {
		return datasource.Datasource{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return ds, nil
}

func (d *Database) DeleteDatasource(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM datasources WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete error: %w", err)
	}
	return nil
}

func (d *Database) UpdateDatasource(ctx context.Context, id string, ds datasource.Datasource) (datasource.Datasource, error) {
	dsRow := DatasourceRow{
		ID:   id,
		Name: sql.NullString{String: ds.Name, Valid: true},
		URL:  sql.NullString{String: ds.URL, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE datasources SET
		name = :name,
		url = :url
		WHERE id = :id`,
		dsRow,
	)
	if err != nil {
		return datasource.Datasource{}, fmt.Errorf("failed to update the datasource: %w", err)
	}
	if err := rows.Close(); err != nil {
		return datasource.Datasource{}, fmt.Errorf("failed to close rows: %w", err)
	}
	return convertDatasourceRowToDatasource(dsRow), nil
}
