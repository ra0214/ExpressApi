package infraestructure

import (
	"expresApi/src/config"
	"expresApi/src/products/domain"
	"fmt"
	"log"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IProduct = (*MySQL)(nil)

func NewMySQL() domain.IProduct {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveProduct(name string, description string, price float64, category string, imageURL string) error {
	query := "INSERT INTO products (name, description, price, category, image_url) VALUES (?, ?, ?, ?, ?)"
	result, err := mysql.conn.ExecutePreparedQuery(query, name, description, price, category, imageURL)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Producto guardado correctamente: Name:%s Description:%s Price:%f Category:%s ImageURL:%s", name, description, price, category, imageURL)
	} else {
		log.Println("[MySQL] - No se insertó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) GetAll() ([]domain.Product, error) {
	query := "SELECT id, name, description, price, category, COALESCE(image_url, '') as image_url FROM products"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("Error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.ImageURL); err != nil {
			return nil, fmt.Errorf("Error al escanear la fila: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterando sobre las filas: %v", err)
	}
	return products, nil
}

func (mysql *MySQL) Delete(id string) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta DELETE: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No se encontró el producto con ID: %s", id)
	}

	log.Printf("[MySQL] - Producto eliminado correctamente con ID: %s", id)
	return nil
}

func (mysql *MySQL) Update(id string, name string, description string, price float64, category string, imageURL string) error {
	query := "UPDATE products SET name = ?, description = ?, price = ?, category = ?, image_url = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, name, description, price, category, imageURL, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta UPDATE: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No se encontró el producto con ID: %s", id)
	}

	log.Printf("[MySQL] - Producto actualizado correctamente con ID: %s", id)
	return nil
}
