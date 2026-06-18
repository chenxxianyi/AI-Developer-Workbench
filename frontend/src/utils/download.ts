/**
 * Download Utilities
 * Trigger browser downloads
 */

/**
 * Download a Blob as a file
 * Creates a temporary Object URL and triggers download
 * @param blob Blob to download
 * @param filename Filename for the download
 */
export async function downloadBlob(blob: Blob, filename: string): Promise<void> {
  const url = URL.createObjectURL(blob)

  try {
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.style.display = 'none'

    document.body.appendChild(a)
    a.click()

    // Wait for click to complete
    await new Promise((resolve) => setTimeout(resolve, 100))

    document.body.removeChild(a)
  } finally {
    // Always release the Object URL to prevent memory leaks
    URL.revokeObjectURL(url)
  }
}

/**
 * Download from a URL
 * @param url URL to download from
 * @param filename Optional filename (will use URL filename if not provided)
 */
export function downloadUrl(url: string, filename?: string): void {
  const a = document.createElement('a')
  a.href = url
  a.download = filename || url.split('/').pop() || 'download'
  a.style.display = 'none'

  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}