import axios from 'axios'
import type { SAMSNode, SAMSEdge, Context, ValidationResponse, ComponentCategory } from '../types'

const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
})

export async function validate(
  nodes: SAMSNode[],
  edges: SAMSEdge[],
  context: Context
): Promise<ValidationResponse> {
  const payload = {
    nodes: nodes.map((n) => ({
      id: n.id,
      type: n.type,
      data: {
        label: n.data.label,
        subtype: n.data.subtype,
        properties: n.data.properties,
      },
    })),
    edges: edges.map((e) => ({
      id: e.id,
      source: e.source,
      target: e.target,
      connectionType: e.data?.connectionType ?? 'sync',
      protocol: e.data?.protocol ?? 'https',
    })),
    context,
  }
  const res = await api.post<ValidationResponse>('/validate', payload)
  return res.data
}

export async function fetchComponents(): Promise<ComponentCategory[]> {
  const res = await api.get<{ components: ComponentCategory[] }>('/components')
  return res.data.components
}
