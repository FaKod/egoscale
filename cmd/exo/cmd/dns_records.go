package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// dnsACmd represents the A command
var dnsACmd = &cobra.Command{
	Use:   "A <domain name>",
	Short: "Add A record type to a domain",
	Long:  `Add an "A" record that points your domain or a subdomain to an IP address.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		addr, err := cmd.Flags().GetString("address")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "A",
			Name:       name,
			Content:    addr,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsACmd)
	dnsACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsACmd.Flags().StringP("address", "a", "", "Example: 127.0.0.1")
	dnsACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsACmd.MarkFlagRequired("address")
}

// AAAACmd represents the AAAA command
var dnsAAAACmd = &cobra.Command{
	Use:   "AAAA <domain name>",
	Short: "Add AAAA record type to a domain",
	Long:  `Add an "AAAA" record that points your domain to an IPv6 address. These records are the same as A records except they use IPv6 addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		addr, err := cmd.Flags().GetString("address")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "AAAA",
			Name:       name,
			Content:    addr,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsAAAACmd)
	dnsAAAACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsAAAACmd.Flags().StringP("address", "a", "", "Example: 2001:0db8:85a3:0000:0000:EA75:1337:BEEF")
	dnsAAAACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsAAAACmd.MarkFlagRequired("address")
}

// ALIASCmd represents the ALIAS command
var dnsALIASCmd = &cobra.Command{
	Use:   "ALIAS <domain name>",
	Short: "Add ALIAS record type to a domain",
	Long: `Add an "ALIAS" record. An ALIAS record is a special record that will
map a domain to another domain transparently. It can be used like a CNAME but
for a name with other records, like the root. When the record is resolved it will
look up the A records for the aliased domain and return those as the records for 
the record name. Note: If you want to redirect to a URL, use a URL record instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		alias, err := cmd.Flags().GetString("alias")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "ALIAS",
			Name:       name,
			Content:    alias,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsALIASCmd)
	dnsALIASCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsALIASCmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsALIASCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsALIASCmd.MarkFlagRequired("alias")
}

// CNAMECmd represents the CNAME command
var dnsCNAMECmd = &cobra.Command{
	Use:   "CNAME <domain name>",
	Short: "Add CNAME record type to a domain",
	Long: `Add a "CNAME" record that aliases a subdomain to another host.
These types of records are used when a server is reached by several names. Only use CNAME records on subdomains.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		alias, err := cmd.Flags().GetString("alias")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "CNAME",
			Name:       name,
			Content:    alias,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsCNAMECmd)
	dnsCNAMECmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsCNAMECmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsCNAMECmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsCNAMECmd.MarkFlagRequired("alias")
	dnsCNAMECmd.MarkFlagRequired("name")
}

// HINFOCmd represents the HINFO command
var dnsHINFOCmd = &cobra.Command{
	Use:   "HINFO <domain name>",
	Short: "Add HINFO record type to a domain",
	Long:  `Add an "HINFO" record is used to describe the CPU and OS of a host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		cpu, err := cmd.Flags().GetString("cpu")
		if err != nil {
			log.Fatal(err)
		}
		os, err := cmd.Flags().GetString("os")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "HINFO",
			Name:       name,
			Content:    fmt.Sprintf("%s %s", cpu, os),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsHINFOCmd)
	dnsHINFOCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsHINFOCmd.Flags().StringP("cpu", "c", "", "Example: IBM-PC/AT")
	dnsHINFOCmd.Flags().StringP("os", "o", "", "The operating system of the machine, example: Linux")
	dnsHINFOCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsHINFOCmd.MarkFlagRequired("cpu")
	dnsHINFOCmd.MarkFlagRequired("os")
}

// MXCmd represents the MX command
var dnsMXCmd = &cobra.Command{
	Use:   "MX <domain name>",
	Short: "Add MX record type to a domain",
	Long: `Add a mail exchange record that points to a mail server or relay.
These types of records are used to describe which servers handle incoming email.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		mailSrv, err := cmd.Flags().GetString("mail-server-host")
		if err != nil {
			log.Fatal(err)
		}
		priority, err := cmd.Flags().GetInt("priority")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "MX",
			Name:       name,
			Content:    mailSrv,
			Prio:       priority,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsMXCmd)
	dnsMXCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>")
	dnsMXCmd.Flags().StringP("mail-server-host", "m", "", "Example: mail-server.example.com")
	dnsMXCmd.Flags().IntP("priority", "p", 0, "Common values are for example 1, 5 or 10")
	dnsMXCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsMXCmd.MarkFlagRequired("mail-server-host")
	dnsMXCmd.MarkFlagRequired("priority")
}

