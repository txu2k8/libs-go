package report

import (
	"fmt"
	"reflect"
)

// This file define some html template for email

// GenHTMLTableString .
func (desc *TestDesc) GenHTMLTableString() (tr string) {
	rType := reflect.TypeOf(desc)
	rVal := reflect.ValueOf(desc)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("desc must be ptr to struct")
	}

	htmlTemplate := `
	<tr id='%d' class='desc'>
		<td colspan='1' align='left' width='15%%'>%s </td>
		<td colspan='1' align='left'> %s</td>
	</tr>
	`
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		v := rVal.Field(i)
		if v.String() != "" {
			tr += fmt.Sprintf(htmlTemplate, i, t.Name, v)
		}
	}
	return
}

// GenHTMLTableString .
func (node *Node) GenHTMLTableString(idx int) (tr string) {
	htmlTemplate := `
	<tr id='%d' class='nodes'>
		<td colspan='1' align='center''>%s</td>
		<td colspan='1' align='center''>%s</td>
		<td colspan='1' align='center''>%s</td>
		<td colspan='1' align='center''>%s</td>
		<td colspan='1' align='center''>%s</td>
		<td colspan='1' align='center''>%s</td>
	</tr>
	`
	tr = fmt.Sprintf(htmlTemplate, idx, node.Name, node.Status, node.IPAddress, node.Roles, node.User, node.Password)
	return
}

// GenNodesHTMLTableString .
func GenNodesHTMLTableString(nodes []Node) (tr string) {
	for idx, node := range nodes {
		tr += node.GenHTMLTableString(idx)
	}
	return
}

// GenHTMLTableString .
func (r *TestResult) GenHTMLTableString(idx int) (tr string) {
	htmlTemplate := `
	<tr id='%d' class='result'>
		<td colspan='1' align='left' class='%s'>%s</td>
		<td colspan='1' align='center' class='%s'>%s</td>
		<td colspan='1' align='center'>%s</td>
		<td colspan='1' align='center'>Loop: %d</td>
	</tr>
	`
	outputTemplate := `
	<tr id='%d' class='output'>
		<td colspan='1' align='left' class='errorMsg'>Message</td>
		<td colspan='3' align='left' class='errorMsg'>%s</td>
	</tr>
	`
	style := "passCase"
	tid := fmt.Sprintf("pt1_%d", idx)
	switch r.Status {
	case FailStatus:
		style = "failCase"
		tid = fmt.Sprintf("ft1_%d", idx)
	case ErrorStatus:
		style = "errorCase"
		tid = fmt.Sprintf("et1_%d", idx)
	case SkipStatus:
		style = "skipCase"
		tid = fmt.Sprintf("st1_%d", idx)
	default:
		style = "passCase"
		tid = fmt.Sprintf("pt1_%d", idx)
	}
	tr = fmt.Sprintf(htmlTemplate, tid, style, r.Name, style, r.Status, r.Elapsed, r.Iteration)
	if r.Output != "" {
		tr += fmt.Sprintf(outputTemplate, tid, r.Output)
	}
	return
}

// GenResultHTMLTableString .
func GenResultHTMLTableString(results []TestResult) (tr string) {
	for idx, r := range results {
		tr += r.GenHTMLTableString(idx)
	}
	return
}
