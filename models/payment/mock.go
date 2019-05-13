package payment

func MockedPayment(paymentID string) Payment {
	return Payment{
		Type:           "Payment",
		ID:             paymentID,
		Version:        1,
		OrganisationID: "an organisation id",
		Attributes: Attributes{
			Amount: "amount",
			BeneficiaryParty: BeneficiaryParty{
				AccountName:       "account name",
				AccountNumber:     "account number",
				AccountNumberCode: "account number code",
				AccountType:       0,
				Address:           "address",
				BankID:            "bank id",
				BankIDCode:        "bank id code",
				Name:              "name",
			},
			ChargesInformation: ChargesInformation{
				BearerCode: "bearer code",
				SenderCharges: []SenderCharge{{
					Amount:   "amount 1",
					Currency: "currency 1",
				}, {
					Amount:   "amount 2",
					Currency: "currency 2",
				}},
				ReceiverChargesAmount:   "receiver charges amount",
				ReceiverChargesCurrency: "receiver charges currency",
			},
			Currency: "currency",
			DebtorParty: DebtorParty{
				AccountName:       "account name",
				AccountNumber:     "account number",
				AccountNumberCode: "account number code",
				Address:           "address",
				BankID:            "bank id",
				BankIDCode:        "bank id code",
				Name:              "name",
			},
			EndToEndReference: "end to end reference",
			FX: FX{
				ContractReference: "contract reference",
				ExchangeRate:      "exchange rage",
				OriginalAmount:    "original amount",
				OriginalCurrency:  "original currency",
			},
			NumericReference:     "numeric reference",
			PaymentID:            "payment id",
			PaymentPurpose:       "payment purpose",
			PaymentScheme:        "payment scheme",
			PaymentType:          "payment type",
			ProcessingDate:       "processing date",
			Reference:            "reference",
			SchemePaymentSubType: "scheme payment sub type",
			SchemePaymentType:    "scheme payment type",
			SponsorParty: SponsorParty{
				AccountNumber: "account number",
				BankID:        "bank id",
				BankIDCode:    "bank id code",
			},
		},
	}
}
