/** @type {import('next').NextConfig} */

module.exports = {
  reactStrictMode: true,
  swcMinify: true,

  async redirects() {
    return [
      {
        source: '/',
        destination: '/staking',
        permanent: true,
      },
    ];
  },
};
