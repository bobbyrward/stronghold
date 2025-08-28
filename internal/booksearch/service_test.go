package booksearch

import "testing"

func TestSearchParametersValidate(t *testing.T) {
    tests := []struct {
        name    string
        params  SearchParameters
        wantErr bool
    }{
        {
            name:    "empty params (neither query nor hash)",
            params:  SearchParameters{},
            wantErr: true,
        },
        {
            name:    "query only",
            params:  SearchParameters{Query: "Dune"},
            wantErr: false,
        },
        {
            name:    "hash only",
            params:  SearchParameters{Hash: "abcdef123456"},
            wantErr: false,
        },
        {
            name:    "both query and hash set",
            params:  SearchParameters{Query: "Dune", Hash: "abcdef123456"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.params.Validate()
            if tt.wantErr && err == nil {
                t.Fatalf("Validate() error = nil, want non-nil error")
            }
            if !tt.wantErr && err != nil {
                t.Fatalf("Validate() error = %v, want nil", err)
            }
        })
    }
}

