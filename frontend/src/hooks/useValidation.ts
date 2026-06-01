import { useEffect, useRef } from 'react'
import { useSAMSStore } from '../store'
import { validate } from '../services/api'
import { useDebounce } from './useDebounce'

export function useValidation() {
  const nodes = useSAMSStore((s) => s.nodes)
  const edges = useSAMSStore((s) => s.edges)
  const context = useSAMSStore((s) => s.context)
  const setResults = useSAMSStore((s) => s.setValidationResults)
  const setValidating = useSAMSStore((s) => s.setIsValidating)

  // Only serialize fields that affect validation — exclude UI state (warningCount, isHighlighted)
  // to prevent syncWarningBadges from triggering an infinite re-validation loop.
  const snapshot = useDebounce(
    JSON.stringify({
      nodes: nodes.map((n) => ({
        id: n.id,
        type: n.type,
        data: { label: n.data.label, subtype: n.data.subtype, properties: n.data.properties },
      })),
      edges: edges.map((e) => ({ id: e.id, source: e.source, target: e.target, data: e.data })),
      context,
    }),
    1000
  )

  const runningRef = useRef(false)

  const runValidation = async (n = nodes, e = edges, c = context) => {
    if (n.length === 0) return
    if (runningRef.current) return
    runningRef.current = true
    setValidating(true)
    try {
      const result = await validate(n, e, c)
      setResults(result)
    } catch {
      // Silently ignore network errors — canvas keeps working
    } finally {
      setValidating(false)
      runningRef.current = false
    }
  }

  // Debounced auto-validation
  useEffect(() => {
    const parsed = JSON.parse(snapshot)
    if (parsed.nodes.length > 0) {
      runValidation(parsed.nodes, parsed.edges, parsed.context)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [snapshot])

  return { runValidation: () => runValidation() }
}
