package rlog

import (
	"ca/pkg/rlog/rpc"
	"context"
	"fmt"
	"sync"
	"time"
)

type Node struct {
	NodeId     NodeId
	Router     rpc.Router
	Cluster    ClusterState
	Acceptor   AcceptorState
	Proposer   ProposerState
	acceptorMu sync.Mutex
}

type RetryPolicy func() time.Duration

func (n *Node) lockAcceptor(f func()) {
	n.acceptorMu.Lock()
	defer n.acceptorMu.Unlock()
	f()
}

// Update : if propose fails, it is due to
// (1) no quorum
// (2) my log is obsolete, call Update to update log
func (n *Node) Update(ctx context.Context) {
	var decidedId LogId
	n.lockAcceptor(func() {
		decidedId = n.Acceptor.Value.DecidedId
	})
	select {
	case <-ctx.Done():
		return
	case <-batchRPC(n.Router, n.Cluster, func(nodeId NodeId) func() (request interface{}) {
		return func() (request interface{}) {
			return &UpdateRequest{
				DecidedId: decidedId,
			}
		}
	}, func(nodeId NodeId) func(response interface{}) {
		return func(response interface{}) {
			if updateResponse, ok := response.(*UpdateResponseDecide); ok && updateResponse != nil {
				// decide
				n.onDecide(&DecideRequest{
					FromId:  updateResponse.FromId,
					Entries: updateResponse.Entries,
				})
			}
			if updateResponse, ok := response.(*UpdateResponseRestore); ok && updateResponse != nil {
				// snapshot
				n.lockAcceptor(func() {
					if startIdOffset := updateResponse.StartId - n.Acceptor.Value.StartId; startIdOffset > 0 {
						fmt.Printf("[node_%d] update start_id from %d to %d\n", n.NodeId, n.Acceptor.Value.StartId, updateResponse.StartId)
						// update start_id
						n.Acceptor.Value.StartId = updateResponse.StartId
						// cut entries to match start_id
						if LogId(len(n.Acceptor.Value.Entries)) <= startIdOffset {
							n.Acceptor.Value.Entries = nil
						} else {
							n.Acceptor.Value.Entries = n.Acceptor.Value.Entries[startIdOffset:]
						}
					}
					if n.Acceptor.Value.DecidedId < updateResponse.DecidedId {
						fmt.Printf("[node_%d] update decided_id from %d to %d\n", n.NodeId, n.Acceptor.Value.DecidedId, updateResponse.DecidedId)
						// update decided_id
						n.Acceptor.Value.DecidedId = updateResponse.DecidedId
						// cope entries to match decided_id
						if len(n.Acceptor.Value.Entries) < len(updateResponse.Entries) {
							n.Acceptor.Value.Entries = updateResponse.Entries
						} else {
							copy(n.Acceptor.Value.Entries, updateResponse.Entries)
						}
						// unmarshal object
						fmt.Printf("[node_%d] unmarshal from snapshot %s\n", n.NodeId, string(updateResponse.Snapshot))
						err := n.Acceptor.Value.Object.Unmarshal(updateResponse.Snapshot)
						if err != nil {
							panic(err)
						}
					}
				})
			}
		}
	}):
	}
}

// Propose : value must be unique among all propose calls
// If Propose is called sequentially and AcceptorState is properly saved into persistent storage
// It can recover at any stage
func (n *Node) Propose(ctx context.Context, command Command, nextTimeout RetryPolicy) error {
	tryCount := 0
	for {
		tryCount++
		if ok := n.proposeOnce(ctx, command); ok {
			return nil
		}
		if n.Cluster.RetryUntilUpdate > 0 && tryCount >= n.Cluster.RetryUntilUpdate {
			n.Update(ctx)
		}
		timeout := nextTimeout()
		timer := time.NewTimer(timeout)
		done := false
		select {
		case <-ctx.Done():
			done = true
		case <-timer.C:
		}
		timer.Stop()
		if done {
			return ctx.Err()
		}
	}
}

