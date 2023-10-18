package compareSubnets

import (
	"bufio"
	"github.com/pterm/pterm"
	"log"
	"os"
	"regexp"
	"strings"
	"vdcL3Migrations/actions"
)

// Subnet representation to include the command and the subnet
type Subnet struct {
	Command string
	Subnet  string
}

func CompareSubnets(pairs [][]string) {
	for _, pair := range pairs {
		// Extracts node name from the pair (e.g. D:\\vdcMigrations\\ld5-1a_route-checks-pre)
		path := pair[0]
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/Subnet comparison:", nodeName)
		pterm.Info.Printf("Comparing subnets: %s and %s\n", pair[1], pair[0])

		// Read file 1
		file1Bytes, err := os.ReadFile(pair[1])
		if err != nil {
			log.Fatal(err)
		}
		file1Content := string(file1Bytes)
		file1Subnets := extractSubnets(file1Content)

		// Read file 2
		file2Bytes, err := os.ReadFile(pair[0])
		if err != nil {
			log.Fatal(err)
		}
		file2Content := string(file2Bytes)
		file2Subnets := extractSubnets(file2Content)

		panel1 := "(*) Existing Subnets on the Edge Router:\n"
		panel1 += "==================================\n"

		// Group subnets by command for file 1
		file1SubnetsByCommand := groupSubnetsByCommand(file1Subnets)
		for command, subnets := range file1SubnetsByCommand {
			panel1 += pterm.FgLightYellow.Sprintf("Routing Table: %s\n", command)
			for _, subnet := range subnets {
				panel1 += pterm.FgLightYellow.Sprintf(" * %s\n", subnet)
			}
		}

		panel2 := "(-) Missing Subnets on BPE:\n"
		panel2 += "==================================\n"
		// Group subnets by command for file 2
		file2SubnetsByCommand := groupSubnetsByCommand(file2Subnets)

		// loop to go through file 1 subnets and compare with file 2
		for command, subnets := range file1SubnetsByCommand {
			file2Subnets, found := file2SubnetsByCommand[command]
			if !found {
				panel2 += pterm.FgLightRed.Sprintf("Routing Table: %s\n", command)
				for _, subnet := range subnets {
					panel2 += pterm.FgLightRed.Sprintf("- %s\n", subnet)
				}
			} else {
				// Compare subnets within the same command
				missingSubnets := findMissingSubnets(file2Subnets, subnets)
				if len(missingSubnets) > 0 {
					panel2 += pterm.FgLightRed.Sprintf("Routing Table: %s\n", command)
					for _, subnet := range missingSubnets {
						panel2 += pterm.FgLightRed.Sprintf("- %s\n", subnet)
					}
				}
			}
		}

		panel3 := "(+) Extra Subnets on BPE:\n"
		panel3 += "==================================\n"

		// Loop through the subnets in file 2
		for command, subnets := range file2SubnetsByCommand {
			file1Subnets, found := file1SubnetsByCommand[command]
			if !found {
				panel3 += pterm.FgLightGreen.Sprintf("Routing Table: %s\n", command)
				for _, subnet := range subnets {
					panel3 += pterm.FgLightGreen.Sprintf("+ %s\n", subnet)
				}
			} else {
				// Compare subnets within the same command
				extraSubnets := findMissingSubnets(file1Subnets, subnets) // The parameters are switched here
				if len(extraSubnets) > 0 {
					panel3 += pterm.FgLightGreen.Sprintf("Routing Table: %s\n", command)
					for _, subnet := range extraSubnets {
						panel3 += pterm.FgLightGreen.Sprintf("+ %s\n", subnet)
					}
				}
			}
		}

		panels := pterm.Panels{
			{{Data: pterm.Sprintf(panel1)}, {Data: pterm.Sprintf(panel2)}, {Data: pterm.Sprintf(panel3)}},
		}
		pterm.DefaultPanel.WithPanels(panels).Render()
	}
}

func extractSubnets(content string) []Subnet {
	subnets := []Subnet{}
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Reg ex to extract subnets and the commands from file
	reCommand := regexp.MustCompile(`show (ip route vrf|route table) (\S+)`)
	reSubnet := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2})`)

	var currentCommand string

	for scanner.Scan() {
		line := scanner.Text()

		commandMatch := reCommand.FindStringSubmatch(line)

		if len(commandMatch) > 0 {
			currentCommand = actions.ReplaceCommand(commandMatch[2])
			continue
		}

		if strings.HasPrefix(line, "") {
			subnetMatch := reSubnet.FindStringSubmatch(line)
			if len(subnetMatch) > 0 {
				subnet := Subnet{
					Command: currentCommand,
					Subnet:  subnetMatch[1],
				}
				subnets = append(subnets, subnet)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return subnets
}

// Group subnets by command
func groupSubnetsByCommand(subnets []Subnet) map[string][]string {
	subnetsByCommand := make(map[string][]string)

	for _, subnet := range subnets {
		subnetsByCommand[subnet.Command] = append(subnetsByCommand[subnet.Command], subnet.Subnet)
	}

	return subnetsByCommand
}

// Find missing subnets
func findMissingSubnets(subnets1 []string, subnets2 []string) []string {
	missingSubnets := []string{}

	subnetsMap := make(map[string]bool)
	for _, subnet := range subnets1 {
		subnetsMap[subnet] = true
	}

	for _, subnet := range subnets2 {
		if _, found := subnetsMap[subnet]; !found {
			missingSubnets = append(missingSubnets, subnet)
		}
	}

	return missingSubnets
}
