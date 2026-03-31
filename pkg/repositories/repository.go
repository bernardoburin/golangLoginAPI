package repositories

import (
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

func login(username, password string) (bool, error) {
	var db = database.CreateDatabaseConnection()
	
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = $1", username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	db.Close()
	
	if storedPassword != password {
		return false, nil
	}

	return true, nil
}