package tmpl

import "fmt"

const table = `
	<table>
		<tbody>
			<tr>
				<th>modulo</th>
				<th>version</th>
				<th>ambiente</th>
				<th>rama</th>
				<th>commit</th>
			</tr>
			<tr>
				<td>assi-ifx-associated-services</td>
				<td>0.0.1-6</td>
				<td>SNAPSHOT - RELEASE</td>
				<td>ESTABLE</td>
				<td>develop - master</td>
			</tr>
		</tbody>
	</table>
`

type Table struct {
	Module      string
	Version     string
	Environment string
	Branch      string
	Commits     string
}

func BuildTitle(value string) string {
	return "<h2>" + value + "</h2>"
}

func BuildDescription(value string) string {
	return "<p>" + value + "</p>"
}

func BuildTable(values []Table) string {
	var table string

	header := `
		<tr>
			<th><strong>modulo</strong></th>
			<th><strong>version</strong></th>
			<th><strong>ambiente</strong></th>
			<th><strong>rama</strong></th>
			<th><strong>commit</strong></th>
		</tr>
	`

	for _, v := range values {
		table += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Module, v.Version, v.Environment, v.Branch, v.Commits)
	}
	return fmt.Sprintf("<h2>Registro de cambios</h2><table><tbody>%s%s</tbody></table>", header, table)
}

func BuildPage(title string, description string, values []Table) string {
	return BuildTitle(title) + BuildDescription(description) + BuildTable(values)
}
