package server

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type Product struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Inventory int        `json:"inventory"`
	Expiry    *time.Time `json:"expiry"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (a *App) getProducts() ([]Product, error) {
	sqlString, _, _ := goqu.
		From("products").
		Select("id", "name", "inventory", "expiry", "created_at", "updated_at").
		ToSQL()

	rows, err := a.db.
		Queryx(sqlString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID, &product.Name, &product.Inventory,
			&product.Expiry, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end scan error: %w", err)
	}

	return products, nil
}

func (a *App) getProduct(productID int) (*Product, error) {
	sqlString, _, _ := goqu.
		From("products").
		Select("id", "name", "inventory", "expiry", "created_at", "updated_at").
		Where(
			goqu.C("id").Eq(productID),
		).
		ToSQL()

	var product Product

	err := a.db.
		QueryRow(sqlString).
		Scan(
			&product.ID, &product.Name, &product.Inventory,
			&product.Expiry, &product.CreatedAt, &product.UpdatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (a *App) isProductNameExist(productName string) (bool, error) {
	sqlString, _, _ := goqu.
		From("products").
		Select("id", "name").
		Where(
			goqu.C("name").Eq(productName),
		).
		ToSQL()

	var product Product

	err := a.db.
		QueryRow(sqlString).
		Scan(
			&product.ID, &product.Name,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (a *App) createProduct(productName string) error {
	sqlString, _, _ := goqu.Insert("products").
		Rows(
			goqu.Record{
				"name": productName,
			},
		).
		ToSQL()

	xres, err := a.db.Exec(sqlString)
	if err != nil {
		return err
	}

	n, err := xres.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("record not affected")
	}

	return nil
}

func (a *App) setProductExpired(productID int) error {
	sqlString, _, _ := goqu.Update("products").
		Set(
			goqu.Record{
				"expiry": time.Now(),
			},
		).
		Where(
			goqu.C("id").Eq(productID),
		).
		ToSQL()

	xres, err := a.db.Exec(sqlString)
	if err != nil {
		return err
	}

	n, err := xres.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("record not affected")
	}

	return nil
}

func (a *App) updateInventory(productID int, amount int) error {
	sqlString, _, _ := goqu.Update("products").
		Set(
			goqu.Record{
				"inventory": amount,
			},
		).
		Where(
			goqu.C("id").Eq(productID),
		).
		ToSQL()

	xres, err := a.db.Exec(sqlString)
	if err != nil {
		return err
	}

	n, err := xres.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("record not affected")
	}

	return nil
}
