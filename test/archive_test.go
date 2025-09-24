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

	// Click on the document to view details
	require.Nil(t, uploadedDocumentTitle.Click())

	// Test title editing functionality
	// Click the Edit button
	editButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Edit",
	})
	require.Nil(t, editButton.Click())

	// Verify the title input field is visible and focused
	titleInput := page.Locator("#title-input")
	isVisible, err = titleInput.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Title input should be visible in edit mode")

	// Verify action buttons have changed to Save/Cancel
	saveButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Save",
	})
	cancelButton := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Cancel",
	})
	isVisible, err = saveButton.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Save button should be visible in edit mode")
	isVisible, err = cancelButton.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Cancel button should be visible in edit mode")

	// Change the title
	newTitle := "Updated Invoice Document"
	require.Nil(t, titleInput.Fill(newTitle))

	// Click Save to submit the changes
	require.Nil(t, saveButton.Click())

	// Wait for page reload after save
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	// Verify the title has been updated
	updatedTitleDisplay := page.GetByRole("heading", playwright.PageGetByRoleOptions{
		Name: newTitle,
	})
	isVisible, err = updatedTitleDisplay.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Updated title should be displayed after save")

	// Verify Edit/Download/Trash buttons are back
	editButton = page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "Edit",
	})
	isVisible, err = editButton.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Edit button should be visible after save")

	// Test file download by clicking on download button

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

	t.Logf("Successfully uploaded file, edited title to '%s', and downloaded file: %s (size: %d bytes)", newTitle, downloadPath, fileInfo.Size())
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

func TestArchiveDocumentUploadAndViewDocumentInformation(t *testing.T) {
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

	// Wait 5 seconds for processing as requested
	page.WaitForTimeout(5000)

	// Verify the uploaded file appears in the document list
	uploadedDocumentTitle := page.GetByText("invoice_0001")
	isVisible, err = uploadedDocumentTitle.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Uploaded document should be visible in the archive")

	// Click on the document to view details
	err = uploadedDocumentTitle.Click()
	require.Nil(t, err)

	// Verify we can see document information
	// Check for document title
	documentTitles, err := page.GetByText("invoice_0001").All()
	require.Nil(t, err)
	require.Len(t, documentTitles, 2)
	isVisible, err = documentTitles[0].IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Document title should be visible")

	// Check for file size information
	fileSizeText := page.GetByText("KB")
	isVisible, err = fileSizeText.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "File size should be visible")

	// Check for download link
	downloadLink := page.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "Download",
	})
	isVisible, err = downloadLink.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Download link should be visible")

	// Check for creation date
	createdText := page.GetByText("Created")
	isVisible, err = createdText.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Creation date should be visible")

	// Check if PDF text was extracted
	extractedText := page.GetByText("INVOICE #0001")
	isVisible, err = extractedText.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Extracted text section should be visible")

	// Check if preview image is showing
	previewImage := page.Locator("img").Filter(playwright.LocatorFilterOptions{
		HasText: "",
	}).First()
	isVisible, err = previewImage.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Preview image should be visible")

	// Alternative check for preview image using more specific selector
	documentPreview := page.Locator("img[alt*='preview'], img[src*='preview']")
	count, err := documentPreview.Count()
	require.Nil(t, err)
	require.Greater(t, count, 0, "At least one preview image should be present")

	t.Logf("Successfully uploaded file, waited 5 seconds, and verified document information, text extraction, and preview image display")
}

func TestExportAllDocuments(t *testing.T) {
	env := NewTestEnvironment()
	go env.StartServer()
	defer env.StopServer()

	page := setupAndLogin(t)
	defer page.Close()

	_, err := page.Goto("http://localhost:8080/archive")
	require.Nil(t, err)

	testPDF1Path := filepath.Join("../testdata/mock_pdfs/invoice_0001.pdf")
	testPDF1AbsPath, err := filepath.Abs(testPDF1Path)
	require.Nil(t, err)

	testPDF2Path := filepath.Join("../testdata/mock_pdfs/contract_SA_0001.pdf")
	testPDF2AbsPath, err := filepath.Abs(testPDF2Path)
	require.Nil(t, err)

	files := []string{testPDF1AbsPath, testPDF2AbsPath}
	fileInput := page.Locator("input[name='documents'][type='file']")
	err = fileInput.SetInputFiles(files)
	require.Nil(t, err)

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	require.Nil(t, err)

	exportButton := page.GetByRole("link", playwright.PageGetByRoleOptions{
		Name: "Export All",
	})
	isVisible, err := exportButton.IsVisible()
	require.Nil(t, err)
	require.True(t, isVisible, "Export All button should be visible")

	download, err := page.ExpectDownload(func() error {
		return exportButton.Click()
	})
	require.Nil(t, err)

	downloadFilename := download.SuggestedFilename()
	require.Contains(t, downloadFilename, "documents-")
	require.Contains(t, downloadFilename, ".zip")

	tempDir := t.TempDir()
	savedPath := filepath.Join(tempDir, downloadFilename)
	err = download.SaveAs(savedPath)
	require.Nil(t, err)

	fileInfo, err := os.Stat(savedPath)
	require.Nil(t, err)
	require.Greater(t, fileInfo.Size(), int64(0), "Downloaded ZIP file should not be empty")

	t.Logf("Successfully exported all documents as ZIP: %s (size: %d bytes)", downloadFilename, fileInfo.Size())
}
