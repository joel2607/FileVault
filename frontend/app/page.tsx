"use client"

import { useState } from "react"
import { ProtectedRoute } from "@/components/protected-route"
import { DashboardLayout } from "@/components/dashboard/dashboard-layout"
import { FileBrowser } from "@/components/dashboard/file-browser"
import { ShareModal } from "@/components/modals/share-modal"
import { Typography, Box } from "@mui/material"
import type { File } from "@/lib/types"
import { SearchBar } from "@/components/dashboard/search-bar"

export default function HomePage() {
  const [shareFile, setShareFile] = useState<File | null>(null)
  const [searchResults, setSearchResults] = useState<File[]>([])
  const [isSearching, setIsSearching] = useState(false)

  const handleShareFile = (file: File) => {
    setShareFile(file)
  }

  const handleSearchResults = (files: File[]) => {
    setSearchResults(files)
    setIsSearching(true)
  }

  const handleClearSearch = () => {
    setSearchResults([])
    setIsSearching(false)
  }

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            My Files
          </Typography>
          <SearchBar onSearchResults={handleSearchResults} onClearSearch={handleClearSearch} />
          <FileBrowser onShareFile={handleShareFile} isSearching={isSearching} searchResults={searchResults} />
        </Box>

        <ShareModal
          open={!!shareFile}
          onClose={() => setShareFile(null)}
          file={shareFile}
          onFileUpdate={() => {}}
        />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
