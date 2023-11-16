package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/swaggo/gin-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "main.go/docs"
)

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func init() {
	var err error
	connStr := "user=postgres dbname=CRUD sslmode=disable password=15522345"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// if err = db.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

}

func main() {
	r := gin.Default()

	// Swagger documentation setup
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))

	// Define routes
	r.GET("/employees", GetEmployees)
	r.GET("/employees/:id", GetEmployee)
	r.POST("/employees", CreateEmployee)
	r.PUT("/employees/:id", UpdateEmployee)
	r.DELETE("/employees/:id", DeleteEmployee)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// @Summary Get a list of employees
// @Produce json
// @Success 200 {array} Employee
// @Router /employees [get]

func GetEmployees(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		employees = append(employees, employee)
	}

	c.JSON(http.StatusOK, employees)
}



// @Summary Get an employee by ID
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} Employee
// @Router /employees/{id} [get]

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	employee := Employee{}
	err := db.QueryRow("SELECT id, name FROM employees WHERE id = $1", id).Scan(&employee.ID, &employee.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Employee not found")
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	c.JSON(http.StatusOK, employee)
}

// @Summary Create a new employee
// @Accept json
// @Produce json
// @Param employee body Employee true "Employee Object"
// @Success 201 {object} Employee
// @Router /employees [post]

func CreateEmployee(c *gin.Context) {
	var employee Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := db.QueryRow("INSERT INTO employees(name) VALUES($1) RETURNING id", employee.Name).Scan(&employee.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// @Summary Update an existing employee
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body Employee true "Employee Object"
// @Router /employees/{id} [put]

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var updatedEmployee Employee
	if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := db.Exec("UPDATE employees SET name = $1 WHERE id = $2", updatedEmployee.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}

// @Summary Delete an employee by ID
// @Param id path int true "Employee ID"
// @Router /employees/{id} [delete]

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, "Employee not found")
		return
	}
	c.Status(http.StatusNoContent)
}
