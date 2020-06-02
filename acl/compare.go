package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
)

var templates = struct {
	report string
}{
	report: `ACL DIFF REPORT {{ .DateTime }}

{{range $id,$value := .Diffs}}
  DEVICE {{ $id }}{{if $value.Unchanged}}
    Incorrect:  {{range $value.Updated}}{{.}}
                {{end}}{{end}}{{if $value.Added}}
    Missing:    {{range $value.Added}}{{.}}
                {{end}}{{end}}{{if $value.Deleted}}
    Unexpected: {{range $value.Deleted}}{{.}}
                {{end}}{{end}}{{end}}
`,
}

type Report struct {
	DateTime types.DateTime
	Diffs    map[uint32]api.Diff
}

func (a *ACL) Compare(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	devices := ctx.Value("devices").([]*uhppote.Device)

	body := struct {
		URL struct {
			ACL    *string `json:"acl"`
			Report *string `json:"report"`
		} `json:"url"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}, fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.URL.ACL == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid download URL",
		}, fmt.Errorf("Missing/invalid download URL")
	}

	uri, err := url.Parse(*body.URL.ACL)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid download URL",
		}, fmt.Errorf("Invalid download URL '%v' (%w)", body.URL.ACL, err)
	}

	if body.URL.Report == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid report URL",
		}, fmt.Errorf("Missing/invalid report URL")
	}

	rpt, err := url.Parse(*body.URL.Report)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid report URL",
		}, fmt.Errorf("Invalid report URL '%v' (%w)", body.URL.Report, err)
	}

	acl, err := a.fetch("acl:compare", uri.String(), devices)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error downloading ACL",
		}, err
	}

	if acl == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error downloading ACL",
		}, fmt.Errorf("Download return nil ACL")
	}

	for k, l := range *acl {
		a.info("acl:compare", fmt.Sprintf("%v  Retrieved %v records", k, len(l)))
	}

	current, err := api.GetACL(impl.Uhppote, devices)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving current ACL",
		}, err
	}

	diff, err := api.Compare(current, *acl)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error comparing current and downloaded ACL's",
		}, err
	}

	var w strings.Builder
	if err := report(diff, templates.report, &w); err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error generating ACL compare report",
		}, err
	}

	filename := time.Now().Format("acl-2006-01-02T150405.rpt")
	if err = a.store("acl:compare", rpt.String(), filename, []byte(w.String())); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error uploading report",
		}, err
	}

	summary := map[uint32]struct {
		Same       int `json:"same"`
		Different  int `json:"different"`
		Missing    int `json:"missing"`
		Extraneous int `json:"extraneous"`
	}{}

	for k, v := range diff {
		a.info("acl:compare", fmt.Sprintf("%v  SUMMARY  same:%v  different:%v  missing:%v  extraneous:%v", k, len(v.Unchanged), len(v.Updated), len(v.Added), len(v.Deleted)))

		summary[k] = struct {
			Same       int `json:"same"`
			Different  int `json:"different"`
			Missing    int `json:"missing"`
			Extraneous int `json:"extraneous"`
		}{
			Same:       len(v.Unchanged),
			Different:  len(v.Updated),
			Missing:    len(v.Added),
			Extraneous: len(v.Deleted),
		}
	}

	return struct {
		URL    string `json:"url"`
		Report map[uint32]struct {
			Same       int `json:"same"`
			Different  int `json:"different"`
			Missing    int `json:"missing"`
			Extraneous int `json:"extraneous"`
		} `json:"report"`
	}{
		URL:    rpt.String(),
		Report: summary,
	}, nil
}

func report(diff map[uint32]api.Diff, format string, w io.Writer) error {
	t, err := template.New("report").Parse(format)
	if err != nil {
		return err
	}

	rpt := Report{
		DateTime: types.DateTime(time.Now()),
		Diffs:    diff,
	}

	return t.Execute(w, rpt)
}
