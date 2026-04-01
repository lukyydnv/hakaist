package license

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	supabaseURL = "https://pikiqnsdfymzlqofrejo.supabase.co"
	supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBpa2lxbnNkZnltemxxb2ZyZWpvIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzQ5MDQ1MDMsImV4cCI6MjA5MDQ4MDUwM30.ijBy9FgjaH5S6r-AL2aDHFIWn0kH6bKeDIzdlaR_pAA"
)

type entry struct {
	Active    bool       `json:"active"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func Check(licenseKey string) {
	licenseKey = strings.TrimSpace(licenseKey)
	if licenseKey == "" {
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/rest/v1/licenses?license_key=eq.%s&select=active,expires_at", supabaseURL, licenseKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Exit(1)
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}

	var rows []entry
	if err := json.Unmarshal(body, &rows); err != nil || len(rows) == 0 {
		os.Exit(1)
	}

	l := rows[0]
	if !l.Active {
		os.Exit(1)
	}
	if l.ExpiresAt != nil && time.Now().After(*l.ExpiresAt) {
		os.Exit(1)
	}
}