package maincofiguration

import (
	"fmt"
	"log"
	"src/pkg/database" // Importa sua configuração de conexão
)

func main() {
	// 1. Conecta ao banco de dados utilizando sua função existente
	db := database.CreateDatabaseConnection()
	defer db.Close()

	// 2. Define as queries SQL
	// Importante: Ordem de criação e inserção respeitando as constraints de FK
	queries := []string{
		// Criação das Tabelas
		`CREATE TABLE IF NOT EXISTS users (
			id       SERIAL PRIMARY KEY,
			name     VARCHAR(100) NOT NULL,
			email    VARCHAR(150) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role     VARCHAR(50)  NOT NULL DEFAULT 'user'
		);`,

		`CREATE TABLE IF NOT EXISTS orders (
			id           SERIAL PRIMARY KEY,
			description  VARCHAR(255) NOT NULL,
			amount       DECIMAL(10, 2) NOT NULL,
			created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id      INTEGER NOT NULL,
			CONSTRAINT fk_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(id) 
				ON DELETE CASCADE
		);`,

		// Inserção de Usuários primeiro (necessário para o ID existir)
		`INSERT INTO users (name, email, password, role) VALUES 
			('João Silva',  'joao@email.com',   'senha123', 'admin'),
			('Maria Souza', 'maria@email.com',  'senha456', 'user'),
			('Carlos Lima', 'carlos@email.com', 'senha789', 'user')
		ON CONFLICT (email) DO NOTHING;`, // Evita erro se rodar o script duas vezes

		// Inserção de Pedidos
		`INSERT INTO orders (description, amount, user_id) VALUES 
			('Notebook Gamer', 4500.00, 1),
			('Mouse Sem Fio', 150.00, 1),
			('Teclado Mecânico', 350.00, 2),
			('Monitor 24pol', 900.00, 3);`,
	}

	// 3. Executa cada query
	for i, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Erro ao executar a operação %d: %v", i+1, err)
		}
	}

	fmt.Println("Tabelas criadas e dados inseridos com sucesso!")
}
