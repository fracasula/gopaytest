package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopaytest/container/mock"
	"gopaytest/models/payment"
	repository "gopaytest/repositories/payments"
	"gopaytest/repositories/payments/paymentsfakes"
	"gopaytest/uuid"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestGetPayments(t *testing.T) {
	router, payments := getMockedPayments(t, 14)

	// TESTING PAGE 1
	req1, err := http.NewRequest("GET", "/v1/payments", nil)
	require.Nil(t, err)

	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)
	require.Equal(t, http.StatusOK, rr1.Code)

	var page1 paymentsPage
	require.Nil(t, json.Unmarshal(rr1.Body.Bytes(), &page1))
	require.Equal(t, paymentsPage{
		Data: payments[0:10],
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments",
		},
	}, page1)
	require.Equal(t, "2", rr1.Header().Get("X-Page-Count"))
	require.Equal(t, "1", rr1.Header().Get("X-Page-Number"))
	require.Equal(t, "10", rr1.Header().Get("X-Page-Size"))
	require.Equal(t, "14", rr1.Header().Get("X-Total-Count"))

	// TESTING PAGE 2
	req2, err := http.NewRequest("GET", "/v1/payments?$page=2", nil)
	require.Nil(t, err)

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	require.Equal(t, http.StatusOK, rr2.Code)

	var page2 paymentsPage
	require.Nil(t, json.Unmarshal(rr2.Body.Bytes(), &page2))
	require.Equal(t, paymentsPage{
		Data: payments[10:],
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments?$page=2",
		},
	}, page2)
	require.Equal(t, "2", rr2.Header().Get("X-Page-Count"))
	require.Equal(t, "2", rr2.Header().Get("X-Page-Number"))
	require.Equal(t, "10", rr2.Header().Get("X-Page-Size"))
	require.Equal(t, "14", rr2.Header().Get("X-Total-Count"))
}

func TestGetPayments_PageSize(t *testing.T) {
	router, payments := getMockedPayments(t, 14)

	req, err := http.NewRequest("GET", "/v1/payments?$size=14", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var page paymentsPage
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, paymentsPage{
		Data: payments,
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments?$size=14",
		},
	}, page)
	require.Equal(t, "1", rr.Header().Get("X-Page-Count"))
	require.Equal(t, "1", rr.Header().Get("X-Page-Number"))
	require.Equal(t, "14", rr.Header().Get("X-Page-Size"))
	require.Equal(t, "14", rr.Header().Get("X-Total-Count"))
}

func TestGetPayments_BadRequest(t *testing.T) {
	router, _ := getMockedPayments(t, 0)

	req1, err := http.NewRequest("GET", "/v1/payments?$size=INVALID", nil)
	require.Nil(t, err)

	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)
	require.Equal(t, http.StatusBadRequest, rr1.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr1.Body.Bytes(), &page))
	require.Equal(t, Error{
		Code:         400,
		Description:  "Invalid page size",
		ReasonPhrase: "Bad Request",
	}, page)

	req2, err := http.NewRequest("GET", "/v1/payments?$page=INVALID", nil)
	require.Nil(t, err)

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	require.Equal(t, http.StatusBadRequest, rr2.Code)

	require.Nil(t, json.Unmarshal(rr2.Body.Bytes(), &page))
	require.Equal(t, Error{
		Code:         400,
		Description:  "Invalid page number",
		ReasonPhrase: "Bad Request",
	}, page)
}

