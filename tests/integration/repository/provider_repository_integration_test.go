package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/model"
	"github.com/koubae/game-hangar/internal/identity/app/modules/auth/repository"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/koubae/game-hangar/tests/integration"
)

func TestProviderRepository_GetProvider(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown()

	tests := []struct {
		id       string
		source   string
		_type    string
		expected model.Provider
	}{
		{
			id:     "provider-username",
			source: "global",
			_type:  "username",
			expected: model.Provider{
				Source:      "global",
				Type:        "username",
				DisplayName: "Username",
				Category:    "managed",
			},
		},
		{
			id:     "provider-email",
			source: "global",

			_type: "email",
			expected: model.Provider{
				Source:      "global",
				Type:        "email",
				DisplayName: "Email",
				Category:    "managed",
			},
		},
		{
			id:     "provider-device",
			source: "global",
			_type:  "device",
			expected: model.Provider{
				Source:      "global",
				Type:        "device",
				DisplayName: "Device",
				Category:    "managed",
			},
		},
		{
			id:     "provider-guest",
			source: "global",
			_type:  "guest",
			expected: model.Provider{
				Source:      "global",
				Type:        "guest",
				DisplayName: "Guest",
				Category:    "anonymous",
			},
		},
		{
			id:     "provider-anonymous",
			source: "global",
			_type:  "anonymous",
			expected: model.Provider{
				Source:      "global",
				Type:        "anonymous",
				DisplayName: "Anonymous",
				Category:    "anonymous",
			},
		},
		{
			id:     "provider-steam",
			source: "global",
			_type:  "steam",
			expected: model.Provider{
				Source:      "global",
				Type:        "steam",
				DisplayName: "Steam",
				Category:    "platform",
			},
		},
		{
			id:     "provider-playstation",
			source: "global",
			_type:  "psn",
			expected: model.Provider{
				Source:      "global",
				Type:        "psn",
				DisplayName: "PlayStation Network",
				Category:    "platform",
			},
		},
		{
			id:     "provider-xbox",
			source: "global",
			_type:  "xbox",
			expected: model.Provider{
				Source:      "global",
				Type:        "xbox",
				DisplayName: "Xbox",
				Category:    "platform",
			},
		},
		{
			id:     "provider-nintendo",
			source: "global",
			_type:  "nintendo",
			expected: model.Provider{
				Source:      "global",
				Type:        "nintendo",
				DisplayName: "Nintendo",
				Category:    "platform",
			},
		},
	}

	providerRepository := repository.NewProviderRepository()
	providerRepository.LoadProviders(context.Background(), connector)
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			provider, err := providerRepository.GetProvider(context.Background(), connector, tt.source, tt._type)
			if err != nil {
				t.Fatalf("Failed to get provider: %v", err)
			}
			if provider == nil {
				t.Fatalf("Provider is nil")
			}

			if provider.Source != tt.expected.Source {
				t.Fatalf("Provider source is not '%s' got: %s\n ", tt.expected.Source, provider.Source)
			}

			if provider.Type != tt.expected.Type {
				t.Fatalf("Provider type is not '%s' got: %s\n ", tt.expected.Type, provider.Type)
			}
			if provider.DisplayName != tt.expected.DisplayName {
				t.Fatalf("Provider display name is not '%s' \n got: %s", tt.expected.DisplayName, provider.DisplayName)
			}
			if provider.Category != tt.expected.Category {
				t.Fatalf("Provider category is not '%s' \n got: %s", tt.expected.Category, provider.Category)
			}
		})
	}
}

func TestProviderRepository_GetProviderFoundWhenCacheMiss(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown()

	source := "global"
	_type := "username"
	providerRepository := repository.NewProviderRepository()
	provider, err := providerRepository.GetProvider(context.Background(), connector, source, _type)
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}
	if provider == nil {
		t.Fatalf("Provider is nil")
	}

	if provider.Source != source {
		t.Fatalf("Provider source is not '%s' got: %s\n ", source, provider.Source)
	}

	if provider.Type != _type {
		t.Fatalf("Provider type is not '%s' got: %s\n ", _type, provider.Type)
	}

	if provider.DisplayName != "Username" {
		t.Fatalf("Provider display name is not '%s' \n got: %s", "Username", provider.DisplayName)
	}
	if provider.Category != "managed" {
		t.Fatalf("Provider category is not '%s' \n got: %s", "managed", provider.Category)
	}
}

func TestProviderRepository_GetProviderNotFound(t *testing.T) {
	_, connector, tearDown := integration.DBWithCleanup(t)
	defer tearDown()

	source := "global"
	_type := "not-exists"
	providerRepository := repository.NewProviderRepository()
	providerRepository.LoadProviders(context.Background(), connector)

	provider, err := providerRepository.GetProvider(context.Background(), connector, source, _type)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			t.Fatalf("Failed to get provider: %v", err)
		}
	}
	if provider != nil {
		t.Fatalf("Provider is not nil")
	}
}
