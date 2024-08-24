<template>
  <div>
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
    <div class="output-container">
      <pre id="output" aria-live="polite">{{ output }}</pre>
    </div>
  </div>
</template>
<script>
import { ref, onMounted } from "vue";

export default {
  setup() {
    const commands = ref([]);
    const customCommand = ref("");
    const output = ref("");

    const handleResponse = async (response) => {
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Kesalahan ${response.status}: ${errorText}`);
      }
      return response.json();
    };

    const loadCommands = async () => {
      try {
        const response = await fetch("/api/list");
        const data = await handleResponse(response);
        commands.value = Array.isArray(data.commands) ? data.commands : [];
        if (!Array.isArray(data.commands)) {
          output.value = "Kesalahan memuat perintah";
        }
      } catch (error) {
        console.error("Kesalahan pengambilan:", error);
        output.value = "Kesalahan memuat perintah";
      }
    };

    const executeCommand = async (cmd) => {
      if (cmd === "custom" && !customCommand.value.trim()) {
        output.value = "Perintah kustom tidak boleh kosong";
        return;
      }

      const url = `/api/execute?value=${encodeURIComponent(
        cmd === "custom" ? "custom" : cmd
      )}`;
      const body =
        cmd === "custom"
          ? `custom_command=${encodeURIComponent(customCommand.value)}`
          : null;

      try {
        const response = await fetch(url, {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: body,
        });
        const result = await handleResponse(response);
        output.value = result.output;
      } catch (error) {
        console.error("Kesalahan pengambilan:", error);
        output.value = "Kesalahan saat mengeksekusi perintah";
      }
    };

    onMounted(async () => {
      await loadCommands();
    });

    return { commands, customCommand, output, executeCommand };
  },
};
</script>

<style>
:root {
  --primary-color: #007bff;
  --primary-hover: #0056b3;
  --secondary-color: #fff;
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --border-color: #ddd;
  --button-padding: 10px;
  --input-padding: 10px;
  --border-radius: 8px;
  --font-size-small: 12px;
  --font-size-medium: 14px;
  --font-family: Arial, sans-serif;
  --font-monospace: "Courier New", Courier, monospace;
  --text-color: #333;
  --gap: 10px;
  --margin: 20px;
}

body {
  font-family: var(--font-family);
  padding: 0;
  background-color: #f4f4f4;
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  min-height: 100vh;
  padding: var(--margin);
  box-sizing: border-box;
}
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap);
  max-width: 70vw;
  width: 100%;
  box-sizing: border-box;
  padding: var(--margin);
  background-color: var(--secondary-color);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 8px var(--shadow-color);
}
.button-container button {
  background-color: var(--primary-color);
  color: var(--secondary-color);
  border: none;
  padding: var(--button-padding);
  border-radius: var(--border-radius);
  cursor: pointer;
  flex: 1 1 calc(33.333% - var(--gap));
  box-sizing: border-box;
  position: relative;
  text-align: center;
  margin: 1px;
  height: auto;
}
.button-container button:hover {
  background-color: var(--primary-hover);
  transform: translateY(-2px);
  box-shadow: 0 6px 10px rgba(0, 0, 0, 0.2);
}
.button-container button:hover::after {
  display: block;
}
.tooltip {
  display: none;
  position: absolute;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%);
  padding: var(--input-padding);
  background-color: var(--text-color);
  color: #fff;
  border-radius: 4px;
  font-size: var(--font-size-medium);
  white-space: nowrap;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s ease;
}
.button-container button:hover .tooltip {
  display: block;
  opacity: 1;
}

.input-container {
  display: flex;
  flex-direction: row;
  align-items: stretch;
  max-width: 70vw;
  width: 100%;
  box-sizing: border-box;
  padding: var(--margin);
  background-color: var(--secondary-color);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 8px var(--shadow-color);
  margin-top: var(--margin);
}
.input-container input {
  flex: 1;
  padding: var(--input-padding);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius) 0 0 var(--border-radius);
  margin-right: -1px;
}
.input-container button {
  background-color: var(--primary-color);
  color: var(--secondary-color);
  border: none;
  padding: var(--button-padding) 20px;
  border-radius: 0 var(--border-radius) var(--border-radius) 0;
  cursor: pointer;
  margin-left: -1px;
}
.output-container {
  display: flex;
  align-items: center;
  max-width: 70vw;
  width: 100%;
  box-sizing: border-box;
  padding: 5px var(--margin);
  background-color: var(--secondary-color);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 8px var(--shadow-color);
  margin-top: var(--margin);
}
#output {
  display: flex;
  background-color: var(--secondary-color);
  color: var(--text-color);
  border-radius: var(--border-radius);
  max-width: 70vw;
  width: 100%;
  overflow-x: auto;
  overflow-y: auto;
  font-family: var(--font-monospace);
  white-space: pre-wrap;
  font-size: var(--font-size-medium);
}

@media (max-width: 1200px) {
  .button-container .input-container {
    max-width: 90%;
  }
}
@media (max-width: 768px) {
  .button-container button {
    flex: 1 1 calc(33.333% - var(--gap));
  }
}
@media (max-width: 480px) {
  .button-container button {
    flex: 1 1 calc(33.333% - var(--gap));
  }
  .button-container {
    max-width: 90vw;
    width: 100%;
  }
  .input-container {
    max-width: 90vw;
    width: 100%;
  }
  .output-container {
    max-width: 90vw;
    width: 100%;
  }
}
</style>
