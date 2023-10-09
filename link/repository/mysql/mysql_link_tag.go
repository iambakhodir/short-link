package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/iambakhodir/short-link/domain"
	"github.com/sirupsen/logrus"
)

type mysqlLinkTagRepository struct {
	Conn *sql.DB
}

func NewMysqlLinkTagRepository(conn *sql.DB) domain.LinkTagRepository {
	return &mysqlLinkTagRepository{Conn: conn}
}

func (m *mysqlLinkTagRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.LinkTag, err error) {
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

	result = make([]domain.LinkTag, 0)
	for rows.Next() {
		t := domain.LinkTag{}
		err = rows.Scan(
			&t.ID,
			&t.LinkId,
			&t.TagId,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlLinkTagRepository) Fetch(ctx context.Context, limit int64) ([]domain.LinkTag, error) {
	query := `SELECT id, link_id, tag_id, created_at, updated_at
				FROM link_tag ORDER BY created_at LIMIT ?`

	res, err := m.fetch(ctx, query, limit)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlLinkTagRepository) GetById(ctx context.Context, id int64) (domain.LinkTag, error) {
	query := `SELECT id, link_id, tag_id, created_at, updated_at
				FROM link_tag where id = ?`

	list, err := m.fetch(ctx, query, id)

	if err != nil {
		return domain.LinkTag{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	} else {
		return domain.LinkTag{}, domain.ErrNotFound
	}
}

func (m *mysqlLinkTagRepository) Update(ctx context.Context, linkTag domain.LinkTag) (int64, error) {
	query := `UPDATE link_tag SET link_id = ?, tag_id = ?, updated_at = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, linkTag.LinkId, linkTag.TagId, linkTag.UpdatedAt, linkTag.ID)

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

	return linkTag.ID, nil
}

func (m *mysqlLinkTagRepository) Store(ctx context.Context, linkTag domain.LinkTag) (int64, error) {
	query := `INSERT link_tag SET link_id = ?, tag_id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, linkTag.LinkId, linkTag.TagId)
	if err != nil {
		mysqlErr, _ := err.(*mysql.MySQLError)
		if mysqlErr.Number == 1062 { // MySQL error code for "Duplicate entry"
			return 0, domain.ErrConflict
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

func (m *mysqlLinkTagRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM link_tag WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
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
