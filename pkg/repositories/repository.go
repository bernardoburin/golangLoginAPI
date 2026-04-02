package repositories

import (
	"database/sql"
	"src/pkg/database"
	"src/pkg/entities"
)

func GetAllUsers() ([]entities.User, error) {
	var db = database.CreateDatabaseConnection()

	rows, err := db.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
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
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE name = $1", username).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, nil
		}
		return entities.User{}, err
	}
	db.Close()
	return user, nil
}