func TestGetPayments_InternalServerError(t *testing.T) {
	repo := &paymentsfakes.FakeRepository{}
	repo.CountPaymentsReturns(0, fmt.Errorf("a database error"))

	router := getMockedRouter(repo)

	// Testing database error on CountPayments
	req1, err := http.NewRequest("GET", "/v1/payments", nil)
	require.Nil(t, err)

	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)
	require.Equal(t, http.StatusInternalServerError, rr1.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr1.Body.Bytes(), &page))
	require.Equal(t, Error{
		Code:         500,
		Description:  "Could not get payments count",
		ReasonPhrase: "Internal Server Error",
	}, page)

	require.Equal(t, 1, repo.CountPaymentsCallCount())
	require.Equal(t, 0, repo.GetPaymentsCallCount())

	// Testing database error on GetPayments
	repo.CountPaymentsReturns(10, nil)
	repo.GetPaymentsReturns(nil, fmt.Errorf("a database error"))

	req2, err := http.NewRequest("GET", "/v1/payments", nil)
	require.Nil(t, err)

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	require.Equal(t, http.StatusInternalServerError, rr2.Code)

	require.Nil(t, json.Unmarshal(rr2.Body.Bytes(), &page))
	require.Equal(t, Error{
		Code:         500,
		Description:  "Could not get payments list",
		ReasonPhrase: "Internal Server Error",
	}, page)

	require.Equal(t, 1, repo.GetPaymentsCallCount())
}

func TestGetPayment(t *testing.T) {
	p1 := payment.MockedPayment(uuid.NewV4String())

	repo := &paymentsfakes.FakeRepository{}
	repo.GetPaymentReturns(&p1, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("GET", "/v1/payments/"+p1.ID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var page paymentPage
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &paymentPage{
		Data: p1,
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments/" + p1.ID,
		},
	}, &page)
}

func TestGetPayment_BadRequest(t *testing.T) {
	p1 := payment.MockedPayment("INVALID UUID")

	repo := &paymentsfakes.FakeRepository{}
	repo.GetPaymentReturns(&p1, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("GET", "/v1/payments/"+p1.ID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Invalid payment ID supplied, must be a valid UUID",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestGetPayment_InternalServerError(t *testing.T) {
	paymentID := uuid.NewV4String()
	repo := &paymentsfakes.FakeRepository{}
	repo.GetPaymentReturns(nil, fmt.Errorf("a database error"))

	router := getMockedRouter(repo)

	req, err := http.NewRequest("GET", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusInternalServerError,
		Description:  fmt.Sprintf("Could not get payment with ID %q", paymentID),
		ReasonPhrase: "Internal Server Error",
	}, &page)
}

func TestGetPayment_NotFound(t *testing.T) {
	paymentID := uuid.NewV4String()
	repo := &paymentsfakes.FakeRepository{}
	repo.GetPaymentReturns(nil, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("GET", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusNotFound,
		Description:  fmt.Sprintf("Could not find payment with ID %q", paymentID),
		ReasonPhrase: "Not Found",
	}, &page)
}

func TestCreatePayment(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)

	repo := repository.NewInMemoryPaymentsRepository()
	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	matches := regexp.MustCompile("payments/(.*)").FindStringSubmatch(rr.Header().Get("Location"))
	require.Len(t, matches, 2)
	p.ID = matches[1]
	p.Version = 1

	repoPayment, err := repo.GetPayment(p.ID)
	require.Nil(t, err)
	require.Equal(t, &p, repoPayment)

	var page paymentPage
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, paymentPage{
		Data: p,
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments/" + p.ID,
		},
	}, page)
}

func TestCreatePayment_BadRequest1(t *testing.T) {
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(nil))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Request body is not a valid Payment",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestCreatePayment_BadRequest2(t *testing.T) {
	p := payment.MockedPayment("")
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
	require.Equal(t, "", rr.Body.String())
}

func TestCreatePayment_BadRequest3(t *testing.T) {
	p := payment.MockedPayment(uuid.NewV4String())
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Payment ID is autogenerated, please leave it empty",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestCreatePayment_BadRequest4(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 1
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Payment version is autogenerated, please leave it empty",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestCreatePayment_InternalServerError(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)

	repo := &paymentsfakes.FakeRepository{}
	repo.CreatePaymentReturns(fmt.Errorf("a database error"))

	router := getMockedRouter(repo)

	req, err := http.NewRequest("POST", "/v1/payments", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusInternalServerError,
		Description:  "Could not create new payment",
		ReasonPhrase: "Internal Server Error",
	}, &page)
}

