package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

func init() {

	for i := A; i <= URL; i++ {

		var cmdUpdateRecord = &cobra.Command{
			Use:   fmt.Sprintf("%s <domain name> <record name | id>", record.String(i)),
			Short: fmt.Sprintf("Update %s record type to a domain", record.String(i)),
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) < 2 {
					cmd.Usage()
					return
				}

				recordID, err := getRecordIDByName(csDNS, args[0], args[1])
				if err != nil {
					log.Fatal(err)
				}

				name, err := cmd.Flags().GetString("name")
				if err != nil {
					log.Fatal(err)
				}
				addr, err := cmd.Flags().GetString("content")
				if err != nil {
					log.Fatal(err)
				}
				ttl, err := cmd.Flags().GetInt("ttl")
				if err != nil {
					log.Fatal(err)
				}

				domain, err := csDNS.GetDomain(args[0])
				if err != nil {
					log.Fatal(err)
				}

				resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
					ID:         recordID,
					DomainID:   domain.ID,
					TTL:        ttl,
					RecordType: record.String(i),
					Name:       name,
					Content:    addr,
				})
				if err != nil {
					log.Fatal(err)
				}
				println(resp.ID)
			},
		}
		cmdUpdateRecord.Flags().StringP("name", "n", "", "Update name")
		cmdUpdateRecord.Flags().StringP("content", "c", "", "Update Content")
		cmdUpdateRecord.Flags().IntP("ttl", "t", 0, "Update ttl")
		cmdUpdateRecord.Flags().IntP("priority", "p", 0, "Update priority")
		dnsUpdateCmd.AddCommand(cmdUpdateRecord)
	}
}

//go:generate stringer -type=record

type record int

const (
	// A record type
	A record = iota
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
