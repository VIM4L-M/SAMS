import { useSAMSStore } from '../../store'
import { ScoreDisplay } from '../validation/ScoreDisplay'
import { useValidation } from '../../hooks/useValidation'

export function Header() {
  const results = useSAMSStore((s) => s.validationResults)
  const isValidating = useSAMSStore((s) => s.isValidating)
  const resetCanvas = useSAMSStore((s) => s.resetCanvas)
  const { runValidation } = useValidation()
  const nodes = useSAMSStore((s) => s.nodes)

  const score = results?.results.score ?? null

  return (
    <header className="h-13 flex items-center justify-between px-5 bg-[#13131a] border-b border-[#1e1e2e] flex-shrink-0 z-10">
      <div className="flex items-center gap-2.5">
        <span className="text-lg">🏗️</span>
        <div>
          <h1 className="text-sm font-bold text-zinc-100 leading-none tracking-wide">SAMS</h1>
          <p className="text-[10px] text-zinc-600 leading-none mt-0.5">Architecture Validator</p>
        </div>
      </div>

      <div className="flex items-center gap-2.5">
        {score !== null && <ScoreDisplay score={score} />}

        <button
          onClick={resetCanvas}
          className="text-xs px-3 py-1.5 rounded-md border border-[#2a2a3a] text-zinc-400 hover:bg-[#1e1e2e] hover:text-zinc-200 transition-colors"
        >
          New
        </button>

        <button disabled title="Coming in V2"
          className="text-xs px-3 py-1.5 rounded-md border border-[#1e1e2e] text-zinc-700 cursor-not-allowed">
          Save
        </button>

        <button disabled title="Coming in V2"
          className="text-xs px-3 py-1.5 rounded-md border border-[#1e1e2e] text-zinc-700 cursor-not-allowed">
          Load
        </button>

        <button
          onClick={runValidation}
          disabled={isValidating || nodes.length === 0}
          className={`text-xs px-4 py-1.5 rounded-md font-semibold transition-all
            ${nodes.length === 0 || isValidating
              ? 'bg-[#1e1e2e] text-zinc-600 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-500 shadow-lg shadow-blue-900/30'
            }`}
        >
          {isValidating ? 'Validating…' : 'Validate'}
        </button>
      </div>
    </header>
  )
}
