package main

import (
	"github.com/friesipayung/cardVault/internal/card"
	"github.com/friesipayung/cardVault/internal/db"
	"github.com/friesipayung/cardVault/internal/nfc"
	"github.com/friesipayung/cardVault/internal/provider"
)

// CardService exposes card operations to the frontend.
type CardService struct {
	db        *db.DB
	cards     *card.Repository
	providers *provider.Repository
	nfc       *nfc.Repository
}

func NewCardService(d *db.DB) *CardService {
	return &CardService{
		db:        d,
		cards:     card.NewRepository(d.DB),
		providers: provider.NewRepository(d.DB),
		nfc:       nfc.NewRepository(d.DB),
	}
}

func (s *CardService) ListCards(search string) ([]card.Card, error) {
	return s.cards.List(search)
}

func (s *CardService) GetCard(id int64) (*card.Card, error) {
	return s.cards.Get(id)
}

func (s *CardService) AddCard(c card.Card) (card.Card, error) {
	id, err := s.cards.Create(&c)
	if err != nil {
		return card.Card{}, err
	}
	got, err := s.cards.Get(id)
	if err != nil {
		return c, nil
	}
	return *got, nil
}

func (s *CardService) UpdateCard(c card.Card) error {
	return s.cards.Update(&c)
}

func (s *CardService) DeleteCard(id int64) error {
	return s.cards.Delete(id)
}

func (s *CardService) TouchCard(id int64) error {
	return s.cards.Touch(id)
}

// ProviderService exposes provider operations.
type ProviderService struct {
	providers *provider.Repository
}

func NewProviderService(d *db.DB) *ProviderService {
	return &ProviderService{providers: provider.NewRepository(d.DB)}
}

func (s *ProviderService) ListProviders(search string) ([]provider.Provider, error) {
	return s.providers.List(search)
}

func (s *ProviderService) GetProvider(id int64) (*provider.Provider, error) {
	return s.providers.Get(id)
}

func (s *ProviderService) AddProvider(p provider.Provider) (provider.Provider, error) {
	id, err := s.providers.Create(&p)
	if err != nil {
		return provider.Provider{}, err
	}
	got, err := s.providers.Get(id)
	if err != nil {
		return p, nil
	}
	return *got, nil
}

func (s *ProviderService) DeleteProvider(id int64) error {
	return s.providers.Delete(id)
}

// NFCService exposes NFC capture/simulation operations.
type NFCService struct {
	nfc   *nfc.Repository
	cards *card.Repository
}

func NewNFCService(d *db.DB) *NFCService {
	return &NFCService{nfc: nfc.NewRepository(d.DB), cards: card.NewRepository(d.DB)}
}

func (s *NFCService) Capture(uid, atqa, sak, ndef string, cardID *int64) (nfc.Capture, error) {
	c := nfc.Capture{UID: uid, ATQA: atqa, SAK: sak, NDEF: ndef, CardID: cardID}
	id, err := s.nfc.Save(&c)
	if err != nil {
		return nfc.Capture{}, err
	}
	// link card nfc_uid if cardID provided
	if cardID != nil {
		card, err := s.cards.Get(*cardID)
		if err == nil {
			card.NFCUID = uid
			_ = s.cards.Update(card)
		}
	}
	_ = id
	return c, nil
}

func (s *NFCService) ListCaptures(cardID int64) ([]nfc.Capture, error) {
	return s.nfc.ListByCard(cardID)
}

func (s *NFCService) Simulate(cardID int64) (string, error) {
	return s.nfc.Simulate(cardID)
}
