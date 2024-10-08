/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wasa-1967862/service/structures"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	SearchUser(nickname string) (int, error)
	PutNewUser(nickname string) (int, error)
	SearchUserFromNick(nickname string, idUser int) ([]structures.User, error)
	ExistsUser(id int) (bool, error)
	CheckBan(users []structures.User, idUser int) ([]structures.User, error)
	PutNewBan(id_banner int, id_banned int) error
	ExistsBan(id_banner int, id_banned int) (bool, error)
	DeleteBan(id_banner int, id_banned int) error
	PutFollowing(follower_id int, followed_id int) error
	ExistsFollowing(follower_id int, followed_id int) (bool, error)
	DeleteFollowing(follower_id int, followed_id int) error
	PutNewNickname(nicknameNew string, idUser int) error
	ObtainIDFromNick(nickname string) (float64, error)
	PostNewPhoto(id_user int, path string, timestamp time.Time) (int, string, error)
	GetNickname(id int) (string, error)
	ExistsPhoto(photoId int) (bool, error)
	PostComment(id_photo int, id_user int, comment string) error
	OwnerPhotoFromIdPhoto(photoId int) (structures.User, error)
	ExistsComment(comment_id int) (bool, error)
	OwnerComment(commentId int, userId int) (bool, error)
	DeleteComment(idComment int) error
	PutLike(idPhoto int, idUser int) error
	ExistsLike(idUser int, photoId int) (bool, error)
	DeleteLike(idUser int, photoId int) error
	GetLikes(photoId int) ([]structures.User, int, error)
	CommentsPhoto(photoId int, idUser int) ([]structures.Comment, error)
	GetPhotoDate(photoId int) (time.Time, error)
	GetPhotoPath(photoId int) (string, error)
	DeletePhoto(photoId int) error
	DeleteCommentPhoto(idPhoto int) error
	DeleteLikePhoto(idPhoto int) error
	GetFollower(userId int) ([]structures.User, int, error)
	GetFollowed(userId int) ([]structures.User, int, error)
	GetPhotosProfileSorted(idProfileSearched int) ([]structures.Photo, error)
	GetStream(idProfile int) ([]structures.Photo, error)
	GetPhotoComplete(photoId int, idUser int) (structures.Photo, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// Creates all the necessary sql tables for the WASAPhoto app.
func createDatabase(db *sql.DB) error {
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			id_user INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			nickname VARCHAR(16) NOT NULL UNIQUE
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			id_photo INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user INTEGER NOT NULL,
			uploadDate DATETIME NOT NULL,
			path VARCHAR(150) NOT NULL,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS  likes (
			id_photo INTEGER NOT NULL,
			id_user INTEGER NOT NULL,
			PRIMARY KEY (id_photo,id_user),
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id_comment INTEGER PRIMARY KEY AUTOINCREMENT,
			id_photo INTEGER NOT NULL,
			id_user INTEGER NOT NULL,
			comment VARCHAR(150) NOT NULL,
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner_id INTEGER NOT NULL,
			banned_id INTEGER NOT NULL,
			PRIMARY KEY (banner_id,banned_id),
			FOREIGN KEY(banner_id) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(banned_id) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS followers(
			follower_id INTEGER NOT NULL,
			followed_id INTEGER NOT NULL,
			PRIMARY KEY (follower_id,followed_id),
			FOREIGN KEY(follower_id) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(followed_id) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
	}

	// Iteration to create all the needed sql schemas
	for i := 0; i < len(tables); i++ {
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil
}
