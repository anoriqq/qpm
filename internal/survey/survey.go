package survey

import "github.com/AlecAivazis/survey/v2"

func AskOneInputRequired(msg, def string) (string, error) {
	var result string

	p := &survey.Input{ Message: msg, Default: def }

	err := survey.AskOne(p, &result, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return result, nil
}

func AskOneConfirm(msg string, def bool) (bool, error) {
	var result bool

	p := &survey.Confirm{ Message: msg, Default: def }

	err := survey.AskOne(p, &result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func AskOnePasswordRequired(msg string) (string, error) {
	var result string

	p := &survey.Password{ Message: msg }

	err := survey.AskOne(p, &result, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return result, nil
}
