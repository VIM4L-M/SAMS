import React, { useCallback, useRef } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  type Connection,
  type Node,
  type Edge,
  BackgroundVariant,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { nodeTypes } from './NodeTypes'
import { edgeTypes } from './EdgeTypes/DefaultEdge'
import { useSAMSStore } from '../../store'
import type { SAMSEdge } from '../../types'

// Generate a collision-proof node id even after loading a saved diagram.
function nextNodeId(existing: { id: string }[]): string {
  let max = 0
  for (const n of existing) {
    const m = /^node_(\d+)$/.exec(n.id)
    if (m) max = Math.max(max, Number(m[1]))
  }
  return `node_${max + 1}`
}

export const Canvas = React.memo(function Canvas() {
  // Store is the single source of truth (controlled mode).
  const nodes = useSAMSStore((s) => s.nodes)
  const edges = useSAMSStore((s) => s.edges)
  const onNodesChange = useSAMSStore((s) => s.onNodesChange)
  const onEdgesChange = useSAMSStore((s) => s.onEdgesChange)
  const addStoreEdge = useSAMSStore((s) => s.addEdge)
  const addStoreNode = useSAMSStore((s) => s.addNode)
  const selectNode = useSAMSStore((s) => s.selectNode)

  const reactFlowWrapper = useRef<HTMLDivElement>(null)
  const [reactFlowInstance, setReactFlowInstance] = React.useState<any>(null)

  const onConnect = useCallback(
    (params: Connection) => {
      const newEdge: SAMSEdge = {
        id: `edge_${Date.now()}`,
        source: params.source,
        target: params.target,
        data: { connectionType: 'sync', protocol: 'https' },
      }
      addStoreEdge(newEdge)
    },
    [addStoreEdge]
  )

  const onNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => selectNode(node.id),
    [selectNode]
  )
  const onPaneClick = useCallback(() => selectNode(null), [selectNode])

  const onDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
  }, [])

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault()
      if (!reactFlowInstance || !reactFlowWrapper.current) return
      const raw = event.dataTransfer.getData('application/sams-component')
      if (!raw) return
      let componentData: any
      try { componentData = JSON.parse(raw) } catch { return }

      const position = reactFlowInstance.screenToFlowPosition({
        x: event.clientX,
        y: event.clientY,
      })

      const id = nextNodeId(useSAMSStore.getState().nodes)
      const newNode = {
        id,
        type: componentData.type,
        position,
        data: {
          label: componentData.label,
          subtype: componentData.id,
          properties: { ...componentData.defaultProperties },
        },
      }
      addStoreNode(newNode)
      selectNode(id)
    },
    [reactFlowInstance, addStoreNode, selectNode]
  )

  const isEmpty = nodes.length === 0

  return (
    <div ref={reactFlowWrapper} className="absolute inset-0">
      <ReactFlow
        nodes={nodes as unknown as Node[]}
        edges={edges as unknown as Edge[]}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={onNodeClick}
        onPaneClick={onPaneClick}
        onInit={setReactFlowInstance}
        onDrop={onDrop}
        onDragOver={onDragOver}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        defaultEdgeOptions={{ type: 'default', data: { connectionType: 'sync', protocol: 'https' } }}
        fitView
        deleteKeyCode="Delete"
        colorMode="dark"
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={24}
          size={1}
          color="#2a2a3a"
        />
        <Controls />
        <MiniMap
          nodeColor={(n) => {
            const map: Record<string, string> = {
              frontend: '#3b82f6', backend: '#a855f7', microservice: '#6366f1',
              database: '#10b981', cache: '#eab308', queue: '#f97316',
              loadbalancer: '#14b8a6', apigateway: '#06b6d4', cdn: '#ec4899', storage: '#64748b',
              monitoring: '#f59e0b',
            }
            return map[n.type ?? ''] ?? '#3a3a50'
          }}
          maskColor="rgba(12,12,16,0.7)"
          style={{ background: '#13131a' }}
        />
      </ReactFlow>

      {isEmpty && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div className="text-center">
            <div className="text-5xl mb-4 opacity-30">🏗️</div>
            <p className="text-zinc-500 text-base font-medium">
              Drag components from the left panel
            </p>
            <p className="text-zinc-600 text-sm mt-1">
              to start designing your architecture
            </p>
          </div>
        </div>
      )}
    </div>
  )
})
