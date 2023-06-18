package quic_tester

type QUICResult struct {
	connectivity                     int
	initDecrypt                      int
	basicStreamReassembly            int
	flowControlAwareStreamReassembly int
	overlappingOffset                int
}

func TestQUIC(ip string, domain string) QUICResult {
	return QUICResult{
		connectivity:                     -5,
		initDecrypt:                      -5,
		basicStreamReassembly:            -5,
		flowControlAwareStreamReassembly: -5,
		overlappingOffset:                -5,
	}
}
