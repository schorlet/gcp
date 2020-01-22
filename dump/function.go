package dump

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"text/template"
)

const body = `
{{define "req"}}
## HTTP headers
{{range $key, $_ := .Header}}
- {{ $key }}: {{ ($.Header.Get $key) }}{{end}}
{{end}}

{{define "jwt"}}
## JWT ClaimSet

` + "```" + `
{{.}}
` + "```" + `
{{end}}
`

var t *template.Template

func Dump(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// the function does not allow unauthenticated access
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// create the html template
	if t == nil {
		var ert error
		t, ert = template.New("body").Parse(body)
		if ert != nil {
			http.Error(w, ert.Error(), http.StatusInternalServerError)
			return
		}
	}

	var out bytes.Buffer

	// output request headers
	err := t.ExecuteTemplate(&out, "req", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// According to the documentation (https://cloud.google.com/functions/docs/securing/authenticating#getting_user_profile_information),
	// there is no need to validate the token, as the token has already been validated by Cloud IAM.

	// remove the Bearer Auth Scheme
	payload := auth[7:]
	// decode the claimset from the JWT payload
	claimset, err := decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// output JWT claimset
	err = t.ExecuteTemplate(&out, "jwt", claimset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=UTF-8")
	out.WriteTo(w)
}

func decode(payload string) (string, error) {
	s := strings.Split(payload, ".")
	if len(s) < 2 {
		return "", errors.New("invalid token")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = json.Indent(&out, decoded, "", "  ")
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
