let currentPage = 0;
let totalPages = 0;
let documentId = "";

function initPreview(docId, pages) {
  documentId = docId;
  totalPages = pages;
  currentPage = 0;
  updatePreview();
}

function nextPage() {
  if (currentPage < totalPages - 1) {
    currentPage++;
    updatePreview();
  }
}

function previousPage() {
  if (currentPage > 0) {
    currentPage--;
    updatePreview();
  }
}

function updatePreview() {
  const img = document.getElementById("previewImage");
  const pageInfo = document.getElementById("pageInfo");
  const prevBtn = document.getElementById("prevBtn");
  const nextBtn = document.getElementById("nextBtn");

  if (img && pageInfo) {
    img.src = `/archive/documents/${documentId}/previews/${currentPage}`;
    pageInfo.textContent = `Page ${currentPage + 1} of ${totalPages}`;
  }

  if (prevBtn) prevBtn.disabled = currentPage === 0;
  if (nextBtn) nextBtn.disabled = currentPage === totalPages - 1;
}

// Initialize when DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  const previewSection = document.getElementById("previewImage");
  if (previewSection) {
    const docId = previewSection.getAttribute("data-doc-id");
    const pages = parseInt(previewSection.getAttribute("data-total-pages"));
    initPreview(docId, pages);

    // Add event listeners for navigation buttons
    const prevBtn = document.getElementById("prevBtn");
    const nextBtn = document.getElementById("nextBtn");

    if (prevBtn) {
      prevBtn.addEventListener("click", previousPage);
    }

    if (nextBtn) {
      nextBtn.addEventListener("click", nextPage);
    }
  }
});
