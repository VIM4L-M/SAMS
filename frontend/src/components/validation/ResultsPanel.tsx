import { useState } from 'react'
import type { MouseEvent } from 'react'
import { useSAMSStore } from '../../store'
import type { ValidationIssue } from '../../types'

const categoryColors: Record<string, string> = {
  security: 'text-red-400',
  performance: 'text-orange-400',
  scalability: 'text-blue-400',
  reliability: 'text-purple-400',
}

const categoryBg: Record<string, string> = {
  security: 'bg-red-500/5 border-red-500/20',
  performance: 'bg-orange-500/5 border-orange-500/20',
  scalability: 'bg-blue-500/5 border-blue-500/20',
  reliability: 'bg-purple-500/5 border-purple-500/20',
}

function IssueRow({ issue, isError }: { issue: ValidationIssue; isError: boolean }) {
  const highlight = useSAMSStore((s) => s.highlightIssue)
  const clearHighlights = useSAMSStore((s) => s.clearHighlights)
  const highlightedNodes = useSAMSStore((s) => s.highlightedNodes)
  const [expanded, setExpanded] = useState(false)

  const isPinned =
    (issue.affectedNodes ?? []).length > 0 &&
    (issue.affectedNodes ?? []).every((id) => highlightedNodes.includes(id))

  const togglePin = (e: MouseEvent) => {
    e.stopPropagation()
    isPinned ? clearHighlights() : highlight(issue.affectedNodes ?? [], issue.affectedEdges ?? [])
  }

  return (
    <div
      onClick={() => setExpanded((e) => !e)}
      className={`border rounded-md mb-1.5 overflow-hidden cursor-pointer transition-all
        ${isPinned ? 'ring-1 ring-red-500/50' : ''}
        ${categoryBg[issue.category]}`}
    >
      <div className="flex items-start gap-2 px-3 py-2">
        <span className="mt-0.5 flex-shrink-0 text-xs">{isError ? '🔴' : '🟡'}</span>
        <div className="flex-1 min-w-0">
          <p className="text-xs font-medium text-zinc-300 leading-snug">{issue.title}</p>
          {expanded && (
            <>
              <p className="text-xs text-zinc-500 mt-1.5 leading-relaxed">{issue.description}</p>
              <p className="text-xs text-emerald-500 mt-1.5">💡 {issue.suggestion}</p>
              {(issue.affectedNodes ?? []).length > 0 && (
                <button
                  onClick={togglePin}
                  className={`mt-2 text-[10px] px-2 py-0.5 rounded border transition-colors
                    ${isPinned
                      ? 'bg-red-500/20 border-red-500/50 text-red-400'
                      : 'bg-transparent border-[#252535] text-zinc-500 hover:border-zinc-500 hover:text-zinc-300'
                    }`}
                >
                  {isPinned ? '✕ Clear highlight' : '⬡ Highlight on canvas'}
                </button>
              )}
            </>
          )}
        </div>
        <span className="text-zinc-700 text-xs flex-shrink-0">{expanded ? '▲' : '▼'}</span>
      </div>
    </div>
  )
}

function CategoryGroup({ label, category, errors, warnings }: {
  label: string; category: string; errors: ValidationIssue[]; warnings: ValidationIssue[]
}) {
  const [open, setOpen] = useState(true)
  const allPassed = errors.length === 0 && warnings.length === 0

  if (allPassed) {
    return (
      <div className="flex items-center gap-2 px-3 py-1.5 rounded-md bg-emerald-500/5 border border-emerald-500/20 mb-1.5">
        <span className="text-xs">✅</span>
        <span className={`text-xs font-medium ${categoryColors[category]}`}>{label}</span>
        <span className="text-xs text-emerald-600 ml-auto">All passed</span>
      </div>
    )
  }

  return (
    <div className="mb-2">
      <button
        onClick={() => setOpen((o) => !o)}
        className="w-full flex items-center gap-2 px-1 py-1 hover:bg-white/5 rounded transition-colors"
      >
        <span className={`text-[10px] font-semibold uppercase tracking-wider ${categoryColors[category]}`}>
          {label}
        </span>
        <span className="text-xs text-zinc-600 ml-auto">
          {errors.length > 0 && `${errors.length} error${errors.length > 1 ? 's' : ''}`}
          {errors.length > 0 && warnings.length > 0 && ' · '}
          {warnings.length > 0 && `${warnings.length} warning${warnings.length > 1 ? 's' : ''}`}
        </span>
        <span className="text-zinc-700 text-xs">{open ? '▾' : '▸'}</span>
      </button>
      {open && (
        <div className="mt-1">
          {errors.map((i) => <IssueRow key={i.ruleId} issue={i} isError />)}
          {warnings.map((i) => <IssueRow key={i.ruleId} issue={i} isError={false} />)}
        </div>
      )}
    </div>
  )
}

export function ResultsPanel() {
  const results = useSAMSStore((s) => s.validationResults)
  const isValidating = useSAMSStore((s) => s.isValidating)
  const nodes = useSAMSStore((s) => s.nodes)

  const empty = (
    <div className="h-full flex items-center justify-center bg-[#0f0f14] border-t border-[#1e1e2e]">
      <p className="text-xs text-zinc-600">
        {nodes.length === 0
          ? 'Draw your architecture and we\'ll validate it automatically'
          : isValidating && !results
            ? 'Validating…'
            : 'Click Validate to run analysis'}
      </p>
    </div>
  )

  if (!results) return empty

  const errors = results.results.errors ?? []
  const warnings = results.results.warnings ?? []
  const categories = ['security', 'performance', 'scalability', 'reliability'] as const

  return (
    <div className="h-full flex flex-col bg-[#0f0f14] border-t border-[#1e1e2e]">
      <div className="px-4 py-2 border-b border-[#1e1e2e] flex items-center justify-between flex-shrink-0">
        <span className="text-[10px] font-semibold text-zinc-600 uppercase tracking-wider">
          Validation Results
          {isValidating && <span className="ml-2 text-blue-500 animate-pulse">updating…</span>}
        </span>
        <span className="text-[10px] text-zinc-700">
          {errors.length + warnings.length} issue{errors.length + warnings.length !== 1 ? 's' : ''}
          {' · '}{results.metadata.rulesChecked} rules
        </span>
      </div>
      <div className="flex-1 overflow-y-auto px-4 py-3">
        {categories.map((cat) => (
          <CategoryGroup
            key={cat}
            label={cat.charAt(0).toUpperCase() + cat.slice(1)}
            category={cat}
            errors={errors.filter((e) => e.category === cat)}
            warnings={warnings.filter((w) => w.category === cat)}
          />
        ))}
      </div>
    </div>
  )
}
