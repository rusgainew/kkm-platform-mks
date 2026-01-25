import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  // Remove rewrites - use API routes instead for better control
  rewrites: async () => ({
    beforeFiles: [],
  }),
};

export default nextConfig;
