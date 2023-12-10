import { extendTheme } from '@chakra-ui/react';

const defaultThemeObject = {
  config: {
    initialColorMode: 'light',
    useSystemColorMode: false,
  },
  styles: {
    global: (props: { colorMode: string }) => ({
      body: {
        background: props.colorMode === 'dark' ? '#000000' : '#000000',
        bgGradient:
          props.colorMode === 'dark'
            ? 'linear(to-r, #000000, #808080)'
            : 'linear(to-l, #000000, #000000)',
        color:
          props.colorMode === 'dark' ? 'rgb(255, 255, 255)' : 'rgb(0, 0, 0)',
      },
    }),
  },
  colors: {
    primary: {
      dark: '#333333',
      light: '#adadad',
      50: '#E6E6E6',
      100: '#CCCCCC',
      200: '#B3B3B3',
      300: '#999999',
      400: '#808080',
      500: '#666666',
      600: '#4D4D4D',
      700: '#333333',
      800: '#1A1A1A',
      900: '#000000',
    },
    complimentary: {
      50: '#FFF2E6',
      100: '#FFE6CC',
      200: '#FFD9B3',
      300: '#FFCC99',
      400: '#FFBF80',
      500: '#FFB266',
      600: '#FFA54D',
      700: '#FF9933',
      800: '#FF8C1A',
      900: '#FF8000',
      1000: '#b35a02',
    },
    background: {
      start: 'rgb(214, 219, 220)',
      end: 'rgb(255, 255, 255)',
      darkStart: 'rgb(0, 0, 0)',
      darkEnd: 'rgb(0, 0, 0)',
    },
    text: {
      light: 'white',
      dark: 'rgb(255, 255, 255)',
    },
    swiper: '#007aff',
    lightText: '#fbfbfb',
    tableBackground: 'transparent',
    tile: {
      start: 'rgb(239, 245, 249)',
      end: 'rgb(228, 232, 233)',
      darkStart: 'rgb(2, 13, 46)',
      darkEnd: 'rgb(2, 5, 19)',
    },
    callout: {
      main: 'rgb(238, 240, 241)',
      border: 'rgb(172, 175, 176)',
      darkMain: 'rgb(20, 20, 20)',
      darkBorder: 'rgb(108, 108, 108)',
    },
    card: {
      main: 'rgb(180, 185, 188)',
      border: 'rgb(131, 134, 135)',
      darkMain: 'rgb(100, 100, 100)',
      darkBorder: 'rgb(200, 200, 200)',
    },
    primaryGlow: {
      light:
        'conic-gradient(from 180deg at 50% 50%, #16abff33 0deg, #0885ff33 55deg, #54d6ff33 120deg, #0071ff33 160deg, transparent 360deg)',
      dark: 'radial-gradient(rgba(1, 65, 255, 0.4), rgba(1, 65, 255, 0))',
    },
    secondaryGlow: {
      light: 'radial-gradient(rgba(255, 255, 255, 1), rgba(255, 255, 255, 0))',
      dark: 'linear-gradient(to bottom right, rgba(1, 65, 255, 0), rgba(1, 65, 255, 0), rgba(1, 65, 255, 0.3))',
    },
  },
  fonts: {
    heading: 'Lato, Poppins',
    body: 'Lato, Poppins',
  },
  textStyles: {
    h1: {
      fontWeight: 'bold',
      fontSize: '2xl',
      letterSpacing: '-0.1rem',
      lineHeight: '1.2',
    },
    h2: {
      fontWeight: 'semibold',
      fontSize: 'xl',
      letterSpacing: '-0.05rem',
      lineHeight: '1.2',
    },
  },
  components: {
    Tooltip: {
      baseStyle: {
        fontSize: '1em',
        bgColor: 'primary.700',
        color: 'primary.50',
        borderRadius: '8px',
        px: '0.75em',
        py: '0.5em',
      },
    },
    Button: {
      baseStyle: {
        fontWeight: 'bold',
      },
      variants: {
        solid: {
          bgColor: 'complimentary.900',
          color: 'white',
        },
        outline: {
          borderColor: 'complimentary.900',
          color: 'primary.600',
          _hover: {
            color: 'complimentary.300',
          },
          _active: {
            bg: 'complimentary.300',
          },
        },
      },
    },
    Box: {
      baseStyle: {
        boxShadow:
          '0 4px 6px rgba(255, 190, 190, 0.1), 0 1px 3px rgba(255, 190, 190, 0.1)',
        width: 'md',
        height: 'md',
      },
    },
    Flex: {
      baseStyle: {
        boxShadow:
          '0 4px 6px rgba(255, 190, 190, 0.1), 0 1px 3px rgba(255, 190, 190, 0.1)',
      },
    },
    Text: {
      baseStyle: {
        color: 'text.light',
      },
    },
  },
};

export const defaultTheme = extendTheme(defaultThemeObject);
