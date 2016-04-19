package main

type TestDef struct {
        Title string `yaml:"title"`
        Iterations int `yaml:"iterations"`
        DockerComposeRoot string `yaml:"docker_compose_root"`
        DockerComposeFile string `yaml:"docker_compose_file"`
        Host string `yaml:"host"`
        Services []string `yaml:"services"`
        Endpoints []Endpoint `yaml:"endpoints"`
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
        Token_key string `yaml:"token_key"`
}

type Endpoint struct {
        Url string `yaml:"url"`
        Auth_method string `yaml:"auth_method"`
        Method string `yaml:"method"`
}


// Authorization: Bearer $TOKEN