package nfc

import (
	"database/sql"
	"time"
)

type Capture struct {
	ID         int64     `json:"id"`
	CardID     *int64    `json:"card_id"`
	UID        string    `json:"uid"`
	ATQA       string    `json:"atqa"`
	SAK        string    `json:"sak"`
	NDEF       string    `json:"ndef"`
	CapturedAt time.Time `json:"captured_at"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) Save(c *Capture) (int64, error) {
	now := time.Now()
	res, err := r.DB.Exec(`INSERT INTO nfc_captures(card_id, uid, atqa, sak, ndef, captured_at) VALUES(?,?,?,?,?,?)`,
		c.CardID, c.UID, c.ATQA, c.SAK, c.NDEF, now)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	c.ID = id
	c.CapturedAt = now
	return id, nil
}

func (r *Repository) ListByCard(cardID int64) ([]Capture, error) {
	rows, err := r.DB.Query(`SELECT id, card_id, uid, COALESCE(atqa,''), COALESCE(sak,''), COALESCE(ndef,''), captured_at
		FROM nfc_captures WHERE card_id=? ORDER BY captured_at DESC`, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Capture
	for rows.Next() {
		var c Capture
		var cid sql.NullInt64
		if err := rows.Scan(&c.ID, &cid, &c.UID, &c.ATQA, &c.SAK, &c.NDEF, &c.CapturedAt); err != nil {
			return nil, err
		}
		if cid.Valid {
			v := cid.Int64
			c.CardID = &v
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Simulate is a placeholder for device-side HCE. On desktop it returns the stored UID/NDEF.
func (r *Repository) Simulate(cardID int64) (string, error) {
	row := r.DB.QueryRow(`SELECT COALESCE(nfc_uid,''), COALESCE(nfc_data,'') FROM cards WHERE id=?`, cardID)
	var uid, data string
	if err := row.Scan(&uid, &data); err != nil {
		return "", err
	}
	return uid, nil
}
