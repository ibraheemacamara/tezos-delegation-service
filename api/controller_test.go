package api

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	"github.com/ibraheemacara/tezos-delegation-service/types"
)

type MockDB struct {
	DelegationsToGet []db.Delegations
}

func (m *MockDB) GetDelegations() ([]db.Delegations, error) {
	return m.DelegationsToGet, nil
}

func (m *MockDB) GetDelegationsByYear(year string) ([]db.Delegations, error) {
	return m.DelegationsToGet, nil
}

func (m *MockDB) InsertDelegations(delegator string, timestamp time.Time, block int32, amount int64) error {
	return nil
}

func (m *MockDB) GetLastBlock() (int32, error) {
	return 0, nil
}

func (m *MockDB) BulkInsertDelegations(delegations []db.Delegations) error {
	return nil
}

type MockDBError struct {
	GetDelegationsError      error
	GetDelegationsByYearErr  error
	InsertDelegationsErr     error
	GetLastBlockErr          error
	BulkInsertDelegationsErr error
}

func (m *MockDBError) GetDelegations() ([]db.Delegations, error) {
	return nil, m.GetDelegationsError
}

func (m *MockDBError) GetDelegationsByYear(year string) ([]db.Delegations, error) {
	return nil, m.GetDelegationsByYearErr
}

func (m *MockDBError) InsertDelegations(delegator string, timestamp time.Time, block int32, amount int64) error {
	return m.InsertDelegationsErr
}

func (m *MockDBError) GetLastBlock() (int32, error) {
	return 0, m.GetLastBlockErr
}

func (m *MockDBError) BulkInsertDelegations(delegations []db.Delegations) error {
	return m.BulkInsertDelegationsErr
}

func TestGetDelegations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	delegator := "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd"
	timestamp := time.Now().UTC()
	block := 109
	amount := 25079312620

	mockDB := &MockDB{
		DelegationsToGet: []db.Delegations{
			{
				Delegator: delegator,
				Timestamp: timestamp,
				Block:     int32(block),
				Amount:    int64(amount),
			},
		},
	}

	expectedDelegations := types.DelegationsResponse{
		Delegations: []types.Delegation{
			{
				Delegator: delegator,
				Timestamp: timestamp,
				Block:     int32(block),
				Amount:    int64(amount),
			},
		},
	}

	controller := NewController(mockDB)

	r := gin.New()
	r.GET("/delegations", controller.GetDelegations)

	req := httptest.NewRequest("GET", "/delegations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var actualDelegations types.DelegationsResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualDelegations)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(actualDelegations.Delegations) != len(expectedDelegations.Delegations) {
		t.Errorf("expected %d delegations, got %d", len(expectedDelegations.Delegations), len(actualDelegations.Delegations))
	}
	for i := range actualDelegations.Delegations {
		if actualDelegations.Delegations[i] != expectedDelegations.Delegations[i] {
			t.Errorf("expected delegation %v, got %v", expectedDelegations.Delegations[i], actualDelegations.Delegations[i])
		}
	}
}

func TestGetDelegationsError(t *testing.T) {
	mockDB := &MockDBError{
		GetDelegationsError: errors.New("test error"),
	}
	controller := NewController(mockDB)

	r := gin.New()
	r.GET("/delegations", controller.GetDelegations)

	req := httptest.NewRequest("GET", "/delegations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
