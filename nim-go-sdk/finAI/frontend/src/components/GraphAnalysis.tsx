import React, { useState } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { Line, Bar, Doughnut } from 'react-chartjs-2';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

export default function GraphAnalysis() {
  const days = Array.from({ length: 30 }, (_, i) => `Day ${i + 1}`);
  const [enlargedGraph, setEnlargedGraph] = useState<string | null>(null);


  // --- 1. Monthly Spending Bar Chart Data ---
  const barData = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
    datasets: [{
      label: 'Total Spending',
      data: [1213.18, 957.80, 1007.35, 1110.74, 893.44, 828.34, 1059.12, 1257.66, 980.38, 1108.97, 1041.25, 920.58],
      backgroundColor: 'rgba(54, 162, 235, 0.7)',
      borderRadius: 8,
    }]
  };


  function BlurBack() {
    return (
      <div onClick={() => setEnlargedGraph(null)} style={{
        position: 'fixed',
        inset: 0,
        zIndex: 9999,
        backgroundColor: 'rgba(255, 255, 255, 0.3)',
        backdropFilter: 'blur(2px)',
        WebkitBackdropFilter: 'blur(2px)', // Safari
        pointerEvents: 'all',
      }}
      />
    )
  }

  function EnlargeButton({ graphId }: { graphId: string }) {
    return (
      <button
        onClick={() => setEnlargedGraph(graphId)}
        style={{
          position: 'absolute',
          top: '16px',
          right: '16px',
          background: '#cdcdcdff',
          color: 'white',
          border: 'none',
          borderRadius: '8px',
          padding: '8px 12px',
          cursor: 'pointer',
          fontSize: '14px',
          fontWeight: '500',
          zIndex: 10,
          transition: 'all 0.2s ease',
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.15)',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.transform = 'scale(1.05)';
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.transform = 'scale(1)';
        }}
      >
        ⤢
      </button>
    );
  }

  const cardStyle: React.CSSProperties = {
    background: '#ffffff',
    padding: '24px',
    borderRadius: '16px',
    boxShadow: '0 12px 40px rgba(0, 0, 0, 0.08)',
    height: '33vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
  }

  const enlargedCardStyle: React.CSSProperties = {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: '80vw',
    height: '80vh',
    background: '#ffffff',
    padding: '40px',
    borderRadius: '20px',
    boxShadow: '0 20px 60px rgba(0, 0, 0, 0.3)',
    zIndex: 10000,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  }

  // --- 2. Category Doughnut Chart Data ---
  const doughnutData = {
    labels: ["Food", "Leisure", "Bills", "Transport", "Shopping"],
    datasets: [{
      data: [872.38, 448.09, 1376.98, 324.12, 510.12],
      backgroundColor: ["#FF6B2C", "#2CA784", "#454A5A", "#2C7AA8", "#A83E60"],
      borderWidth: 0,
      hoverOffset: 8,
    }]
  };

  // --- 3. Average Daily Spending Line Chart Data ---
  const dailyAvgData = {
    labels: days,
    datasets: [
      {
        label: '30-Day Average',
        data: [45, 42, 50, 48, 46, 47, 49, 44, 50, 48, 46, 45, 47, 49, 46, 44, 48, 47, 46, 45, 50, 49, 48, 46, 45, 44, 47, 48, 46, 45],
        borderColor: '#2C7AA8',
        backgroundColor: 'rgba(44,122,168,0.1)',
        fill: true,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Current Average',
        data: Array(30).fill(47),
        borderColor: '#FF6B2C',
        borderDash: [6, 4],
        tension: 0.4,
        pointRadius: 0,
      }
    ]
  };

  // --- 4. Budget Tracking Line Chart Data ---
  const budgetValue = 600;
  const spendingTrend = [
    100, 120, 150, 180, 200, 230, 250, 270, 300, 320,
    350, 370, 400, 420, 430, 450, 470, 480, 500, 520,
    550, 560, 580, 600, 620, 650, 670, 690, 720, 750
  ];

  const budgetData = {
    labels: days,
    datasets: [
      {
        label: 'User Spending',
        data: spendingTrend,
        // The 'segment' property is the key to restoring the gradient effect
        segment: {
          borderColor: (ctx: any) => {
            const val = (ctx.p0.parsed.y + ctx.p1.parsed.y) / 2;
            if (val < budgetValue * 0.8) return "#28a745"; // Green
            if (val <= budgetValue) return "#FF8C00";      // Orange
            return "#dc3545";                             // Red
          },
        },
        borderWidth: 3,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Budget Limit',
        data: Array(30).fill(budgetValue),
        borderColor: '#6c757d',
        borderDash: [6, 4],
        pointRadius: 0,
      }
    ]
  };

  const commonOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' as const },
    },
  };

  const doughnutOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' as const },
    },
    cutout: '68%',
  };


  return (
<<<<<<< HEAD
    <div style={{ padding: '40px', backgroundColor: '#f4f3f2', overflow: 'hidden' }}>
      {enlargedGraph && <BlurBack />}

      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))',