// NAPTRCmd represents the NAPTR command
var dnsNAPTRCmd = &cobra.Command{
	Use:   "NAPTR <domain name>",
	Short: "Add NAPTR record type to a domain",
	Long: `Add an "NAPTR" record to provide a means to map a resource that is not in
the domain name syntax to a label that is. More information can be found in RFC 2915.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		order, err := cmd.Flags().GetInt("order")
		if err != nil {
			log.Fatal(err)
		}
		preference, err := cmd.Flags().GetInt("preference")
		if err != nil {
			log.Fatal(err)
		}

		flags := ""
		//flags
		s, err := cmd.Flags().GetBool("s")
		if err != nil {
			log.Fatal(err)
		}
		a, err := cmd.Flags().GetBool("a")
		if err != nil {
			log.Fatal(err)
		}
		u, err := cmd.Flags().GetBool("u")
		if err != nil {
			log.Fatal(err)
		}
		p, err := cmd.Flags().GetBool("p")
		if err != nil {
			log.Fatal(err)
		}

		if s {
			flags += "s"
		}
		if a {
			flags += "a"
		}
		if u {
			flags += "u"
		}
		if p {
			flags += "p"
		}

		service, err := cmd.Flags().GetString("service")
		if err != nil {
			log.Fatal(err)
		}
		regex, err := cmd.Flags().GetString("regex")
		if err != nil {
			log.Fatal(err)
		}
		replacement, err := cmd.Flags().GetString("replacement")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "NAPTR",
			Name:       name,
			Content:    fmt.Sprintf("%d %d %q %q %q %q", order, preference, flags, service, regex, replacement),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsNAPTRCmd)
	dnsNAPTRCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsNAPTRCmd.Flags().IntP("order", "o", 0, "Used to determine the processing order, lowest first.")
	dnsNAPTRCmd.Flags().IntP("preference", "", 0, "Used to give weight to records with the same value in the 'order' field, low to high.")
	dnsNAPTRCmd.Flags().StringP("service", "", "", "Service")
	dnsNAPTRCmd.Flags().StringP("regex", "", "", "The substituion expression.")
	dnsNAPTRCmd.Flags().StringP("replacement", "", "", "The next record to look up, which must be a fully-qualified domain name.")

	//flags
	dnsNAPTRCmd.Flags().BoolP("s", "", false, "Flag indicates the next lookup is for an SRV.")
	dnsNAPTRCmd.Flags().BoolP("a", "", false, "Flag indicates the next lookup is for an A or AAAA record.")
	dnsNAPTRCmd.Flags().BoolP("u", "", false, "Flag indicates the next record is the output of the regular expression as a URI.")
	dnsNAPTRCmd.Flags().BoolP("p", "", false, "Flag indicates that processing should continue in a protocol-specific fashion.")

	dnsNAPTRCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsNAPTRCmd.MarkFlagRequired("order")
	dnsNAPTRCmd.MarkFlagRequired("preference")
	dnsNAPTRCmd.MarkFlagRequired("service")
	dnsNAPTRCmd.MarkFlagRequired("regex")
	dnsNAPTRCmd.MarkFlagRequired("replacement")
}

// NSCmd represents the NS command
var dnsNSCmd = &cobra.Command{
	Use:   "NS <domain name>",
	Short: "Add NS record type to a domain",
	Long: `Add an "NS" record the delegates a domain to another name server.
You may only delegate subdomains (for example subdomain.yourdomain.com).`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		mailSrv, err := cmd.Flags().GetString("mail-server-host")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "NS",
			Name:       name,
			Content:    mailSrv,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsNSCmd)
	dnsNSCmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsNSCmd.Flags().StringP("mail-server-host", "m", "", "Example: 'ns1.example.com'")
	dnsNSCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsNSCmd.MarkFlagRequired("name")
	dnsNSCmd.MarkFlagRequired("mail-server-host")
}

// POOLCmd represents the POOL command
var dnsPOOLCmd = &cobra.Command{
	Use:   "POOL <domain name>",
	Short: "Add POOL record type to a domain",
	Long: `Add a "POOL" record that aliases a subdomain to another host as
part of a pool of available CNAME records. This is a DNSimple custom record type.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		alias, err := cmd.Flags().GetString("alias")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "POOL",
			Name:       name,
			Content:    alias,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsPOOLCmd)
	dnsPOOLCmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsPOOLCmd.Flags().StringP("alias", "a", "", "Alias for: Example: 'some-other-site.com'")
	dnsPOOLCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsPOOLCmd.MarkFlagRequired("name")
	dnsPOOLCmd.MarkFlagRequired("alias")
}

