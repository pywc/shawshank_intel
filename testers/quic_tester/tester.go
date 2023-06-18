package quic_tester

type QUICResult struct {
	Connectivity                     int `json:"connectivity,omitempty"`
	InitDecrypt                      int `json:"init_decrypt,omitempty"`
	BasicStreamReassembly            int `json:"basic_stream_reassembly,omitempty"`
	FlowControlAwareStreamReassembly int `json:"flow_control_aware_stream_reassembly,omitempty"`
	OverlappingOffset                int `json:"overlapping_offset,omitempty"`
}

func TestQUIC(ip string, domain string) QUICResult {
	return QUICResult{
		Connectivity:                     -5,
		InitDecrypt:                      -5,
		BasicStreamReassembly:            -5,
		FlowControlAwareStreamReassembly: -5,
		OverlappingOffset:                -5,
	}
}
