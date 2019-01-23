package payment

import (
	"math/big"

	"github.com/gofrs/uuid"
)

// Pymt contains all the information which a single payment has.
type Pymt struct {
	PymtUpsert
	ID      uuid.UUID `json:"id"`
	Version uint32    `json:"version"`
}

// PymtUpsert contains the information required to create or update a payment.
type PymtUpsert struct {
	Type           string `json:"type"`
	OrganisationID string `json:"organisation_id"`
	Attributes     struct {
		Amount               *big.Float `json:"amount"`
		Currency             string     `json:"currency"`
		Reference            string     `json:"reference"`
		EndToEndReference    string     `json:"end_to_end_reference"`
		NumericReference     string     `json:"numeric_reference"`
		PaymentID            string     `json:"payment_id"`
		PaymentPurpose       string     `json:"payment_purpose"`
		PaymentScheme        string     `json:"payment_scheme"`
		PaymentType          string     `json:"payment_type"`
		ProcessingDate       string     `json:"processing_date"`
		SchemePaymentSubType string     `json:"scheme_payment_sub_type"`
		SchemePaymentType    string     `json:"scheme_payment_type"`
		BeneficiaryParty     Party      `json:"beneficiary_party"`
		DebtorParty          Party      `json:"debtor_party"`
		SponsorParty         Party      `json:"sponsor_party"`
		ChargesInformation   struct {
			BearerCode    string `json:"bearer_code"`
			SenderCharges []struct {
				Amount   *big.Float `json:"amount"`
				Currency string     `json:"currency"`
			} `json:"sender_charges"`
			ReceiverChargesAmount   *big.Float `json:"receiver_charges_amount"`
			ReceiverChargesCurrency string     `json:"receiver_charges_currency"`
		} `json:"charges_information"`
		Fx struct {
			ContractReference string `json:"contract_reference"`
			ExchangeRate      string `json:"exchange_rate"`
			OriginalAmount    string `json:"original_amount"`
			OriginalCurrency  string `json:"original_currency"`
		} `json:"fx"`
	} `json:"attributes"`
}

// Party contains the information of each single party involved in a payment.
type Party struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}
