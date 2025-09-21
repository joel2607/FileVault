/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  async rewrites() {
    return [
      {
        source: '/query',
        // The destination is where your docker-compose exposes the backend port
        destination: 'http://localhost:8080/query', 
      },
    ];
  },
};

export default nextConfig
