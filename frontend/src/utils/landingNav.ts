const baseClass = 'inline-flex min-h-10 items-center px-4 py-2 rounded-lg text-base font-bold transition-smooth cursor-pointer'
const inactiveClass = 'text-text-secondary hover:text-accent hover:bg-surface'
const activeClass = 'text-accent bg-accent-soft shadow-sm'

export function getLandingNavLinkClass(currentHash: string, itemHash: string): string {
  const normalizedCurrentHash = currentHash || ''
  const isActive = normalizedCurrentHash === itemHash

  return `${baseClass} ${isActive ? activeClass : inactiveClass}`
}
