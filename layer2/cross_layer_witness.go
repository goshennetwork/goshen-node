package layer2

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/utils"
)

var MessageSentEventID = crypto.Keccak256Hash([]byte("MessageSent(uint64,address,address,bytes32,bytes)"))

type MessageSentEvent struct {
	MessageIndex uint64
	Target       web3.Address
	Sender       web3.Address
	Message      []byte
	MmrRoot      [32]byte
}

func convertHashes(hash []common.Hash) []web3.Hash {
	var topics []web3.Hash
	for _, h := range hash {
		topics = append(topics, web3.Hash(h))
	}

	return topics
}

func MustParseMessageSentEvent(log *types.Log) *MessageSentEvent {
	res, err := ParseMessageSentEvent(log)
	utils.Ensure(err)

	return res
}

func ParseMessageSentEvent(log *types.Log) (*MessageSentEvent, error) {
	web3Log := &web3.Log{
		Address: web3.Address(log.Address),
		Topics:  convertHashes(log.Topics),
		Data:    log.Data,
	}
	evt := abi.MustNewEvent("MessageSent(uint64 indexed _messageIndex,address indexed _target,address indexed _sender, bytes32 _mmrRoot, bytes _message)")

	utils.EnsureTrue(evt.ID() == web3.Hash(MessageSentEventID))

	args, err := evt.ParseLog(web3Log)
	if err != nil {
		return nil, err
	}
	var evtItem MessageSentEvent
	err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
	if err != nil {
		return nil, err
	}
	return &evtItem, nil
}
