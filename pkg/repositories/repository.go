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

func CreateOrder(order entities.Order) error {
	db := database.CreateDatabaseConnection()
	defer db.Close()

	_, err := db.Exec("INSERT INTO orders (description, amount, user_id) VALUES ($1, $2, $3)",
		order.Description, order.Amount, order.UserID)
	return err
}

func GetOrdersByUserID(userID int) ([]entities.Order, error) {
	db := database.CreateDatabaseConnection()
	defer db.Close()

	rows, err := db.Query("SELECT id, description, amount, created_at, user_id FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(&order.ID, &order.Description, &order.Amount, &order.CreatedAt, &order.UserID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func DeleteOrder(orderID int) error {
	db := database.CreateDatabaseConnection()
	defer db.Close()

	_, err := db.Exec("DELETE FROM orders WHERE id = $1", orderID)
	return err
}
