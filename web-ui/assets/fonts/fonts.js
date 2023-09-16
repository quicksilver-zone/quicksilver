import localFont from 'next/font/local';

export const nexa = localFont({
  src: [
    {
      path: './Nexa-Bold.woff2',
      weight: '800',
      style: "normal"
    },
    {
      path: './Nexa-Light.woff2',
      weight: '200',
      style: "normal"
    }
  ],
  variable: '--font-nexa',
  fallback: ['ui-sans-serif'],
});

export const quicksand = localFont({
  src: [
    {
      path: './Quicksand-Bold.woff2',
      weight: '800',
      style: "normal"
    },
    {
      path: './Quicksand-Light.woff2',
      weight: '200',
      style: "normal"
    },
    {
      path: './Quicksand-Medium.woff2',
      weight: '600',
      style: "normal"
    },
    {
      path: './Quicksand-Regular.woff2',
      weight: '400',
      style: "normal"
    }
  ],
  variable: '--font-quicksand',
  fallback: ['ui-sans-serif'],
});