// proposeOnce : propose once
func (n *Node) proposeOnce(ctx context.Context, command Command) bool {
	var decidedValue Value
	n.lockAcceptor(func() {
		decidedValue = n.Acceptor.Value.decided().copyEntries()
	})

	onResponseMu := sync.Mutex{}
	n.Proposer.CountId++
	proposalId := MakeProposalId(n.NodeId, n.Proposer.CountId)
	// PREPARE
	prepareResponseMap := map[NodeId]*PrepareResponse{}
	{
		select {
		case <-ctx.Done():
			return false
		case <-batchRPC(n.Router, n.Cluster, func(nodeId NodeId) func() (request interface{}) {
			return func() (request interface{}) {
				return &PrepareRequest{
					ProposalId: proposalId,
					FromId:     decidedValue.DecidedId,
				}
			}
		}, func(receiver NodeId) func(response interface{}) {
			return func(response interface{}) {
				if prepareResponse, ok := response.(*PrepareResponse); ok && prepareResponse != nil {
					onResponseMu.Lock()
					defer onResponseMu.Unlock()
					if prepareResponse.Success {
						prepareResponseMap[receiver] = prepareResponse
					}
					n.processHeader(receiver, prepareResponse.Promise, prepareResponse.AcceptId)
				}
			}
		}):
		}

		if len(prepareResponseMap) <= len(n.Cluster.AddressBook)/2 {
			return false
		}
		fmt.Printf("[node_%d] prepare ok command=%+v,proposal_id=%+v\n", n.NodeId, command, proposalId)
	}
	// PROCESS PREPARE
	var addonEntries []Command
	{
		maxAccepted := ProposalId(0)
		for _, prepareResponse := range prepareResponseMap {
			if maxAccepted < prepareResponse.Accepted {
				maxAccepted = prepareResponse.Accepted
				addonEntries = prepareResponse.Entries
			}
		}
		existed := false
		for _, e := range addonEntries {
			if e == command {
				existed = true
				break
			}
		}
		if !existed {
			addonEntries = append(addonEntries, command)
		}
		decidedValue = decidedValue.append(addonEntries)
		decidedValue.DecidedId += LogId(len(addonEntries))
	}
	// ACCEPT
	{
		acceptResponseMap := map[NodeId]*AcceptResponse{}
		select {
		case <-ctx.Done():
			return false
		case <-batchRPC(n.Router, n.Cluster, func(id NodeId) func() (request interface{}) {
			return func() (request interface{}) {
				fromId := decidedValue.DecidedId
				if lid, ok := n.Proposer.AcceptIdMap[id]; ok {
					fromId = lid
				}
				if fromId < decidedValue.StartId {
					// ignore obsolete node
					return nil
				}
				return &AcceptRequest{
					ProposalId: proposalId,
					FromId:     fromId,
					Entries:    decidedValue.tail(fromId),
				}
			}
		}, func(receiver NodeId) func(response interface{}) {
			return func(response interface{}) {
				if acceptResponse, ok := response.(*AcceptResponse); ok && acceptResponse != nil {
					onResponseMu.Lock()
					defer onResponseMu.Unlock()
					if acceptResponse.Success {
						acceptResponseMap[receiver] = acceptResponse
					}
					n.processHeader(receiver, acceptResponse.Promise, acceptResponse.AcceptId)
				}
			}
		}):
		}

		if len(acceptResponseMap) <= len(n.Cluster.AddressBook)/2 {
			return false
		}
		fmt.Printf("[node_%d] accept ok command=%+v,proposal_id=%+v,next_id=%d,entries=%+v\n", n.NodeId, command, proposalId, decidedValue.DecidedId, addonEntries)
	}
	// DECIDE
	go batchRPC(n.Router, n.Cluster, func(id NodeId) func() (request interface{}) {
		return func() (request interface{}) {
			fromId := decidedValue.DecidedId
			if lid, ok := n.Proposer.AcceptIdMap[id]; ok {
				fromId = lid
			}
			if fromId < decidedValue.StartId {
				// ignore obsolete node
				return nil
			}
			return &DecideRequest{
				FromId:  fromId,
				Entries: decidedValue.tail(fromId),
			}
		}
	}, nil)
	return true
}

