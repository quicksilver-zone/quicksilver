/** @type {import('next').NextConfig} */

module.exports = {
  reactStrictMode: true,
  swcMinify: true,
  basePath: '/quicksilver',
  assetPrefix: '/quicksilver',
  async redirects() {
    return [
      {
        source: '/',
        destination: '/staking',
        permanent: true,
      },
    ];
  },
  typescript: {
    // !! WARN !! //
    // There are no fatal errors in this project, this option is used as a workaround due to the amalgamation of packages we are using //
    // This option will be removed once all dependencies are updated to use the lates versions //
    ignoreBuildErrors: true,
  },
};
