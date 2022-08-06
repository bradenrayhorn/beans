/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8000/api/:path*'
      }
    ]
  },
  async redirects() {
    return [
      {
        source: '/',
        destination: '/app',
        permanent: true,
      }
    ]
  }
}

module.exports = nextConfig
