package infraestructure

import (
	"expresApi/src/config"
	"expresApi/src/users/domain"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.IUser = (*MySQL)(nil)

func NewMySQL() domain.IUser {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) SaveUser(userName string, email string, password string, estado bool) error {
	query := "INSERT INTO user (userName, email, password, estado, created_at) VALUES (?, ?, ?, ?, NOW())"
	result, err := mysql.conn.ExecutePreparedQuery(query, userName, email, password, estado)
	if err != nil {
		return fmt.Errorf("error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario creado correctamente: Username:%s Email:%s Estado:%t", userName, email, estado)
	} else {
		log.Println("[MySQL] - No se insertó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) GetAll() ([]domain.User, error) {
	query := "SELECT id, userName, email, password, estado, created_at FROM user"
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("Error al ejecutar la consulta SELECT: %v", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Estado, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("Error al escanear la fila: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterando sobre las filas: %v", err)
	}
	return users, nil
}

func (mysql *MySQL) UpdateUser(id int32, userName string, email string, password string) error {
	query := "UPDATE user SET userName = ?, email = ?, password = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, userName, email, password, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario actualizado correctamente: ID: %d Username:%s Email: %s", id, userName, email)
	} else {
		log.Println("[MySQL] - No se actualizó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) DeleteUser(id int32) error {
	query := "DELETE FROM user WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Usuario eliminado correctamente: ID: %d", id)
	} else {
		log.Println("[MySQL] - No se eliminó ninguna fila")
	}
	return nil
}

func (mysql *MySQL) GetUserByCredentials(userName string) (*domain.User, error) {
	query := "SELECT id, userName, email, password, estado, created_at FROM user WHERE userName = ?"
	row, err := mysql.conn.FetchRow(query, userName)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}

	var user domain.User
	err = row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Estado, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	return &user, nil
}

func (mysql *MySQL) ToggleUserStatus(id int32) (bool, error) {
	// Primero obtenemos el estado actual
	query := "SELECT estado FROM user WHERE id = ?"
	row, err := mysql.conn.FetchRow(query, id)
	if err != nil {
		return false, fmt.Errorf("error al obtener el usuario: %v", err)
	}

	var currentStatus bool
	err = row.Scan(&currentStatus)
	if err != nil {
		return false, fmt.Errorf("usuario no encontrado")
	}

	// Cambiamos el estado
	newStatus := !currentStatus
	updateQuery := "UPDATE user SET estado = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(updateQuery, newStatus, id)
	if err != nil {
		return false, fmt.Errorf("error al actualizar el estado: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Estado del usuario actualizado: ID: %d, Nuevo estado: %t", id, newStatus)
		return newStatus, nil
	} else {
		return false, fmt.Errorf("no se actualizó ninguna fila")
	}
}

// ToggleUserStatusController maneja la solicitud HTTP para cambiar el estado de un usuario
func (mysql *MySQL) ToggleUserStatusController(c *gin.Context) {
	// Obtener el ID del usuario desde los parámetros de la URL
	idParam := c.Param("id")

	// Convertir el ID a int32
	var id int32
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(400, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	// Llamar al método para cambiar el estado
	newStatus, err := mysql.ToggleUserStatus(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Error al cambiar el estado del usuario",
			"details": err.Error(),
		})
		return
	}

	// Responder con el nuevo estado
	c.JSON(200, gin.H{
		"message":    "Estado del usuario actualizado correctamente",
		"user_id":    id,
		"new_status": newStatus,
	})
}
