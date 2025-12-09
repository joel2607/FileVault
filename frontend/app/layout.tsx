import type React from "react"
import type { Metadata } from "next"
import { GeistSans } from "geist/font/sans"
import { GeistMono } from "geist/font/mono"
import { Analytics } from "@vercel/analytics/next"
import { CustomThemeProvider } from "@/contexts/theme-context"
import { CustomApolloProvider } from "@/contexts/apollo-provider"
import { Suspense } from "react"
import "./globals.css"

export const metadata: Metadata = {
  title: "FileVault - File Management",
  description: "Modern file management application",
  generator: "v0.app",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  console.log("Browser sees endpoint as:", process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT);
  return (
    <html lang="en">
      <body className={`font-sans ${GeistSans.variable} ${GeistMono.variable}`}>
        <Suspense fallback={null}>
          <CustomApolloProvider>
            <CustomThemeProvider>{children}</CustomThemeProvider>
          </CustomApolloProvider>
          <Analytics />
        </Suspense>
      </body>
    </html>
  )
}
