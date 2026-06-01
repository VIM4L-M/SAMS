export interface NodeData {
  label: string
  subtype: string
  icon?: string
  properties: Record<string, boolean>
  [key: string]: unknown // required for React Flow generic compatibility
}

export interface SAMSNode {
  id: string
  type: string
  position: { x: number; y: number }
  data: NodeData
}

export interface SAMSEdge {
  id: string
  source: string
  target: string
  data?: {
    connectionType: 'sync' | 'async'
    protocol: 'http' | 'https' | 'grpc' | 'amqp'
  }
}

export interface Context {
  trafficLevel: 'low' | 'medium' | 'high' | 'massive'
  readWriteRatio: 'read_heavy' | 'balanced' | 'write_heavy'
  userBase: 'local' | 'regional' | 'global'
  teamSize: 'solo' | 'small' | 'medium' | 'large'
  stage: 'early' | 'growing' | 'scale'
}

export interface ValidationIssue {
  ruleId: string
  category: 'security' | 'performance' | 'scalability' | 'reliability'
  title: string
  description: string
  affectedNodes: string[]
  affectedEdges: string[]
  suggestion: string
}

export interface PassedRule {
  ruleId: string
  category: string
  title: string
}

export interface ValidationResponse {
  results: {
    errors: ValidationIssue[]
    warnings: ValidationIssue[]
    passed: PassedRule[]
    score: number
  }
  metadata: {
    rulesChecked: number
    timeMs: number
  }
}

export interface ComponentItem {
  id: string
  type: string
  label: string
  icon: string
  defaultProperties: Record<string, boolean>
}

export interface ComponentCategory {
  category: string
  categoryLabel: string
  categoryIcon: string
  items: ComponentItem[]
}
