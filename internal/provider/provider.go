package provider

import (
	"database/sql"
)

type Provider struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	LogoURL       string `json:"logo_url"`
	BarcodeFormat string `json:"barcode_format"`
	Country       string `json:"country"`
	Notes         string `json:"notes"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) List(search string) ([]Provider, error) {
	q := `SELECT id, name, COALESCE(logo_url,''), COALESCE(barcode_format,''), COALESCE(country,''), COALESCE(notes,'') FROM providers`
	args := []any{}
	if search != "" {
		q += ` WHERE name LIKE ? OR country LIKE ?`
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	q += ` ORDER BY name ASC`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Provider
	for rows.Next() {
		var p Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.LogoURL, &p.BarcodeFormat, &p.Country, &p.Notes); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) Get(id int64) (*Provider, error) {
	row := r.DB.QueryRow(`SELECT id, name, COALESCE(logo_url,''), COALESCE(barcode_format,''), COALESCE(country,''), COALESCE(notes,'') FROM providers WHERE id=?`, id)
	var p Provider
	if err := row.Scan(&p.ID, &p.Name, &p.LogoURL, &p.BarcodeFormat, &p.Country, &p.Notes); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Create(p *Provider) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO providers(name, logo_url, barcode_format, country, notes) VALUES(?,?,?,?,?)`,
		p.Name, p.LogoURL, p.BarcodeFormat, p.Country, p.Notes)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.DB.Exec("DELETE FROM providers WHERE id=?", id)
	return err
}
