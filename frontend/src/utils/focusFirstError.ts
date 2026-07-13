/**
 * focusFirstError — focus the first invalid field after a failed submit.
 *
 * Looks for elements marked invalid (data-invalid="true" or [aria-invalid="true"])
 * within the given root, or falls back to the first [data-error] anchor, and
 * focuses it. Returns true when an element was focused.
 *
 * Used by tool pages to keep keyboard users oriented after validation errors.
 */
export function focusFirstError(root: ParentNode = document): boolean {
  const selector = '[data-invalid="true"], [aria-invalid="true"], [data-error="true"]'
  const el = root.querySelector<HTMLElement>(selector)
  if (el) {
    el.focus()
    if (typeof el.scrollIntoView === 'function') {
      el.scrollIntoView({ block: 'nearest' })
    }
    return true
  }
  return false
}
