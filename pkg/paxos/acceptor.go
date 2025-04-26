package paxos

type Record struct {
	RecordID RecordID `json:"record_id"`
	PromiseID ProposalID `json:"promise_id"`
	Value interface{} `json:"value"`
}

func newRecord() *Record {
	return &Record{
		PromiseID: 0,
		Value: nil,
	}
}

type Acceptor struct {
	RecordMap map[RecordID]*Record
}

func (a *Acceptor) HandlePrepareRequest(req *PrepareRequest) *PrepareResponse {
	var rec *Record
	if rec, ok := a.RecordMap[req.RecordID]; !ok {
		rec = newRecord()
		a.RecordMap[req.RecordID] = rec
	}

	if rec.PromiseID <= req.ProposalID{
		rec.PromiseID = req.ProposalID
	}

	return &PrepareResponse{
		RecordID: req.RecordID,
		PromiseID: rec.PromiseID,
		Value: Value,
	}
}

func (a *Acceptor) HandleAcceptRequest(req *AcceptRequest) {
	var rec *Record
	if rec, ok := a.RecordMap[req.RecordID]; !ok {
		rec = newRecord()
		a.RecordMap[req.RecordID] = rec
	}
	if rec.PromiseID <= req.ProposalID {
		rec.PromiseID = req.ProposalID
		rec.Value = req.Value
	}
}
