// @vitest-environment jsdom

import { describe, expect, it } from 'vitest'
import { focusFirstError } from './focusFirstError'

function mount(html: string): HTMLElement {
  document.body.innerHTML = ''
  const root = document.createElement('div')
  root.innerHTML = html
  document.body.appendChild(root)
  return root
}

describe('focusFirstError', () => {
  it('focuses the first element marked data-invalid', () => {
    const root = mount(`
      <input id="a" />
      <input id="b" data-invalid="true" />
      <input id="c" aria-invalid="true" />
    `)

    const focused = focusFirstError(root)

    expect(focused).toBe(true)
    expect(document.activeElement).toBe(root.querySelector('#b'))
  })

  it('focuses the first aria-invalid element when no data-invalid is present', () => {
    const root = mount(`
      <input id="a" />
      <input id="b" aria-invalid="true" />
    `)

    focusFirstError(root)

    expect(document.activeElement).toBe(root.querySelector('#b'))
  })

  it('returns false when there are no invalid elements', () => {
    const root = mount(`<input id="a" /><input id="b" />`)

    expect(focusFirstError(root)).toBe(false)
  })

  it('focuses a [data-error="true"] anchor as a fallback', () => {
    const root = mount(`
      <div id="err" data-error="true" tabindex="-1"></div>
    `)

    focusFirstError(root)

    expect(document.activeElement).toBe(root.querySelector('#err'))
  })
})
