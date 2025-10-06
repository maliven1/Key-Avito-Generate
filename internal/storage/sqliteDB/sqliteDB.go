package sqliteDB

import (
	"avito/internal/logger"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func CreateDB(storagePath string) (*Storage, bool, error) {
	const op = "storage.sqliteDB.New"
	log := logger.Log.With(
		slog.String("op", op),
	)
	var install bool
	if _, err := os.Stat(storagePath); err != nil {
		if os.IsNotExist(err) {
			log.Info("База данных будет создана")
			install = true
		} else {
			log.Error("не получилось проверить файл", err)
			os.Exit(1)
		}
	}
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		log.Error("BD not found", err)
		os.Exit(1)
	}
	if install {
		db, err := sql.Open("sqlite", storagePath)
		if err != nil {
			return nil, install, fmt.Errorf("%s: %w", op, err)
		}

		stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS KEY(
	id INTEGER PRIMARY KEY,
	Hkey TEXT ,
	Mkey TEXT ,
	Lkey TEXT )`)
		if err != nil {
			return nil, install, fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, install, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{db: db}, install, nil
}

func (s *Storage) SaveHighKey(key string) (int64, error) {
	const op = "storage.sqliteDB.SaveHighKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("INSERT INTO KEY(Hkey) VALUES(?)")
	if err != nil {
		log.Error("Save Hkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.Exec(key)
	if err != nil {
		log.Error("Read Hkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("failed to get last insert id", err)
		return 0, fmt.Errorf("%s: failed to get last insert id %w", op, err)
	}
	return id, nil
}

func (s *Storage) SaveMiddleKey(key string) (int64, error) {
	const op = "storage.sqliteDB.SaveMiddleKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("INSERT INTO KEY(Mkey) VALUES(?)")
	if err != nil {
		log.Error("Save Mkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.Exec(key)
	if err != nil {
		log.Error("Read Mkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("failed to get last insert id", err)
		return 0, fmt.Errorf("%s: failed to get last insert id %w", op, err)
	}
	return id, nil
}

func (s *Storage) SaveLowKey(key string) (int64, error) {
	const op = "storage.sqliteDB.SaveLowKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("INSERT INTO KEY(Lkey) VALUES(?)")
	if err != nil {
		log.Error("Save Lkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.Exec(key)
	if err != nil {
		log.Error("Read Lkey error", err)
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("Failed to get last insert id", err)
		return 0, fmt.Errorf("%s: failed to get last insert id %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetHightKey() ([]string, error) {
	const op = "storage.sqlite.GetHightKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("SELECT Hkey FROM KEY")
	if err != nil {
		log.Error("prepare statement:", err)
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resKey []string
	var key string
	row, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		log.Error("execute statement:", err)
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	for row.Next() {
		err := row.Scan(&key)
		if err != nil {
			continue
		}
		resKey = append(resKey, key)
	}
	return resKey, nil
}

func (s *Storage) GetMiddleKey() ([]string, error) {
	const op = "storage.sqlite.GetMiddleKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("SELECT Mkey FROM KEY")
	if err != nil {
		log.Error("prepare statement:", err)
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resKey []string
	var key string
	row, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		log.Error("execute statement:", err)
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	for row.Next() {
		err := row.Scan(&key)
		if err != nil {
			continue
		}
		resKey = append(resKey, key)
	}
	return resKey, nil
}

func (s *Storage) GetLowKey() ([]string, error) {
	const op = "storage.sqlite.GetLowKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	stmt, err := s.db.Prepare("SELECT Lkey FROM KEY")
	if err != nil {
		log.Error("prepare statement:", err)
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resKey []string
	var key string
	row, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		log.Error("execute statement:", err)
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	for row.Next() {
		err := row.Scan(&key)
		if err != nil {
			continue
		}
		resKey = append(resKey, key)
	}
	return resKey, nil
}

func (s *Storage) DeleteHightKey(key string) error {
	const op = "storage.sqlite.DeleteHightKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	query := "DELETE FROM KEY WHERE Hkey = ?"
	_, err := s.db.Query(query, key)
	if err != nil {
		log.Error("delete key error: ", err)
		return fmt.Errorf("%s: delete key error: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteMeddleKey(key string) error {
	const op = "storage.sqlite.DeleteHightKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	query := "DELETE FROM KEY WHERE Mkey = ?"
	_, err := s.db.Query(query, key)
	if err != nil {
		log.Error("delete key error: ", err)
		return fmt.Errorf("%s: delete key error: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteLowKey(key string) error {
	const op = "storage.sqlite.DeleteHightKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	query := "DELETE FROM KEY WHERE Lkey = ?"
	_, err := s.db.Query(query, key)
	if err != nil {
		log.Error("delete key error: ", err)
		return fmt.Errorf("%s: delete key error: %w", op, err)
	}
	return nil
}
