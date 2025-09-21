"use client"

import { useState } from "react"
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  TextField,
  Switch,
  FormControlLabel,
  Tabs,
  Tab,
  Alert,
  InputAdornment,
  IconButton,
  Chip,
} from "@mui/material"
import { ContentCopy, Public, Lock, Person } from "@mui/icons-material"
import { useMutation } from "@apollo/client"
import { SET_FILE_PUBLIC_MUTATION, SET_FILE_PRIVATE_MUTATION } from "@/lib/graphql/mutations"
import type { File } from "@/lib/types"

interface ShareModalProps {
  open: boolean
  onClose: () => void
  file: File | null
  onFileUpdate: () => void
}

export function ShareModal({ open, onClose, file, onFileUpdate }: ShareModalProps) {
  const [activeTab, setActiveTab] = useState(0)
  const [userEmail, setUserEmail] = useState("")
  const [error, setError] = useState("")
  const [copySuccess, setCopySuccess] = useState(false)

  const [setFilePublic] = useMutation(SET_FILE_PUBLIC_MUTATION)
  const [setFilePrivate] = useMutation(SET_FILE_PRIVATE_MUTATION)

  if (!file) return null

  const shareUrl = `${window.location.origin}/shared/${file.id}`

  const handleVisibilityToggle = async (isPublic: boolean) => {
    try {
      setError("")
      if (isPublic) {
        await setFilePublic({ variables: { fileID: file.id } })
      } else {
        await setFilePrivate({ variables: { fileID: file.id } })
      }
      onFileUpdate()
    } catch (error: any) {
      setError(error.message)
    }
  }

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(shareUrl)
      setCopySuccess(true)
      setTimeout(() => setCopySuccess(false), 2000)
    } catch (error) {
      setError("Failed to copy link to clipboard")
    }
  }

  const handleShareWithUser = () => {
    // This would typically call the shareFileWithUser mutation
    console.log("Share with user:", userEmail)
    setUserEmail("")
  }

  const handleClose = () => {
    setActiveTab(0)
    setUserEmail("")
    setError("")
    setCopySuccess(false)
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Share "{file.fileName}"</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {copySuccess && (
          <Alert severity="success" sx={{ mb: 2 }}>
            Link copied to clipboard!
          </Alert>
        )}

        <Tabs value={activeTab} onChange={(_, newValue) => setActiveTab(newValue)} sx={{ mb: 3 }}>
          <Tab label="Share with People" icon={<Person />} />
          <Tab label="Get Link" icon={<Public />} />
        </Tabs>

        {activeTab === 0 && (
          <Box>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Share this file with specific people by entering their email address.
            </Typography>

            <TextField
              fullWidth
              label="Email address"
              value={userEmail}
              onChange={(e) => setUserEmail(e.target.value)}
              margin="normal"
              placeholder="Enter email to share with"
            />

            <Button variant="contained" onClick={handleShareWithUser} disabled={!userEmail.trim()} sx={{ mt: 2 }}>
              Share
            </Button>

            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle2" gutterBottom>
                Current Access
              </Typography>
              <Chip
                icon={file.isPublic ? <Public /> : <Lock />}
                label={file.isPublic ? "Public" : "Private"}
                color={file.isPublic ? "success" : "default"}
              />
            </Box>
          </Box>
        )}

        {activeTab === 1 && (
          <Box>
            <Box sx={{ mb: 3 }}>
              <FormControlLabel
                control={<Switch checked={file.isPublic} onChange={(e) => handleVisibilityToggle(e.target.checked)} />}
                label={
                  <Box>
                    <Typography variant="body1">{file.isPublic ? "Public" : "Private"}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {file.isPublic ? "Anyone with the link can view this file" : "Only you can access this file"}
                    </Typography>
                  </Box>
                }
              />
            </Box>

            {file.isPublic && (
              <Box>
                <Typography variant="subtitle2" gutterBottom>
                  Shareable Link
                </Typography>
                <TextField
                  fullWidth
                  value={shareUrl}
                  InputProps={{
                    readOnly: true,
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton onClick={handleCopyLink}>
                          <ContentCopy />
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                  margin="normal"
                />

                <Box sx={{ mt: 2 }}>
                  <Typography variant="body2" color="text.secondary">
                    Download count: {file.downloadCount}
                  </Typography>
                </Box>
              </Box>
            )}
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Close</Button>
      </DialogActions>
    </Dialog>
  )
}
