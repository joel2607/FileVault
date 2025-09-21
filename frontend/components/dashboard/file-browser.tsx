"use client"

import type React from "react"

import { useState } from "react"
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Alert,
  Chip,
} from "@mui/material"
import {
  Folder as FolderIcon,
  MoreVert,
  Edit,
  Download,
  Share,
  Delete,
  Public,
  CreateNewFolder,
} from "@mui/icons-material"
import { useQuery, useMutation } from "@apollo/client"
import { ROOT_QUERY, FOLDER_QUERY } from "@/lib/graphql/queries"
import {
  UPDATE_FILE_MUTATION,
  UPDATE_FOLDER_MUTATION,
  DELETE_FILE_MUTATION,
  DELETE_FOLDER_MUTATION,
  GENERATE_DOWNLOAD_URL_MUTATION,
  CREATE_FOLDER_MUTATION,
} from "@/lib/graphql/mutations"
import type { File, Folder } from "@/lib/types"
import { DashboardBreadcrumbs } from "./breadcrumbs"

interface FileBrowserProps {
  onShareFile: (file: File) => void
}

export function FileBrowser({ onShareFile }: FileBrowserProps) {
  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>()
  const [currentPath, setCurrentPath] = useState<Folder[]>([])
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedItem, setSelectedItem] = useState<File | Folder | null>(null)
  const [renameDialogOpen, setRenameDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [createFolderDialogOpen, setCreateFolderDialogOpen] = useState(false)
  const [newName, setNewName] = useState("")
  const [newFolderName, setNewFolderName] = useState("")
  const [error, setError] = useState("")

  // Queries
  const {
    data: rootData,
    loading: rootLoading,
    refetch: refetchRoot,
  } = useQuery(ROOT_QUERY, {
    skip: !!currentFolderId,
  })

  const {
    data: folderData,
    loading: folderLoading,
    refetch: refetchFolder,
  } = useQuery(FOLDER_QUERY, {
    variables: { id: currentFolderId! },
    skip: !currentFolderId,
  })

  // Mutations
  const [updateFile] = useMutation(UPDATE_FILE_MUTATION)
  const [updateFolder] = useMutation(UPDATE_FOLDER_MUTATION)
  const [deleteFile] = useMutation(DELETE_FILE_MUTATION)
  const [deleteFolder] = useMutation(DELETE_FOLDER_MUTATION)
  const [generateDownloadUrl] = useMutation(GENERATE_DOWNLOAD_URL_MUTATION)
  const [createFolder] = useMutation(CREATE_FOLDER_MUTATION)

  const currentData = currentFolderId ? folderData?.folder : rootData?.root
  const loading = currentFolderId ? folderLoading : rootLoading

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>, item: File | Folder) => {
    setAnchorEl(event.currentTarget)
    setSelectedItem(item)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedItem(null)
  }

  const handleFolderClick = (folder: Folder) => {
    setCurrentFolderId(folder.id)
    setCurrentPath([...currentPath, folder])
  }

  const handleNavigate = (folderId?: string) => {
    if (!folderId) {
      // Navigate to root
      setCurrentFolderId(undefined)
      setCurrentPath([])
    } else {
      // Navigate to specific folder in path
      const folderIndex = currentPath.findIndex((f) => f.id === folderId)
      if (folderIndex !== -1) {
        setCurrentFolderId(folderId)
        setCurrentPath(currentPath.slice(0, folderIndex + 1))
      }
    }
  }

  const handleRename = () => {
    if (!selectedItem) return
    setNewName("fileName" in selectedItem ? selectedItem.fileName : selectedItem.folderName)
    setRenameDialogOpen(true)
    handleMenuClose()
  }

  const handleDownload = async () => {
    if (!selectedItem || !("fileName" in selectedItem)) return

    try {
      const { data } = await generateDownloadUrl({
        variables: { fileID: selectedItem.id },
      })

      if (data?.generateDownloadUrl) {
        const link = document.createElement("a")
        link.href = data.generateDownloadUrl
        link.download = selectedItem.fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      }
    } catch (error) {
      console.error("Download failed:", error)
    }

    handleMenuClose()
  }

  const handleShare = () => {
    if (!selectedItem || !("fileName" in selectedItem)) return
    onShareFile(selectedItem)
    handleMenuClose()
  }

  const handleDelete = () => {
    setDeleteDialogOpen(true)
    handleMenuClose()
  }

  const confirmRename = async () => {
    if (!selectedItem || !newName.trim()) return

    try {
      if ("fileName" in selectedItem) {
        await updateFile({
          variables: {
            input: {
              id: selectedItem.id,
              name: newName.trim(),
            },
          },
        })
      } else {
        await updateFolder({
          variables: {
            input: {
              id: selectedItem.id,
              name: newName.trim(),
            },
          },
        })
      }

      // Refetch current data
      if (currentFolderId) {
        refetchFolder()
      } else {
        refetchRoot()
      }

      setRenameDialogOpen(false)
      setNewName("")
    } catch (error: any) {
      setError(error.message)
    }
  }

  const confirmDelete = async () => {
    if (!selectedItem) return

    try {
      if ("fileName" in selectedItem) {
        await deleteFile({
          variables: { id: selectedItem.id },
        })
      } else {
        await deleteFolder({
          variables: { id: selectedItem.id },
        })
      }

      // Refetch current data
      if (currentFolderId) {
        refetchFolder()
      } else {
        refetchRoot()
      }

      setDeleteDialogOpen(false)
    } catch (error: any) {
      setError(error.message)
    }
  }

  const handleCreateFolder = async () => {
    if (!newFolderName.trim()) return

    try {
      await createFolder({
        variables: {
          input: {
            name: newFolderName.trim(),
            parentFolderID: currentFolderId,
          },
        },
      })

      // Refetch current data
      if (currentFolderId) {
        refetchFolder()
      } else {
        refetchRoot()
      }

      setCreateFolderDialogOpen(false)
      setNewFolderName("")
    } catch (error: any) {
      setError(error.message)
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

  if (loading) {
    return <Typography>Loading...</Typography>
  }

  return (
    <Box>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        <DashboardBreadcrumbs currentPath={currentPath} onNavigate={handleNavigate} />
        <Button variant="outlined" startIcon={<CreateNewFolder />} onClick={() => setCreateFolderDialogOpen(true)}>
          New Folder
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError("")}>
          {error}
        </Alert>
      )}

      <Grid container spacing={2}>
        {/* Folders */}
        {currentData?.folders?.map((folder: Folder) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={folder.id}>
            <Card
              sx={{
                cursor: "pointer",
                "&:hover": { elevation: 4 },
                position: "relative",
              }}
            >
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                  <FolderIcon
                    sx={{ mr: 1, fontSize: 40, color: "primary.main" }}
                    onClick={() => handleFolderClick(folder)}
                  />
                  <Box sx={{ flexGrow: 1 }} onClick={() => handleFolderClick(folder)}>
                    <Typography variant="subtitle2" noWrap>
                      {folder.folderName}
                    </Typography>
                  </Box>
                  <IconButton size="small" onClick={(e) => handleMenuClick(e, folder)}>
                    <MoreVert />
                  </IconButton>
                </Box>
                <Box sx={{ display: "flex", gap: 1 }}>
                  {folder.isPublic && <Chip icon={<Public />} label="Public" size="small" color="success" />}
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}

        {/* Files */}
        {currentData?.files?.map((file: File) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={file.id}>
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
                  <IconButton size="small" onClick={(e) => handleMenuClick(e, file)}>
                    <MoreVert />
                  </IconButton>
                </Box>
                <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                  {file.isPublic && <Chip icon={<Public />} label="Public" size="small" color="success" />}
                  {file.downloadCount > 0 && (
                    <Chip label={`${file.downloadCount} downloads`} size="small" variant="outlined" />
                  )}
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Context Menu */}
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
        <MenuItem onClick={handleRename}>
          <Edit sx={{ mr: 1 }} />
          Rename
        </MenuItem>
        {selectedItem && "fileName" in selectedItem && (
          <>
            <MenuItem onClick={handleDownload}>
              <Download sx={{ mr: 1 }} />
              Download
            </MenuItem>
            <MenuItem onClick={handleShare}>
              <Share sx={{ mr: 1 }} />
              Share
            </MenuItem>
          </>
        )}
        <MenuItem onClick={handleDelete} sx={{ color: "error.main" }}>
          <Delete sx={{ mr: 1 }} />
          Delete
        </MenuItem>
      </Menu>

      {/* Rename Dialog */}
      <Dialog open={renameDialogOpen} onClose={() => setRenameDialogOpen(false)}>
        <DialogTitle>Rename Item</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="New Name"
            fullWidth
            variant="outlined"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRenameDialogOpen(false)}>Cancel</Button>
          <Button onClick={confirmRename} variant="contained">
            Rename
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "
            {selectedItem && ("fileName" in selectedItem ? selectedItem.fileName : selectedItem?.folderName)}"? This
            action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={confirmDelete} variant="contained" color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Folder Dialog */}
      <Dialog open={createFolderDialogOpen} onClose={() => setCreateFolderDialogOpen(false)}>
        <DialogTitle>Create New Folder</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Folder Name"
            fullWidth
            variant="outlined"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateFolderDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleCreateFolder} variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
