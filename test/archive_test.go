package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestFileUploadAndDownload(t *testing.T) {
	env := NewTestEnvironment()
	go env.StartServer()
	defer env.StopServer()

	page := setupAndLogin(t)
	defer page.Close()

	// Navigate to archive page
	_, err := page.Goto("http://localhost:8080/archive")
	require.Nil(t, err)

	// Verify we're on the archive page
	archiveHeading := page.GetByText("Archive")
	isVisible, err := archiveHeading.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible)

	// Test file upload
	testPDFPath := filepath.Join("../testdata/mock_pdfs/invoice_0001.pdf")
	testPDFAbsPath, err := filepath.Abs(testPDFPath)
	require.Nil(t, err)

	// Verify test file exists
	_, err = os.Stat(testPDFAbsPath)
	require.Nil(t, err, "Test PDF file should exist at %s", testPDFAbsPath)

	// Find the file input and upload the file
	fileInput := page.Locator("input[name='documents'][type='file']")
	err = fileInput.SetInputFiles(testPDFAbsPath)
	require.Nil(t, err)

	// Wait for the page to reload after upload (form auto-submits)
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Verify the uploaded file appears in the document list
	uploadedDocumentTitle := page.GetByText("invoice_0001")
	isVisible, err = uploadedDocumentTitle.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Uploaded document should be visible in the archive")
	// Test file download by clicking on the document
	// Set up download handler
	require.Nil(t, uploadedDocumentTitle.Click())

	downloadLink := page.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "Download",
	})
	download, err := page.ExpectDownload(func() error {
		return downloadLink.Click()
	})
	require.Nil(t, err)

	// Verify download properties
	downloadPath := download.SuggestedFilename()
	require.Equal(t, "invoice_0001.pdf", downloadPath, "Downloaded file should have correct filename")

	// Save the download to verify it's a valid file
	tempDir := t.TempDir()
	savedPath := filepath.Join(tempDir, "downloaded_invoice_0001.pdf")
	err = download.SaveAs(savedPath)
	require.Nil(t, err)

	// Verify the downloaded file exists and has content
	fileInfo, err := os.Stat(savedPath)
	require.Nil(t, err)
	require.Greater(t, fileInfo.Size(), int64(0), "Downloaded file should not be empty")

	t.Logf("Successfully uploaded and downloaded file: %s (size: %d bytes)", downloadPath, fileInfo.Size())
}

func TestFileUploadInFolder(t *testing.T) {
	env := NewTestEnvironment()
	go env.StartServer()
	defer env.StopServer()

	page := setupAndLogin(t)
	defer page.Close()

	// Navigate to archive page
	_, err := page.Goto("http://localhost:8080/archive")
	require.Nil(t, err)

	// Create a new folder first
	newFolderButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "New Folder",
	})
	err = newFolderButton.Click()
	require.Nil(t, err)

	// Fill in folder name in the modal
	folderNameInput := page.Locator("input[name='name']")
	err = folderNameInput.Fill("Test Folder")
	require.Nil(t, err)

	// Submit the folder creation form
	createButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Create",
	})
	err = createButton.Click()
	require.Nil(t, err)

	// Wait for page reload
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Navigate into the newly created folder
	folderCard := page.Locator("a.card.card-compact").Filter(playwright.LocatorFilterOptions{
		HasText: "Test Folder",
	})
	err = folderCard.Click()
	require.Nil(t, err)

	// Verify we're in the folder (breadcrumb should show "Test Folder")
	breadcrumb := page.GetByText("Test Folder")
	isVisible, err := breadcrumb.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible)

	// Upload a file to this folder
	testPDFPath := filepath.Join("../testdata/mock_pdfs/contract_SA_0001.pdf")
	testPDFAbsPath, err := filepath.Abs(testPDFPath)
	require.Nil(t, err)

	fileInput := page.Locator("input[name='documents'][type='file']")
	err = fileInput.SetInputFiles(testPDFAbsPath)
	require.Nil(t, err)

	// Wait for upload to complete
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Verify the file appears in the folder
	uploadedDocumentTitle := page.GetByText("contract_SA_0001")
	isVisible, err = uploadedDocumentTitle.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Uploaded document should be visible in the archive")

	t.Logf("Successfully uploaded file to folder: Test Folder")
}
