"use client"

import { useEffect, useState } from "react"
import { fetchStatus, fetchHistory, fetchConfig, ServiceStatus, ReconcileEvent } from "@/lib/api"

export default function Home() {
  const [services, setServices] = useState<ServiceStatus[]>([])
  const [history, setHistory] = useState<ReconcileEvent[]>([])
  const [config, setConfig] = useState<string>("")
  const [lastSync, setLastSync] = useState<string>("never")

  const refresh = async () => {
    const [s, h, c] = await Promise.all([fetchStatus(), fetchHistory(), fetchConfig()])
    setServices(s || [])
    setHistory(h || [])
    setConfig(c?.config || "")
    setLastSync(new Date().toLocaleTimeString())
  }

  useEffect(() => {
    refresh()
    const interval = setInterval(refresh, 10000)
    return () => clearInterval(interval)
  }, [])

  const totalPods = services.reduce((a, s) => a + s.ready_replicas, 0)
  const healthy = services.filter(s => s.status === "healthy").length
  const degraded = services.filter(s => s.status === "degraded").length

  return (
    <main className="min-h-screen bg-[#1E1E2E] text-[#CDD6F4] p-6 font-mono">
      
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#CBA6F7]">FlowOps</h1>
          <p className="text-[#A6ADC8] text-sm">GitOps Deployment Engine</p>
        </div>
        <div className="text-right">
          <p className="text-[#A6ADC8] text-xs">last sync</p>
          <p className="text-[#A6E3A1] text-sm font-bold">{lastSync}</p>
        </div>
      </div>

      {/* Stats bar */}
      <div className="grid grid-cols-4 gap-4 mb-6">
        {[
          { label: "Total Pods", value: totalPods, color: "text-[#89B4FA]" },
          { label: "Healthy", value: healthy, color: "text-[#A6E3A1]" },
          { label: "Degraded", value: degraded, color: "text-[#F38BA8]" },
          { label: "Services", value: services.length, color: "text-[#CBA6F7]" },
        ].map(stat => (
          <div key={stat.label} className="bg-[#181825] rounded-xl p-4 border border-[#45475A]">
            <p className="text-[#A6ADC8] text-xs mb-1">{stat.label}</p>
            <p className={`text-3xl font-bold ${stat.color}`}>{stat.value}</p>
          </div>
        ))}
      </div>

      {/* Services grid */}
      <h2 className="text-[#89B4FA] font-bold mb-3 text-sm uppercase tracking-widest">Services</h2>
      <div className="grid grid-cols-5 gap-3 mb-6">
        {services.map(svc => (
          <div key={svc.name} className="bg-[#181825] rounded-xl p-4 border border-[#45475A]">
            <div className="flex items-center justify-between mb-2">
              <p className="text-[#CBA6F7] font-bold text-sm">{svc.name}</p>
              <span className={`text-xs px-2 py-0.5 rounded-full font-bold ${
                svc.status === "healthy"
                  ? "bg-[#A6E3A1]/20 text-[#A6E3A1]"
                  : "bg-[#F38BA8]/20 text-[#F38BA8]"
              }`}>
                {svc.status}
              </span>
            </div>
            <p className="text-[#A6ADC8] text-xs mb-1">
              {svc.ready_replicas}/{svc.desired_replicas} pods
            </p>
            <p className="text-[#6C7086] text-xs truncate">{svc.image}</p>
          </div>
        ))}
      </div>

      {/* Bottom split */}
      <div className="grid grid-cols-2 gap-4">
        
        {/* Reconciliation history */}
        <div className="bg-[#181825] rounded-xl p-4 border border-[#45475A]">
          <h2 className="text-[#89B4FA] font-bold mb-3 text-sm uppercase tracking-widest">
            Reconciliation History
          </h2>
          {history.length === 0 ? (
            <p className="text-[#6C7086] text-sm">No changes recorded yet.</p>
          ) : (
            <div className="space-y-2">
              {[...history].reverse().map((event, i) => (
                <div key={i} className="flex items-start gap-3 text-xs border-b border-[#313244] pb-2">
                  <span className="text-[#6C7086] shrink-0">
                    {new Date(event.timestamp).toLocaleTimeString()}
                  </span>
                  <span className="text-[#CBA6F7] shrink-0">{event.service_name}</span>
                  <span className="text-[#A6ADC8]">
                    {event.what_changed}: 
                    <span className="text-[#F38BA8]"> {event.old_value}</span>
                    <span className="text-[#A6ADC8]"> → </span>
                    <span className="text-[#A6E3A1]">{event.new_value}</span>
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Config */}
        <div className="bg-[#181825] rounded-xl p-4 border border-[#45475A]">
          <h2 className="text-[#89B4FA] font-bold mb-3 text-sm uppercase tracking-widest">
            Live Config
          </h2>
          <pre className="text-[#CDD6F4] text-xs leading-5 overflow-auto max-h-64 bg-[#11111B] p-3 rounded-lg">
            {config}
          </pre>
        </div>

      </div>
    </main>
  )
}