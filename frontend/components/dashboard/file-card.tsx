"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { Card, CardContent, Typography, IconButton, Box, Chip, Menu, MenuItem } from "@mui/material"
import { MoreVert, Edit, Download, Share, Delete, Public, MoveUp } from "@mui/icons-material"
import { useSubscription, useMutation } from "@apollo/client"
import { FILE_DOWNLOAD_COUNT_SUBSCRIPTION } from "@/lib/graphql/subscriptions"
import { GENERATE_DOWNLOAD_URL_MUTATION } from "@/lib/graphql/mutations"
import type { File } from "@/lib/types"

interface FileCardProps {
  file: File
  onRename: (file: File) => void
  onShare: (file: File) => void
  onDelete: (file: File) => void
  onMove: (file: File) => void
}

export function FileCard({ file, onRename, onShare, onDelete, onMove }: FileCardProps) {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [downloadCount, setDownloadCount] = useState(file.downloadCount)
  const [downloading, setDownloading] = useState(false)

  const [generateDownloadUrl] = useMutation(GENERATE_DOWNLOAD_URL_MUTATION)

  const { data: subscriptionData } = useSubscription(FILE_DOWNLOAD_COUNT_SUBSCRIPTION, {
    variables: { fileID: file.id },
  })

  useEffect(() => {
    if (subscriptionData?.fileDownloadCount) {
      setDownloadCount(subscriptionData.fileDownloadCount.downloadCount)
    }
  }, [subscriptionData])

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleDownload = async () => {
    setDownloading(true)
    handleMenuClose()

    try {
      const { data } = await generateDownloadUrl({
        variables: { fileID: file.id },
      })

      if (data?.generateDownloadUrl) {
        const link = document.createElement("a")
        link.href = `${data.generateDownloadUrl}`
        link.download = file.fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      }
    } catch (error) {
      console.error("Download failed:", error)
    } finally {
      setDownloading(false)
    }
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes"
    const k = 1024
    const sizes = ["Bytes", "KB", "MB", "GB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
  }

  const getFileIcon = (mimeType: string) => {
    if (mimeType.startsWith("image/")) return "🖼️"
    if (mimeType.startsWith("video/")) return "🎥"
    if (mimeType.startsWith("audio/")) return "🎵"
    if (mimeType.includes("pdf")) return "📄"
    if (mimeType.includes("document") || mimeType.includes("word")) return "📝"
    if (mimeType.includes("spreadsheet") || mimeType.includes("excel")) return "📊"
    return "📄"
  }

  return (
    <>
      <Card sx={{ "&:hover": { elevation: 4 }, position: "relative" }}>
        <CardContent>
          <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
            <Box sx={{ mr: 1, fontSize: 32 }}>{getFileIcon(file.mimeType)}</Box>
            <Box sx={{ flexGrow: 1, minWidth: 0 }}>
              <Typography variant="subtitle2" noWrap>
                {file.fileName}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {formatFileSize(file.size)}
              </Typography>
            </Box>
            <IconButton size="small" onClick={handleMenuClick} disabled={downloading}>
              <MoreVert />
            </IconButton>
          </Box>
          <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
            {file.isPublic && <Chip icon={<Public />} label="Public" size="small" color="success" />}
            {downloadCount > 0 && (
              <Chip
                label={`${downloadCount} download${downloadCount !== 1 ? "s" : ""}`}
                size="small"
                variant="outlined"
              />
            )}
            {downloading && <Chip label="Downloading..." size="small" color="primary" />}
          </Box>
        </CardContent>
      </Card>

      {/* Context Menu */}
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
        <MenuItem
          onClick={() => {
            onRename(file)
            handleMenuClose()
          }}
        >
          <Edit sx={{ mr: 1 }} />
          Rename
        </MenuItem>
        <MenuItem onClick={handleDownload} disabled={downloading}>
          <Download sx={{ mr: 1 }} />
          {downloading ? "Downloading..." : "Download"}
        </MenuItem>
        <MenuItem
          onClick={() => {
            onMove(file)
            handleMenuClose()
          }}
        >
          <MoveUp sx={{ mr: 1 }} />
          Move
        </MenuItem>
        <MenuItem
          onClick={() => {
            onShare(file)
            handleMenuClose()
          }}
        >
          <Share sx={{ mr: 1 }} />
          Share
        </MenuItem>
        <MenuItem
          onClick={() => {
            onDelete(file)
            handleMenuClose()
          }}
          sx={{ color: "error.main" }}
        >
          <Delete sx={{ mr: 1 }} />
          Delete
        </MenuItem>
      </Menu>
    </>
  )
}
