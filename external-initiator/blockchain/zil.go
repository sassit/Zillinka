package blockchain

import (
	"encoding/json"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/external-initiator/store"
	"github.com/smartcontractkit/external-initiator/subscriber"
)

const ZIL = "zilliqa"

type zilFilterQuery struct {
	ServiceName string
	Addresses   []string
}

// The zilManager implements the subscriber.JsonManager interface and allows
// for interacting with ZIL nodes over RPC or WS.
type zilManager struct {
	fq *zilFilterQuery
	p  subscriber.Type
}

/*
Example query from: https://dev.zilliqa.com/docs/dev/dev-tools-websockets/#subscribe-event-log

{
  "query":"EventLog",
  "addresses":[
    "0x0000000000000000000000000000000000000000",
    "0x1111111111111111111111111111111111111111"
  ]
}
*/
type ZilEventLogQueryRequest struct {
	Query     string   `json:"query"`
	Addresses []string `json:"addresses"`
}

/*
Example response from:
{
   "type": "Notification",
   "values": [
      {
         "query": "EventLog",
         "value": [
            {
               "address": "afccafdc1ce8249cec35a0b432e329ce1bfac179",
               "event_logs": [
                  {
                     "_eventname": "request",
                     "params": [
                        {
                           "type": "String",
                           "value": "TEST",
                           "vname": "oracleId"
                        },
                        {
                           "type": "Uint32",
                           "value": "0",
                           "vname": "requestId"
                        },
                        {
                           "type": "ByStr20",
                           "value": "0x1a8ba23182e4686fb8121a310111d03b55c91b46",
                           "vname": "initiator"
                        },
                        {
                           "type": "String",
                           "value": "kaub",
                           "vname": "argument"
                        }
                     ]
                  }
               ]
            }
         ]
      }
   ]
}
*/

type ZilEventLogQueryResponse struct {
	Type   string                  `json:"type"`
	Values []ZilEventLogQueryValue `json:"values"`
}

type ZilEventLogQueryValue struct {
	Query string `json:"query"`
	Value []struct {
		Address   string `json:"address"`
		EventLogs []struct {
			Eventname string `json:"_eventname"`
			Params    []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
				Vname string `json:"vname"`
			} `json:"params"`
		} `json:"event_logs"`
	} `json:"value"`
}

// createZilManager creates a new instance of zilManager with the provided
// connection type and store.ZilSubscription config.
func createZilManager(p subscriber.Type, config store.Subscription) zilManager {
	var addresses []string
	for _, a := range config.Zilliqa.Accounts {
		addresses = append(addresses, a)
	}
	return zilManager{
		fq: &zilFilterQuery{
			ServiceName: config.Zilliqa.ServiceName,
			Addresses:   addresses,
		},
		p: p,
	}
}

// GetTriggerJson generates a JSON payload to the ZIL node
// using the config in zilManager.
func (z zilManager) GetTriggerJson() []byte {
	logger.Debugw("Getting trigger json", "ExpectsMock", ExpectsMock)

	queryCall := ZilEventLogQueryRequest{
		Query:     "EventLog",
		Addresses: z.fq.Addresses,
	}

	logger.Debug("addresses: ", z.fq.Addresses)

	switch z.p {
	case subscriber.WS:
		//logger.Debug("Addresses from filter query: ", queryCall.Addresses)
		//if len(z.fq.Addresses) == 0 || z.fq.Addresses == nil {
		//	return nil
		//}
		bytes, err := json.Marshal(queryCall)
		if err != nil {
			return nil
		}
		logger.Debug("Payload", string(bytes))
		return bytes
	default:
		logger.Errorw(ErrSubscriberType.Error(), "type", z.p)
		return nil
	}
}

// GetTestJson generates a JSON payload to test
// the connection to the ZIL node.
//
// If zilManager is using WebSocket:
// Returns nil.
//
// If zilManager is using RPC:
// Sends a request to get the latest block number.
func (z zilManager) GetTestJson() []byte {
	logger.Debugw("Get test json", "ExpectsMock", ExpectsMock)
	return nil
}

// ParseTestResponse parses the response from the
// ZIL node after sending GetTestJson(), and returns
// the error from parsing, if any.
//
// If zilManager is using WebSocket:
// Returns nil.
func (z zilManager) ParseTestResponse(data []byte) error {
	logger.Debugw("Parsing test response", "ExpectsMock", ExpectsMock)
	return nil
}

// ParseResponse parses the response from the
// ZIL node, and returns a slice of subscriber.Events
// and if the parsing was successful.
func (z zilManager) ParseResponse(data []byte) ([]subscriber.Event, bool) {
	logger.Debugw("Parsing response", "ExpectsMock", ExpectsMock)

	var msg ZilEventLogQueryResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		logger.Error("failed parsing message: ", string(data))
		return nil, false
	}

	if msg.Type != "Notification" || len(msg.Values) == 0 {
		logger.Error("invalid message: ", msg)
		return nil, false
	}

	var events []subscriber.Event

	switch z.p {
	case subscriber.WS:
		for _, v := range msg.Values {
			event, err := json.Marshal(v)
			if err == nil {
				events = append(events, event)
			}
		}
	default:
		logger.Errorw(ErrSubscriberType.Error(), "type", z.p)
		return nil, false
	}

	return events, true
}
