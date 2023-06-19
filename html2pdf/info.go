package main

type config struct {
	Config map[string]map[string]any `json:"config,omitempty"`
}

// infoT as received via the command line
type infoT struct {
	ApiURL      string `json:"api_url"`
	ExternalURL string `json:"external_url"`
	ApiToken    string `json:"api_user_access_token"`
	Config      struct {
		System config            `json:"system,omitempty"`
		Plugin map[string]config `json:"plugin,omitempty"`
	}
}
