package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/iambakhodir/short-link/domain"
	"github.com/sirupsen/logrus"
)

type mysqlTagsRepository struct {
	Conn *sql.DB
}

func NewMysqlTagsRepository(conn *sql.DB) domain.TagsRepository {
	return &mysqlTagsRepository{Conn: conn}
}

func (m *mysqlTagsRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Tags, err error) {
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

	result = make([]domain.Tags, 0)
	for rows.Next() {
		t := domain.Tags{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlTagsRepository) Fetch(ctx context.Context, limit int64) ([]domain.Tags, error) {
	query := `SELECT id, name, created_at, updated_at
				FROM tags ORDER BY created_at LIMIT ?`

	res, err := m.fetch(ctx, query, limit)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlTagsRepository) GetById(ctx context.Context, id int64) (domain.Tags, error) {
	query := `SELECT id, name, created_at, updated_at
				FROM tags where id = ?`

	list, err := m.fetch(ctx, query, id)

	if err != nil {
		return domain.Tags{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	} else {
		return domain.Tags{}, domain.ErrNotFound
	}
}

func (m *mysqlTagsRepository) GetByName(ctx context.Context, name string) (domain.Tags, error) {
	query := `SELECT id, name, created_at, updated_at
				FROM tags where name = ?`

	list, err := m.fetch(ctx, query, name)

	if err != nil {
		return domain.Tags{}, err
	}

	if len(list) > 0 {
		return list[0], nil
	} else {
		return domain.Tags{}, domain.ErrNotFound
	}
}

func (m *mysqlTagsRepository) FetchByLinkId(ctx context.Context, linkId int64) ([]domain.Tags, error) {
	query := `SELECT t.id, t.name, t.created_at, t.updated_at
				FROM tags as t LEFT JOIN link_tag as lt 
				    ON t.id = lt.tag_id where lt.link_id = ?`

	return m.fetch(ctx, query, linkId)
}

func (m *mysqlTagsRepository) Update(ctx context.Context, tags domain.Tags) (int64, error) {
	query := `UPDATE tags SET name = ?, updated_at = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, tags.Name, tags.ID)

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

	return tags.ID, nil
}

func (m *mysqlTagsRepository) Store(ctx context.Context, tags domain.Tags) (int64, error) {
	query := `INSERT tags SET name = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, tags.Name)
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
func (m *mysqlTagsRepository) FirstOrCreate(ctx context.Context, tags domain.Tags) (int64, error) {
	tag, err := m.GetByName(ctx, tags.Name)
	if err == nil {
		return tag.ID, nil
	}

	query := `INSERT tags SET name = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, tags.Name)
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

func (m *mysqlTagsRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tags WHERE id = ?`

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
