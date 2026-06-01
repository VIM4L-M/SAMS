import { create } from 'zustand'
import type { SAMSNode, SAMSEdge, Context, ValidationResponse, NodeData } from '../types'

const defaultContext: Context = {
  trafficLevel: 'medium',
  readWriteRatio: 'balanced',
  userBase: 'regional',
  teamSize: 'small',
  stage: 'early',
}

interface SAMSStore {
  nodes: SAMSNode[]
  edges: SAMSEdge[]
  selectedNodeId: string | null

  context: Context

  validationResults: ValidationResponse | null
  isValidating: boolean
  highlightedNodes: string[]
  highlightedEdges: string[]

  addNode: (node: SAMSNode) => void
  removeNode: (id: string) => void
  updateNode: (id: string, data: Partial<NodeData>) => void
  setNodes: (nodes: SAMSNode[]) => void
  setEdges: (edges: SAMSEdge[]) => void
  addEdge: (edge: SAMSEdge) => void
  removeEdge: (id: string) => void
  selectNode: (id: string | null) => void
  updateContext: (ctx: Partial<Context>) => void
  setValidationResults: (results: ValidationResponse) => void
  setIsValidating: (val: boolean) => void
  highlightIssue: (nodeIds: string[], edgeIds: string[]) => void
  clearHighlights: () => void
  resetCanvas: () => void
}

export const useSAMSStore = create<SAMSStore>((set) => ({
  nodes: [],
  edges: [],
  selectedNodeId: null,
  context: defaultContext,
  validationResults: null,
  isValidating: false,
  highlightedNodes: [],
  highlightedEdges: [],

  addNode: (node) => set((s) => ({ nodes: [...s.nodes, node] })),
  removeNode: (id) =>
    set((s) => ({
      nodes: s.nodes.filter((n) => n.id !== id),
      edges: s.edges.filter((e) => e.source !== id && e.target !== id),
      selectedNodeId: s.selectedNodeId === id ? null : s.selectedNodeId,
    })),
  updateNode: (id, data) =>
    set((s) => ({
      nodes: s.nodes.map((n) =>
        n.id === id ? { ...n, data: { ...n.data, ...data } } : n
      ),
    })),
  setNodes: (nodes) => set({ nodes }),
  setEdges: (edges) => set({ edges }),
  addEdge: (edge) => set((s) => ({ edges: [...s.edges, edge] })),
  removeEdge: (id) => set((s) => ({ edges: s.edges.filter((e) => e.id !== id) })),
  selectNode: (id) => set({ selectedNodeId: id }),
  updateContext: (ctx) => set((s) => ({ context: { ...s.context, ...ctx } })),
  setValidationResults: (results) => set({ validationResults: results }),
  setIsValidating: (val) => set({ isValidating: val }),
  highlightIssue: (nodeIds, edgeIds) =>
    set({ highlightedNodes: nodeIds, highlightedEdges: edgeIds }),
  clearHighlights: () => set({ highlightedNodes: [], highlightedEdges: [] }),
  resetCanvas: () =>
    set({
      nodes: [],
      edges: [],
      selectedNodeId: null,
      validationResults: null,
      highlightedNodes: [],
      highlightedEdges: [],
    }),
}))
