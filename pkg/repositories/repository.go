package repositories

import (
	"database/sql"
	"src/pkg/database"
	"src/pkg/entities"
)

func GetAllUsers() ([]entities.User, error) {
	var db = database.CreateDatabaseConnection()

	rows, err := db.Query("SELECT id, name, email, password, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	db.Close()
	return users, nil
}

func GetUser(username string) (entities.User, error) {
	var db = database.CreateDatabaseConnection()

	var user entities.User
	err := db.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = $1", username).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, nil
		}
		return entities.User{}, err
	}
	db.Close()
	return user, nil
}

func CreateUser(user entities.User) error {
	var db = database.CreateDatabaseConnection()

	_, err := db.Exec("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
