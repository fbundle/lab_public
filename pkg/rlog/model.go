package rlog

import (
	"github.com/khanh-nguyen-code/go_util/pkg/rlog/rpc"
	"time"
)

type NodeId uint8 // at most 256 nodes
type CountId uint64

type ProposalId uint64 // CountId<<8 + NodeId // at most 2^(64-8) proposals each node

func MakeProposalId(nodeId NodeId, countId CountId) ProposalId {
	return ProposalId(uint64(countId)<<8 + uint64(nodeId))
}

func (pid ProposalId) CountId() CountId {
	return CountId(pid >> 8)
}

type LogId uint64

type Command = string

type Object interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Commit(...Command)
}

type Value struct {
	Object    Object    `json:"object"`
	StartId   LogId     `json:"start_id"`
	DecidedId LogId     `json:"decided_id"`
	Entries   []Command `json:"entries"`
}

func (v Value) decided() Value {
	return Value{
		Object:    v.Object,
		StartId:   v.StartId,
		DecidedId: v.DecidedId,
		Entries:   v.Entries[:v.DecidedId-v.StartId],
	}
}

func (v Value) copyEntries() Value {
	out := v // shallow copy
	out.Entries = make([]Command, len(v.Entries))
	copy(out.Entries, v.Entries)
	return out
}

func (v Value) append(entries []Command) Value {
	v.Entries = append(v.Entries, entries...)
	return v
}

func (v *Value) compact(block LogId, ratio LogId) LogId {
	beforeStartId := v.StartId
	if block > 0 && ratio >= 2 {
		for LogId(len(v.Entries)) > block*ratio {
			v.StartId += block
			v.Entries = v.Entries[block:]
		}
	}
	return beforeStartId
}
func (v *Value) tail(from LogId) []Command {
	return v.Entries[from-v.StartId:]
}

type AcceptorState struct {
	Promised ProposalId `json:"promised"`
	Accepted ProposalId `json:"accepted"`
	Value    Value      `json:"value"`
}

func (hs *AcceptorState) decideAndCompactEntries(entries []Command, compactionBlock LogId, compactionRatio LogId) LogId {
	hs.Value = hs.Value.decided().append(entries)
	hs.Value.DecidedId += LogId(len(entries))
	hs.Value.Object.Commit(entries...)
	return hs.Value.compact(compactionBlock, compactionRatio)
}
func (hs *AcceptorState) acceptEntries(accepted ProposalId, entries []Command) {
	hs.Value = hs.Value.decided().append(entries)
	hs.Accepted = accepted
}

type ProposerState struct {
	CountId     CountId          `json:"count_id"`
	AcceptIdMap map[NodeId]LogId `json:"accept_id_map"`
}

type ClusterState struct {
	AddressBook      map[NodeId]rpc.Address `json:"address_book"`
	RpcTimeout       time.Duration          `json:"rpc_timeout"`
	CompactionBlock  LogId                  `json:"compaction_block"`
	CompactionRatio  LogId                  `json:"compaction_ratio"`
	RetryUntilUpdate int                    `json:"retry_until_update"`
}
