export function setupEditor(
  textareaId: string,
  onEvaluate: (code: string) => void
) {
  const textarea = document.getElementById(textareaId) as HTMLTextAreaElement;
  if (!textarea || !(textarea instanceof HTMLTextAreaElement)) {
    throw new Error(`Textarea with id ${textareaId} not found`);
  }

  const STORAGE_KEY = `editor_content_${textareaId}`;

  // Load saved content
  const savedContent = localStorage.getItem(STORAGE_KEY);
  if (savedContent) {
    textarea.value = savedContent;
  }

  textarea.spellcheck = false;

  // Save content on change with debounce
  let saveTimeout: number | null = null;
  const saveContent = () => {
    if (saveTimeout) {
      window.clearTimeout(saveTimeout);
    }
    saveTimeout = window.setTimeout(() => {
      localStorage.setItem(STORAGE_KEY, textarea.value);
    }, 500); // Save 500ms after last change
  };

  textarea.addEventListener("input", saveContent);

  // Handle shift+enter and tab
  textarea.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter" && e.shiftKey) {
      e.preventDefault();
      const selectedText = textarea.value.substring(
        textarea.selectionStart,
        textarea.selectionEnd
      );
      if (selectedText) {
        // If there's selected text, evaluate that
        onEvaluate(selectedText.trim());
      } else {
        // Otherwise, get and evaluate the current line
        const value = textarea.value;
        const start = value.lastIndexOf("\n", textarea.selectionStart - 1) + 1;
        const end = value.indexOf("\n", textarea.selectionStart);
        const currentLine = value.substring(
          start,
          end === -1 ? value.length : end
        );
        onEvaluate(currentLine.trim());
      }
    }
    // Handle tab key
    if (e.key === "Tab") {
      e.preventDefault();
      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;
      // Insert tab character
      textarea.value =
        textarea.value.substring(0, start) +
        "\t" +
        textarea.value.substring(end);
      // Move cursor after tab
      textarea.selectionStart = textarea.selectionEnd = start + 1;
      // Save after tab insertion
      saveContent();
    }
  });

  // Clean up on page unload
  window.addEventListener("unload", () => {
    if (saveTimeout) {
      window.clearTimeout(saveTimeout);
      localStorage.setItem(STORAGE_KEY, textarea.value);
    }
  });
}
