package payment

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayment(t *testing.T) {
	var payment Payment
	err := json.Unmarshal([]byte(jsonData), &payment)
	require.Nil(t, err, "Could not unmarshal JSON data")

	require.Equal(t, Payment{
		Type:           "Payment",
		ID:             "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
		Version:        0,
		OrganisationID: "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		Attributes: Attributes{
			Amount: "100.21",
			BeneficiaryParty: BeneficiaryParty{
				AccountName:       "W Owens",
				AccountNumber:     "31926819",
				AccountNumberCode: "BBAN",
				AccountType:       0,
				Address:           "1 The Beneficiary Localtown SE2",
				BankID:            "403000",
				BankIDCode:        "GBDSC",
				Name:              "Wilfred Jeremiah Owens",
			},
			ChargesInformation: ChargesInformation{
				BearerCode: "SHAR",
				SenderCharges: []SenderCharge{
					{
						Amount:   "5.00",
						Currency: "GBP",
					},
					{
						Amount:   "10.00",
						Currency: "USD",
					},
				},
				ReceiverChargesAmount:   "1.00",
				ReceiverChargesCurrency: "USD",
			},
			Currency: "GBP",
			DebtorParty: DebtorParty{
				AccountName:       "EJ Brown Black",
				AccountNumber:     "GB29XABC10161234567801",
				AccountNumberCode: "IBAN",
				Address:           "10 Debtor Crescent Sourcetown NE1",
				BankID:            "203301",
				BankIDCode:        "GBDSC",
				Name:              "Emelia Jane Brown",
			},
			EndToEndReference: "Wil piano Jan",
			FX: FX{
				ContractReference: "FX123",
				ExchangeRate:      "2.00000",
				OriginalAmount:    "200.42",
				OriginalCurrency:  "USD",
			},
			NumericReference:     "1002001",
			PaymentID:            "123456789012345678",
			PaymentPurpose:       "Paying for goods/services",
			PaymentScheme:        "FPS",
			PaymentType:          "Credit",
			ProcessingDate:       "2017-01-18",
			Reference:            "Payment for Em's piano lessons",
			SchemePaymentSubType: "InternetBanking",
			SchemePaymentType:    "ImmediatePayment",
			SponsorParty: SponsorParty{
				AccountNumber: "56781234",
				BankID:        "123123",
				BankIDCode:    "GBDSC",
			},
		},
	}, payment)
}

var jsonData = `{
	"type": "Payment",
	"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
	"version": 0,
	"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	"attributes": {
		"amount": "100.21",
		"beneficiary_party": {
			"account_name": "W Owens",
			"account_number": "31926819",
			"account_number_code": "BBAN",
			"account_type": 0,
			"address": "1 The Beneficiary Localtown SE2",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Wilfred Jeremiah Owens"
		},
		"charges_information": {
			"bearer_code": "SHAR",
			"sender_charges": [
				{
					"amount": "5.00",
					"currency": "GBP"
				},
				{
					"amount": "10.00",
					"currency": "USD"
				}
			],
			"receiver_charges_amount": "1.00",
			"receiver_charges_currency": "USD"
		},
		"currency": "GBP",
		"debtor_party": {
			"account_name": "EJ Brown Black",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "10 Debtor Crescent Sourcetown NE1",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "Emelia Jane Brown"
		},
		"end_to_end_reference": "Wil piano Jan",
		"fx": {
			"contract_reference": "FX123",
			"exchange_rate": "2.00000",
			"original_amount": "200.42",
			"original_currency": "USD"
		},
		"numeric_reference": "1002001",
		"payment_id": "123456789012345678",
		"payment_purpose": "Paying for goods/services",
		"payment_scheme": "FPS",
		"payment_type": "Credit",
		"processing_date": "2017-01-18",
		"reference": "Payment for Em's piano lessons",
		"scheme_payment_sub_type": "InternetBanking",
		"scheme_payment_type": "ImmediatePayment",
		"sponsor_party": {
			"account_number": "56781234",
			"bank_id": "123123",
			"bank_id_code": "GBDSC"
		}
	}
}`
