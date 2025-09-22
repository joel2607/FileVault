"use client"

import { useState, useEffect } from "react"
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
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Avatar,
  Autocomplete,
  CircularProgress,
} from "@mui/material"
import { ContentCopy, Public, Lock, Person, Delete } from "@mui/icons-material"
import { useMutation, useQuery, useLazyQuery } from "@apollo/client"
import {
  SET_FILE_PUBLIC_MUTATION,
  SET_FILE_PRIVATE_MUTATION,
  SHARE_FILE_WITH_USER_MUTATION,
  REMOVE_FILE_ACCESS_MUTATION,
} from "@/lib/graphql/mutations"
import { GET_FILE_ACCESS_QUERY, SEARCH_USERS_QUERY } from "@/lib/graphql/queries"
import type { File } from "@/lib/types"
import { useDebounce } from "@/hooks/use-debounce"

interface ShareModalProps {
  open: boolean
  onClose: () => void
  file: File | null
  onFileUpdate: () => void
}

interface User {
  id: string
  username: string
  email: string
}

interface FileAccess {
  id: string
  user: User
  file: {
    id: string
    fileName: string
  }
}

export function ShareModal({ open, onClose, file, onFileUpdate }: ShareModalProps) {
  const [activeTab, setActiveTab] = useState(0)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [error, setError] = useState("")
  const [copySuccess, setCopySuccess] = useState(false)
  const [searchTerm, setSearchTerm] = useState("")
  const debouncedSearchTerm = useDebounce(searchTerm, 300)

  // Queries and mutations
  const { data: accessData, refetch: refetchAccess } = useQuery(GET_FILE_ACCESS_QUERY, {
    variables: { fileID: file?.id },
    skip: !file?.id || !open,
  })

  const [searchUsers, { data: usersData, loading: searchingUsers }] = useLazyQuery(SEARCH_USERS_QUERY)
  const [setFilePublic] = useMutation(SET_FILE_PUBLIC_MUTATION)
  const [setFilePrivate] = useMutation(SET_FILE_PRIVATE_MUTATION)
  const [shareFileWithUser] = useMutation(SHARE_FILE_WITH_USER_MUTATION)
  const [removeFileAccess] = useMutation(REMOVE_FILE_ACCESS_MUTATION)

  useEffect(() => {
    if (debouncedSearchTerm.length > 2) {
      searchUsers({ variables: { query: debouncedSearchTerm } })
    }
  }, [debouncedSearchTerm, searchUsers])

  if (!file) return null

  const shareUrl = `${window.location.origin}/shared/${file.id}`
  const fileAccess: FileAccess[] = accessData?.getFileAccess || []

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

  const handleShareWithUser = async () => {
    if (!selectedUser) return

    try {
      setError("")
      await shareFileWithUser({
        variables: {
          fileID: file.id,
          userID: selectedUser.id,
        },
      })

      setSelectedUser(null)
      setSearchTerm("")
      refetchAccess()
    } catch (error: any) {
      setError(error.message)
    }
  }

  const handleRevokeAccess = async (userId: string) => {
    try {
      setError("")
      await removeFileAccess({
        variables: {
          fileID: file.id,
          userID: userId,
        },
      })
      refetchAccess()
    } catch (error: any) {
      setError(error.message)
    }
  }

  const handleClose = () => {
    setActiveTab(0)
    setSelectedUser(null)
    setSearchTerm("")
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

            <Autocomplete
              options={usersData?.searchUsers || []}
              getOptionLabel={(option) => `${option.username} (${option.email})`}
              value={selectedUser}
              onChange={(_, newValue) => setSelectedUser(newValue)}
              inputValue={searchTerm}
              onInputChange={(_, newInputValue) => setSearchTerm(newInputValue)}
              loading={searchingUsers}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Search users by email"
                  margin="normal"
                  fullWidth
                  InputProps={{
                    ...params.InputProps,
                    endAdornment: (
                      <>
                        {searchingUsers ? <CircularProgress color="inherit" size={20} /> : null}
                        {params.InputProps.endAdornment}
                      </>
                    ),
                  }}
                />
              )}
              renderOption={(props, option) => (
                <Box component="li" {...props}>
                  <Avatar sx={{ mr: 2, width: 32, height: 32 }}>{option.username.charAt(0).toUpperCase()}</Avatar>
                  <Box>
                    <Typography variant="body2">{option.username}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {option.email}
                    </Typography>
                  </Box>
                </Box>
              )}
            />

            <Button variant="contained" onClick={handleShareWithUser} disabled={!selectedUser} sx={{ mt: 2 }}>
              Share
            </Button>

            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle2" gutterBottom>
                People with Access ({fileAccess.length})
              </Typography>

              {fileAccess.length > 0 ? (
                <List dense>
                  {fileAccess.map((access) => (
                    <ListItem key={access.id}>
                      <Avatar sx={{ mr: 2, width: 32, height: 32 }}>
                        {access.user.username.charAt(0).toUpperCase()}
                      </Avatar>
                      <ListItemText primary={access.user.username} secondary={access.user.email} />
                      <ListItemSecondaryAction>
                        <IconButton edge="end" onClick={() => handleRevokeAccess(access.user.id)} color="error">
                          <Delete />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))}
                </List>
              ) : (
                <Typography variant="body2" color="text.secondary">
                  No users have been granted access to this file.
                </Typography>
              )}
            </Box>

            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle2" gutterBottom>
                Current Visibility
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
