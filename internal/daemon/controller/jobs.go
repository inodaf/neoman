package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

func (c *controller) GetJobsStatus(w http.ResponseWriter, r *http.Request) {
	jobs, err := c.jobRepository.GetAll()
	if err != nil {
		slog.Error("failed to get jobs", "error", err)
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	var inProgress, pending, failed int
	var sb strings.Builder

	for _, job := range jobs {
		switch job.Status {
		case domain.JobStatusInProgress:
			inProgress++
		case domain.JobStatusPending:
			pending++
		case domain.JobStatusFailed:
			failed++
		}
	}

	sb.WriteString("Running: ")
	sb.WriteString(strconv.Itoa(inProgress))
	sb.WriteString("\nPending: ")
	sb.WriteString(strconv.Itoa(pending))
	sb.WriteString("\nFailed: ")
	sb.WriteString(strconv.Itoa(failed))
	sb.WriteString("\n\n")

	for _, job := range jobs {
		switch job.Status {
		case domain.JobStatusFailed:
			errMsg := "unknown error"
			if job.LastError != nil {
				errMsg = *job.LastError
			}
			sb.WriteString(job.ID)
			fmt.Fprintf(&sb, ": %s\n", job.Status)
			sb.WriteString(errMsg)
			sb.WriteString("\n")
		default:
			sb.WriteString(job.ID)
			fmt.Fprintf(&sb, ": %s\n", job.Status)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))
}
