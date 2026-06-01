import React from 'react'
import { Handle, Position } from '@xyflow/react'
import { useSAMSStore } from '../../../store'

interface BaseNodeProps {
  id: string
  data: {
    label: string
    subtype: string
    properties: Record<string, boolean>
  }
  icon: string
  iconBg: string       // e.g. "bg-blue-500/20"
  accentBorder: string // e.g. "border-l-blue-500"
}

export const BaseNode = React.memo(function BaseNode({
  id,
  data,
  icon,
  iconBg,
  accentBorder,
}: BaseNodeProps) {
  const isHighlighted = useSAMSStore((s) => s.highlightedNodes.includes(id))

  const warningCount = useSAMSStore((s) => {
    if (!s.validationResults) return 0
    let count = 0
    for (const issue of [
      ...(s.validationResults.results.errors ?? []),
      ...(s.validationResults.results.warnings ?? []),
    ]) {
      if (issue.affectedNodes?.includes(id)) count++
    }
    return count
  })

  const label = typeof data.label === 'string' ? data.label : ''
  const subtype = typeof data.subtype === 'string' ? data.subtype : ''

  return (
    <div
      className={`
        relative flex items-center gap-2.5 pl-1 pr-3 py-2.5 rounded-lg
        bg-[#1e1e2e] border border-[#2a2a3a] border-l-4 ${accentBorder}
        shadow-lg min-w-[130px] cursor-pointer select-none
        transition-all duration-150
        ${isHighlighted
          ? '!border-red-500 ring-2 ring-red-500/30 shadow-red-900/30'
          : 'hover:border-[#3a3a50] hover:shadow-xl'
        }
      `}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="!w-2.5 !h-2.5 !bg-[#3a3a50] !border-2 !border-[#1e1e2e] !-left-1.5"
      />
      <Handle
        type="source"
        position={Position.Right}
        className="!w-2.5 !h-2.5 !bg-[#3a3a50] !border-2 !border-[#1e1e2e] !-right-1.5"
      />

      {warningCount > 0 && (
        <span className="absolute -top-2 -right-2 bg-red-500 text-white text-[10px] font-bold rounded-full w-4.5 h-4.5 min-w-[18px] min-h-[18px] flex items-center justify-center z-10 shadow-lg">
          {warningCount}
        </span>
      )}

      <div className={`${iconBg} rounded-md w-8 h-8 flex items-center justify-center text-base flex-shrink-0`}>
        {icon}
      </div>

      <div className="flex flex-col min-w-0">
        <span className="text-xs font-semibold text-zinc-200 leading-tight truncate">
          {label}
        </span>
        {subtype && subtype !== label.toLowerCase().replace(/\s+/g, '-') && (
          <span className="text-[10px] text-zinc-500 leading-tight truncate mt-0.5">
            {subtype}
          </span>
        )}
      </div>
    </div>
  )
})
