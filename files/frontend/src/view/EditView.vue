<template>
  <!-- <div class="container"> -->
  <div class="button-group">
    <button class="btn btn-secondary" @click="backTohome">Back</button>
    <button type="submit" form="edit-form" class="btn btn-primary">Save</button>
  </div>
  <form id="edit-form" @submit.prevent="saveChanges">
    <input type="hidden" name="file" :value="fileName" />
    <div class="form-group editor-container">
      <div id="line-numbers" class="line-numbers" ref="linenumber"></div>
      <textarea
        name="content"
        id="content"
        rows="20"
        class="form-control editor-content"
        aria-label="File content"
        ref="content"
        v-model="fileData.content"
        @input="updateLineNumbers"
      ></textarea>
    </div>
  </form>
  <!-- </div> -->
</template>

<script>
import axios from "axios";
import { useToast } from "vue-toastification";
export default {
  name: "FileEditor",
  data() {
    return {
      fileData: {
        fileName: "",
        content: "",
      },
    };
  },
  methods: {
    backTohome() {
      this.$router.push({
        name: "Home",
      });
    },
    async fetchFileContent() {
      try {
        const filepath = this.$route.params.filepath;
        const response = await axios.get(`/api/files/view`, {
          params: { filepath: filepath },
        });
        this.fileData = response.data;
      } catch (error) {
        console.error("Error fetching file content:", error);
      }
    },
    async saveChanges() {
      const toast = useToast();
      try {
        const response = await axios.post(`/api/files/save`, this.fileData);
        this.$router.push({ name: "Home" });
        toast.success(response.data.message, this.getToastOptions());
      } catch (error) {
        let errorMessage = "Error saving file.";

        // Determine specific error message
        if (error.response) {
          errorMessage = error.response.data.message || errorMessage;
        } else if (error.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }
        toast.error(errorMessage, this.getToastOptions());
      }
    },
    getToastOptions() {
      return {
        position: "top-right",
        timeout: 1989,
        closeOnClick: true,
        pauseOnFocusLoss: true,
        pauseOnHover: true,
        draggable: true,
        draggablePercent: 0.6,
        showCloseButtonOnHover: true,
        hideProgressBar: true,
        closeButton: "button",
        icon: true,
        rtl: false,
      };
    },
    updateLineNumbers() {
      const textarea = this.$refs.content;
      const lineNumbersEle = this.$refs.linenumber;

      const textareaStyles = window.getComputedStyle(textarea);
      [
        "fontFamily",
        "fontSize",
        "fontWeight",
        "letterSpacing",
        "lineHeight",
        "padding",
      ].forEach((property) => {
        lineNumbersEle.style[property] = textareaStyles[property];
      });

      const parseValue = (v) =>
        v.endsWith("px") ? parseInt(v.slice(0, -2), 10) : 0;

      const font = `${textareaStyles.fontSize} ${textareaStyles.fontFamily}`;
      const paddingLeft = parseValue(textareaStyles.paddingLeft);
      const paddingRight = parseValue(textareaStyles.paddingRight);

      const canvas = document.createElement("canvas");
      const context = canvas.getContext("2d");
      context.font = font;

      const calculateNumLines = (str) => {
        const textareaWidth =
          textarea.getBoundingClientRect().width - paddingLeft - paddingRight;
        const words = str.split(" ");
        let lineCount = 0;
        let currentLine = "";
        for (let i = 0; i < words.length; i++) {
          const wordWidth = context.measureText(words[i] + " ").width;
          const lineWidth = context.measureText(currentLine).width;

          if (lineWidth + wordWidth > textareaWidth) {
            lineCount++;
            currentLine = words[i] + " ";
          } else {
            currentLine += words[i] + " ";
          }
        }

        if (currentLine.trim() !== "") {
          lineCount++;
        }

        return lineCount;
      };

      const calculateLineNumbers = () => {
        const lines = textarea.value.split("\n");
        const numLines = lines.map((line) => calculateNumLines(line));

        let lineNumbers = [];
        let i = 1;
        while (numLines.length > 0) {
          const numLinesOfSentence = numLines.shift();
          lineNumbers.push(i);
          if (numLinesOfSentence > 1) {
            Array(numLinesOfSentence - 1)
              .fill("")
              .forEach(() => lineNumbers.push(""));
          }
          i++;
        }

        return lineNumbers;
      };

      const displayLineNumbers = () => {
        const lineNumbers = calculateLineNumbers();
        lineNumbersEle.innerHTML = Array.from(
          {
            length: lineNumbers.length,
          },
          (_, i) => `<div>${lineNumbers[i] || "&nbsp;"}</div>`
        ).join("");
      };

      textarea.addEventListener("input", () => {
        displayLineNumbers();
      });

      displayLineNumbers();

      const ro = new ResizeObserver(() => {
        const rect = textarea.getBoundingClientRect();
        lineNumbersEle.style.height = `${rect.height}px`;
        displayLineNumbers();
      });
      ro.observe(textarea);

      textarea.addEventListener("scroll", () => {
        lineNumbersEle.scrollTop = textarea.scrollTop;
      });
    },
  },
  mounted() {
    this.fetchFileContent();
    this.updateLineNumbers();
  },
};
</script>

<style scoped>
body {
  /* font-family: Arial, sans-serif; */
  background-color: #f8f9fa; /* Default light background */
  color: #495057; /* Default light text color */
  transition: background-color 0.3s, color 0.3s;
  margin: 0;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
/* * {
  box-sizing: border-box;
} */
.container {
  max-width: 800px;
  margin-top: 50px;
  padding: 20px;
  background-color: white; /* Default light container */
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  transition: background-color 0.3s, color 0.3s;
  position: relative; /* Allow absolute positioning of toggle */
}
h1 {
  font-size: 28px;
  margin-bottom: 30px;
  color: #343a40; /* Default light heading color */
  text-align: center;
}
.form-control {
  border: none;
  outline: none;
  width: 100%;
  transition: border-color 0.3s;
}
.form-control:focus {
  border-color: #007bff; /* Change border color on focus */
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
}
.line-numbers {
  border-right: 1px solid rgb(203 213 225);
  padding: 0.5rem;
  text-align: right;
  overflow: hidden;
}
.editor-container {
  display: flex;
  border: 1px solid #ced4da;
  border-radius: 0.5rem;
  overflow: hidden;
  max-height: 32rem;
  padding: 0.5rem;
  font-family: "Courier New", Courier, monospace;
  font-size: 14px;
}
.button-group {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between; /* Space between buttons */
}
.btn {
  width: 48%; /* Make buttons equal width */
}
.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
  transition: background-color 0.3s, border-color 0.3s;
}
.btn-primary:hover {
  background-color: #0056b3;
  border-color: #004085;
}
.btn-secondary {
  background-color: #6c757d;
  border-color: #6c757d;
  transition: background-color 0.3s, border-color 0.3s;
}
.btn-secondary:hover {
  background-color: #5a6268;
  border-color: #545b62;
}
</style>
