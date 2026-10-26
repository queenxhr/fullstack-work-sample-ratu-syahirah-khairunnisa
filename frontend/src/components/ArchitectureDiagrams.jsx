import React, { useState } from 'react';
import { DollarSign, Layers, PieChart, TrendingDown, Server, Cpu, CheckSquare } from 'lucide-react';

export default function ArchitectureDiagrams() {
  // Soal 4: TCO Calculator State
  const [monolithCost, setMonolithCost] = useState(150);
  const [eksNodes, setEksNodes] = useState(3);
  const [rdsInstances, setRdsInstances] = useState(3);
  const [devopsSalaries, setDevopsSalaries] = useState(2500); // 1 DevOps engineer monthly USD equivalent

  // Calculated Microservices Cost
  const eksControlPlane = 73; // AWS EKS control plane monthly
  const nodeCost = eksNodes * 80; // ~$80/month for t3.xlarge
  const rdsCost = rdsInstances * 95; // ~$95/month for db.t3.medium
  const networkingObservability = 220; // NAT Gateway, ALB, Datadog/Cloudwatch
  const totalMicroservicesInfra = eksControlPlane + nodeCost + rdsCost + networkingObservability;
  const totalMicroservicesTCO = totalMicroservicesInfra + devopsSalaries;

  // Soal 5: LTV / CAC Simulator State
  const [arpu, setArpu] = useState(250000); // IDR per customer per month
  const [churnRate, setChurnRate] = useState(6); // 6% monthly churn
  const [cac, setCac] = useState(750000); // IDR acquisition cost

  const calculatedLTV = Math.round(arpu / (churnRate / 100));
  const ltvToCacRatio = (calculatedLTV / cac).toFixed(2);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      {/* Soal #4: TCO Analysis */}
      <div className="glass-panel" style={{ padding: '1.75rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
          <DollarSign size={24} color="var(--accent-cyan)" />
          <div>
            <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>
              Soal #4: Evaluasi TCO (Total Cost of Ownership) &amp; Trade-off Monolith vs Microservices
            </h2>
            <div
              title="Kalkulator membandingkan biaya monolith dengan biaya infrastruktur cloud-native dan tenaga operasional."
              style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}
            >
              Kalkulator TCO
            </div>
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem', marginBottom: '1.5rem' }}>
          {/* Controls */}
          <div style={{ padding: '1.25rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)' }}>
            <h3 style={{ fontSize: '0.95rem', fontWeight: 600, marginBottom: '1rem', color: 'var(--accent-cyan)' }}>
              Parameter Infrastruktur Microservices:
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem', fontSize: '0.85rem' }}>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>EKS Worker Nodes (Multi-AZ):</span>
                  <strong>{eksNodes} nodes (${nodeCost}/mo)</strong>
                </div>
                <input
                  type="range"
                  min="2"
                  max="10"
                  value={eksNodes}
                  onChange={(e) => setEksNodes(Number(e.target.value))}
                  style={{ width: '100%' }}
                />
              </div>

              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>Database Per Service (RDS):</span>
                  <strong>{rdsInstances} DBs (${rdsCost}/mo)</strong>
                </div>
                <input
                  type="range"
                  min="2"
                  max="8"
                  value={rdsInstances}
                  onChange={(e) => setRdsInstances(Number(e.target.value))}
                  style={{ width: '100%' }}
                />
              </div>

              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>Gaji DevOps/Platform Engineer:</span>
                  <strong>${devopsSalaries}/mo</strong>
                </div>
                <input
                  type="range"
                  min="0"
                  max="5000"
                  step="250"
                  value={devopsSalaries}
                  onChange={(e) => setDevopsSalaries(Number(e.target.value))}
                  style={{ width: '100%' }}
                />
              </div>
            </div>
          </div>

          {/* Results Comparison */}
          <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between', padding: '1.25rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)' }}>
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: '0.75rem', borderBottom: '1px solid var(--border-color)', marginBottom: '0.75rem' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Biaya Monolith Saat Ini:</span>
                <span className="mono" style={{ fontSize: '1.15rem', fontWeight: 700, color: 'var(--accent-emerald)' }}>
                  ${monolithCost}/bulan
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: '0.75rem', borderBottom: '1px solid var(--border-color)', marginBottom: '0.75rem' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Infra Cloud Microservices:</span>
                <span className="mono" style={{ fontSize: '1.15rem', fontWeight: 700, color: 'var(--accent-rose)' }}>
                  ${totalMicroservicesInfra}/bulan
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: '0.75rem', borderBottom: '1px solid var(--border-color)', marginBottom: '0.75rem' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Total TCO (Infra + Human Effort):</span>
                <span className="mono" style={{ fontSize: '1.35rem', fontWeight: 800, color: 'var(--accent-rose)' }}>
                  ${totalMicroservicesTCO}/bulan
                </span>
              </div>
            </div>

            <div style={{ padding: '0.75rem', background: 'rgba(244, 63, 94, 0.1)', border: '1px solid rgba(244, 63, 94, 0.3)', borderRadius: 'var(--radius-sm)', fontSize: '0.8rem', color: '#fb7185' }}>
              <strong>Kesimpulan Finansial:</strong> Migrasi langsung ke Microservices menyebabkan lonjakan biaya <strong>{(totalMicroservicesTCO / monolithCost).toFixed(1)}x lipat</strong> tanpa peningkatan revenue yang sebanding.
              Rekomendasi terbaik untuk startup tahap awal-menengah adalah <strong>Modular Monolith</strong>.
            </div>
          </div>
        </div>
      </div>

      {/* Soal #5: LTV/CAC & Sprint Allocation */}
      <div className="glass-panel" style={{ padding: '1.75rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
          <PieChart size={24} color="var(--accent-cyan)" />
          <div>
            <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>
              Soal #5: Product Business Metrics (LTV / CAC) &amp; Sprint Capacity Allocation
            </h2>
            <div
              title="Simulasi hubungan churn dengan LTV/CAC dan contoh pembagian kapasitas sprint 70/20/10."
              style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}
            >
              LTV/CAC &amp; alokasi sprint
            </div>
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
          {/* LTV Simulator */}
          <div style={{ padding: '1.25rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)' }}>
            <h3 style={{ fontSize: '0.95rem', fontWeight: 600, marginBottom: '0.75rem', color: 'var(--accent-purple)' }}>
              Simulasi Dampak Downtime Terhadap Churn &amp; LTV:
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', fontSize: '0.825rem' }}>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>Customer Monthly Churn Rate (%):</span>
                  <strong style={{ color: churnRate > 5 ? 'var(--accent-rose)' : 'var(--accent-emerald)' }}>{churnRate}%</strong>
                </div>
                <input
                  type="range"
                  min="2"
                  max="20"
                  step="0.5"
                  value={churnRate}
                  onChange={(e) => setChurnRate(Number(e.target.value))}
                  style={{ width: '100%' }}
                />
              </div>

              <div style={{ padding: '0.75rem', background: 'var(--bg-primary)', borderRadius: 'var(--radius-sm)', marginTop: '0.5rem' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>Customer Lifetime Value (LTV):</span>
                  <strong className="mono" style={{ color: 'var(--accent-cyan)' }}>Rp {calculatedLTV.toLocaleString('id-ID')}</strong>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                  <span>Customer Acquisition Cost (CAC):</span>
                  <strong className="mono">Rp {cac.toLocaleString('id-ID')}</strong>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid var(--border-color)', paddingTop: '0.5rem', marginTop: '0.25rem' }}>
                  <span>Rasio LTV / CAC:</span>
                  <strong className="mono" style={{ fontSize: '1.1rem', color: Number(ltvToCacRatio) >= 3 ? 'var(--accent-emerald)' : 'var(--accent-rose)' }}>
                    {ltvToCacRatio}x {Number(ltvToCacRatio) >= 3 ? '✅ (Sehat)' : '⚠️ (Berbahaya)'}
                  </strong>
                </div>
              </div>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                *Downtime dan bug checkout menaikkan churn rate dari 4% ke 12%, memotong LTV hingga 66% dan membuat pengeluaran CAC marketing rugi.
              </div>
            </div>
          </div>

          {/* 70/20/10 Sprint Framework */}
          <div style={{ padding: '1.25rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)' }}>
            <h3 style={{ fontSize: '0.95rem', fontWeight: 600, marginBottom: '0.75rem', color: 'var(--accent-cyan)' }}>
              Framework Alokasi Sprint 70 / 20 / 10:
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', fontSize: '0.825rem' }}>
              <div style={{ padding: '0.625rem', background: 'rgba(56, 189, 248, 0.1)', borderLeft: '4px solid var(--accent-cyan)', borderRadius: '4px' }}>
                <strong>70% Kapasitas: New Business Features</strong>
                <div style={{ color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
                  3 fitur baru yang diminta Product Manager untuk mengejar akuisisi pelanggan (CAC target).
                </div>
              </div>
              <div style={{ padding: '0.625rem', background: 'rgba(244, 63, 94, 0.1)', borderLeft: '4px solid var(--accent-rose)', borderRadius: '4px' }}>
                <strong>20% Kapasitas: Tech Debt &amp; Reliability Refactoring</strong>
                <div style={{ color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
                  Mengatasi memory leak, mengoptimasi query deadlock, dan mencegah catastrophic server downtime.
                </div>
              </div>
              <div style={{ padding: '0.625rem', background: 'rgba(168, 85, 247, 0.1)', borderLeft: '4px solid var(--accent-purple)', borderRadius: '4px' }}>
                <strong>10% Kapasitas: Unplanned Bugs &amp; Maintenance</strong>
                <div style={{ color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
                  Buffer untuk production hotfixes, library security patches, dan eksplorasi teknis.
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
