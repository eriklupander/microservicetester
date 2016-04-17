package main

type TestDef struct {
        Title string `yaml:"title"`
        Iterations int `yaml:"iterations"`
        Host string `yaml:"host"`
        Services []string `yaml:"services"`
        OAuth OAuthDef `yaml:"oauth"`
}

type OAuthDef struct {
        Url string `yaml:"url"`
        Client_id string `yaml:"client_id"`
        Client_password string `yaml:"client_password"`
        Scope string `yaml:"scope"`
        Grant_type string `yaml:"grant_type"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
}