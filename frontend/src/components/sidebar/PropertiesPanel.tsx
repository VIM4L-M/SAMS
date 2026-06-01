import { useSAMSStore } from '../../store'
import type { Context } from '../../types'

const NODE_PROPERTIES: Record<string, { key: string; label: string }[]> = {
  database: [
    { key: 'replication', label: 'Replication' },
    { key: 'backups', label: 'Backups' },
    { key: 'connectionPooling', label: 'Connection Pooling' },
    { key: 'encryptionAtRest', label: 'Encryption at Rest' },
    { key: 'sharding', label: 'Sharding' },
  ],
  cache: [
    { key: 'persistence', label: 'Persistence' },
    { key: 'clustering', label: 'Clustering' },
  ],
  backend: [
    { key: 'stateless', label: 'Stateless' },
    { key: 'healthChecks', label: 'Health Checks' },
    { key: 'circuitBreaker', label: 'Circuit Breaker' },
  ],
  microservice: [
    { key: 'stateless', label: 'Stateless' },
    { key: 'healthChecks', label: 'Health Checks' },
    { key: 'circuitBreaker', label: 'Circuit Breaker' },
  ],
  loadbalancer: [
    { key: 'healthChecks', label: 'Health Checks' },
    { key: 'sslTermination', label: 'SSL Termination' },
  ],
}

function Toggle({ checked, onChange }: { checked: boolean; onChange: () => void }) {
  return (
    <button
      onClick={onChange}
      className={`relative w-9 h-5 rounded-full transition-colors duration-200 focus:outline-none
        ${checked ? 'bg-blue-600' : 'bg-zinc-700'}`}
    >
      <span className={`absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform duration-200
        ${checked ? 'translate-x-4' : 'translate-x-0'}`} />
    </button>
  )
}

function Radio<T extends string>({
  label, options, value, onChange,
}: {
  label: string; options: { value: T; label: string }[]; value: T; onChange: (v: T) => void
}) {
  return (
    <div className="mb-4">
      <p className="text-[10px] font-semibold text-zinc-600 uppercase tracking-wider mb-2">{label}</p>
      <div className="flex flex-wrap gap-1.5">
        {options.map((o) => (
          <button
            key={o.value}
            onClick={() => onChange(o.value)}
            className={`text-xs px-2.5 py-1 rounded-md border transition-colors
              ${value === o.value
                ? 'bg-blue-600 text-white border-blue-600'
                : 'bg-transparent text-zinc-500 border-[#252535] hover:border-zinc-500 hover:text-zinc-300'
              }`}
          >
            {o.label}
          </button>
        ))}
      </div>
    </div>
  )
}

export function PropertiesPanel() {
  const selectedId = useSAMSStore((s) => s.selectedNodeId)
  const nodes = useSAMSStore((s) => s.nodes)
  const updateNode = useSAMSStore((s) => s.updateNode)
  const context = useSAMSStore((s) => s.context)
  const updateContext = useSAMSStore((s) => s.updateContext)

  const selectedNode = nodes.find((n) => n.id === selectedId)

  const toggleProperty = (key: string) => {
    if (!selectedNode) return
    updateNode(selectedNode.id, {
      properties: {
        ...(selectedNode.data.properties as Record<string, boolean>),
        [key]: !(selectedNode.data.properties as Record<string, boolean>)[key],
      },
    })
  }

  const nodeProps = selectedNode ? (NODE_PROPERTIES[selectedNode.type] ?? []) : []
  const relevantProps = nodeProps.filter((p) =>
    Object.prototype.hasOwnProperty.call(selectedNode?.data.properties ?? {}, p.key)
  )

  return (
    <div className="h-full flex flex-col bg-[#13131a] border-l border-[#1e1e2e]">
      <div className="px-4 py-3 border-b border-[#1e1e2e]">
        <h2 className="text-xs font-semibold text-zinc-300 tracking-wide">Properties</h2>
      </div>

      <div className="flex-1 overflow-y-auto px-4 py-4">
        {selectedNode ? (
          <>
            <div className="mb-5">
              <p className="text-[10px] font-semibold text-zinc-600 uppercase tracking-widest mb-3">
                {selectedNode.data.label as string}
              </p>
              {relevantProps.length > 0 ? (
                relevantProps.map((prop) => (
                  <div key={prop.key} className="flex items-center justify-between mb-3">
                    <span className="text-xs text-zinc-400">{prop.label}</span>
                    <Toggle
                      checked={!!(selectedNode.data.properties as Record<string, boolean>)[prop.key]}
                      onChange={() => toggleProperty(prop.key)}
                    />
                  </div>
                ))
              ) : (
                <p className="text-xs text-zinc-600">No configurable properties</p>
              )}
            </div>
            <div className="border-t border-[#1e1e2e] mb-5" />
          </>
        ) : (
          <p className="text-xs text-zinc-600 text-center mt-6 mb-6">
            Click a node to configure it
          </p>
        )}

        <p className="text-[10px] font-semibold text-zinc-600 uppercase tracking-widest mb-3">
          Context
        </p>

        <Radio label="Traffic Level" value={context.trafficLevel}
          onChange={(v) => updateContext({ trafficLevel: v as Context['trafficLevel'] })}
          options={[{ value: 'low', label: 'Low' }, { value: 'medium', label: 'Medium' },
            { value: 'high', label: 'High' }, { value: 'massive', label: 'Massive' }]} />
        <Radio label="Read / Write Ratio" value={context.readWriteRatio}
          onChange={(v) => updateContext({ readWriteRatio: v as Context['readWriteRatio'] })}
          options={[{ value: 'read_heavy', label: 'Read Heavy' }, { value: 'balanced', label: 'Balanced' },
            { value: 'write_heavy', label: 'Write Heavy' }]} />
        <Radio label="User Base" value={context.userBase}
          onChange={(v) => updateContext({ userBase: v as Context['userBase'] })}
          options={[{ value: 'local', label: 'Local' }, { value: 'regional', label: 'Regional' },
            { value: 'global', label: 'Global' }]} />
        <Radio label="Team Size" value={context.teamSize}
          onChange={(v) => updateContext({ teamSize: v as Context['teamSize'] })}
          options={[{ value: 'solo', label: 'Solo' }, { value: 'small', label: 'Small' },
            { value: 'medium', label: 'Medium' }, { value: 'large', label: 'Large' }]} />
        <Radio label="Stage" value={context.stage}
          onChange={(v) => updateContext({ stage: v as Context['stage'] })}
          options={[{ value: 'early', label: 'Early' }, { value: 'growing', label: 'Growing' },
            { value: 'scale', label: 'Scale' }]} />
      </div>
    </div>
  )
}
