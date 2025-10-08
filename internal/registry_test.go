package internal

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	client := &godo.Client{}

	t.Run("no services specified", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		s := server.NewMCPServer("test-server", "0.0.1")

		err := Register(logger, s, client)
		require.NoError(t, err)

		logOutput := buf.String()
		require.Contains(t, logOutput, "no services specified, loading all supported services")

		for svc := range supportedServices {
			require.Contains(t, logOutput, fmt.Sprintf("Registering tool and resources for service: %s", svc))
		}
	})

	t.Run("one service specified", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		s := server.NewMCPServer("test-server", "0.0.1")

		err := Register(logger, s, client, "droplets")
		require.NoError(t, err)

		logOutput := buf.String()
		require.Contains(t, logOutput, "Registering tool and resources for service: droplets")
		require.NotContains(t, logOutput, "Registering tool and resources for service: networking")
		require.NotContains(t, logOutput, "no services specified, loading all supported services")
	})

	t.Run("multiple services specified", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		s := server.NewMCPServer("test-server", "0.0.1")

		err := Register(logger, s, client, "droplets", "networking")
		require.NoError(t, err)

		logOutput := buf.String()
		require.Contains(t, logOutput, "Registering tool and resources for service: droplets")
		require.Contains(t, logOutput, "Registering tool and resources for service: networking")
		require.NotContains(t, logOutput, "Registering tool and resources for service: apps")
	})

	t.Run("unsupported service", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		s := server.NewMCPServer("test-server", "0.0.1")

		err := Register(logger, s, client, "non-existent-service")
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported service: non-existent-service")
	})
}