import "htmx.org";

// Search functionality
document.addEventListener("DOMContentLoaded", function () {
  const searchInput = document.getElementById("search-input");
  const searchResults = document.getElementById("search-results");

  if (searchInput && searchResults) {
    // Show search results when input gets focus and has content
    searchInput.addEventListener("focus", function () {
      if (searchInput.value.trim() !== "") {
        searchResults.classList.remove("hidden");
      }
    });

    // Hide search results when clicking outside
    document.addEventListener("click", function (event) {
      if (
        !searchInput.contains(event.target) &&
        !searchResults.contains(event.target)
      ) {
        searchResults.classList.add("hidden");
      }
    });

    // Show results after HTMX request completes
    searchInput.addEventListener("htmx:afterRequest", function (event) {
      if (searchInput.value.trim() !== "") {
        searchResults.classList.remove("hidden");
      } else {
        searchResults.classList.add("hidden");
      }
    });

    // Clear results when input is empty
    searchInput.addEventListener("input", function () {
      if (searchInput.value.trim() === "") {
        searchResults.classList.add("hidden");
      }
    });

    // Hide results when a search result is clicked
    searchResults.addEventListener("click", function () {
      searchResults.classList.add("hidden");
      searchInput.blur();
    });
  }
});
