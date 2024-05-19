import { ChakraProps, chakra, keyframes } from '@chakra-ui/react';
import { JSX, SVGProps } from 'react';

const scaleKeyframes = keyframes`
  0%, 100% {
    transform: scale(1);
  }
  30% {
    transform: scale(1.2);
  }
  50% {
    transform: scale(0.9);
  }
  70% {
    transform: scale(1.1);
  }
`;

const opacityKeyframes = keyframes`
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
`;

const rotateKeyframes = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
`;

const KadoIconContent = (
  props: JSX.IntrinsicAttributes &
    Omit<SVGProps<SVGSVGElement>, 'color' | 'height' | 'width' | 'size' | 'htmlTranslate' | 'as' | 'viewBox' | 'ref' | 'children'> & {
      htmlTranslate?: 'no' | 'yes' | undefined;
      showAnimation?: boolean;
      orange?: boolean;
    } & Omit<ChakraProps, never> & { as?: 'svg' | undefined },
) => {
  const KadoIconLoading = chakra('svg', {
    baseStyle: {
      display: 'inline-block',
      lineHeight: '1em',
      flexShrink: 0,
      color: 'currentColor',
      transition: 'transform 0.3s ease',
      '.segment': {
        animation: `${scaleKeyframes} 1s infinite, ${opacityKeyframes} 2s infinite, ${rotateKeyframes} 5s infinite linear`,
        transformOrigin: 'center center',
      },
    },
  });

  const KadoIcon = chakra('svg', {
    baseStyle: {
      display: 'inline-block',
      lineHeight: '1em',
      flexShrink: 0,
      color: 'currentColor',
      transition: 'transform 0.3s ease',
    },
  });
  {
    if (!props.showAnimation && !props.orange) {
      return (
        <KadoIcon viewBox="0 0 260 260" {...props}>
          <g className="segment" opacity="0.75">
            <path
              fill="#808080"
              d="M209.53,161.73a115,115,0,0,0,6.3-55.36c-4.48-32.06-21.62-56.69-53.53-74.17-40.85-22.38-96-14.49-129,18.43C15,68.91,3.92,99.22,18.62,120.49c10.66,15.42,31.23,21,43.33,35.32,15.14,17.88,14,45.49,29,63.46,13.33,15.93,37.32,20.27,57,13.74C164.46,227.55,194.25,203.45,209.53,161.73Z"
            />
          </g>
          <g className="segment" opacity="0.8">
            <path
              fill="#9ea3b4"
              d="M148,233c16.47-5.46,46.26-29.56,61.55-71.28,6.47-17.68,10-34.7,6.29-55.36C203.56,38.13,120,125.89,105.59,134.8c-28.44,17.58-36.29,56.12-16,82.71.47.61.93,1.2,1.4,1.76C104.29,235.2,128.28,239.54,148,233Z"
            />
          </g>
          <g className="segment" opacity="0.7">
            <path
              fill="#cedcf8"
              d="M174,42.78a115,115,0,0,0-55.36-6.3C86.55,41,61.92,58.1,44.44,90c-22.37,40.85-14.49,96,18.43,129,18.28,18.3,48.59,29.39,69.87,14.69,15.41-10.65,21-31.22,35.31-43.33,17.88-15.14,45.49-14,63.46-29,15.93-13.33,20.27-37.31,13.74-57C239.79,87.85,215.69,58.06,174,42.78Z"
            />
          </g>
          <path
            fill="#fff"
            d="M137.37,177.57l-16.23-32.81-10.9,10-4.86,22.86H86.79l17.76-84h18.72l-7,32.69,35.77-32.69h25.83l-42.17,38.14,23.92,45.84Z"
          />
        </KadoIcon>
      );
    }
    {
      if (props.showAnimation) {
        return (
          <KadoIconLoading viewBox="0 0 260 260" {...props}>
            <g className="segment" opacity="0.75">
              <path
                fill="#5493f7"
                d="M209.53,161.73a115,115,0,0,0,6.3-55.36c-4.48-32.06-21.62-56.69-53.53-74.17-40.85-22.38-96-14.49-129,18.43C15,68.91,3.92,99.22,18.62,120.49c10.66,15.42,31.23,21,43.33,35.32,15.14,17.88,14,45.49,29,63.46,13.33,15.93,37.32,20.27,57,13.74C164.46,227.55,194.25,203.45,209.53,161.73Z"
              />
            </g>
            <g className="segment" opacity="0.8">
              <path
                fill="#2043b5"
                d="M148,233c16.47-5.46,46.26-29.56,61.55-71.28,6.47-17.68,10-34.7,6.29-55.36C203.56,38.13,120,125.89,105.59,134.8c-28.44,17.58-36.29,56.12-16,82.71.47.61.93,1.2,1.4,1.76C104.29,235.2,128.28,239.54,148,233Z"
              />
            </g>
            <g className="segment" opacity="0.7">
              <path
                fill="#3573ec"
                d="M174,42.78a115,115,0,0,0-55.36-6.3C86.55,41,61.92,58.1,44.44,90c-22.37,40.85-14.49,96,18.43,129,18.28,18.3,48.59,29.39,69.87,14.69,15.41-10.65,21-31.22,35.31-43.33,17.88-15.14,45.49-14,63.46-29,15.93-13.33,20.27-37.31,13.74-57C239.79,87.85,215.69,58.06,174,42.78Z"
              />
            </g>
            <path
              fill="#fff"
              d="M137.37,177.57l-16.23-32.81-10.9,10-4.86,22.86H86.79l17.76-84h18.72l-7,32.69,35.77-32.69h25.83l-42.17,38.14,23.92,45.84Z"
            />
          </KadoIconLoading>
        );
      }
    }
    {
      if (props.orange) {
        return (
          <KadoIcon viewBox="0 0 260 260" {...props}>
            <g className="segment" opacity="0.75">
              <path
                fill="rgba(255, 168, 81, 0.9)"
                d="M209.53,161.73a115,115,0,0,0,6.3-55.36c-4.48-32.06-21.62-56.69-53.53-74.17-40.85-22.38-96-14.49-129,18.43C15,68.91,3.92,99.22,18.62,120.49c10.66,15.42,31.23,21,43.33,35.32,15.14,17.88,14,45.49,29,63.46,13.33,15.93,37.32,20.27,57,13.74C164.46,227.55,194.25,203.45,209.53,161.73Z"
              />
            </g>
            <g className="segment" opacity="0.8">
              <path
                fill="rgba(255, 128, 0, 0.9)"
                d="M148,233c16.47-5.46,46.26-29.56,61.55-71.28,6.47-17.68,10-34.7,6.29-55.36C203.56,38.13,120,125.89,105.59,134.8c-28.44,17.58-36.29,56.12-16,82.71.47.61.93,1.2,1.4,1.76C104.29,235.2,128.28,239.54,148,233Z"
              />
            </g>
            <g className="segment" opacity="0.7">
              <path
                fill="rgba(175, 88, 0, 0.9)"
                d="M174,42.78a115,115,0,0,0-55.36-6.3C86.55,41,61.92,58.1,44.44,90c-22.37,40.85-14.49,96,18.43,129,18.28,18.3,48.59,29.39,69.87,14.69,15.41-10.65,21-31.22,35.31-43.33,17.88-15.14,45.49-14,63.46-29,15.93-13.33,20.27-37.31,13.74-57C239.79,87.85,215.69,58.06,174,42.78Z"
              />
            </g>
            <path
              fill="#fff"
              d="M137.37,177.57l-16.23-32.81-10.9,10-4.86,22.86H86.79l17.76-84h18.72l-7,32.69,35.77-32.69h25.83l-42.17,38.14,23.92,45.84Z"
            />
          </KadoIcon>
        );
      }
    }
  }
};

export default KadoIconContent;
