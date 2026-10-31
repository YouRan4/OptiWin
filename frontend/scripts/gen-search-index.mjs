import { readFileSync, readdirSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))

const routes = {
  'SecurityPage.vue': ['/security', 'nav.security'],
  'PerformancePage.vue': ['/performance', 'nav.performance'],
  'PersonalizationPage.vue': ['/personalization', 'nav.personalization'],
  'UtilitiesPage.vue': ['/utilities', 'nav.utilities'],
  'UpdatesPage.vue': ['/updates', 'nav.updates'],
}

const viewsDir = join(here, '..', 'src', 'views')
const locale = JSON.parse(readFileSync(join(here, '..', 'src', 'locales', 'zh.json'), 'utf8'))
const lookup = (key) => key.split('.').reduce((o, k) => (o == null ? undefined : o[k]), locale)

const entries = []
for (const file of readdirSync(viewsDir).sort()) {
  const route = routes[file]
  if (!route) continue
  const [path, page] = route
  const lines = readFileSync(join(viewsDir, file), 'utf8').split('\n')
  let group = ''
  const seen = new Set()
  lines.forEach((line, i) => {
    const title = line.match(/header-title"?>\{\{ i18n\('([\w.]+)'\)/)
    if (title) group = title[1]
    if (!line.includes('row-label')) return
    const inline = line.match(/i18n\('([\w.]+)'\)/)
    const next = inline ? null : (lines[i + 1] || '').match(/i18n\('([\w.]+)'\)/)
    const label = (inline || next)?.[1]
    if (!label || seen.has(label) || !group) return
    seen.add(label)
    entries.push({ path, page, group, label })
  })
}

const rows = entries.map(({ path, page, group, label }) => {
  const desc = `${label}Desc`
  const tail = lookup(desc) ? `, desc: "${desc}"` : ''
  return `  { path: "${path}", page: "${page}", group: "${group}", label: "${label}"${tail} },`
})

const out = [
  'export interface SearchEntry {',
  '  path: string',
  '  page: string',
  '  group: string',
  '  label: string',
  '  desc?: string',
  '}',
  '',
  'export const searchIndex: SearchEntry[] = [',
  ...rows,
  ']',
  '',
].join('\n')

writeFileSync(join(here, '..', 'src', 'searchIndex.ts'), out)
console.log(`searchIndex.ts: ${entries.length} entries`)
