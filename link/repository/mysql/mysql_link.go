package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/iambakhodir/short-link/domain"
	"github.com/sirupsen/logrus"
	"time"
)

type mysqlLinkRepository struct {
	Conn *sql.DB
}

func NewMysqlLinkRepository(conn *sql.DB) domain.LinkRepository {
	return &mysqlLinkRepository{Conn: conn}
}

func (m *mysqlLinkRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Link, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Link, 0)
	for rows.Next() {
		t := domain.Link{}
		err = rows.Scan(
			&t.ID,
			&t.UserId,
			&t.Alias,
			&t.Target,
			&t.Description,
			&t.UpdatedAt,
			&t.CreatedAt,
			&t.DeletedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlLinkRepository) Fetch(ctx context.Context, limit int64) ([]domain.Link, error) {
	query := `SELECT id, user_id, alias, target, description, created_at, updated_at, deleted_at
				FROM link ORDER BY created_at LIMIT ?`

	res, err := m.fetch(ctx, query, limit)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlLinkRepository) GetById(ctx context.Context, id int64) (domain.Link, error) {
	query := `SELECT id, user_id, alias, target, description, created_at, updated_at, deleted_at
				FROM link where id = ?`

	list, err := m.fetch(ctx, query, id)

	if err != nil {
		return domain.Link{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	} else {
		return domain.Link{}, domain.ErrNotFound
	}
}

func (m *mysqlLinkRepository) Update(ctx context.Context, link domain.Link) (int64, error) {
	query := `UPDATE link SET alias = ?, target = ?, user_id =?, deleted_at = ?, updated_at = ?, description=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, link.Alias, link.Target, link.UserId, link.DeletedAt, link.UpdatedAt, link.Description, link.ID)

	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	if affect != 1 {
		return 0, fmt.Errorf("Total Affected: %d", affect)
	}

	return link.ID, nil
}

func (m *mysqlLinkRepository) GetByAlias(ctx context.Context, alias string) (domain.Link, error) {
	query := `SELECT id, user_id, alias, target, description, created_at, updated_at, deleted_at
				FROM link where alias = ?`

	list, err := m.fetch(ctx, query, alias)

	if err != nil {
		return domain.Link{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	} else {
		return domain.Link{}, domain.ErrNotFound
	}
}

func (m *mysqlLinkRepository) Store(ctx context.Context, link domain.Link) (int64, error) {
	query := `INSERT link SET alias = ?, target = ?, user_id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, link.Alias, link.Target, link.UserId)
	if err != nil {
		mysqlErr, _ := err.(*mysql.MySQLError)
		if mysqlErr.Number == 1062 { // MySQL error code for "Duplicate entry"
			return 0, domain.ErrLinkIsExists
		}

		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected != 1 {
		return 0, fmt.Errorf("Total affected: %d", rowsAffected)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *mysqlLinkRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE link SET deleted_at = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("Total Affected: %d", rowsAffected)
	}

	return nil
}
