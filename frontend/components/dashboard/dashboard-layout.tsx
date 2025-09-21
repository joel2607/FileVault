"use client"

import type React from "react"

import { Box, Toolbar } from "@mui/material"
import { DashboardHeader } from "./header"
import { Sidebar } from "./sidebar"
import { UploadModal } from "@/components/modals/upload-modal"
import { useState } from "react"

const DRAWER_WIDTH = 280

interface DashboardLayoutProps {
  children: React.ReactNode
  currentFolderId?: string
  onUploadComplete?: () => void
}

export function DashboardLayout({ children, currentFolderId, onUploadComplete }: DashboardLayoutProps) {
  const [uploadModalOpen, setUploadModalOpen] = useState(false)

  const handleUploadClick = () => {
    setUploadModalOpen(true)
  }

  const handleUploadComplete = () => {
    onUploadComplete?.()
    setUploadModalOpen(false)
  }

  return (
    <Box sx={{ display: "flex" }}>
      <DashboardHeader />
      <Sidebar onUploadClick={handleUploadClick} />

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${DRAWER_WIDTH}px)` },
        }}
      >
        <Toolbar />
        {children}
      </Box>

      <UploadModal
        open={uploadModalOpen}
        onClose={() => setUploadModalOpen(false)}
        currentFolderId={currentFolderId}
        onUploadComplete={handleUploadComplete}
      />
    </Box>
  )
}
