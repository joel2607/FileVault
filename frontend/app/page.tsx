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

  const handleShareFile = (file: File) => {
    setShareFile(file)
  }

  return (
    <ProtectedRoute>
      <DashboardLayout currentFolderId={currentFolderId}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            My Files
          </Typography>
          <FileBrowser onShareFile={handleShareFile} />
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
