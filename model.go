package awesomeProject
import (
	"database/sql"
	_ "errors"
)

type Car struct {
	ID    int   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Year  float64 `json:"year"`
}

func (car *Car) getCar(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM cars WHERE id=$1",
		car.ID).Scan(&car.Name, &car.Price,&car.Year)
}

func getCars(db *sql.DB, start, count int) ([]Car, error) {
	rows, err := db.Query(
		"SELECT id, name,price ,year FROM cars",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cars := []Car{}

	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Name, &car.Price,&car.Year); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (car *Car) updateCar(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE cars SET name=$1, price=$2 ,year=$3 WHERE id=$4",
			car.Name, car.Price,&car.Year ,car.ID)

	return err
}

func (car *Car) deleteCar(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM cars WHERE id=$1", car.ID)

	return err
}

func (car *Car) createCar(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO cars(name, price ,year ) VALUES($1, $2 ,$3) RETURNING id",
		car.Name, car.Price ,&car.Year ).Scan(&car.ID)

	if err != nil {
		return err
	}

	return nil
}