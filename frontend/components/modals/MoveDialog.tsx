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
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from "@mui/material"
import { Folder as FolderIcon, ArrowBack as ArrowBackIcon } from "@mui/icons-material"
import { useQuery } from "@apollo/client"
import { ROOT_QUERY, FOLDER_QUERY } from "@/lib/graphql/queries"
import type { File, Folder } from "@/lib/types"

interface MoveDialogProps {
  open: boolean
  onClose: () => void
  onMove: (destinationFolderId: string | null) => void
  itemToMove: File | Folder | null
}

export function MoveDialog({ open, onClose, onMove, itemToMove }: MoveDialogProps) {
  const [currentFolderId, setCurrentFolderId] = useState<string | null>(null)
  const [folderHistory, setFolderHistory] = useState<(string | null)[]>([])

  const { data, loading } = useQuery(currentFolderId ? FOLDER_QUERY : ROOT_QUERY, {
    variables: { id: currentFolderId },
    skip: !open,
  })

  const handleFolderClick = (folderId: string) => {
    setFolderHistory([...folderHistory, currentFolderId])
    setCurrentFolderId(folderId)
  }

  const handleBackClick = () => {
    const previousFolderId = folderHistory[folderHistory.length - 1]
    setFolderHistory(folderHistory.slice(0, -1))
    setCurrentFolderId(previousFolderId)
  }

  const handleMoveClick = () => {
    onMove(currentFolderId)
  }

  const folders = currentFolderId ? data?.folder?.folders : data?.root?.folders

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>
        <Box sx={{ display: "flex", alignItems: "center" }}>
          {folderHistory.length > 0 && (
            <IconButton onClick={handleBackClick} sx={{ mr: 1 }}>
              <ArrowBackIcon />
            </IconButton>
          )}
          Move Item
        </Box>
      </DialogTitle>
      <DialogContent>
        {loading ? (
          <Typography>Loading...</Typography>
        ) : (
          <List>
            {folders?.map((folder: Folder) => (
              <ListItem button key={folder.id} onClick={() => handleFolderClick(folder.id)}>
                <ListItemIcon>
                  <FolderIcon />
                </ListItemIcon>
                <ListItemText primary={folder.folderName} />
              </ListItem>
            ))}
          </List>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={handleMoveClick}
          variant="contained"
          disabled={
            !itemToMove ||
            ("parentFolderId" in itemToMove && itemToMove.parentFolderId === currentFolderId) ||
            ("parentFolderId" in itemToMove && !itemToMove.parentFolderId && currentFolderId === null)
          }
        >
          Move Here
        </Button>
      </DialogActions>
    </Dialog>
  )
}
