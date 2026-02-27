package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type WorkflowRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	_ = godotenv.Load()

	owner := os.Getenv("OWNER")
	repo := os.Getenv("REPO")
	token := os.Getenv("GITHUB_TOKEN")

	if owner == "" || repo == "" || token == "" {
		log.Fatal("OWNER, REPO, and GITHUB_TOKEN must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := &http.Client{}

	runs, err := fetchWorkflowRuns(ctx, client, owner, repo, token)
	if err != nil {
		log.Fatalf("failed to fetch workflow runs: %v", err)
	}

	fmt.Printf("Total Count: %d\n", runs.TotalCount)

	for _, run := range runs.WorkflowRuns {
		if err := deleteWorkflowRun(ctx, client, owner, repo, token, run); err != nil {
			log.Printf("error deleting %s (%d): %v", run.Name, run.ID, err)
		}
	}
}

func fetchWorkflowRuns(ctx context.Context, client *http.Client, owner, repo, token string) (*WorkflowRuns, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(req, token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var runs WorkflowRuns
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return nil, err
	}

	return &runs, nil
}

func deleteWorkflowRun(ctx context.Context, client *http.Client, owner, repo, token string, run WorkflowRun) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d", owner, repo, run.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	setHeaders(req, token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("204 - %s (%d)\n", run.Name, run.ID)
		return nil
	}

	return fmt.Errorf("status %d", resp.StatusCode)
}

func setHeaders(req *http.Request, token string) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}
