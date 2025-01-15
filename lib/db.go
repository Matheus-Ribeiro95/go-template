package LIB

import (
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

type DB_INSERT_RESPONSE struct {
	Id   string
	Name string
}

type DB_DELETE_RESPONSE struct {
	Id string
}

type DB_UPDATE_RESPONSE struct {
	Id   string
	Name string
}

func SelectAll(session *gocql.Session, logger *zap.Logger) any {
	logger.Info("Select All Users")
	var query = session.Query("SELECT id, name FROM local.users")

	if rows, err := query.Iter().SliceMap(); err == nil {
		return rows
	} else {
		logger.Error("select all local.users: ", zap.Error(err))
		return nil
	}
}

func InsertQuery(name string, session *gocql.Session, logger *zap.Logger) (DB_INSERT_RESPONSE, bool) {
	logger.Info("Inserting User")

	var userId gocql.UUID
	var err error
	if userId, err = gocql.RandomUUID(); err != nil {
		logger.Error("insert local.users", zap.Error(err))
		return DB_INSERT_RESPONSE{}, false
	}

	if err = session.Query("INSERT INTO local.users (id, name) VALUES (?, ?)", userId, name).Exec(); err != nil {
		logger.Error("insert local.users", zap.Error(err))
		return DB_INSERT_RESPONSE{}, false
	}

	return DB_INSERT_RESPONSE{
		Id:   userId.String(),
		Name: name,
	}, true
}

func DeleteQuery(id string, session *gocql.Session, logger *zap.Logger) (DB_DELETE_RESPONSE, bool) {
	logger.Info("Deleting User")

	var (
		err  error
		q    int
		qNew int
	)

	err = session.Query("SELECT COUNT(*) FROM local.users WHERE id = ?", id).Scan(&q)
	if err != nil {
		logger.Error("delete local.users", zap.Error(err))
		return DB_DELETE_RESPONSE{}, false
	}

	if err = session.Query("DELETE FROM local.users WHERE id = ?", id).Exec(); err != nil {
		logger.Error("delete local.users", zap.Error(err))
		return DB_DELETE_RESPONSE{}, false
	}

	err = session.Query("SELECT COUNT(*) FROM local.users WHERE id = ?", id).Scan(&qNew)
	if err != nil {
		logger.Error("delete local.users", zap.Error(err))
		return DB_DELETE_RESPONSE{}, false
	}
	if q == qNew {
		logger.Error("delete local.users", zap.String("message", "User Id not found"))
		return DB_DELETE_RESPONSE{}, false
	}

	return DB_DELETE_RESPONSE{
		Id: id,
	}, true
}

func UpdateQuery(id string, name string, session *gocql.Session, logger *zap.Logger) (DB_UPDATE_RESPONSE, bool) {
	logger.Info("Updating User")

	var (
		err error
		q   int
	)

	err = session.Query("SELECT COUNT(id) AS q FROM local.users WHERE id = ?", id).Scan(&q)
	if err != nil {
		logger.Error("update local.users", zap.Error(err))
		return DB_UPDATE_RESPONSE{}, false
	}
	if q == 0 {
		logger.Error("update local.users", zap.String("message", "User Id not found"))
		return DB_UPDATE_RESPONSE{}, false
	}

	if err := session.Query("UPDATE local.users SET name = ? WHERE id = ?", name, id).Exec(); err != nil {
		logger.Error("update local.users", zap.Error(err))
		return DB_UPDATE_RESPONSE{}, false
	}
	return DB_UPDATE_RESPONSE{
		Id:   id,
		Name: name,
	}, true
}
