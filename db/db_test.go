package db

import (
	"testing"
	"time"
)

// MockDB implements DBInterface for unit testing
type MockDB struct {
	InsertCalled     bool
	InsertedArgs     []interface{}
	DelegationsToGet []Delegations
	InsertErr        error
	GetErr           error
}

func (m *MockDB) InsertDelegations(delegator string, timestamp time.Time, block int32, amount int64) error {
	m.InsertCalled = true
	m.InsertedArgs = []interface{}{delegator, timestamp, block, amount}
	return m.InsertErr
}

func (m *MockDB) GetDelegations() ([]Delegations, error) {
	return m.DelegationsToGet, m.GetErr
}

func TestInsertAndGetDelegations_Mock(t *testing.T) {
	mock := &MockDB{
		DelegationsToGet: []Delegations{{Delegator: "tz1mock", Timestamp: time.Now(), Block: 42, Amount: 99}},
	}

	// Insert
	err := mock.InsertDelegations("tz1mock", time.Now(), 42, 99)
	if err != nil {
		t.Fatalf("InsertDelegations failed: %v", err)
	}
	if !mock.InsertCalled {
		t.Error("InsertDelegations was not called")
	}

	// Retrieve
	delegations, err := mock.GetDelegations()
	if err != nil {
		t.Fatalf("GetDelegations failed: %v", err)
	}
	if len(delegations) != 1 || delegations[0].Delegator != "tz1mock" {
		t.Errorf("Delegations not mocked as expected: %+v", delegations)
	}
}