func TestDeletePayment(t *testing.T) {
	paymentID := uuid.NewV4String()
	repo := &paymentsfakes.FakeRepository{}
	repo.DeletePaymentReturns(true, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("DELETE", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Equal(t, "", rr.Body.String())
}

func TestDeletePayment_BadRequest(t *testing.T) {
	paymentID := "INVALID UUID"
	repo := &paymentsfakes.FakeRepository{}
	repo.DeletePaymentReturns(true, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("DELETE", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Invalid payment ID supplied, must be a valid UUID",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestDeletePayment_InternalServerError(t *testing.T) {
	paymentID := uuid.NewV4String()
	repo := &paymentsfakes.FakeRepository{}
	repo.DeletePaymentReturns(false, fmt.Errorf("a database error"))

	router := getMockedRouter(repo)

	req, err := http.NewRequest("DELETE", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusInternalServerError,
		Description:  fmt.Sprintf("Could not delete payment with ID %q", paymentID),
		ReasonPhrase: "Internal Server Error",
	}, &page)
}

func TestDeletePayment_NotFound(t *testing.T) {
	paymentID := uuid.NewV4String()
	repo := &paymentsfakes.FakeRepository{}
	repo.DeletePaymentReturns(false, nil)

	router := getMockedRouter(repo)

	req, err := http.NewRequest("DELETE", "/v1/payments/"+paymentID, nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusNotFound,
		Description:  fmt.Sprintf("Could not find payment with ID %q", paymentID),
		ReasonPhrase: "Not Found",
	}, &page)
}

func TestUpdatePayment(t *testing.T) {
	repo := repository.NewInMemoryPaymentsRepository()
	router := getMockedRouter(repo)

	paymentID := uuid.NewV4String()
	p := payment.MockedPayment(paymentID)
	p.Version = 1
	require.Nil(t, repo.CreatePayment(p))

	u := p
	u.ID = ""
	u.Version = 1
	u.OrganisationID = "a fancy company"
	u.Attributes.Amount = "500.00"
	u.Attributes.Currency = "BTC"
	body, err := json.Marshal(u)
	require.Nil(t, err)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	updatedPayment, err := repo.GetPayment(p.ID)
	require.Nil(t, err)
	require.Equal(t, 2, updatedPayment.Version)
	require.Equal(t, paymentID, updatedPayment.ID)
	require.Equal(t, "a fancy company", updatedPayment.OrganisationID)
	require.Equal(t, "500.00", updatedPayment.Attributes.Amount)
	require.Equal(t, "BTC", updatedPayment.Attributes.Currency)

	var page paymentPage
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, paymentPage{
		Data: *updatedPayment,
		Links: pageLinks{
			Self: "https://api.test.gopaytest.tech/v1/payments/" + paymentID,
		},
	}, page)
}

func TestUpdatePayment_BadRequest1(t *testing.T) {
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)
	paymentID := uuid.NewV4String()

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(nil))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Request body is not a valid Payment",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestUpdatePayment_BadRequest2(t *testing.T) {
	p := payment.MockedPayment(uuid.NewV4String())
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+p.ID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
	require.Equal(t, "", rr.Body.String())
}

func TestUpdatePayment_BadRequest3(t *testing.T) {
	// the paymentID below should be empty to make clear that the user cannot change the ID of an existing resource
	// here we're testing this behaviour
	p := payment.MockedPayment(uuid.NewV4String())
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+p.ID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Payment ID is autogenerated, please leave it empty",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestUpdatePayment_BadRequest4(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 1
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/INVALID-UUID", bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusBadRequest,
		Description:  "Invalid payment ID supplied, must be a valid UUID",
		ReasonPhrase: "Bad Request",
	}, &page)
}

