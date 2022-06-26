package db

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	OpenDatabase() error
	CreateTable() error
	DisplayAllTodos() ([]Todo, error)
	GetTodo(id string) (Todo, error)
	InsertTodo(name, description, tasks string) error
	UpdateTodo(id int, todo Todo) error
	DeleteTodo(id string) error
}

type data struct {
	logger       *log.Entry
	dbDriver     string
	dbPath       string
	dbConnection *sql.DB
}

func NewDatabase() Database {
	return &data{
		logger:   log.WithFields(log.Fields{"package": "db"}),
		dbDriver: "sqlite3",
		dbPath:   "./sqlite-database.db",
	}
}

func (d *data) OpenDatabase() error {
	database, databaseOpenErr := sql.Open(d.dbDriver, d.dbPath)
	if databaseOpenErr != nil {
		return databaseOpenErr
	}

	d.dbConnection = database
	return database.Ping()
}

func (d *data) CreateTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"description" TEXT,
		"tasks" TEXT
	  );`

	statement, err := d.dbConnection.Prepare(createTableSQL)
	if err != nil {
		d.logger.Fatal(err.Error())
		return err
	}

	statement.Exec()
	d.logger.Debugf("todos table created")
	return nil
}

func (d *data) InsertTodo(name, description, tasks string) error {
	statement, err := d.dbConnection.Prepare(`INSERT INTO todos(name, description, tasks) VALUES (?, ?, ?)`)
	defer statement.Close()
	if err != nil {
		d.logger.Fatalln(err)
	}

	_, err = statement.Exec(name, description, tasks)
	if err != nil {
		d.logger.Fatalln(err)
		return err
	}

	d.logger.Debugln("Inserted todo successfully")
	return nil
}

func (d *data) DisplayAllTodos() ([]Todo, error) {
	result := make([]Todo, 0)

	row, err := d.dbConnection.Query("SELECT * FROM todos ORDER BY id")
	defer row.Close()

	if err != nil {
		d.logger.Fatal(err)
		return result, err
	}

	for row.Next() {
		var id int
		var name, description, tasks string
		row.Scan(&id, &name, &description, &tasks)

		result = append(result, Todo{Id: id, Name: name, Description: description, Tasks: tasks})
		d.logger.Println(id, "[", name, "] ", description, "—", tasks)
	}

	return result, err
}

func (d *data) GetTodo(todoId string) (Todo, error) {
	queryResult, err := d.dbConnection.Query(fmt.Sprintf(`SELECT * from todos WHERE id = %s`, todoId))
	defer queryResult.Close()

	if err != nil {
		d.logger.Fatalln(err)
		return Todo{}, err
	}

	found := queryResult.Next()
	if found {
		var id int
		var name, description, tasks string
		queryResult.Scan(&id, &name, &description, &tasks)

		d.logger.Println(name, description, tasks)
		return Todo{id, name, description, tasks}, nil
	}

	d.logger.Info("todo not found")
	return Todo{}, fmt.Errorf("todo not found")
}

func (d *data) UpdateTodo(id int, todo Todo) error {
	statement, errPrepare := d.dbConnection.Prepare("UPDATE todos set name=?, description=?, tasks=? where id=?")
	defer statement.Close()
	if errPrepare != nil {
		d.logger.Fatalln(errPrepare)
		return errPrepare
	}

	result, errExec := statement.Exec(todo.Name, todo.Description, todo.Tasks, id)
	if errExec != nil {
		d.logger.Fatalln(errExec)
		return errExec
	}

	affect, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		d.logger.Fatalln(errRowsAffected)
		return errRowsAffected
	}

	log.Infof("%d rows affected", affect)
	return nil
}

func (d *data) DeleteTodo(id string) error {
	statement, err := d.dbConnection.Prepare(`DELETE from todos WHERE id = ?`)
	defer statement.Close()
	if err != nil {
		d.logger.Fatalln(err)
		return err
	}

	_, err = statement.Exec(id)
	if err != nil {
		d.logger.Fatalln(err)
		return err
	}

	d.logger.Debugln("Todo deleted successfully")
	return nil
}
