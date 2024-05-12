/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.templ", "**/*.go"],
  theme: {
    colors: {
      'vesper': {
        'bg': '#101010',
        'text': '#ffffff',
        'highlight': '#e6b99d'
      }
    }
  },
  plugins: [],
};