package web

import (
	"context"
	"embed"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"unterlagen/features/administration"
	"unterlagen/features/archive"
	"unterlagen/features/common"
	"unterlagen/features/search"
	"unterlagen/platform/configuration"
	"unterlagen/platform/web/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

//go:generate go tool templ generate
//go:generate npm run build:css
//go:generate npm run build:js
//go:embed public
var public embed.FS

const startMessage = `
 _    _       _            _
| |  | |     | |          | |
| |  | |_ __ | |_ ___ _ __| | __ _  __ _  ___ _ __
| |  | | '_ \| __/ _ \ '__| |/ _  |/ _  |/ _ \ '_  \
| |__| | | | | ||  __/ |  | | (_| | (_| |  __/ | | |
 \____/|_| |_|\__\___|_|  |_|\__,_|\__, |\___|_| |_|
                                    __/ |
                                   |___/
                SERVER STARTED
`

type Server struct {
	administration *administration.Administration
	archive        *archive.Archive
	search         *search.Search
	sessionStore   sessions.Store
	internal       *http.Server
}

func (server *Server) Start() {
	go func() {
		println(startMessage)
		if err := server.internal.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (server *Server) buildPublicHandler() http.HandlerFunc {
	publicFs, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix("/public/", http.FileServer(http.FS(publicFs))).ServeHTTP
}

func (server *Server) setup(w http.ResponseWriter, r *http.Request) {
	if server.administration.AdminExists() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	notifications := server.buildNotifications(r, w)
	templates.Setup(notifications).Render(r.Context(), w)
}

func (server *Server) handleSetup(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	username := r.FormValue("username")
	if username == "" {
		session.AddFlash("Username is required", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		session.AddFlash("Password is required", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	err := server.administration.CreateUser(username, password, administration.UserRoleAdmin)
	if err != nil {
		session.AddFlash(err.Error(), "error")
		session.Save(r, w)
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	session.AddFlash("Setup successful! Please log in.", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (server *Server) login(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	_, ok := session.Values["username"]
	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	notifications := server.buildNotifications(r, w)
	templates.Login(notifications).Render(r.Context(), w)
}

func (server *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	username := r.FormValue("username")
	if username == "" {
		session.AddFlash("Username is required", "error")
		session.Save(r, w)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		session.AddFlash("Password is required", "error")
		session.Save(r, w)
		return
	}

	user, err := server.administration.GetUser(username)
	if err != nil {
		session.AddFlash("That didn't work. Please try again.", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if !user.IsValidPassword(password) {
		session.AddFlash("That didn't work. Please try again.", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	notifications := server.buildNotifications(r, w)
	templates.Home(notifications, server.isAdmin(r)).Render(r.Context(), w)
}

func (server *Server) getArchive(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	folderID := "root"
	folderIDs := r.URL.Query()["folderID"]
	if len(folderIDs) > 0 {
		folderID = folderIDs[0]
	}

	// Verify user owns the folder they're trying to access
	_, err := server.archive.GetFolder(folderID, user)
	if err != nil {
		slog.Error("failed to get folder", slog.String("folderID", folderID), slog.String("user", user), slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	// Check if we should show trashed documents
	showTrashed := false
	if _, exists := r.URL.Query()["showTrashed"]; exists {
		showTrashed = true
	}

	documents, err := server.archive.GetDocumentsInFolder(folderID, user)
	if err != nil {
		slog.Error("failed to get documents in folder", slog.String("folderID", folderID), slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	folders, err := server.archive.GetFolderChildren(folderID, user)
	if err != nil {
		slog.Error("failed to get children of folder", slog.String("folderID", folderID), slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	hierarchy, err := server.archive.GetFolderHierarchy(folderID, user)
	if err != nil {
		slog.Error("failed to get hierarchy of folder", slog.String("folderID", folderID), slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	notifications := server.buildNotifications(r, w)
	templates.Archive(folderID, documents, folders, hierarchy, notifications, server.isAdmin(r), showTrashed).Render(r.Context(), w)
}

func (server *Server) handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	parentFolderID := r.PostFormValue("parentFolderID")
	session := server.getSession(r)
	username := server.getAuthenticatedUser(r)
	if name == "" {
		session.AddFlash("Name is required", "error")
		session.Save(r, w)
		return
	}

	if parentFolderID == "" {
		parentFolderID = archive.FolderRootID
	}

	// Verify user owns the parent folder
	_, err := server.archive.GetFolder(parentFolderID, username)
	if err != nil {
		slog.Error("failed to verify parent folder ownership", slog.String("error", err.Error()))
		session.AddFlash("You don't have permission to create folders here", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/archive", http.StatusFound)
		return
	}

	err = server.archive.CreateFolder(name, parentFolderID, username)
	if err != nil {
		slog.Error("failed to create folder", slog.String("error", err.Error()))
		server.createGenericErrorNotification()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/archive?folderID=%s", parentFolderID), http.StatusFound)
}

func (server *Server) handleSynchronize(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	folderID := r.PostFormValue("folderID")
	if folderID == "" {
		folderID = archive.FolderRootID
	}

	userID := server.getAuthenticatedUser(r)
	err := server.archive.Synchronize(userID)
	if err != nil {
		slog.Error("failed to synchronize archive", slog.String("error", err.Error()))
		server.createGenericErrorNotification()
		return
	}

	// Add success notification that synchronization was started
	session.AddFlash("Archive synchronization started", "success")
	session.Save(r, w)

	// Redirect back to the archive page with the current folder
	http.Redirect(w, r, fmt.Sprintf("/archive?folderID=%s", folderID), http.StatusFound)
}

func (server *Server) handleUploadDocument(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	username := server.getAuthenticatedUser(r)
	folderID := r.PostFormValue("folderID")
	if folderID == "" {
		folderID = archive.FolderRootID
	}

	// Verify user owns the target folder
	_, err := server.archive.GetFolder(folderID, username)
	if err != nil {
		slog.Error("failed to verify folder ownership for upload", slog.String("error", err.Error()))
		session.AddFlash("You don't have permission to upload to this folder", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/archive", http.StatusFound)
		return
	}

	err = r.ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		session.AddFlash("Failed to parse uploaded files", "error")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/archive?folderID=%s", folderID), http.StatusFound)
		return
	}

	files := r.MultipartForm.File["documents"]
	if len(files) == 0 {
		session.AddFlash("No files selected", "error")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/archive?folderID=%s", folderID), http.StatusFound)
		return
	}

	uploadedCount := 0
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			slog.Error("failed to open uploaded file", slog.String("filename", fileHeader.Filename), slog.String("error", err.Error()))
			continue
		}

		err = server.archive.UploadDocument(fileHeader.Filename, uint64(fileHeader.Size), folderID, username, file)
		file.Close()

		if err != nil {
			slog.Error("failed to upload document", slog.String("filename", fileHeader.Filename), slog.String("error", err.Error()))
			continue
		}

		uploadedCount++
	}

	if uploadedCount == 0 {
		session.AddFlash("Failed to upload any documents", "error")
	} else if uploadedCount == len(files) {
		if len(files) == 1 {
			session.AddFlash("Document uploaded successfully", "success")
		} else {
			session.AddFlash(fmt.Sprintf("%d documents uploaded successfully", uploadedCount), "success")
		}
	} else {
		session.AddFlash(fmt.Sprintf("%d of %d documents uploaded successfully", uploadedCount, len(files)), "warning")
	}

	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/archive?folderID=%s", folderID), http.StatusFound)
}

func (server *Server) getDocumentDetails(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	if documentID == "" {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	document, err := server.archive.GetDocument(documentID, user)
	if err != nil {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	notifications := server.buildNotifications(r, w)
	templates.DocumentDetails(document, notifications, server.isAdmin(r)).Render(r.Context(), w)
}

func (server *Server) downloadDocument(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	if documentID == "" {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	document, err := server.archive.GetDocument(documentID, user)
	if err != nil {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", document.Filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", document.Filesize))
	err = server.archive.DownloadDocument(documentID, user, func(r io.Reader) error {
		_, err := io.Copy(w, r)
		return err
	})

	if err != nil {
		slog.Error("failed to download document",
			slog.String("documentID", documentID),
			slog.String("user", user),
			slog.String("error", err.Error()))
	}
}

func (server *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	if documentID == "" {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	session := server.getSession(r)
	err := server.archive.TrashDocument(documentID, user)
	if err != nil {
		slog.Error("failed to delete document", slog.String("error", err.Error()))
		session.AddFlash("Failed to delete document", "error")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/archive/documents/%s", documentID), http.StatusFound)
		return
	}

	session.AddFlash("Document trashed successfully", "success")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/archive/documents/%s", documentID), http.StatusFound)
}

func (server *Server) handleRestoreDocument(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	if documentID == "" {
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	session := server.getSession(r)
	err := server.archive.RestoreDocument(documentID, user)
	if err != nil {
		slog.Error("failed to restore document", slog.String("error", err.Error()))
		session.AddFlash("Failed to restore document", "error")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/archive/documents/%s", documentID), http.StatusFound)
		return
	}

	session.AddFlash("Document restored successfully", "success")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/archive/documents/%s", documentID), http.StatusFound)
}

func (server *Server) getDocumentPreview(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	pageNumberStr := chi.URLParam(r, "page")

	if documentID == "" || pageNumberStr == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	err = server.archive.GetDocumentPreview(documentID, user, pageNumber, func(r io.Reader) error {
		_, err := io.Copy(w, r)
		return err
	})

	if err != nil {
		slog.Error("failed to get document preview",
			slog.String("documentID", documentID),
			slog.String("user", user),
			slog.Int("page", pageNumber),
			slog.String("error", err.Error()))
		http.Error(w, "Preview not found", http.StatusNotFound)
	}
}

func (server *Server) getDocumentPreviewComponent(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	documentID := chi.URLParam(r, "id")
	pageNumberStr := chi.URLParam(r, "page")

	if documentID == "" || pageNumberStr == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	document, err := server.archive.GetDocument(documentID, user)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	if pageNumber < 0 || pageNumber >= len(document.PreviewFilepaths) {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	templates.DocumentPreviewComponent(document, pageNumber).Render(r.Context(), w)
}

func (server *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (server *Server) getSearch(w http.ResponseWriter, r *http.Request) {
	notifications := server.buildNotifications(r, w)
	isAdmin := server.isAdmin(r)

	page := templates.PageSearch

	// Start with empty results
	var results []archive.Document

	templates.Search(notifications, page, isAdmin, results).Render(r.Context(), w)
}

func (server *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	user := server.getAuthenticatedUser(r)
	query := r.URL.Query().Get("q")

	if query == "" {
		templates.EmptySearchResults().Render(r.Context(), w)
		return
	}

	hits, err := server.search.SearchDocuments(query, user, 20)
	if err != nil {
		templates.EmptySearchResults().Render(r.Context(), w)
		return
	}

	var documentIDs []string
	for _, hit := range hits {
		documentIDs = append(documentIDs, hit.DocumentID)
	}

	documents, err := server.archive.GetDocuments(documentIDs, user)
	if err != nil {
		log.Printf("Error getting documents: %v", err)
		templates.EmptySearchResults().Render(r.Context(), w)
		return
	}

	templates.SearchResults(documents).Render(r.Context(), w)
}

func (server *Server) profile(w http.ResponseWriter, r *http.Request) {
	username := server.getAuthenticatedUser(r)

	user, err := server.administration.GetUser(username)
	if err != nil {
		slog.Error("failed to get user", slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Profile Page</h1><p>Welcome, " + user.Username + "</p><p>This is a placeholder profile page.</p>"))
}

func (server *Server) admin(w http.ResponseWriter, r *http.Request) {
	notifications := server.buildNotifications(r, w)
	currentTab := r.URL.Query().Get("tab")
	if currentTab == "" {
		currentTab = "general"
	}

	settings, err := server.administration.Get()
	if err != nil {
		slog.Error("failed to get settings", slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}
	users, err := server.administration.GetAllUsers()
	if err != nil {
		slog.Error("failed to get users", slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	// Get page for tasks pagination
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	tasks, totalTasks, totalPages, err := server.administration.GetTasksPaginated(page)
	if err != nil {
		slog.Error("failed to get tasks", slog.String("error", err.Error()))
		templates.ErrorServer("").Render(r.Context(), w)
		return
	}

	// Check if there are any completed tasks across all pages
	hasCompletedTasks, err := server.administration.HasCompletedTasks()
	if err != nil {
		slog.Error("failed to check for completed tasks", slog.String("error", err.Error()))
		// Non-critical error, continue with hasCompletedTasks = false
		hasCompletedTasks = false
	}

	properties := templates.TaskTabProperties{
		Tasks:             tasks,
		CurrentPage:       page,
		TotalPages:        totalPages,
		TotalTasks:        totalTasks,
		HasCompletedTasks: hasCompletedTasks,
	}

	runtimeInfo := server.administration.GetRuntimeInfo()
	templates.Administration(notifications, currentTab, settings, users, properties, runtimeInfo).Render(r.Context(), w)
}

func (server *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	registrationEnabled := r.FormValue("registrationEnabled") == "true"

	err := server.administration.UpdateRegistrationEnabled(registrationEnabled)
	if err != nil {
		slog.Error("failed to update settings", slog.String("error", err.Error()))
		session.AddFlash("Failed to update settings", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/admin?tab=general", http.StatusFound)
		return
	}

	session.AddFlash("Settings updated successfully", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/admin?tab=general", http.StatusFound)
}

func (server *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	username := r.FormValue("username")
	if username == "" {
		session.AddFlash("username is required", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/admin?tab=users", http.StatusFound)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		session.AddFlash("password is required", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/admin?tab=users", http.StatusFound)
		return
	}

	err := server.administration.CreateUser(username, password, administration.UserRoleUser)
	if err != nil {
		session.AddFlash(err.Error(), "error")
		session.Save(r, w)
		http.Redirect(w, r, "/admin?tab=users", http.StatusFound)
		return
	}

	session.AddFlash("User created successfully", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/admin?tab=users", http.StatusFound)
}

func (server *Server) handleForceGC(w http.ResponseWriter, r *http.Request) {
	runtime.GC()

	session := server.getSession(r)
	session.AddFlash("Garbage collection triggered successfully", "success")
	session.Save(r, w)

	http.Redirect(w, r, "/admin?tab=runtime", http.StatusFound)
}

func (server *Server) handleClearCompletedTasks(w http.ResponseWriter, r *http.Request) {
	session := server.getSession(r)
	err := server.administration.ClearCompletedTasks()
	if err != nil {
		slog.Error("failed to clear completed tasks", slog.String("error", err.Error()))
		session.AddFlash("Failed to clear completed tasks", "error")
		session.Save(r, w)
		http.Redirect(w, r, "/admin?tab=tasks", http.StatusFound)
		return
	}

	session.AddFlash("Completed tasks cleared successfully", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/admin?tab=tasks", http.StatusFound)
}

func (server *Server) buildNotifications(r *http.Request, w http.ResponseWriter) []templates.Notification {
	var notifications []templates.Notification
	session := server.getSession(r)

	errMessages := session.Flashes("error")
	if len(errMessages) > 0 {
		for _, message := range errMessages {
			messageStr := message.(string)
			notifications = append(notifications, templates.Notification{
				Type:    templates.NotificationError,
				Message: messageStr,
			})
			slog.Error(messageStr)
		}
		delete(session.Values, "error")
	}

	warnMessages := session.Flashes("warning")
	if len(warnMessages) > 0 {
		for _, message := range warnMessages {
			notifications = append(notifications, templates.Notification{
				Type:    templates.NotificationWarning,
				Message: message.(string),
			})
		}
		delete(session.Values, "warning")
	}

	infoMessages := session.Flashes("info")
	if len(infoMessages) > 0 {
		for _, message := range infoMessages {
			notifications = append(notifications, templates.Notification{
				Type:    templates.NotificationInfo,
				Message: message.(string),
			})
		}
		delete(session.Values, "info")
	}

	successMessages := session.Flashes("success")
	if len(successMessages) > 0 {
		for _, message := range successMessages {
			notifications = append(notifications, templates.Notification{
				Type:    templates.NotificationSuccess,
				Message: message.(string),
			})
		}
		delete(session.Values, "success")
	}

	session.Save(r, w)
	return notifications
}

func (server *Server) createGenericErrorNotification() templates.Notification {
	return templates.Notification{
		Type:    templates.NotificationError,
		Message: "Something went wrong. Please try again later.",
	}
}

func (server *Server) requireSetup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.administration.AdminExists() {
			http.Redirect(w, r, "/setup", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (server *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := server.sessionStore.Get(r, "unterlagen")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.isAdmin(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) getSession(r *http.Request) *sessions.Session {
	session, err := server.sessionStore.Get(r, "unterlagen")
	if err != nil {
		return sessions.NewSession(server.sessionStore, "unterlagen")
	}
	return session
}

func (server *Server) getAuthenticatedUser(r *http.Request) string {
	session, err := server.sessionStore.Get(r, "unterlagen")
	if err != nil {
		panic(err)
	}

	username, ok := session.Values["username"]
	if !ok {
		panic("no username in session found")
	}

	return username.(string)
}

func (server *Server) isAdmin(r *http.Request) bool {
	session := server.getSession(r)
	return session.Values["role"].(administration.UserRole) == administration.UserRoleAdmin
}

func NewServer(
	administration *administration.Administration,
	archive *archive.Archive,
	search *search.Search,
	shutdown *common.Shutdown,
	configuration configuration.Configuration,
) *Server {
	sessionStore := sessions.NewCookieStore([]byte(configuration.Server.SessionKey))
	for _, role := range administration.UserRoles() {
		gob.Register(role)
	}

	server := &Server{
		administration: administration,
		archive:        archive,
		search:         search,
		sessionStore:   sessionStore,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/setup", server.setup)
	router.Post("/setup", server.handleSetup)

	router.Group(func(router chi.Router) {
		router.Use(server.requireSetup)
		router.Get("/login", server.login)
		router.Post("/login", server.handleLogin)

		router.Group(func(router chi.Router) {
			router.Use(server.requireLogin)
			router.Get("/", server.home)
			router.Get("/profile", server.profile)
			router.Post("/logout", server.handleLogout)
			router.Get("/archive", server.getArchive)
			router.Post("/archive/folders", server.handleCreateFolder)
			router.Post("/archive/synchronize", server.handleSynchronize)
			router.Post("/archive/documents", server.handleUploadDocument)
			router.Get("/archive/documents/{id}", server.getDocumentDetails)
			router.Get("/archive/documents/{id}/download", server.downloadDocument)
			router.Get("/archive/documents/{id}/previews/{page}", server.getDocumentPreview)
			router.Get("/archive/documents/{id}/preview-component/{page}", server.getDocumentPreviewComponent)
			router.Post("/archive/documents/{id}/delete", server.handleDeleteDocument)
			router.Post("/archive/documents/{id}/restore", server.handleRestoreDocument)
			router.Get("/search", server.getSearch)
			router.Get("/search/execute", server.handleSearch)

			router.Group(func(router chi.Router) {
				router.Use(server.requireAdmin)
				router.Get("/admin", server.admin)
				router.Post("/admin/settings", server.handleUpdateSettings)
				router.Post("/admin/users", server.handleCreateUser)
				router.Post("/admin/runtime/gc", server.handleForceGC)
				router.Post("/admin/tasks/clear-completed", server.handleClearCompletedTasks)
			})
		})
	})
	router.Get("/public/*", server.buildPublicHandler())

	server.internal = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	shutdown.AddCallback(func() {
		if err := server.internal.Shutdown(context.Background()); err != nil {
			slog.Error("server shutdown failed", slog.String("error", err.Error()))
		} else {
			slog.Info("shutdown server")
		}
	})

	return server
}
