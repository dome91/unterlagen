package test

import (
	"path/filepath"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestSearchAfterUpload(t *testing.T) {
	// Setup: Initialize environment, start server, and log in
	env := NewTestEnvironment()
	go env.StartServer()
	defer env.
		StopServer()

	page := setupAndLogin(t)
	defer page.Close()

	// Navigate to the archive page
	_, err := page.Goto("http://localhost:8080/archive")
	require.Nil(t, err)

	// Define paths for the files to be uploaded
	invoicePDFPath, err := filepath.Abs("../testdata/mock_pdfs/invoice_0001.pdf")
	require.Nil(t, err)
	manualPDFPath, err := filepath.Abs("../testdata/mock_pdfs/manual_0001.pdf")
	require.Nil(t, err)
	presentationPDFPath, err := filepath.Abs("../testdata/mock_pdfs/presentation_001.pdf")
	require.Nil(t, err)

	// Upload the files
	fileInput := page.Locator("input[name='documents'][type='file']")
	err = fileInput.SetInputFiles([]string{invoicePDFPath, manualPDFPath, presentationPDFPath})
	require.Nil(t, err)

	// Wait for the upload to complete and the page to settle
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Verify that the success message is visible
	successMessage := page.GetByText("3 documents uploaded successfully")
	isVisible, err := successMessage.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Success message for 3 uploaded documents should be visible")

	// Navigate to the search page
	searchLink := page.GetByRole("link", playwright.PageGetByRoleOptions{Name: "Search"})
	require.Nil(t, searchLink.Click())
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Helper function to perform a search and verify the result
	searchAndVerify := func(searchTerm, expectedDocument string) {
		t.Run("SearchFor"+expectedDocument, func(t *testing.T) {
			searchInput := page.GetByRole("textbox", playwright.PageGetByRoleOptions{Name: "Search documents..."})
			require.Nil(t, searchInput.Fill(searchTerm))

			documentLink := page.GetByRole("link", playwright.PageGetByRoleOptions{Name: expectedDocument})
			// Necessary due to the dynamic nature of HTMX loading the results
			err := documentLink.WaitFor(playwright.LocatorWaitForOptions{
				State: playwright.WaitForSelectorStateAttached,
			})
			require.Nil(t, err)

			isVisible, err := documentLink.IsVisible()
			require.Nil(t, err)
			require.True(t, isVisible, "Document '%s' should be found when searching for '%s'", expectedDocument, searchTerm)
		})
	}

	// Perform searches for each of the uploaded documents
	searchAndVerify("invoice", "invoice_0001")
	searchAndVerify("manual", "manual_0001")
	searchAndVerify("presentation", "presentation_001")
}
