<template>
  <div>
    <!-- Button Container -->
    <div class="button-container">
      <button
        v-for="command in commands"
        :key="command.Name"
        @click="executeCommand(command.Value)"
        aria-label="Execute command"
      >
        {{ command.Name }}
        <span class="tooltip">{{ command.Description }}</span>
      </button>
    </div>

    <!-- Input Container -->
    <div class="input-container">
      <input
        type="text"
        v-model="customCommand"
        id="customCommandInput"
        placeholder="Masukkan perintah kustom"
        aria-label="Custom command input"
      />
      <button
        id="executeCustomButton"
        @click="executeCommand('custom')"
        :aria-disabled="!customCommand.trim()"
        :disabled="!customCommand.trim()"
        aria-label="Execute custom command"
      >
        Execute
      </button>
    </div>

    <!-- Output Container -->
    <div class="output-container">
      <pre id="output" aria-live="polite">{{ output }}</pre>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: 'CommandExecutor',
  data() {
    return {
      commands: [],
      customCommand: "",
      output: "",
      isLoading: false,
      error: null,
    };
  },
  computed: {
    isCustomCommandValid() {
      return this.customCommand.trim().length > 0;
    },
  },
  methods: {
    async loadCommands() {
      if (this.commands.length) return;

      this.isLoading = true;
      this.error = null;

      try {
        const { data } = await axios.get("/api/command/list");
        if (Array.isArray(data)) {
          this.commands = data;
        } else {
          throw new Error("Data tidak valid");
        }
      } catch (error) {
        this.handleError("Kesalahan memuat perintah:", error);
      } finally {
        this.isLoading = false;
      }
    },
    async executeCommand(cmd) {
      if (cmd === "custom" && !this.isCustomCommandValid) {
        this.output = "Perintah kustom tidak boleh kosong";
        return;
      }

      this.isLoading = true;
      this.error = null;
      this.output = "";

      const url = "/api/command/execute";
      const params = new URLSearchParams({
        value: cmd === "custom" ? "custom" : cmd,
        ...(cmd === "custom" && { custom_command: this.customCommand }),
      });

      try {
        const { data } = await axios.post(url, params, {
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        });
        this.output = data.message;
      } catch (error) {
        this.handleError("Kesalahan saat mengeksekusi perintah:", error);
      } finally {
        this.isLoading = false;
      }
    },
    handleError(prefix, error) {
      console.error(prefix, error);
      this.error = error.response?.data?.message 
        || error.message 
        || "Terjadi kesalahan yang tidak diketahui";
    },
  },
  mounted() {
    this.loadCommands();
  },
};
</script>

<style>
:root {
  --primary-color: #00509e; /* Biru gelap yang lebih cerah */
  --primary-hover: #003d6b; /* Biru gelap lebih tua untuk hover */
  --secondary-color: #e1e8f0; /* Abu-abu-biru sangat terang untuk kontras */
  --text-color: #ffffff; /* Teks putih untuk kontras */
  --shadow-color: rgba(0, 0, 0, 0.2); /* Bayangan lebih lembut */
  --border-color: #0074d9; /* Warna border biru cerah */
  --button-padding: 12px;
  --input-padding: 12px;
  --border-radius: 6px;
  --font-size-small: 12px;
  --font-size-medium: 14px;
  --font-size-large: 16px;
  --font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  --font-monospace: "Courier New", Courier, monospace;
  --gap: 12px;
  --margin: 20px;
  --button-gradient-start: #0074d9; /* Biru cerah untuk gradient tombol */
  --button-gradient-end: #00509e; /* Biru gelap untuk gradient tombol */
  --tooltip-background: #333; /* Latar belakang tooltip gelap */
  --tooltip-text: #fff; /* Teks tooltip putih */
  --tooltip-border: #555; /* Border tooltip gelap */
}

/* General Styles */
body {
  font-family: var(--font-family);
  padding: 0;
  background-color: var(
    --secondary-color
  ); /* Latar belakang abu-abu-biru sangat terang */
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  min-height: 100vh;
  padding: var(--margin);
  box-sizing: border-box;
}

/* Button Container */
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap);
  max-width: 80vw;
  width: 100%;
  box-sizing: border-box;
  padding: var(--margin);
  background-color: var(--secondary-color);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 10px var(--shadow-color);
}

/* Button Styles */
.button-container button {
  background: linear-gradient(
    145deg,
    var(--button-gradient-start),
    var(--button-gradient-end)
  );
  color: var(--text-color);
  border: none;
  padding: var(--button-padding);
  border-radius: var(--border-radius);
  cursor: pointer;
  flex: 1 1 calc(25% - var(--gap));
  box-sizing: border-box;
  position: relative;
  text-align: center;
  margin: 1px;
  height: auto;
  font-size: var(--font-size-medium);
  transition: all 0.3s ease;
}

.button-container button:hover {
  background: linear-gradient(
    145deg,
    var(--button-gradient-end),
    var(--button-gradient-start)
  );
  transform: translateY(-6px);
  box-shadow: 0 8px 20px var(--shadow-color);
}

.button-container button:active {
  transform: translateY(2px);
  box-shadow: 0 4px 10px var(--shadow-color);
}

/* Tooltip Styles */
.tooltip {
  display: none;
  position: absolute;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%);
  padding: var(--input-padding);
  background-color: var(--tooltip-background);
  color: var(--tooltip-text);
  border: 1px solid var(--tooltip-border);
  border-radius: 4px;
  font-size: var(--font-size-small);
  white-space: nowrap;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.button-container button:hover .tooltip {
  display: block;
  opacity: 1;
  transform: translateX(-50%) translateY(-10px);
}

/* Input and Output Containers */
.input-container,
.output-container {
  display: flex;
  align-items: center;
  max-width: 80vw;
  width: 100%;
  box-sizing: border-box;
  padding: var(--margin);
  background-color: var(--secondary-color);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 10px var(--shadow-color);
  margin-top: var(--margin);
}

/* Input Styles */
.input-container input {
  flex: 1;
  padding: var(--input-padding);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius) 0 0 var(--border-radius);
  margin-right: -1px;
  font-size: var(--font-size-medium);
  background-color: #ffffff; /* Background input putih untuk keterbacaan */
}

/* Button Styles in Input Container */
.input-container button {
  background: linear-gradient(
    145deg,
    var(--button-gradient-start),
    var(--button-gradient-end)
  );
  color: var(--text-color);
  border: none;
  padding: var(--button-padding) 20px;
  border-radius: 0 var(--border-radius) var(--border-radius) 0;
  cursor: pointer;
  margin-left: -1px;
  font-size: var(--font-size-medium);
  transition: all 0.3s ease;
}

.input-container button:hover {
  background: linear-gradient(
    145deg,
    var(--button-gradient-end),
    var(--button-gradient-start)
  );
}

/* Output Styles */
#output {
  display: block;
  background-color: var(--secondary-color);
  color: #333333;
  border-radius: var(--border-radius);
  min-height: 120px;
  max-height: 250px;
  overflow-y: auto;
  overflow-x: auto;
  font-family: var(--font-monospace);
  white-space: pre-wrap;
  font-size: var(--font-size-medium);
}

@media (max-width: 768px) {
  .button-container button {
    flex: 1 1 calc(33.333% - var(--gap));
  }
  .button-container,
  .input-container,
  .output-container {
    max-width: 90vw;
  }
}

@media (max-width: 480px) {
  .button-container button {
    flex: 1 1 calc(33.333% - var(--gap));
  }
}
</style>