func (n *Node) onPrepare(request *PrepareRequest) (response *PrepareResponse) {
	if request == nil {
		return nil
	}
	n.lockAcceptor(func() {
		response = &PrepareResponse{
			Promise:  n.Acceptor.Promised,
			AcceptId: n.Acceptor.Value.DecidedId,
			Success:  false,
		}
		if request.ProposalId <= n.Acceptor.Promised {
			return
		}
		n.Acceptor.Promised = request.ProposalId
		if offset := request.FromId - n.Acceptor.Value.StartId; 0 <= offset && offset <= LogId(len(n.Acceptor.Value.Entries)) {
			response.Success = true
			response.Accepted = n.Acceptor.Accepted
			response.Entries = n.Acceptor.Value.tail(request.FromId)
			fmt.Printf("[node_%d] promise to request=%+v,response=%+v\n", n.NodeId, *request, *response)
		} else {
			// ignore obsolete node
		}
	})
	return response
}

func (n *Node) onAccept(request *AcceptRequest) (response *AcceptResponse) {
	if request == nil {
		return nil
	}
	n.lockAcceptor(func() {
		response = &AcceptResponse{
			Promise:  n.Acceptor.Promised,
			AcceptId: n.Acceptor.Value.DecidedId,
			Success:  false,
		}
		if request.ProposalId < n.Acceptor.Promised {
			return
		}
		n.Acceptor.Promised = request.ProposalId

		if offset := n.Acceptor.Value.DecidedId - request.FromId; 0 <= offset {
			// do not ignore obsolete node
			response.Success = true
			var entries []Command
			if offset < LogId(len(request.Entries)) {
				entries = request.Entries[offset:] // guarantee all commands before offset are identical with node log
			}
			n.Acceptor.acceptEntries(request.ProposalId, entries)
			fmt.Printf("[node_%d] accept to request=%+v,response=%+v\n", n.NodeId, *request, *response)
		}
	})
	return response
}
func (n *Node) onDecide(request *DecideRequest) (response *DecideResponse) {
	if request == nil {
		return nil
	}
	var entries []Command
	n.lockAcceptor(func() {
		response = &DecideResponse{}
		if entriesOffset := n.Acceptor.Value.DecidedId - request.FromId; entriesOffset >= 0 {
			if entriesOffset < LogId(len(request.Entries)) {
				entries = request.Entries[entriesOffset:]
				beforeStartId := n.Acceptor.decideAndCompactEntries(entries, n.Cluster.CompactionBlock, n.Cluster.CompactionRatio)
				if beforeStartId < n.Acceptor.Value.StartId {
					fmt.Printf("[node_%d] compact log from %d to %d\n", n.NodeId, beforeStartId, n.Acceptor.Value.StartId)
				}
				fmt.Printf("[node_%d] decide ok %+v\n", n.NodeId, entries)
			}
		}
	})
	return response
}

// onUpdate : return either UpdateResponse or UpdateResponseDecide or UpdateResponseRestore
func (n *Node) onUpdate(request *UpdateRequest) (response interface{}) {
	if request == nil {
		return nil
	}
	n.lockAcceptor(func() {
		decidedValue := n.Acceptor.Value.decided()
		if decidedValue.DecidedId > request.DecidedId {
			if decidedValue.StartId > request.DecidedId {
				// snapshot
				snapshot, err := decidedValue.Object.Marshal()
				if err != nil {
					panic(err)
				}
				response = &UpdateResponseRestore{
					Snapshot:  snapshot,
					StartId:   decidedValue.StartId,
					DecidedId: decidedValue.DecidedId,
					Entries:   decidedValue.Entries,
				}
			} else {
				// decide
				response = &UpdateResponseDecide{
					FromId:  request.DecidedId,
					Entries: decidedValue.tail(request.DecidedId),
				}
			}
		} else {
			// nothing
			response = &UpdateResponse{}
		}
	})
	return response
}
func (n *Node) Handler(call *rpc.Call) {
	response := func() interface{} {
		switch request := call.Request.(type) {
		case *PrepareRequest:
			return n.onPrepare(request)
		case *AcceptRequest:
			return n.onAccept(request)
		case *DecideRequest:
			return n.onDecide(request)
		case *UpdateRequest:
			return n.onUpdate(request)
		default:
			return nil
		}
	}()
	call.Write(response)
}

func (n *Node) processHeader(receiver NodeId, promise ProposalId, acceptId LogId) {
	if n.Proposer.CountId < promise.CountId() {
		n.Proposer.CountId = promise.CountId()
	}
	n.Proposer.AcceptIdMap[receiver] = acceptId
}
