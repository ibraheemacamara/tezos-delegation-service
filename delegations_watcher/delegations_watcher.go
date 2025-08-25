package delegationswatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	events "github.com/dipdup-net/go-lib/tzkt/events"
	"github.com/ibraheemacara/tezos-delegation-service/config"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	"github.com/ibraheemacara/tezos-delegation-service/httpclient"
	"github.com/ibraheemacara/tezos-delegation-service/types"
	log "github.com/sirupsen/logrus"
)

type DelegationsWatcher struct {
	config     config.Config
	httpClient httpclient.HttpInterface
	db         db.DBInterface
	tzktClient TzktClient
}

type TzktClient interface {
	Connect(ctx context.Context) error
	SubscribeToHead() error
	Listen() <-chan events.Message
}

func NewDelegationsWatcher(config config.Config, httpClient httpclient.HttpInterface, db db.DBInterface) *DelegationsWatcher {
	wsUrl := fmt.Sprintf("%s/v1/ws", config.Tzkt.Url)
	return &DelegationsWatcher{
		config:     config,
		httpClient: httpClient,
		db:         db,
		tzktClient: events.NewTzKT(wsUrl),
	}
}

func (dw *DelegationsWatcher) Start() {
	log.Info("Delegations watcher started")

	//first we get the last block recorded in the database
	lastBlock, err := dw.db.GetLastBlock()
	if err != nil {
		log.Error("Failed to get last block from database: ", err)
		return
	}

	var allDelegations []types.TzktDelegationsResponse
	if lastBlock == 0 {
		log.Info("No blocks recorded in the database, query all delegations from tzkt ...")

		allDelegations, err = getDelegations(dw.config.Tzkt.Url, 0, dw.httpClient)
		if err != nil {
			log.Error("Failed to get delegations from tzkt: ", err)
			return
		}

	} else {
		log.Infof("Last block recorded in the database: %v, getting delegations from last block to current state", lastBlock)
		//get delegations from last block to current state
		allDelegations, err = getDelegationsFromLevel(dw.config.Tzkt.Url, lastBlock, dw.httpClient)
		if err != nil {
			log.Error("Failed to get delegations from tzkt: ", err)
			return
		}
	}

	//all past delegations are retrived from tzkt, start watching for new blocks
	go dw.WatchNewBlocks(context.Background())

	//insert all delegations into database
	log.Infof("Number of delegations: %v, inserting into database", len(allDelegations))
	err = bulkInsertDelegations(dw.db, allDelegations)
	if err != nil {
		log.Error("Failed to insert delegations into database: ", err)
		return
	}

	log.Infof("All %v delegations inserted into database", len(allDelegations))

	// go dw.WatchNewBlocks()
}

func (dw *DelegationsWatcher) WatchNewBlocks(ctx context.Context) {
	log.Info("Start watching for new blocks...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := dw.tzktClient.Connect(ctx); err != nil {
			log.Errorf("Failed to connect to tzkt: %v, retrying in 5 seconds", err)
			select {
			case <-time.After(5 * time.Second):
				continue
			case <-ctx.Done():
				return
			}
		}
		log.Info("Connected to tzkt events hub")

		//subscribe to head events
		if err := dw.tzktClient.SubscribeToHead(); err != nil {
			log.Errorf("Failed to subscribe to head events: %v", err)
		}

		//process received messages
		for msg := range dw.tzktClient.Listen() {
			if msg.Channel == events.ChannelHead {
				log.Info("Received head event")
				raw, err := json.Marshal(msg.Body)
				if err != nil {
					log.Errorf("Failed to marshal head event: %v", err)
					continue
				}

				var head map[string]any
				err = json.Unmarshal(raw, &head)
				if err != nil {
					log.Errorf("Failed to unmarshal head event: %v", err)
					continue
				}

				if head["level"] == nil {
					continue
				}
				level := head["level"].(float64)

				log.Infof("New block received: %v, getting delegations", level)

				delegationsResponse, err := getDelegations(dw.config.Tzkt.Url, int32(level), dw.httpClient)
				if err != nil {
					log.Errorf("Failed to get delegations from tzkt: %v", err)
					continue
				}
				if len(delegationsResponse) == 0 {
					log.Infof("No delegations found for block: %v", level)
					continue
				}
				log.Infof("Number of delegations: %v, inserting into database; this may take a while", len(delegationsResponse))
				err = bulkInsertDelegations(dw.db, delegationsResponse)
				if err != nil {
					log.Errorf("Failed to insert delegations for block %v into database: %v", level, err)
					continue
				}

				log.Infof("All %v delegations inserted into database for block %v", len(delegationsResponse), level)

			}
		}

		// Reconnect logic
		select {
		case <-ctx.Done():
			return
		default:
			log.Println("disconnected, retrying in 5s…")
			time.Sleep(5 * time.Second)
		}

	}
}

func getDelegations(tzktUrl string, level int32, httpClient httpclient.HttpInterface) ([]types.TzktDelegationsResponse, error) {
	limit := 10000
	offset := 0
	allDelegations := []types.TzktDelegationsResponse{}
	for {
		var url string
		if level == 0 {
			url = fmt.Sprintf("%s/v1/operations/delegations?limit=%d&offset=%d", tzktUrl, limit, offset)
		} else {
			url = fmt.Sprintf("%s/v1/operations/delegations?limit=%d&offset=%d&level=%d", tzktUrl, limit, offset, level)
		}

		data, err := httpClient.Get(url)
		if err != nil {
			return nil, err
		}

		var delegations []types.TzktDelegationsResponse
		err = json.Unmarshal(data, &delegations)
		if err != nil {
			return nil, err
		}

		if len(delegations) == 0 {
			break
		}
		allDelegations = append(allDelegations, delegations...)

		offset += limit

	}

	return allDelegations, nil
}

func getDelegationsFromLevel(tzktUrl string, fromLevel int32, httpClient httpclient.HttpInterface) ([]types.TzktDelegationsResponse, error) {
	limit := 10000
	offset := 0
	allDelegations := []types.TzktDelegationsResponse{}
	for {
		url := fmt.Sprintf("%s/v1/operations/delegations?limit=%d&offset=%d&level.gt=%d", tzktUrl, limit, offset, fromLevel)
		data, err := httpClient.Get(url)
		if err != nil {
			return nil, err
		}

		var delegations []types.TzktDelegationsResponse
		err = json.Unmarshal(data, &delegations)
		if err != nil {
			return nil, err
		}

		if len(delegations) == 0 {
			break
		}
		allDelegations = append(allDelegations, delegations...)

		offset += limit

	}

	return allDelegations, nil
}

func bulkInsertDelegations(dbInterface db.DBInterface, delegationsResponse []types.TzktDelegationsResponse) error {
	delegations := make([]db.Delegations, len(delegationsResponse))
	for i, delegation := range delegationsResponse {
		delegations[i] = db.Delegations{
			Delegator: delegation.Sender.Address,
			Timestamp: delegation.Timestamp,
			Block:     delegation.Level,
			Amount:    delegation.Amount,
		}
	}
	err := dbInterface.BulkInsertDelegations(delegations)
	if err != nil {
		return err
	}
	return nil
}