// SPFCmd represents the SPF command
var dnsSPFCmd = &cobra.Command{
	Use:   "SPF <domain name>",
	Short: "Add SPF record type to a domain",
	Long: `Add an "SPF" record to indicate what hosts and addresses are allowed to send mail from your domain.
When creating an SPF record we will automatically create a corresponding TXT record
for you as some older mail exchanges require a TXT version of the SPF record.

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("SPF called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSPFCmd)
}

// SRVCmd represents the SRV command
var dnsSRVCmd = &cobra.Command{
	Use:   "SRV <domain name>",
	Short: "Add SRV record type to a domain",
	Long:  `Add an "SRV" record to specify the location of servers for a specific service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}

		if name != "" {
			name = fmt.Sprintf(".%s", name)
		}

		symbName, err := cmd.Flags().GetString("symbolic-name")
		if err != nil {
			log.Fatal(err)
		}
		protocol, err := cmd.Flags().GetString("protocol")
		if err != nil {
			log.Fatal(err)
		}
		prio, err := cmd.Flags().GetInt("priority")
		if err != nil {
			log.Fatal(err)
		}
		weight, err := cmd.Flags().GetInt("weight")
		if err != nil {
			log.Fatal(err)
		}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}
		target, err := cmd.Flags().GetString("target")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "SRV",
			Name:       fmt.Sprintf("_%s._%s%s", symbName, protocol, name),
			Content:    fmt.Sprintf("%d %s %s", weight, port, target),
			Prio:       prio,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSRVCmd)
	dnsSRVCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsSRVCmd.Flags().StringP("symbolic-name", "s", "", "This will be a symbolic name for the service, like 'sip'. It might also be called Service at other DNS providers.")
	dnsSRVCmd.Flags().StringP("protocol", "p", "", "This will usually be 'TCP' or 'UDP'.")
	dnsSRVCmd.Flags().IntP("priority", "", 0, "Priority")
	dnsSRVCmd.Flags().IntP("weight", "w", 0, "A relative weight for 'SRV' records with the same priority.")
	dnsSRVCmd.Flags().StringP("port", "P", "", "The 'TCP' or 'UDP' port on which the service is found.")
	dnsSRVCmd.Flags().StringP("target", "", "", "The canonical hostname of the machine providing the service.")
	dnsSRVCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsSRVCmd.MarkFlagRequired("ttl")
	dnsSRVCmd.MarkFlagRequired("symbolic-name")
	dnsSRVCmd.MarkFlagRequired("protocol")
	dnsSRVCmd.MarkFlagRequired("priority")
	dnsSRVCmd.MarkFlagRequired("weight")
	dnsSRVCmd.MarkFlagRequired("port")
	dnsSRVCmd.MarkFlagRequired("target")
}

// SSHFPCmd represents the SSHFP command
var dnsSSHFPCmd = &cobra.Command{
	Use:   "SSHFP <domain name>",
	Short: "Add SSHFP record type to a domain",
	Long:  `Edit an "SSHFP" record to share your SSH fingerprint with others.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		algo, err := cmd.Flags().GetInt("algorithm")
		if err != nil {
			log.Fatal(err)
		}
		fingerIDType, err := cmd.Flags().GetInt("fingerprint-type")
		if err != nil {
			log.Fatal(err)
		}
		fingerprint, err := cmd.Flags().GetString("fingerprint")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "SSHFP",
			Name:       name,
			Content:    fmt.Sprintf("%d %d %s", algo, fingerIDType, fingerprint),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSSHFPCmd)
	dnsSSHFPCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsSSHFPCmd.Flags().IntP("algorithm", "a", 0, "RSA(1) | DSA(2) | ECDSA(3) | ED25519(4)")
	dnsSSHFPCmd.Flags().IntP("fingerprint-type", "", 0, "SHA1(1) | SHA256(2)")
	dnsSSHFPCmd.Flags().StringP("fingerprint", "f", "", "Fingerprint")
	dnsSSHFPCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsSSHFPCmd.MarkFlagRequired("ttl")
	dnsSSHFPCmd.MarkFlagRequired("algorithm")
	dnsSSHFPCmd.MarkFlagRequired("fingerprint-type")
}

// TXTCmd represents the TXT command
var dnsTXTCmd = &cobra.Command{
	Use:   "TXT <domain name>",
	Short: "Add TXT record type to a domain",
	Long: `Add a "TXT" record. This is useful for domain records that are not covered by
the standard record types. For example, Google uses this type of record for domain verification.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		content, err := cmd.Flags().GetString("content")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "TXT",
			Name:       name,
			Content:    content,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsTXTCmd)
	dnsTXTCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsTXTCmd.Flags().StringP("content", "c", "", "Content record")
	dnsTXTCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsTXTCmd.MarkFlagRequired("ttl")
	dnsTXTCmd.MarkFlagRequired("content")
}

// URLCmd represents the URL command
var dnsURLCmd = &cobra.Command{
	Use:   "URL <domain name>",
	Short: "Add URL record type to a domain",
	Long: `Add an URL redirection record that points your domain to a URL.
This type of record uses an HTTP redirect to redirect visitors from a domain to a web site.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		destURL, err := cmd.Flags().GetString("destination-url")
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

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "URL",
			Name:       name,
			Content:    destURL,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsURLCmd)
	dnsURLCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsURLCmd.Flags().StringP("destination-url", "d", "", "Example: https://www.example.com")
	dnsURLCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsURLCmd.MarkFlagRequired("ttl")
	dnsURLCmd.MarkFlagRequired("destination-url")
}