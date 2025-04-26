package codec

import (
	"ca/pkg/rlog"
	"ca/pkg/rlog_util/codec/proto/gen/rlog_pb"
	"errors"
	"google.golang.org/protobuf/proto"
)

var (
	errTypeUnknown = errors.New("type unknown")
)

func Marshal(r interface{}) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	var m proto.Message
	switch r := r.(type) {
	case *rlog.PrepareRequest:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_PREPARE_REQUEST,
			PrepareRequest: &rlog_pb.PrepareRequest{
				ProposalId: uint64(r.ProposalId),
				FromId:     uint64(r.FromId),
			},
		}
	case *rlog.PrepareResponse:
		var accepted uint64
		var entries []string
		if r.Success {
			accepted = uint64(r.Accepted)
			entries = r.Entries
		}
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_PREPARE_RESPONSE,
			PrepareResponse: &rlog_pb.PrepareResponse{
				Promise:  uint64(r.Promise),
				AcceptId: uint64(r.AcceptId),
				Success:  r.Success,
				Accepted: &accepted,
				Entries:  entries,
			},
		}
	case *rlog.AcceptRequest:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_ACCEPT_REQUEST,
			AcceptRequest: &rlog_pb.AcceptRequest{
				ProposalId: uint64(r.ProposalId),
				FromId:     uint64(r.FromId),
				Entries:    r.Entries,
			},
		}
	case *rlog.AcceptResponse:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_ACCEPT_RESPONSE,
			AcceptResponse: &rlog_pb.AcceptResponse{
				Promise:  uint64(r.Promise),
				AcceptId: uint64(r.AcceptId),
				Success:  r.Success,
			},
		}
	case *rlog.DecideRequest:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_DECIDE_REQUEST,
			DecideRequest: &rlog_pb.DecideRequest{
				FromId:  uint64(r.FromId),
				Entries: r.Entries,
			},
		}
	case *rlog.DecideResponse:
		m = &rlog_pb.Message{
			Type:           rlog_pb.Type_DECIDE_RESPONSE,
			DecideResponse: &rlog_pb.DecideResponse{},
		}
	case *rlog.UpdateRequest:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_UPDATE_REQUEST,
			UpdateRequest: &rlog_pb.UpdateRequest{
				DecidedId: uint64(r.DecidedId),
			},
		}
	case *rlog.UpdateResponse:
		m = &rlog_pb.Message{
			Type:           rlog_pb.Type_UPDATE_RESPONSE,
			UpdateResponse: &rlog_pb.UpdateResponse{},
		}
	case *rlog.UpdateResponseDecide:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_UPDATE_RESPONSE_DECIDE,
			UpdateResponseDecide: &rlog_pb.UpdateResponseDecide{
				FromId:  uint64(r.FromId),
				Entries: r.Entries,
			},
		}
	case *rlog.UpdateResponseRestore:
		m = &rlog_pb.Message{
			Type: rlog_pb.Type_UPDATE_RESPONSE_RESTORE,
			UpdateResponseRestore: &rlog_pb.UpdateResponseRestore{
				Snapshot:  r.Snapshot,
				StartId:   proto.Uint64(uint64(r.StartId)),
				DecidedId: proto.Uint64(uint64(r.DecidedId)),
				Entries:   r.Entries,
			},
		}
	default:
		return nil, errTypeUnknown
	}
	return proto.Marshal(m)
}

func Unmarshal(b []byte) (interface{}, error) {
	if len(b) == 0 {
		return nil, nil
	}
	m := &rlog_pb.Message{}
	err := proto.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	switch m.GetType() {
	case rlog_pb.Type_PREPARE_REQUEST:
		r := m.GetPrepareRequest()
		return &rlog.PrepareRequest{
			ProposalId: rlog.ProposalId(r.GetProposalId()),
			FromId:     rlog.LogId(r.GetFromId()),
		}, nil
	case rlog_pb.Type_PREPARE_RESPONSE:
		r := m.GetPrepareResponse()
		return &rlog.PrepareResponse{
			Promise:  rlog.ProposalId(r.GetPromise()),
			AcceptId: rlog.LogId(r.GetAcceptId()),
			Success:  r.GetSuccess(),
			Accepted: rlog.ProposalId(r.GetAccepted()),
			Entries:  r.GetEntries(),
		}, nil
	case rlog_pb.Type_ACCEPT_REQUEST:
		r := m.GetAcceptRequest()
		return &rlog.AcceptRequest{
			ProposalId: rlog.ProposalId(r.GetProposalId()),
			FromId:     rlog.LogId(r.GetFromId()),
			Entries:    r.GetEntries(),
		}, nil
	case rlog_pb.Type_ACCEPT_RESPONSE:
		r := m.GetAcceptResponse()
		return &rlog.AcceptResponse{
			Promise:  rlog.ProposalId(r.GetPromise()),
			AcceptId: rlog.LogId(r.GetAcceptId()),
			Success:  r.GetSuccess(),
		}, nil
	case rlog_pb.Type_DECIDE_REQUEST:
		r := m.GetDecideRequest()
		return &rlog.DecideRequest{
			FromId:  rlog.LogId(r.GetFromId()),
			Entries: r.GetEntries(),
		}, nil
	case rlog_pb.Type_DECIDE_RESPONSE:
		return &rlog.DecideResponse{}, nil
	case rlog_pb.Type_UPDATE_REQUEST:
		r := m.GetUpdateRequest()
		return &rlog.UpdateRequest{
			DecidedId: rlog.LogId(r.DecidedId),
		}, nil
	case rlog_pb.Type_UPDATE_RESPONSE:
		return &rlog.UpdateResponse{}, nil
	case rlog_pb.Type_UPDATE_RESPONSE_DECIDE:
		r := m.GetUpdateResponseDecide()
		return &rlog.UpdateResponseDecide{
			FromId:  rlog.LogId(r.GetFromId()),
			Entries: r.GetEntries(),
		}, nil
	case rlog_pb.Type_UPDATE_RESPONSE_RESTORE:
		r := m.GetUpdateResponseRestore()
		return &rlog.UpdateResponseRestore{
			Snapshot:  r.GetSnapshot(),
			StartId:   rlog.LogId(r.GetStartId()),
			DecidedId: rlog.LogId(r.GetDecidedId()),
			Entries:   r.GetEntries(),
		}, nil
	default:
		return nil, errTypeUnknown
	}
}
