package test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	env := NewTestEnvironment()
	go env.StartServer()
	defer env.StopServer()
	setupAndLogin(t)
}

func setupAndLogin(t *testing.T) playwright.Page {
	username := "admin"
	password := "admin"

	page, err := browser.NewPage()
	require.Nil(t, err)
	_, err = page.Goto("http://localhost:8080")
	require.Nil(t, err)

	submitMask := func(submitButtonName string) {
		usernameInput := page.GetByRole("textbox", playwright.PageGetByRoleOptions{
			Name: "Username",
		})
		require.Nil(t, usernameInput.Fill(username))

		passwordInput := page.GetByRole("textbox", playwright.PageGetByRoleOptions{
			Name: "Password",
		})
		require.Nil(t, passwordInput.Fill(password))

		submitButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
			Name: submitButtonName,
		})
		require.Nil(t, submitButton.Click())
	}

	textExists := func(text string) {
		signInText := page.GetByText(text)
		isVisible, err := signInText.IsVisible()
		require.Nil(t, err)
		require.True(t, isVisible)
	}

	setup := func() {
		submitMask("Create")
		textExists("Sign in to your account")
	}

	signin := func() {
		submitMask("Sign in")
		textExists("Welcome to Unterlagen")
	}

	setup()
	signin()
	return page
}
