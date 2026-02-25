package stock // thats the namespace in go

import (
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) ListProducts() error {
	tx, err := r.DB.Begin() // wrapping in db.begin and commit will result in consistent atomic call
	if err != nil {
		return err
	}

	defer tx.Rollback() //on failure, roll back

	rows, err := r.DB.Query("SELECT sku, name, stock FROM products ORDER BY sku")
	if err != nil {
		return err
	}
	defer rows.Close() // defer turns out to be some kind of scheduler, like a flag. it runs on exit.

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0) // thats the frame for the tui
	fmt.Fprintln(w, "SKU\tNAME\tSTOCK")
	for rows.Next() {
		var sku, name string
		var stock int
		if err := rows.Scan(&sku, &name, &stock); err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", sku, name, stock)
	}
	w.Flush() // so it turns out go is terribly lazy and it had been operating like a stringbuilder, this flush is where it writes it to the screen
	return tx.Commit()
}

func (r *Repository) UpdateStock(id, sku string, qty int, isIncrement bool, note string) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be > 0")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // another flag to roll back anything in case of an error. independent from the above code

	// 1. Idempotency Check
	var exists int
	err = tx.QueryRow("SELECT 1 FROM orders WHERE id = $1", id).Scan(&exists) // seems like this dollar sign is how you pass value in psql
	if err == nil {
		fmt.Println("Skipped/Already applied")
		return nil
	}

	// 2. Update Product Stock
	var updateQuery string
	if isIncrement {
		updateQuery = "UPDATE products SET stock = stock + $1 WHERE sku = $2 RETURNING stock"
	} else {
		updateQuery = "UPDATE products SET stock = stock - $1 WHERE sku = $2 RETURNING stock"
	}

	var newStock int
	err = tx.QueryRow(updateQuery, qty, sku).Scan(&newStock) // also this tx.queryrow is nice, looks weird tho, another lazy loading i guess
	if err != nil {
		return fmt.Errorf("failed to update stock (check SKU or sufficient stock): %v", err)
	}

	// 3. Insert Order
	_, err = tx.Exec(`
	INSERT INTO orders (id, sku, quantity, is_increment, note)
	VALUES ($1, $2, $3, $4, $5)`,
		id, sku, qty, isIncrement, note)
	if err != nil {
		return err
	}
	// yapp, exec is where it gets executed
	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Success! New stock for %s: %d\n", sku, newStock)
	return nil
}

func (r *Repository) Report(topN, lowStockThreshold int) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// 1. Start Transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var totalSKUs int
	err = tx.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalSKUs)
	if err != nil {
		return fmt.Errorf("failed to count SKUs: %v", err)
	}

	var totalStock int
	err = tx.QueryRow("SELECT COALESCE(SUM(stock), 0) FROM products").Scan(&totalStock)
	if err != nil {
		return fmt.Errorf("failed to sum stock: %v", err)
	}

	fmt.Printf("Total SKU count: %d\n", totalSKUs)
	fmt.Printf("Total stock units: %d\n\n", totalStock)

	fmt.Printf("--- Top %d Products by Stock ---\n", topN)
	rows, err := tx.Query("SELECT sku, name, stock FROM products ORDER BY stock DESC LIMIT $1", topN) // the equivalent of linq .Take(x)
	if err != nil {
		return fmt.Errorf("failed to get top products: %v", err)
	}
	// i was told rows must be closed now while defer could work as well

	fmt.Fprintln(w, "SKU\tNAME\tSTOCK")
	for rows.Next() {
		var sku, name string
		var stock int
		if err := rows.Scan(&sku, &name, &stock); err != nil {
			rows.Close() // Clean up on error
			return err
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", sku, name, stock)
	}
	w.Flush()
	rows.Close() // Close explicitly
	fmt.Println()

	fmt.Printf("--- Low Stock (<= %d) ---\n", lowStockThreshold)
	lowRows, err := tx.Query("SELECT sku, name, stock FROM products WHERE stock <= $1 ORDER BY stock ASC", lowStockThreshold)
	if err != nil {
		return fmt.Errorf("failed to get low stock products: %v", err)
	}

	fmt.Fprintln(w, "SKU\tNAME\tSTOCK")
	for lowRows.Next() {
		var sku, name string
		var stock int
		if err := lowRows.Scan(&sku, &name, &stock); err != nil {
			lowRows.Close()
			return err
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", sku, name, stock)
	}
	w.Flush()
	lowRows.Close()

	return tx.Commit()
}
