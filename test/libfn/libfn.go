package libfn

import "regexp"

type SecretData struct {
	Data []byte
	Name string
}

func GetSecrets(projectId string, secretsRegexp *regexp.Regexp) ([]SecretData, error) {
	return nil, nil
}

type SecretsInterface interface {
	GetSecrets(projectId string, secretsRegexp *regexp.Regexp) ([]SecretData, error)
}
