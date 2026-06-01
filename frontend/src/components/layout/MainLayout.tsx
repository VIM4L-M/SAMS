import { Header } from './Header'
import { ComponentLibrary } from '../sidebar/ComponentLibrary'
import { Canvas } from '../canvas/Canvas'
import { PropertiesPanel } from '../sidebar/PropertiesPanel'
import { ResultsPanel } from '../validation/ResultsPanel'

export function MainLayout() {
  return (
    <div className="h-screen flex flex-col overflow-hidden bg-[#0c0c10]">
      <Header />
      <div className="flex-1 flex overflow-hidden">
        <div className="w-[250px] flex-shrink-0 overflow-hidden">
          <ComponentLibrary />
        </div>
        <div className="flex-1 flex flex-col overflow-hidden">
          <div className="flex-1 relative overflow-hidden">
            <Canvas />
          </div>
          <div className="h-[280px] flex-shrink-0 overflow-hidden">
            <ResultsPanel />
          </div>
        </div>
        <div className="w-[300px] flex-shrink-0 overflow-hidden">
          <PropertiesPanel />
        </div>
      </div>
    </div>
  )
}
