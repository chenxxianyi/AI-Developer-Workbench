import { describe, expect, it } from 'vitest'
import { getLandingNavLinkClass } from './landingNav'

describe('landing navigation active state', () => {
  it('highlights home when the route has no section hash', () => {
    expect(getLandingNavLinkClass('', '')).toContain('text-accent')
    expect(getLandingNavLinkClass('', '#tools')).toContain('text-text-secondary')
  })

  it('highlights the matching section link when the route has a section hash', () => {
    expect(getLandingNavLinkClass('#tools', '#tools')).toContain('text-accent')
    expect(getLandingNavLinkClass('#workflow', '#workflow')).toContain('text-accent')
    expect(getLandingNavLinkClass('#features', '#features')).toContain('text-accent')
    expect(getLandingNavLinkClass('#tools', '#workflow')).toContain('text-text-secondary')
  })
})
