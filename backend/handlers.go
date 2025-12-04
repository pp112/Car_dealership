package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB *sql.DB
}

func (s *Server) routes() *gin.Engine {
	r := gin.Default()
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	r.GET("/api/search", s.handleSearch)
	r.GET("/api/car/:id", s.handleGetCar)
	r.GET("/api/all", s.handleAllCars)

	return r
}

// ---------------------- SEARCH --------------------------

func (s *Server) handleSearch(c *gin.Context) {
	field := c.Query("field")
	q := c.Query("q")

	if field == "" || q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field and q required"})
		return
	}

	var rows *sql.Rows
	var err error

	switch field {
	case "brand":
		rows, err = s.DB.Query(`
            SELECT cars.id, brands.name, models.name, cars.year, cars.price
            FROM cars
            JOIN brands ON cars.brand_id = brands.id
            JOIN models ON cars.model_id = models.id
            WHERE LOWER(brands.name) LIKE $1
            ORDER BY cars.id
        `, "%"+strings.ToLower(q)+"%")

	case "model":
		rows, err = s.DB.Query(`
            SELECT cars.id, brands.name, models.name, cars.year, cars.price
            FROM cars
            JOIN brands ON cars.brand_id = brands.id
            JOIN models ON cars.model_id = models.id
            WHERE LOWER(models.name) LIKE $1
            ORDER BY cars.id
        `, "%"+strings.ToLower(q)+"%")

	case "year":
		rows, err = s.queryNumeric("year", q)

	case "price":
		rows, err = s.queryNumeric("price", q)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown field"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	cars := []Car{}
	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cars = append(cars, car)
	}

	c.JSON(http.StatusOK, cars)
}

// ---------------------- NUMERIC SEARCH --------------------------

func (s *Server) queryNumeric(col string, q string) (*sql.Rows, error) {

	op := "="
	valStr := q

	if strings.HasPrefix(q, ">=") || strings.HasPrefix(q, "<=") ||
		strings.HasPrefix(q, ">") || strings.HasPrefix(q, "<") ||
		strings.HasPrefix(q, "=") {

		if strings.HasPrefix(q, ">=") || strings.HasPrefix(q, "<=") {
			op = q[:2]
			valStr = strings.TrimSpace(q[2:])
		} else {
			op = q[:1]
			valStr = strings.TrimSpace(q[1:])
		}
	} else {
		op = "="
		valStr = strings.TrimSpace(q)
	}

	sqlStr := `
        SELECT cars.id, brands.name, models.name, cars.year, cars.price
        FROM cars
        JOIN brands ON cars.brand_id = brands.id
        JOIN models ON cars.model_id = models.id
        WHERE ` + col + " " + op + " $1 ORDER BY cars.id"

	return s.DB.Query(sqlStr, valStr)
}

// ---------------------- GET CAR BY ID --------------------------

func (s *Server) handleGetCar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	row := s.DB.QueryRow(`
        SELECT cars.id, brands.name, models.name, cars.year, cars.price
        FROM cars
        JOIN brands ON cars.brand_id = brands.id
        JOIN models ON cars.model_id = models.id
        WHERE cars.id = $1
    `, id)

	var car Car
	err = row.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, car)
}

// ---------------------- GET ALL CARS --------------------------

func (s *Server) handleAllCars(c *gin.Context) {
	rows, err := s.DB.Query(`
        SELECT cars.id, brands.name, models.name, cars.year, cars.price
        FROM cars
        JOIN brands ON cars.brand_id = brands.id
        JOIN models ON cars.model_id = models.id
        ORDER BY cars.id
    `)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	cars := []Car{}
	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cars = append(cars, car)
	}

	c.JSON(http.StatusOK, cars)
}
