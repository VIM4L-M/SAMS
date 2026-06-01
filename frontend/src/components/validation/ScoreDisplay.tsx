function scoreColor(score: number) {
  if (score <= 40) return { text: 'text-red-400', bg: 'bg-red-500/10 border-red-500/30', label: 'Needs Work' }
  if (score <= 70) return { text: 'text-yellow-400', bg: 'bg-yellow-500/10 border-yellow-500/30', label: 'Improving' }
  if (score <= 90) return { text: 'text-emerald-400', bg: 'bg-emerald-500/10 border-emerald-500/30', label: 'Good' }
  return { text: 'text-blue-400', bg: 'bg-blue-500/10 border-blue-500/30', label: 'Production Ready' }
}

export function ScoreDisplay({ score }: { score: number }) {
  const { text, bg, label } = scoreColor(score)
  return (
    <div className={`flex items-center gap-2 px-3 py-1.5 rounded-md border ${bg}`}>
      <span className={`text-base font-bold ${text}`}>{score}</span>
      <div className="flex flex-col">
        <span className="text-[9px] text-zinc-600 leading-none">Score</span>
        <span className={`text-[10px] font-semibold ${text} leading-tight`}>{label}</span>
      </div>
    </div>
  )
}