func TestUpdatePayment_InternalServerError1(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	fakeError := &paymentsfakes.FakeUpdateError{}
	fakeError.ErrorReturns("a database error")
	fakeError.IsConflictReturns(false)
	repo.UpdatePaymentReturns(false, fakeError)
	paymentID := uuid.NewV4String()
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusInternalServerError,
		Description:  fmt.Sprintf("Could not update payment with ID %q", paymentID),
		ReasonPhrase: "Internal Server Error",
	}, &page)

	require.Equal(t, 1, repo.UpdatePaymentCallCount())
	require.Equal(t, 0, repo.GetPaymentCallCount())
}

func TestUpdatePayment_InternalServerError2(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	repo.UpdatePaymentReturns(true, nil)
	repo.GetPaymentReturns(nil, fmt.Errorf("a database error"))
	paymentID := uuid.NewV4String()
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusInternalServerError,
		Description:  fmt.Sprintf("Could not get payment with ID %q", paymentID),
		ReasonPhrase: "Internal Server Error",
	}, &page)

	require.Equal(t, 1, repo.UpdatePaymentCallCount())
	require.Equal(t, 1, repo.GetPaymentCallCount())
}

func TestUpdatePayment_NotFound1(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	repo.UpdatePaymentReturns(false, nil)
	paymentID := uuid.NewV4String()
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusNotFound,
		Description:  fmt.Sprintf("Could not find payment with ID %q", paymentID),
		ReasonPhrase: "Not Found",
	}, &page)

	require.Equal(t, 1, repo.UpdatePaymentCallCount())
	require.Equal(t, 0, repo.GetPaymentCallCount())
}

func TestUpdatePayment_NotFound2(t *testing.T) {
	p := payment.MockedPayment("")
	p.Version = 0
	body, err := json.Marshal(p)
	require.Nil(t, err)
	repo := &paymentsfakes.FakeRepository{}
	repo.UpdatePaymentReturns(true, nil)
	repo.GetPaymentReturns(nil, nil)
	paymentID := uuid.NewV4String()
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusNotFound,
		Description:  fmt.Sprintf("Could not find payment with ID %q", paymentID),
		ReasonPhrase: "Not Found",
	}, &page)

	require.Equal(t, 1, repo.UpdatePaymentCallCount())
	require.Equal(t, 1, repo.GetPaymentCallCount())
}

func TestUpdatePayment_Conflict(t *testing.T) {
	paymentID := uuid.NewV4String()
	p := payment.MockedPayment("")
	p.Version = 1
	body, err := json.Marshal(p)
	require.Nil(t, err)

	fakeError := &paymentsfakes.FakeUpdateError{}
	fakeError.ErrorReturns("a conflict error")
	fakeError.IsConflictReturns(true)

	repo := &paymentsfakes.FakeRepository{}
	repo.UpdatePaymentReturns(false, fakeError)
	router := getMockedRouter(repo)

	req, err := http.NewRequest("PUT", "/v1/payments/"+paymentID, bytes.NewBuffer(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusConflict, rr.Code)

	var page Error
	require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &page))
	require.Equal(t, &Error{
		Code:         http.StatusConflict,
		Description:  fmt.Sprintf("Could not update payment with ID %q", paymentID),
		ReasonPhrase: "Conflict",
	}, &page)

	require.Equal(t, 1, repo.UpdatePaymentCallCount())
	require.Equal(t, 0, repo.GetPaymentCallCount())
}

func getMockedRouter(repo repository.Repository) *chi.Mux {
	fakeContainer := mock.NewMockedContainer()
	fakeContainer.BaseURLReturns("https://api.test.gopaytest.tech")
	fakeContainer.PaymentsRepositoryReturns(repo)

	return NewRouter(fakeContainer)
}

func getMockedPayments(t *testing.T, noOfPayments int) (*chi.Mux, []payment.Payment) {
	repo := repository.NewInMemoryPaymentsRepository()
	payments := make([]payment.Payment, noOfPayments)
	for i := 0; i < noOfPayments; i++ {
		paymentID := fmt.Sprintf("payment %03d", i)
		payments[i] = payment.MockedPayment(paymentID)
		require.Nil(t, repo.CreatePayment(payments[i]))
	}

	return getMockedRouter(repo), payments
}
