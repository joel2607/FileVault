"use client"

import { useState } from "react"
import { ProtectedRoute } from "@/components/protected-route"
import { DashboardLayout } from "@/components/dashboard/dashboard-layout"
import { FileBrowser } from "@/components/dashboard/file-browser"
import { ShareModal } from "@/components/modals/share-modal"
import { Typography, Box } from "@mui/material"
import type { File } from "@/lib/types"

export default function HomePage() {
  const [shareFile, setShareFile] = useState<File | null>(null)
  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>()
  const [refreshKey, setRefreshKey] = useState(0)

  const handleShareFile = (file: File) => {
    setShareFile(file)
  }

  const handleUploadComplete = () => {
    setRefreshKey((prev) => prev + 1)
  }

  const handleFileUpdate = () => {
    setRefreshKey((prev) => prev + 1)
  }

  return (
    <ProtectedRoute>
      <DashboardLayout currentFolderId={currentFolderId} onUploadComplete={handleUploadComplete}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            My Files
          </Typography>
          <FileBrowser key={refreshKey} onShareFile={handleShareFile} />
        </Box>

        <ShareModal
          open={!!shareFile}
          onClose={() => setShareFile(null)}
          file={shareFile}
          onFileUpdate={handleFileUpdate}
        />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
