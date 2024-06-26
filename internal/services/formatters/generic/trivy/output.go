// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trivy

import (
	"fmt"
	"strings"
)

type trivyOutput struct {
	ArtifactName string               `json:"artifactName"`
	ArtifactType string               `json:"artifactType"`
	Results      []*trivyOutputResult `json:"results"`
}

type trivyOutputResult struct {
	Target            string                   `json:"target"`
	Class             string                   `json:"class"`
	Type              string                   `json:"type"`
	Vulnerabilities   []*trivyVulnerability    `json:"vulnerabilities"`
	Misconfigurations []*trivyMisconfiguration `json:"misconfigurations"`
}

type trivyMisconfiguration struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Message     string   `json:"message"`
	Resolution  string   `json:"resolution"`
	References  []string `json:"references"`
	Severity    string   `json:"severity"`
	IacMetadata struct {
		Resource  string `json:"resource"`
		Provider  string `json:"provider"`
		Service   string `json:"service"`
		StartLine int    `json:"startline"`
		EndLine   int    `json:"endline"`
	} `json:"iacmetadata"`
}

type trivyVulnerability struct {
	VulnerabilityID  string   `json:"vulnerabilityID"`
	PkgName          string   `json:"pkgName"`
	PkgPath          string   `json:"pkgPath"`
	InstalledVersion string   `json:"installedVersion"`
	FixedVersion     string   `json:"fixedVersion"`
	SeveritySource   string   `json:"severitySource"`
	PrimaryURL       string   `json:"primaryURL"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Severity         string   `json:"severity"`
	CweIDs           []string `json:"cweIDs"`
}

func (v *trivyVulnerability) getDetails() string {
	details := v.getBaseDetailsWithoutCWEs()

	if len(v.CweIDs) > 0 {
		details = v.getDetailsWithCWEs(details)
	}
	details = strings.TrimRight(details, "\n")
	return strings.ReplaceAll(details, "\n\n", "")
}

func (v *trivyVulnerability) getBaseDetailsWithoutCWEs() (details string) {
	if v.Description != "" {
		details += v.Description + "\n"
	}

	if v.PrimaryURL != "" {
		details += fmt.Sprintf("PrimaryURL: %s.\n", v.PrimaryURL)
	}

	return details
}

func (v *trivyVulnerability) getInstalledVersionAndUpdateVersion() (code string) {
	if v.InstalledVersion != "" && v.FixedVersion != "" {
		return fmt.Sprintf("\tInstalled Version: %q\n\tUpdate to Version: %q for fix this issue.",
			v.InstalledVersion, v.FixedVersion)
	}
	return ""
}

// nolint:gomnd // magic number "2" is not necessary to check
func (v *trivyVulnerability) getDetailsWithCWEs(details string) string {
	details += "Cwe Links: "

	for _, ID := range v.CweIDs {
		idAfterSplit := strings.SplitAfter(ID, "-")
		if len(idAfterSplit) >= 2 {
			details += v.addCWELinkInDetails(details, idAfterSplit[1])
		}
	}

	return strings.TrimRight(details, ",")
}

func (v *trivyVulnerability) addCWELinkInDetails(details, cweID string) string {
	basePath := "https://cwe.mitre.org/data/definitions/"

	cweLink := basePath + cweID + ".html"
	if !strings.Contains(details, cweLink) {
		return fmt.Sprintf("(%s),", cweLink)
	}

	return ""
}
