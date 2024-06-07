package main

// infoT as received via the command line
type infoT struct {
	ApiURL      string `json:"api_url"`
	ExternalURL string `json:"external_url"`
	ApiToken    string `json:"api_user_access_token"`
}
