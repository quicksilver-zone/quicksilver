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
        permanent: true, // Change to false if temporary redirect
      },
    ];
  },
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  },
};
