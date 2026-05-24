// Package fieldgeo provides a log-entry transformer that resolves IP address
// fields to human-readable region labels using caller-supplied CIDR mappings.
//
// Usage:
//
//	tr := fieldgeo.New([]fieldgeo.Rule{
//		{
//			SrcField:  "client_ip",
//			DestField: "region",
//			Regions: map[string]string{
//				"10.0.0.0/8":  "internal",
//				"8.8.8.0/24":  "google",
//			},
//		},
//	})
//
//	annotated := tr.Apply(entry)
//
// If the IP in SrcField does not match any CIDR, the entry is returned
// unchanged. Invalid CIDRs in the Regions map are silently ignored at
// construction time.
package fieldgeo
