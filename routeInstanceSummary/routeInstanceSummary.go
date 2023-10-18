package routeInstanceSummary

import "encoding/xml"

type RpcReply struct {
	XMLName             xml.Name `xml:"rpc-reply"`
	Text                string   `xml:",chardata"`
	Junos               string   `xml:"junos,attr"`
	InstanceInformation struct {
		Text         string `xml:",chardata"`
		Xmlns        string `xml:"xmlns,attr"`
		Style        string `xml:"style,attr"`
		InstanceCore []struct {
			Text         string `xml:",chardata"`
			InstanceName string `xml:"instance-name"`
			InstanceType string `xml:"instance-type"`
			InstanceRib  []struct {
				Text              string `xml:",chardata"`
				IribName          string `xml:"irib-name"`
				IribActiveCount   string `xml:"irib-active-count"`
				IribHolddownCount string `xml:"irib-holddown-count"`
				IribHiddenCount   string `xml:"irib-hidden-count"`
			} `xml:"instance-rib"`
		} `xml:"instance-core"`
	} `xml:"instance-information"`
	Cli struct {
		Text   string `xml:",chardata"`
		Banner string `xml:"banner"`
	} `xml:"cli"`
}
