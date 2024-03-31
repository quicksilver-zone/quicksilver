const MillionLint = require('@million/lint');
/** @type {import('next').NextConfig} */

module.exports = MillionLint.next()({
  reactStrictMode: true,
  swcMinify: true,
  async redirects() {
    return [{
      source: '/',
      destination: '/staking',
      permanent: true
    }];
  }
});