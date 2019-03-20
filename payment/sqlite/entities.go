package sqlite

import (
	"encoding/json"

	"github.com/ifraixedes/go-payments-api-example/payment"
	"go.fraixed.es/errors"
)

// pymtData groups several payment fields which are stored as blob in the DB.
type pymtData struct {
	Type string `json:"type"`
	payment.Attrs
}

// Serialize returns blob to store pd in the DB.
func (pd *pymtData) Serialize() ([]byte, error) {
	b, err := json.Marshal(pd)
	if err != nil {
		return nil, errors.Wrap(err, payment.ErrUnexpectedSysError)
	}

	return b, nil
}

// Deserialize initializes pd from b. b is usually the bob stored in the DB.
func (pd *pymtData) Deserialize(b []byte) error {
	err := json.Unmarshal(b, pd)
	if err != nil {
		return errors.Wrap(err, ErrInvalidFormatBlob)
	}

	return nil
}

// Init initializes the pd from p.
func (pd *pymtData) Init(p payment.PymtUpsert) {
	pd.Type = p.Type
	pd.Attrs = p.Attributes
}

// Set sets p with the values hold by pd. If p is nill, it panics.
func (pd *pymtData) Set(p *payment.PymtUpsert) {
	p.Type = pd.Type
	p.Attributes = pd.Attrs
}
