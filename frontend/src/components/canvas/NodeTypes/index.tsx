import React from 'react'
import type { NodeProps } from '@xyflow/react'
import { BaseNode } from './BaseNode'
import type { NodeData } from '../../../types'

type SAMSNodeProps = NodeProps & { data: NodeData }

const FrontendNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="🖥️" iconBg="bg-blue-500/20" accentBorder="border-l-blue-500" />
))
const BackendNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="⚙️" iconBg="bg-purple-500/20" accentBorder="border-l-purple-500" />
))
const MicroserviceNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="🔧" iconBg="bg-indigo-500/20" accentBorder="border-l-indigo-500" />
))
const DatabaseNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="🗄️" iconBg="bg-emerald-500/20" accentBorder="border-l-emerald-500" />
))
const CacheNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="⚡" iconBg="bg-yellow-500/20" accentBorder="border-l-yellow-500" />
))
const QueueNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="📨" iconBg="bg-orange-500/20" accentBorder="border-l-orange-500" />
))
const LoadBalancerNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="⚖️" iconBg="bg-teal-500/20" accentBorder="border-l-teal-500" />
))
const ApiGatewayNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="🌐" iconBg="bg-cyan-500/20" accentBorder="border-l-cyan-500" />
))
const CdnNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="🚀" iconBg="bg-pink-500/20" accentBorder="border-l-pink-500" />
))
const StorageNode = React.memo((p: SAMSNodeProps) => (
  <BaseNode {...p} icon="💾" iconBg="bg-slate-500/20" accentBorder="border-l-slate-400" />
))

export const nodeTypes = {
  frontend: FrontendNode,
  backend: BackendNode,
  microservice: MicroserviceNode,
  database: DatabaseNode,
  cache: CacheNode,
  queue: QueueNode,
  loadbalancer: LoadBalancerNode,
  apigateway: ApiGatewayNode,
  cdn: CdnNode,
  storage: StorageNode,
}