=======
    <div style={{ padding: '40px', backgroundColor: '#f4f3f2', 'overflow': 'hidden', maxHeight:'73vh' }}>
      
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', 
>>>>>>> 751f30bd2cbcce0060cc1fc6a94116e1df6bc157
        gap: '24px',
        maxWidth: '1200px',
        margin: '0 auto'
      }}>

        {/* Monthly Spending Bar Chart */}
        <div style={cardStyle}>
          <EnlargeButton graphId="monthly" />
          <Bar
            data={barData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Monthly Spending Overview' } } }}
          />
        </div>

        {/* Category Doughnut Chart */}
        <div style={cardStyle}>
          <EnlargeButton graphId="category" />
          <div style={{ width: '80%', margin: '0 auto', height: '30vh' }}>
            <Doughnut
              data={doughnutData}
              options={{ ...doughnutOptions, cutout: '68%', plugins: { ...doughnutOptions.plugins, title: { display: true, text: 'Spending by Category' } } }}
            />
          </div>
        </div>

        {/* Daily Average Line Chart */}
        <div style={cardStyle}>
          <EnlargeButton graphId="daily" />
          <Line
            data={dailyAvgData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Average Daily Spending' } } }}
          />
        </div>

        {/* Budget vs Spending Line Chart */}
        <div style={cardStyle}>
          <EnlargeButton graphId="budget" />
          <Line
            data={budgetData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Spending vs Budget' } } }}
          />
        </div>

      </div>

      {/* Enlarged Graph Overlays */}
      {enlargedGraph === 'monthly' && (
        <div style={enlargedCardStyle}>
          <button
            onClick={() => setEnlargedGraph(null)}
            style={{
              position: 'absolute',
              top: '20px',
              right: '20px',
              background: '#cdcdcdff',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 16px',
              cursor: 'pointer',
              fontSize: '16px',
              fontWeight: '600',
              zIndex: 10001,
            }}
          >
            ✕
          </button>
          <Bar
            data={barData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Monthly Spending Overview', font: { size: 20 } } } }}
          />
        </div>
      )}

      {enlargedGraph === 'category' && (
        <div style={enlargedCardStyle}>
          <button
            onClick={() => setEnlargedGraph(null)}
            style={{
              position: 'absolute',
              top: '20px',
              right: '20px',
              background: '#cdcdcdff',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 16px',
              cursor: 'pointer',
              fontSize: '16px',
              fontWeight: '600',
              zIndex: 10001,
            }}
          >
            ✕
          </button>
          <div style={{ width: '60%', height: '70vh', margin: '0 auto' }}>
            <Doughnut
              data={doughnutData}
              options={{ ...doughnutOptions, plugins: { ...doughnutOptions.plugins, title: { display: true, text: 'Spending by Category', font: { size: 20 } } } }}
            />
          </div>
        </div>
      )}

      {enlargedGraph === 'daily' && (
        <div style={enlargedCardStyle}>
          <button
            onClick={() => setEnlargedGraph(null)}
            style={{
              position: 'absolute',
              top: '20px',
              right: '20px',
              background: '#cdcdcdff',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 16px',
              cursor: 'pointer',
              fontSize: '16px',
              fontWeight: '600',
              zIndex: 10001,
            }}
          >
            ✕
          </button>
          <Line
            data={dailyAvgData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Average Daily Spending', font: { size: 20 } } } }}
          />
        </div>
      )}

      {enlargedGraph === 'budget' && (
        <div style={enlargedCardStyle}>
          <button
            onClick={() => setEnlargedGraph(null)}
            style={{
              position: 'absolute',
              top: '20px',
              right: '20px',
              background: '#cdcdcdff',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 16px',
              cursor: 'pointer',
              fontSize: '16px',
              fontWeight: '600',
              zIndex: 10001,
            }}
          >
            ✕
          </button>
          <Line
            data={budgetData}
            options={{ ...commonOptions, plugins: { ...commonOptions.plugins, title: { display: true, text: 'Spending vs Budget', font: { size: 20 } } } }}
          />
        </div>
      )}

    </div>
  );
};