<template>
  <div id="app">
    <header>
      <h1>Sistem Informasi</h1>
    </header>

    <nav class="navbar">
      <ul>
        <li><a href="#general-info">Informasi Umum</a></li>
        <li><a href="#memory-stats">Statistik Memori</a></li>
        <li><a href="#disk-stats">Statistik Disk</a></li>
        <li><a href="#network-stats">Statistik Jaringan</a></li>
        <li><a href="#cpu-stats">Statistik CPU</a></li>
      </ul>
    </nav>

    <main>
      <div class="info-container">
        <section v-if="!data && !error" class="loading-section">
          <div class="loading-spinner"></div>
          <p>Loading data, please wait...</p>
        </section>

        <section v-if="data" id="general-info" class="info-section">
          <h2>Informasi Umum</h2>
          <div class="info-card">
            <p><strong>OS:</strong> {{ data.goos }}</p>
            <p><strong>Arsitektur:</strong> {{ data.goarch }}</p>
            <p><strong>Jumlah CPU:</strong> {{ data.num_cpu }}</p>
            <p>
              <strong>Terakhir Boot:</strong> {{ formatDate(data.last_boot) }}
            </p>
            <p><strong>Uptime:</strong> {{ data.uptime }}</p>
            <p><strong>Jumlah Proses:</strong> {{ data.num_processes }}</p>
          </div>
        </section>

        <section v-if="data" id="memory-stats" class="info-section">
          <h2>Statistik Memori</h2>
          <div class="info-card">
            <p><strong>Total Memori:</strong> {{ data.mem_stats.total_mb }}</p>
            <p><strong>Memori Gratis:</strong> {{ data.mem_stats.free_mb }}</p>
            <p>
              <strong>Memori Terpakai:</strong> {{ data.mem_stats.used_mb }}
            </p>
            <p>
              <strong>Persentase Penggunaan Memori:</strong>
              {{ data.mem_stats.usage_percent.toFixed(2) }}%
            </p>
          </div>
        </section>

        <section v-if="data" id="disk-stats" class="info-section">
          <h2>Statistik Disk</h2>
          <div class="info-card">
            <p><strong>Total Disk:</strong> {{ data.disk_stats.total }}</p>
            <p><strong>Disk Gratis:</strong> {{ data.disk_stats.free }}</p>
          </div>
        </section>

        <section v-if="data" id="network-stats" class="info-section">
          <h2>Statistik Jaringan</h2>
          <ul class="network-list">
            <li
              v-for="(adapter, index) in data.network_stats.adapters"
              :key="index"
              class="network-item"
            >
              <i class="fas fa-network-wired"></i> <strong>Nama:</strong>
              {{ adapter.name }}<br />
              <strong>Alamat IP:</strong> {{ adapter.ip_address }}<br />
              <strong>Upload:</strong> {{ adapter.upload }}<br />
              <strong>Download:</strong> {{ adapter.download }}
            </li>
          </ul>
        </section>

        <section v-if="data" id="cpu-stats" class="info-section">
          <h2>Statistik CPU</h2>
          <div class="info-card">
            <p>
              <strong>Penggunaan CPU:</strong>
              {{ data.cpu_stats.usage_percent.toFixed(2) }}%
            </p>
            <p><strong>Jumlah Cores:</strong> {{ data.cpu_stats.cores }}</p>
          </div>
          <!-- Optional: Add CPU usage chart -->
          <div class="chart-container">
            <LineChart :data="cpuChartData" />
          </div>
        </section>

        <section v-if="error">
          <h2>Error</h2>
          <div class="error-message">{{ error }}</div>
        </section>
      </div>
    </main>
  </div>
</template>

<script>
import axios from "axios";
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  CategoryScale,
  LinearScale,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  CategoryScale,
  LinearScale
);

export default {
  name: "App",
  components: {
    LineChart: Line,
  },
  data() {
    return {
      data: null,
      error: null,
      cpuChartData: {
        labels: ["1m", "5m", "10m", "30m", "1h", "6h", "12h"],
        datasets: [
          {
            label: "Penggunaan CPU",
            backgroundColor: "rgba(66, 185, 131, 0.2)",
            borderColor: "rgba(66, 185, 131, 1)",
            data: [2, 3, 4, 5, 6, 5, 4], // Dummy data, replace with actual
          },
        ],
      },
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
      try {
        const response = await axios.get("/api/device-stats");
        this.data = response.data;
        this.updateChartData();
      } catch (error) {
        this.error = "Error fetching data from API";
        console.error(error);
      }
    },
    updateChartData() {
      // Update chart data based on the API response
      // Example for CPU chart: update with real data
      this.cpuChartData.datasets[0].data = this.data.cpu_usage_data; // Example
    },
    formatDate(dateString) {
      const options = {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        timeZoneName: "short",
      };
      return new Date(dateString).toLocaleString("id-ID", options);
    },
  },
};
</script>

<style scoped>
/* Root Variables */
:root {
  --primary-color: #42b983;
  --secondary-color: #f4f4f9;
  --text-color: #333;
  --error-background: #f8d7da;
  --error-text: #721c24;
}

/* General Styles */
body {
  font-family: "Roboto", sans-serif;
  color: var(--text-color);
  background: linear-gradient(135deg, #f4f4f9 0%, #e0e0e0 100%);
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

header {
  background-color: var(--primary-color);
  color: white;
  padding: 20px;
  text-align: center;
}

nav.navbar {
  background-color: var(--primary-color);
  padding: 10px;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.navbar ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  justify-content: center;
}

.navbar li {
  margin: 0 15px;
}

.navbar a {
  color: white;
  text-decoration: none;
  font-weight: 500;
}

.navbar a:hover {
  text-decoration: underline;
}

main {
  padding: 20px;
}

/* Flexbox Layout */
.info-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

/* Info Section */
.info-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

/* Info Card */
.info-card {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #ddd;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin: 0 auto;
  text-align: left;
}

.info-card p {
  margin: 10px 0;
}

h1 {
  font-size: 2em;
  margin-bottom: 20px;
}

h2 {
  font-size: 1.5em;
  color: var(--primary-color);
  border-bottom: 2px solid var(--primary-color);
  padding-bottom: 10px;
  margin-bottom: 20px;
}

p {
  font-size: 1em;
  line-height: 1.5;
}

/* Network List Styles */
.network-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.network-item {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #ddd;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-bottom: 10px;
}

.network-item i {
  margin-right: 5px;
}

/* Error Styles */
.error-message {
  background: var(--error-background);
  border-left: 4px solid var(--error-text);
  color: var(--error-text);
  padding: 15px;
  border-radius: 5px;
  max-width: 600px;
  margin: 20px auto;
}

/* Loading Section */
.loading-section {
  text-align: center;
}

.loading-spinner {
  border: 8px solid #f3f3f3;
  border-top: 8px solid var(--primary-color);
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* Chart Container */
.chart-container {
  width: 100%;
  max-width: 800px;
  margin: 20px auto;
}

.chart-container canvas {
  width: 100% !important;
  height: auto !important;
}

/* Button and Input Styles */
button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 5px;
  padding: 10px 20px;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

button:hover {
  background-color: #369c6a;
}

input {
  border: 1px solid #ddd;
  border-radius: 5px;
  padding: 10px;
  width: 100%;
  max-width: 400px;
}

input:focus {
  border-color: var(--primary-color);
  outline: none;
}

/* Responsive Design */
@media (max-width: 768px) {
  .info-container {
    grid-template-columns: 1fr;
  }

  .navbar ul {
    flex-direction: column;
  }

  .navbar li {
    margin: 10px 0;
  }
}
</style>
