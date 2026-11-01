package card

import (
	"database/sql"
	"time"
)

type Card struct {
	ID             int64      `json:"id"`
	ProviderID     *int64     `json:"provider_id"`
	ProviderName   string     `json:"provider_name"`
	CardNumber     string     `json:"card_number"`
	CardHolderName string     `json:"card_holder_name"`
	BarcodeType    string     `json:"barcode_type"`
	BarcodeValue   string     `json:"barcode_value"`
	NFCUID         string     `json:"nfc_uid"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastUsedAt     *time.Time `json:"last_used_at"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) List(search string) ([]Card, error) {
	q := `SELECT id, provider_id, provider_name, card_number, card_holder_name,
		COALESCE(barcode_type,''), COALESCE(barcode_value,''), COALESCE(nfc_uid,''), COALESCE(notes,''),
		created_at, updated_at, last_used_at FROM cards`
	args := []any{}
	if search != "" {
		q += ` WHERE provider_name LIKE ? OR card_number LIKE ? OR card_holder_name LIKE ? OR notes LIKE ?`
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	q += ` ORDER BY updated_at DESC`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Card
	for rows.Next() {
		var c Card
		var pid sql.NullInt64
		var last sql.NullTime
		if err := rows.Scan(&c.ID, &pid, &c.ProviderName, &c.CardNumber, &c.CardHolderName,
			&c.BarcodeType, &c.BarcodeValue, &c.NFCUID, &c.Notes,
			&c.CreatedAt, &c.UpdatedAt, &last); err != nil {
			return nil, err
		}
		if pid.Valid {
			v := pid.Int64
			c.ProviderID = &v
		}
		if last.Valid {
			t := last.Time
			c.LastUsedAt = &t
		}
		out = append(out, c)
	}
	_ = rows
	return out, rows.Err()
}

func (r *Repository) Get(id int64) (*Card, error) {
	row := r.DB.QueryRow(`SELECT id, provider_id, provider_name, card_number, card_holder_name,
		COALESCE(barcode_type,''), COALESCE(barcode_value,''), COALESCE(nfc_uid,''), COALESCE(notes,''),
		created_at, updated_at, last_used_at FROM cards WHERE id=?`, id)
	var c Card
	var pid sql.NullInt64
	var last sql.NullTime
	if err := row.Scan(&c.ID, &pid, &c.ProviderName, &c.CardNumber, &c.CardHolderName,
		&c.BarcodeType, &c.BarcodeValue, &c.NFCUID, &c.Notes,
		&c.CreatedAt, &c.UpdatedAt, &last); err != nil {
		return nil, err
	}
	if pid.Valid {
		v := pid.Int64
		c.ProviderID = &v
	}
	if last.Valid {
		t := last.Time
		c.LastUsedAt = &t
	}
	return &c, nil
}

func (r *Repository) Create(c *Card) (int64, error) {
	now := time.Now()
	res, err := r.DB.Exec(`INSERT INTO cards(provider_id, provider_name, card_number, card_holder_name,
		barcode_type, barcode_value, nfc_uid, notes, created_at, updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?)`,
		c.ProviderID, c.ProviderName, c.CardNumber, c.CardHolderName,
		c.BarcodeType, c.BarcodeValue, c.NFCUID, c.Notes, now, now)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	c.ID = id
	c.CreatedAt = now
	c.UpdatedAt = now
	return id, nil
}

func (r *Repository) Update(c *Card) error {
	_, err := r.DB.Exec(`UPDATE cards SET provider_id=?, provider_name=?, card_number=?, card_holder_name=?,
		barcode_type=?, barcode_value=?, nfc_uid=?, notes=?, updated_at=? WHERE id=?`,
		c.ProviderID, c.ProviderName, c.CardNumber, c.CardHolderName,
		c.BarcodeType, c.BarcodeValue, c.NFCUID, c.Notes, time.Now(), c.ID)
	return err
}

func (r *Repository) Delete(id int64) error {
	_, err := r.DB.Exec("DELETE FROM cards WHERE id=?", id)
	return err
}

func (r *Repository) Touch(id int64) error {
	_, err := r.DB.Exec("UPDATE cards SET last_used_at=? WHERE id=?", time.Now(), id)
	return err
}
