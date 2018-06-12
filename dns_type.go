package egoscale

//go:generate stringer -type=record

// Record represent record type
type Record int

const (
	// A record type
	A Record = iota
	// AAAA record type
	AAAA
	// ALIAS record type
	ALIAS
	// CNAME record type
	CNAME
	// HINFO record type
	HINFO
	// MX record type
	MX
	// NAPTR record type
	NAPTR
	// NS record type
	NS
	// POOL record type
	POOL
	// SPF record type
	SPF
	// SRV record type
	SRV
	// SSHFP record type
	SSHFP
	// TXT record type
	TXT
	// URL record type
	URL
)
