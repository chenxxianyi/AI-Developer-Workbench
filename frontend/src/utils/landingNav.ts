const baseClass = 'px-3 py-2 rounded-lg text-sm font-semibold transition-smooth cursor-pointer'
const inactiveClass = 'text-text-secondary hover:text-accent hover:bg-surface-muted'
const activeClass = 'text-accent bg-accent-soft font-semibold'

export function getLandingNavLinkClass(currentHash: string, itemHash: string): string {
  const normalizedCurrentHash = currentHash || ''
  const isActive = normalizedCurrentHash === itemHash

  return `${baseClass} ${isActive ? activeClass : inactiveClass}`
}
