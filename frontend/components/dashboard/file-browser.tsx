"use client"

import type React from "react"

import { useState, useEffect } from "react"
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
import { Folder as FolderIcon, MoreVert, Edit, Delete, Public, CreateNewFolder, MoveUp } from "@mui/icons-material"
import { useQuery, useMutation } from "@apollo/client"
import { ROOT_QUERY, FOLDER_QUERY } from "@/lib/graphql/queries"
import {
  UPDATE_FILE_MUTATION,
  UPDATE_FOLDER_MUTATION,
  DELETE_FILE_MUTATION,
  DELETE_FOLDER_MUTATION,
  CREATE_FOLDER_MUTATION,
} from "@/lib/graphql/mutations"
import type { File, Folder } from "@/lib/types"
import { DashboardBreadcrumbs } from "./breadcrumbs"
import { SearchBar } from "./search-bar"
import { FileCard } from "./file-card"
import { MoveDialog } from "../modals/MoveDialog"

interface FileBrowserProps {
  onShareFile: (file: File) => void
  refresh: number
}

export function FileBrowser({ onShareFile, refresh }: FileBrowserProps) {
  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>()
  const [currentPath, setCurrentPath] = useState<Folder[]>([])
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedItem, setSelectedItem] = useState<File | Folder | null>(null)
  const [renameDialogOpen, setRenameDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [createFolderDialogOpen, setCreateFolderDialogOpen] = useState(false)
  const [moveDialogOpen, setMoveDialogOpen] = useState(false)
  const [newName, setNewName] = useState("")
  const [newFolderName, setNewFolderName] = useState("")
  const [error, setError] = useState("")
  const [searchResults, setSearchResults] = useState<File[] | null>(null)
  const [isSearching, setIsSearching] = useState(false)

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

  useEffect(() => {
    if (currentFolderId) {
      refetchFolder()
    } else {
      refetchRoot()
    }
  }, [refresh, currentFolderId, refetchFolder, refetchRoot])

  // Mutations
  const [updateFile] = useMutation(UPDATE_FILE_MUTATION)
  const [updateFolder] = useMutation(UPDATE_FOLDER_MUTATION)
  const [deleteFile] = useMutation(DELETE_FILE_MUTATION)
  const [deleteFolder] = useMutation(DELETE_FOLDER_MUTATION)
  const [createFolder] = useMutation(CREATE_FOLDER_MUTATION)

  const currentData = currentFolderId ? folderData?.folder : rootData?.root
  const loading = currentFolderId ? folderLoading : rootLoading

  const handleSearchResults = (files: File[]) => {
    setSearchResults(files)
    setIsSearching(true)
  }

  const handleClearSearch = () => {
    setSearchResults(null)
    setIsSearching(false)
  }

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>, item: Folder) => {
    setAnchorEl(event.currentTarget)
    setSelectedItem(item)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedItem(null)
  }

  const handleFolderClick = (folder: Folder) => {
    if (isSearching) {
      handleClearSearch()
    }
    setCurrentFolderId(folder.id)
    setCurrentPath([...currentPath, folder])
  }

  const handleNavigate = (folderId?: string) => {
    if (isSearching) {
      handleClearSearch()
    }

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

  const handleRenameFile = (file: File) => {
    setSelectedItem(file)
    setNewName(file.fileName)
    setRenameDialogOpen(true)
  }

  const handleRenameFolder = () => {
    if (!selectedItem) return
    setNewName("folderName" in selectedItem ? selectedItem.folderName : "")
    setRenameDialogOpen(true)
    setAnchorEl(null)
  }

  const handleDeleteFile = (file: File) => {
    setSelectedItem(file)
    setDeleteDialogOpen(true)
  }

  const handleDeleteFolder = () => {
    if (!selectedItem) return
    setDeleteDialogOpen(true)
    setAnchorEl(null)
  }

  const handleMove = (item: File | Folder) => {
    setSelectedItem(item)
    setMoveDialogOpen(true)
  }

  const closeRenameDialog = () => {
    setRenameDialogOpen(false)
    setSelectedItem(null)
  }

  const closeDeleteDialog = () => {
    setDeleteDialogOpen(false)
    setSelectedItem(null)
  }

  const closeMoveDialog = () => {
    setMoveDialogOpen(false)
    setSelectedItem(null)
  }

  const confirmRename = async () => {
    if (!selectedItem || !newName.trim()) return

    try {
      if ("fileName" in selectedItem) {
        await updateFile({
          variables: {
            input: {
              id: selectedItem.id,
              fileName: newName.trim(),
            },
          },
        })
      } else {
        await updateFolder({
          variables: {
            input: {
              id: selectedItem.id,
              folderName: newName.trim(),
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

      closeRenameDialog()
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

      closeDeleteDialog()
    } catch (error: any) {
      setError(error.message)
    }
  }

  const confirmMove = async (destinationFolderId: string | null) => {
    if (!selectedItem) return

    try {
      if ("fileName" in selectedItem) {
        await updateFile({
          variables: {
            input: {
              id: selectedItem.id,
              parentFolderID: destinationFolderId,
            },
          },
        })
      } else {
        await updateFolder({
          variables: {
            input: {
              id: selectedItem.id,
              parentFolderID: destinationFolderId,
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

      closeMoveDialog()
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
            folderName: newFolderName.trim(),
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

  if (loading) {
    return <Typography>Loading...</Typography>
  }

  const displayFiles = isSearching ? searchResults : currentData?.files
  const displayFolders = isSearching ? [] : currentData?.folders

  return (
    <Box>
      <SearchBar onSearchResults={handleSearchResults} onClearSearch={handleClearSearch} />

      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        {isSearching ? (
          <Typography variant="h6">Search Results ({searchResults?.length || 0} files found)</Typography>
        ) : (
          <DashboardBreadcrumbs currentPath={currentPath} onNavigate={handleNavigate} />
        )}

        {!isSearching && (
          <Button variant="outlined" startIcon={<CreateNewFolder />} onClick={() => setCreateFolderDialogOpen(true)}>
            New Folder
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError("")}>
          {error}
        </Alert>
      )}

      <Grid container spacing={2}>
        {/* Folders - only show when not searching */}
        {displayFolders?.map((folder: Folder) => (
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

        {/* Files - using enhanced FileCard component */}
        {displayFiles?.map((file: File) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={file.id}>
            <FileCard
              file={file}
              onRename={handleRenameFile}
              onShare={onShareFile}
              onDelete={handleDeleteFile}
              onMove={handleMove}
            />
          </Grid>
        ))}
      </Grid>

      {isSearching && (!searchResults || searchResults.length === 0) && (
        <Box sx={{ textAlign: "center", py: 4 }}>
          <Typography variant="body1" color="text.secondary">
            No files found matching your search criteria.
          </Typography>
        </Box>
      )}

      {/* Folder Context Menu */}
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
        <MenuItem onClick={handleRenameFolder}>
          <Edit sx={{ mr: 1 }} />
          Rename
        </MenuItem>
        <MenuItem onClick={() => handleMove(selectedItem!)}>
          <MoveUp sx={{ mr: 1 }} />
          Move
        </MenuItem>
        <MenuItem onClick={handleDeleteFolder} sx={{ color: "error.main" }}>
          <Delete sx={{ mr: 1 }} />
          Delete
        </MenuItem>
      </Menu>

      {/* Rename Dialog */}
      <Dialog open={renameDialogOpen} onClose={closeRenameDialog}>
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
          <Button onClick={closeRenameDialog}>Cancel</Button>
          <Button onClick={confirmRename} variant="contained">
            Rename
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={closeDeleteDialog}>
        <DialogTitle>Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "
            {selectedItem && ("fileName" in selectedItem ? selectedItem.fileName : selectedItem?.folderName)}"? This
            action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDeleteDialog}>Cancel</Button>
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

      {/* Move Dialog */}
      <MoveDialog open={moveDialogOpen} onClose={closeMoveDialog} onMove={confirmMove} itemToMove={selectedItem} />
    </Box>
  )
}
