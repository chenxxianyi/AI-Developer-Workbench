/**
 * Clipboard Utilities
 * Copy text to clipboard
 */

/**
 * Copy text to clipboard
 * @param text Text to copy
 * @returns Promise that resolves when copy is complete
 */
export async function copyToClipboard(text: string): Promise<void> {
  try {
    // Modern API (async)
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text)
      return
    }

    // Fallback for older browsers
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)

    textarea.select()
    textarea.setSelectionRange(0, textarea.value.length)

    const success = document.execCommand('copy')
    document.body.removeChild(textarea)

    if (!success) {
      throw new Error('Failed to copy text')
    }
  } catch (error) {
    console.error('Clipboard copy failed:', error)
    throw new Error('复制失败，请手动复制')
  }
}