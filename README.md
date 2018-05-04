Go Based AWS SecretManager

Simple command line for fetching secrets from AWS Secrets Manager

#### Command line Arguments
Currently supported
```
  -c string
        Full path to credentials file (default "~\.aws\credentials")
  -k    Skip profile check and just use default for use when no cred file and default will work
  -p string
        Profile to use (default "default")
  -s string
        Secret To Fetch (default "secret")
  -v string
        Version of secret To Fetch (default "version")
```