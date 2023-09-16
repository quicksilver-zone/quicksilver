/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './app/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        nexa: ['var(--font-nexa)'],
        quicksand: ['var(--font-quicksand)']
      },
      backgroundImage: {
        "background": "url('../assets/imgs/background.svg')"
      }
    },
  },
  plugins: [],
}

