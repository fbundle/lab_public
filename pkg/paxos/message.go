package paxos

type ProposalID = uint64
type RecordID = uint64

type PrepareRequest struct {
	RecordID RecordID `json:"record_id"`
	ProposalID ProposalID `json:"proposal_id"`
}

type PrepareResponse struct {
	RecordID RecordID `json:"record_id"`
	PromiseID ProposalID `json:promise_id`
	Value interface{} `json:value`
}

type AcceptRequest struct {
	RecordID RecordID `json:"record_id"`
	ProposalID ProposalID `json:"proposal_id"`
	Value interface{} `json:"value"`
}
