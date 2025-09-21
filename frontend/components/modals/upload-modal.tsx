"use client"

import type React from "react"

import { useState, useCallback } from "react"
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  LinearProgress,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Alert,
} from "@mui/material"
import { CloudUpload, Delete, AttachFile } from "@mui/icons-material"
import { useMutation } from "@apollo/client"
import { UPLOAD_FILES_MUTATION } from "@/lib/graphql/mutations"

interface UploadModalProps {
  open: boolean
  onClose: () => void
  currentFolderId?: string
  onUploadComplete: () => void
}

export function UploadModal({ open, onClose, currentFolderId, onUploadComplete }: UploadModalProps) {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([])
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [error, setError] = useState("")

  const [uploadFiles] = useMutation(UPLOAD_FILES_MUTATION)

  const handleFileSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || [])
    setSelectedFiles((prev) => [...prev, ...files])
  }, [])

  const handleDrop = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    const files = Array.from(event.dataTransfer.files)
    setSelectedFiles((prev) => [...prev, ...files])
  }, [])

  const handleDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault()
  }, [])

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index))
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes"
    const k = 1024
    const sizes = ["Bytes", "KB", "MB", "GB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
  }

  const handleUpload = async () => {
    if (selectedFiles.length === 0) return

    setUploading(true)
    setError("")
    setUploadProgress(0)

    try {
      // Simulate progress for demo purposes
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval)
            return prev
          }
          return prev + 10
        })
      }, 200)

      await uploadFiles({
        variables: {
          files: selectedFiles,
          parentFolderID: currentFolderId,
        },
      })

      clearInterval(progressInterval)
      setUploadProgress(100)

      setTimeout(() => {
        onUploadComplete()
        handleClose()
      }, 500)
    } catch (error: any) {
      setError(error.message)
      setUploading(false)
      setUploadProgress(0)
    }
  }

  const handleClose = () => {
    if (!uploading) {
      setSelectedFiles([])
      setUploadProgress(0)
      setError("")
      onClose()
    }
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>Upload Files</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {/* Drop Zone */}
        <Box
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          sx={{
            border: "2px dashed",
            borderColor: "primary.main",
            borderRadius: 2,
            p: 4,
            textAlign: "center",
            mb: 2,
            cursor: "pointer",
            "&:hover": {
              backgroundColor: "action.hover",
            },
          }}
          onClick={() => document.getElementById("file-input")?.click()}
        >
          <CloudUpload sx={{ fontSize: 48, color: "primary.main", mb: 2 }} />
          <Typography variant="h6" gutterBottom>
            Drag and drop files here
          </Typography>
          <Typography variant="body2" color="text.secondary">
            or click to select files
          </Typography>
          <input id="file-input" type="file" multiple hidden onChange={handleFileSelect} />
        </Box>

        {/* Selected Files List */}
        {selectedFiles.length > 0 && (
          <Box>
            <Typography variant="subtitle1" gutterBottom>
              Selected Files ({selectedFiles.length})
            </Typography>
            <List dense>
              {selectedFiles.map((file, index) => (
                <ListItem
                  key={index}
                  secondaryAction={
                    !uploading && (
                      <IconButton edge="end" onClick={() => removeFile(index)}>
                        <Delete />
                      </IconButton>
                    )
                  }
                >
                  <AttachFile sx={{ mr: 1 }} />
                  <ListItemText primary={file.name} secondary={formatFileSize(file.size)} />
                </ListItem>
              ))}
            </List>
          </Box>
        )}

        {/* Upload Progress */}
        {uploading && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" gutterBottom>
              Uploading... {uploadProgress}%
            </Typography>
            <LinearProgress variant="determinate" value={uploadProgress} />
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={uploading}>
          Cancel
        </Button>
        <Button onClick={handleUpload} variant="contained" disabled={selectedFiles.length === 0 || uploading}>
          {uploading ? "Uploading..." : `Upload ${selectedFiles.length} file${selectedFiles.length !== 1 ? "s" : ""}`}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
