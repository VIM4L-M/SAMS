import React from 'react'
import { BaseEdge, EdgeLabelRenderer, getBezierPath } from '@xyflow/react'
import { useSAMSStore } from '../../../store'

interface DefaultEdgeProps {
  id: string
  sourceX: number; sourceY: number; targetX: number; targetY: number
  sourcePosition: any; targetPosition: any
  data?: { connectionType?: string; protocol?: string }
}

export const DefaultEdge = React.memo(function DefaultEdge({
  id, sourceX, sourceY, targetX, targetY, sourcePosition, targetPosition, data,
}: DefaultEdgeProps) {
  const isHighlighted = useSAMSStore((s) => s.highlightedEdges.includes(id))

  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX, sourceY, sourcePosition, targetX, targetY, targetPosition,
  })

  const isAsync = data?.connectionType === 'async'
  const isHttp = data?.protocol === 'http'

  const stroke = isHighlighted ? '#ef4444' : isHttp ? '#f97316' : '#4a4a6a'

  return (
    <>
      <BaseEdge
        path={edgePath}
        style={{
          stroke,
          strokeWidth: isHighlighted ? 2.5 : 1.5,
          strokeDasharray: isAsync ? '6 3' : undefined,
        }}
        markerEnd="url(#arrow)"
      />
      {data?.protocol && (
        <EdgeLabelRenderer>
          <div
            style={{ transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)` }}
            className="absolute pointer-events-none"
          >
            <span className="text-[9px] bg-[#1e1e2e] border border-[#2a2a3a] text-zinc-500 px-1.5 py-0.5 rounded">
              {data.protocol}
            </span>
          </div>
        </EdgeLabelRenderer>
      )}
    </>
  )
})

export const edgeTypes = { default: DefaultEdge }
