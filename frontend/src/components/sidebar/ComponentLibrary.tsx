import React, { useState, useEffect } from 'react'
import { fetchComponents } from '../../services/api'
import type { ComponentCategory, ComponentItem } from '../../types'

const categoryIcons: Record<string, string> = {
  frontend: '🖥️', backend: '⚙️', microservice: '🔧', database: '🗄️',
  cache: '⚡', queue: '📨', loadbalancer: '⚖️', apigateway: '🌐',
  cdn: '🚀', storage: '💾', network: '🌐', monitoring: '📊',
}

const accentColors: Record<string, string> = {
  frontend: 'hover:border-blue-500/50 hover:bg-blue-500/5',
  backend: 'hover:border-purple-500/50 hover:bg-purple-500/5',
  microservice: 'hover:border-indigo-500/50 hover:bg-indigo-500/5',
  database: 'hover:border-emerald-500/50 hover:bg-emerald-500/5',
  cache: 'hover:border-yellow-500/50 hover:bg-yellow-500/5',
  queue: 'hover:border-orange-500/50 hover:bg-orange-500/5',
  loadbalancer: 'hover:border-teal-500/50 hover:bg-teal-500/5',
  apigateway: 'hover:border-cyan-500/50 hover:bg-cyan-500/5',
  cdn: 'hover:border-pink-500/50 hover:bg-pink-500/5',
  storage: 'hover:border-slate-400/50 hover:bg-slate-500/5',
  monitoring: 'hover:border-amber-500/50 hover:bg-amber-500/5',
}

function DraggableComponent({ item }: { item: ComponentItem }) {
  const icon = categoryIcons[item.type] ?? '📦'
  const accent = accentColors[item.type] ?? 'hover:border-zinc-500 hover:bg-zinc-800'

  const onDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData('application/sams-component', JSON.stringify(item))
    e.dataTransfer.effectAllowed = 'move'
  }

  return (
    <div
      draggable
      onDragStart={onDragStart}
      className={`
        flex items-center gap-2.5 px-3 py-2 rounded-md
        bg-[#1a1a24] border border-[#252535]
        cursor-grab active:cursor-grabbing
        transition-all duration-150 select-none
        ${accent}
      `}
    >
      <span className="text-sm">{icon}</span>
      <span className="text-xs font-medium text-zinc-300">{item.label}</span>
    </div>
  )
}

function CategorySection({ category, defaultOpen = true }: {
  category: ComponentCategory; defaultOpen?: boolean
}) {
  const [open, setOpen] = useState(defaultOpen)
  return (
    <div className="mb-3">
      <button
        onClick={() => setOpen((o) => !o)}
        className="w-full flex items-center justify-between px-1 py-1.5 text-[10px] font-semibold text-zinc-600 uppercase tracking-widest hover:text-zinc-400 transition-colors"
      >
        <span className="flex items-center gap-1.5">
          <span>{category.categoryIcon}</span>
          {category.categoryLabel}
        </span>
        <span>{open ? '▾' : '▸'}</span>
      </button>
      {open && (
        <div className="flex flex-col gap-1.5 mt-1">
          {category.items.map((item) => (
            <DraggableComponent key={item.id} item={item} />
          ))}
        </div>
      )}
    </div>
  )
}

export function ComponentLibrary() {
  const [categories, setCategories] = useState<ComponentCategory[]>([])

  useEffect(() => {
    fetchComponents().then(setCategories).catch(() => {})
  }, [])

  return (
    <div className="h-full flex flex-col bg-[#13131a] border-r border-[#1e1e2e]">
      <div className="px-4 py-3 border-b border-[#1e1e2e]">
        <h2 className="text-xs font-semibold text-zinc-300 tracking-wide">Components</h2>
        <p className="text-[10px] text-zinc-600 mt-0.5">Drag onto canvas</p>
      </div>
      <div className="flex-1 overflow-y-auto px-3 py-3">
        {categories.length === 0 ? (
          <p className="text-xs text-zinc-600 text-center mt-8">Loading…</p>
        ) : (
          categories.map((cat, i) => (
            <CategorySection key={cat.category} category={cat} defaultOpen={i < 3} />
          ))
        )}
      </div>
    </div>
  )
}
