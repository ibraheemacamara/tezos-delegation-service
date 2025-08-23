package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibraheemacara/tezos-delegation-service/db"
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

	expectedDelegations := []db.Delegations{
		{
			Delegator: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
			Timestamp: time.Now().UTC(),
			Block:     109,
			Amount:    25079312620,
		},
	}
	mockDB := &MockDB{
		DelegationsToGet: expectedDelegations,
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

	var actualDelegations []db.Delegations
	err := json.Unmarshal(w.Body.Bytes(), &actualDelegations)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(actualDelegations) != len(expectedDelegations) {
		t.Errorf("expected %d delegations, got %d", len(expectedDelegations), len(actualDelegations))
	}
	for i := range actualDelegations {
		fmt.Println(actualDelegations[i] == expectedDelegations[i])

		if actualDelegations[i] != expectedDelegations[i] {
			t.Errorf("expected delegation %v, got %v", expectedDelegations[i], actualDelegations[i])
		}
	}
}

// func TestGetDelegationsByYear(t *testing.T) {
// 	var year int = 2018
// 	expectedDelegations := []db.Delegations{
// 		{
// 			Delegator: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
// 			Timestamp: time.Date(2018, time.June, 30, 19, 30, 27, 0, time.UTC),
// 			Block:     109,
// 			Amount:    25079312620,
// 		},
// 		{
// 			Delegator: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZece",
// 			Timestamp: time.Date(2018, time.July, 30, 19, 30, 27, 0, time.UTC),
// 			Block:     109,
// 			Amount:    25079312620,
// 		},
// 	}
// 	mockDB := &MockDB{
// 		DelegationsToGet: expectedDelegations,
// 	}

// 	controller := NewController(mockDB)

// 	r := gin.New()
// 	r.GET("/delegations/:year", controller.GetDelegationsByYear)

// 	path := fmt.Sprintf("/delegations/%d", year)

// 	req := httptest.NewRequest("GET", path, nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	if w.Code != 200 {
// 		t.Fatalf("expected status 200, got %d", w.Code)
// 	}

// 	var actualDelegations []db.Delegations
// 	err := json.Unmarshal(w.Body.Bytes(), &actualDelegations)
// 	if err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if len(actualDelegations) != len(expectedDelegations) {
// 		t.Errorf("expected %d delegations, got %d", len(expectedDelegations), len(actualDelegations))
// 	}
// 	for i := range actualDelegations {
// 		actualYear := actualDelegations[i].Timestamp.Year()
// 		if actualYear != year {
// 			t.Errorf("expected year %d, got %d", year, actualYear)
// 		}
// 		if actualDelegations[i] != expectedDelegations[i] {
// 			t.Errorf("expected delegation %v, got %v", expectedDelegations[i], actualDelegations[i])
// 		}
// 	}
// }

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

// errors responses
// func TestGetDelegationsByYearError(t *testing.T) {
// 	mockDB := &MockDBError{
// 		GetDelegationsByYearErr: errors.New("test error"),
// 	}
// 	controller := NewController(mockDB)

// 	r := gin.New()
// 	r.GET("/delegations/:year", controller.GetDelegationsByYear)

// 	path := fmt.Sprintf("/delegations/%d", 2018)

// 	req := httptest.NewRequest("GET", path, nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	if w.Code != 500 {
// 		t.Fatalf("expected status 500, got %d", w.Code)
// 	}
// }
