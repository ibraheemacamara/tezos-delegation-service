package delegationswatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	events "github.com/dipdup-net/go-lib/tzkt/events"
	"github.com/ibraheemacara/tezos-delegation-service/config"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	"github.com/ibraheemacara/tezos-delegation-service/httpclient"
	"github.com/ibraheemacara/tezos-delegation-service/types"
)

type MockDB struct {
	db.DBInterface
	LastBlock        int32
	BulkInsertCalled bool
	BulkInsertDelegs []db.Delegations
}

func (m *MockDB) GetLastBlock() (int32, error) {
	return m.LastBlock, nil
}
func (m *MockDB) BulkInsertDelegations(delegations []db.Delegations) error {
	m.BulkInsertCalled = true
	m.BulkInsertDelegs = delegations
	return nil
}

type MockHTTPClient struct {
	httpclient.HttpInterface
	callCount int
}

func (m *MockHTTPClient) Get(url string) ([]byte, error) {
	m.callCount++
	if m.callCount == 1 {
		resp := []types.TzktDelegationsResponse{
			{
				Sender:    types.Address{Address: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd"},
				Timestamp: time.Now().UTC(),
				Level:     1,
				Amount:    1000,
			},
		}
		return json.Marshal(resp)
	}
	resp := []types.TzktDelegationsResponse{}
	return json.Marshal(resp)
}

type MockDBError struct {
	GetLastBlockErr    error
	BulkInsertDelegErr error
}

func (m *MockDBError) GetLastBlock() (int32, error) {
	return 0, m.GetLastBlockErr
}
func (m *MockDBError) BulkInsertDelegations([]db.Delegations) error {
	return m.BulkInsertDelegErr
}

// Unused methods
func (m *MockDBError) GetDelegations() ([]db.Delegations, error)               { return nil, nil }
func (m *MockDBError) GetDelegationsByYear(string) ([]db.Delegations, error)   { return nil, nil }
func (m *MockDBError) InsertDelegations(string, time.Time, int32, int64) error { return nil }

type MockTzkt struct {
	msgChan chan events.Message
}

func (m *MockTzkt) Connect(ctx context.Context) error { return nil }
func (m *MockTzkt) SubscribeToHead() error            { return nil }
func (m *MockTzkt) Listen() <-chan events.Message     { return m.msgChan }

func TestGetDelegations(t *testing.T) {
	url := "http://fake-tzkt"
	httpClient := &MockHTTPClient{}
	delegations, err := getDelegations(url, 0, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(delegations) != 1 {
		t.Errorf("expected 1 delegation, got %d", len(delegations))
	}
	if delegations[0].Sender.Address != "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd" {
		t.Errorf("expected delegation from tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd, got %s", delegations[0].Sender.Address)
	}
	if delegations[0].Level != 1 {
		t.Errorf("expected level 1, got %d", delegations[0].Level)
	}
	if delegations[0].Amount != 1000 {
		t.Errorf("expected amount 1000, got %d", delegations[0].Amount)
	}
}

func TestGetDelegationsFromLevel(t *testing.T) {
	url := "http://fake-tzkt"
	httpClient := &MockHTTPClient{}
	delegations, err := getDelegationsFromLevel(url, 0, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(delegations) != 1 {
		t.Errorf("expected 1 delegation, got %d", len(delegations))
	}
}

// we set offset to 1 so first call returns empty
func TestGetDelegations_Empty(t *testing.T) {
	httpClient := &MockHTTPClient{callCount: 1}
	delegations, err := getDelegations("http://fake-tzkt", 0, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	if len(delegations) != 0 {
		t.Errorf("expected 0 delegations, got %d", len(delegations))
	}
}

func TestBulkInsertDelegations(t *testing.T) {
	err := bulkInsertDelegations(&MockDB{}, []types.TzktDelegationsResponse{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBulkInsertDelegationsError(t *testing.T) {
	err := bulkInsertDelegations(&MockDBError{BulkInsertDelegErr: fmt.Errorf("bulk insert error")}, []types.TzktDelegationsResponse{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestWatchBlocksDelegationsInserted(t *testing.T) {
	msgChan := make(chan events.Message, 1)
	mockTzkt := &MockTzkt{msgChan: msgChan}
	mockDB := &MockDB{}
	cfg := config.Config{}
	cfg.Tzkt.Url = "http://fake-tzkt"

	watcher := &DelegationsWatcher{
		config:     cfg,
		httpClient: &MockHTTPClient{},
		db:         mockDB,
		tzktClient: mockTzkt,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		watcher.WatchNewBlocks(ctx)
		close(done)
	}()

	msgChan <- events.Message{
		Channel: events.ChannelHead,
		Body:    map[string]interface{}{"level": float64(1)},
	}
	close(msgChan)
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	fmt.Println("BulkInsertCalled:", mockDB.BulkInsertCalled)
	if !mockDB.BulkInsertCalled {
		t.Errorf("expected BulkInsertDelegations to be called")
	}
}

func TestWatchBlocksNoDelegations(t *testing.T) {
	msgChan := make(chan events.Message, 1)
	mockTzkt := &MockTzkt{msgChan: msgChan}
	mockDB := &MockDB{}
	cfg := config.Config{}
	cfg.Tzkt.Url = "http://fake-tzkt"

	httpClient := &MockHTTPClient{callCount: 1}

	watcher := &DelegationsWatcher{
		config:     cfg,
		httpClient: httpClient,
		db:         mockDB,
		tzktClient: mockTzkt,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan <- events.Message{
		Channel: events.ChannelHead,
		Body:    map[string]interface{}{"level": float64(1)},
	}
	close(msgChan)
	time.Sleep(10 * time.Millisecond)
	cancel()

	watcher.WatchNewBlocks(ctx)

	if mockDB.BulkInsertCalled {
		t.Errorf("expected BulkInsertDelegations NOT to be called")
	}
}

func TestWatchBlocksDBInsertError(t *testing.T) {
	msgChan := make(chan events.Message, 1)
	mockTzkt := &MockTzkt{msgChan: msgChan}
	mockDB := &MockDBError{BulkInsertDelegErr: fmt.Errorf("fail")}
	cfg := config.Config{}
	cfg.Tzkt.Url = "http://fake-tzkt"

	watcher := &DelegationsWatcher{
		config:     cfg,
		httpClient: &MockHTTPClient{},
		db:         mockDB,
		tzktClient: mockTzkt,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan <- events.Message{
		Channel: events.ChannelHead,
		Body:    map[string]interface{}{"level": float64(1)},
	}
	close(msgChan)
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Should not panic even if DB insert fails
	watcher.WatchNewBlocks(ctx)
}
