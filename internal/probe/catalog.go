package probe

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func GenerateCatalog(sourcePath, implementedPath string) (string, error) {
	var source SourceManifest
	if err := ReadJSON(sourcePath, &source); err != nil {
		return "", err
	}
	var implemented ImplementedManifest
	if err := ReadJSON(implementedPath, &implemented); err != nil {
		return "", err
	}
	if source.Source.Commit != implemented.SourceCommit {
		return "", fmt.Errorf("probe: catalog manifests use different source commits")
	}
	mappings := make(map[string]Mapping, len(implemented.Mappings))
	for _, mapping := range implemented.Mappings {
		if _, duplicate := mappings[mapping.Contract]; duplicate {
			return "", fmt.Errorf("probe: duplicate catalog mapping %q", mapping.Contract)
		}
		mappings[mapping.Contract] = mapping
	}
	var builder strings.Builder
	builder.WriteString("# API 对齐索引\n\n")
	builder.WriteString("此文件由 `go run ./cmd/bpi-probe api-doc` 自动生成，请勿手工编辑。\n\n")
	fmt.Fprintf(&builder, "基准：`bpi-rs` `%s`，提交 `%s`。\n\n", source.Source.CargoVersion, source.Source.Commit)
	fmt.Fprintf(&builder, "对齐进度：**%d/%d 条契约**，覆盖 **%d 个领域**。\n\n", len(mappings), source.Summary.PromotedContracts, source.Summary.DomainClients)
	builder.WriteString("| 风险级别 | 契约数 |\n| --- | ---: |\n")
	for _, risk := range []string{"public-read", "authenticated-read", "private-read", "login-session"} {
		fmt.Fprintf(&builder, "| `%s` | %d |\n", risk, source.Summary.RiskCounts[risk])
	}
	builder.WriteString("\n")
	domains := append([]SourceDomain(nil), source.Domains...)
	sort.Slice(domains, func(i, j int) bool { return domains[i].Name < domains[j].Name })
	for _, domain := range domains {
		fmt.Fprintf(&builder, "## %s\n\n", domain.Name)
		builder.WriteString("| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |\n")
		builder.WriteString("| --- | --- | --- | --- | --- | --- |\n")
		contracts := append([]SourceContract(nil), domain.Contracts...)
		sort.Slice(contracts, func(i, j int) bool { return contracts[i].Name < contracts[j].Name })
		for _, contract := range contracts {
			mapping, ok := mappings[contract.Name]
			if !ok {
				return "", fmt.Errorf("probe: catalog is missing mapping %q", contract.Name)
			}
			request := "FLOW"
			if contract.Method != "" {
				request = contract.Method + " " + catalogEndpoint(contract.URL)
			}
			fmt.Fprintf(&builder, "| `%s` | `%s` | `%s` | `%s` | `%s` | `%s` |\n",
				escapeTable(contract.Name), escapeTable(contract.Risk), escapeTable(request),
				escapeTable(mapping.Method), escapeTable(mapping.Params), escapeTable(mapping.Model))
		}
		builder.WriteString("\n")
	}
	if len(mappings) != source.Summary.PromotedContracts {
		return "", fmt.Errorf("probe: catalog has %d mappings, want %d", len(mappings), source.Summary.PromotedContracts)
	}
	return strings.TrimRight(builder.String(), "\n") + "\n", nil
}

func catalogEndpoint(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return rawURL
	}
	return parsed.Host + parsed.EscapedPath()
}

func escapeTable(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}
