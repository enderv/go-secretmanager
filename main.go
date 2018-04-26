package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/alyu/configparser"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	sourceProfile := flag.String("p", "default", "Profile to use")
	secret := flag.String("s", "secret", "Secret To Fetch")
	version := flag.String("v", "version", "Version of secret To Fetch")
	credFile := flag.String("c", filepath.Join(getCredentialPath(), ".aws", "credentials"), "Full path to credentials file")
	flag.Parse()
	if *secret == "secret" {
		fmt.Println("You must specify a secret name to fetch")
		return
	}
	//Get Current Credentials
	exists, err := checkProfileExists(credFile, sourceProfile)
	if err != nil || !exists {
		fmt.Println(err.Error())
		return
	}
	sess := CreateSession(sourceProfile)
	getSecret(sess, secret, version)
}

// getMFACode prompts for MFA Token input
// It returns the value as a string and any error
func getMFACode() (string, error) {
	var mfa string
	fmt.Print("Enter MFA Token: ")
	reader := bufio.NewReader(os.Stdin)
	mfa, err := reader.ReadString('\n')
	if err != nil {
		return mfa, errors.New("failed to get token")
	}
	return strings.TrimSpace(mfa), nil
}

// CreateSession Creates AWS Session with specified profile
func CreateSession(profileName *string) *session.Session {
	profileNameValue := *profileName
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profileNameValue,
		Config:  aws.Config{Region: aws.String("us-east-1")},
	}))
	return sess
}

// getCredentialPath returns the users home directory path as a string
func getCredentialPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// checkProfileExists takes path to the credentials file and profile name to search for
// Returns bool and any errors
func checkProfileExists(credFile *string, profileName *string) (bool, error) {
	config, err := configparser.Read(*credFile)
	if err != nil {
		fmt.Println("Could not find credentials file")
		fmt.Println(err.Error())
		return false, err
	}
	section, err := config.Section(*profileName)
	if err != nil {
		fmt.Println("Could not find profile in credentials file")
		return false, nil
	}
	if !section.Exists("aws_access_key_id") {
		fmt.Println("Could not find access key in profile")
		return false, nil
	}

	return true, nil
}

func getSecret(sess *session.Session, secretName *string, secretVersion *string) {
	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("MyTestDatabaseSecret"),
		VersionStage: aws.String("AWSPREVIOUS"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
