package rlog

type PrepareRequest struct {
	ProposalId ProposalId `json:"proposal_id"`
	FromId     LogId      `json:"from_id"`
}

type PrepareResponse struct {
	// header
	Promise  ProposalId `json:"promise"`
	AcceptId LogId      `json:"accept_id"`
	Success  bool       `json:"success"`
	// body
	Accepted ProposalId `json:"accepted"` // if success=true
	Entries  []Command  `json:"entries"`  // if success=true
}

type AcceptRequest struct {
	ProposalId ProposalId `json:"proposal_id"`
	FromId     LogId      `json:"from_id"`
	Entries    []Command  `json:"entries"`
}

type AcceptResponse struct {
	// header
	Promise  ProposalId `json:"promise"`
	AcceptId LogId      `json:"accept_id"`
	Success  bool       `json:"success"`
	// body
}

type DecideRequest struct {
	FromId  LogId     `json:"from_id"`
	Entries []Command `json:"entries"`
}

type DecideResponse struct{}

type UpdateRequest struct {
	DecidedId LogId `json:"decided_id"`
}

type UpdateResponse struct {
}

type UpdateResponseDecide struct {
	FromId  LogId     `json:"from_id"`
	Entries []Command `json:"entries"`
}

type UpdateResponseRestore struct {
	Snapshot  []byte    `json:"snapshot"`
	StartId   LogId     `json:"start_id"`
	DecidedId LogId     `json:"accept_id"`
	Entries   []Command `json:"entries"`
}
