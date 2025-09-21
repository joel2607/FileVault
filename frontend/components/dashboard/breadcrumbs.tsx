"use client"

import { Breadcrumbs, Link, Typography, Box } from "@mui/material"
import { Home, NavigateNext } from "@mui/icons-material"
import type { Folder } from "@/lib/types"

interface DashboardBreadcrumbsProps {
  currentPath: Folder[]
  onNavigate: (folderId?: string) => void
}

export function DashboardBreadcrumbs({ currentPath, onNavigate }: DashboardBreadcrumbsProps) {
  return (
    <Box sx={{ mb: 2 }}>
      <Breadcrumbs separator={<NavigateNext fontSize="small" />} aria-label="breadcrumb">
        <Link
          component="button"
          variant="body1"
          onClick={() => onNavigate()}
          sx={{
            display: "flex",
            alignItems: "center",
            textDecoration: "none",
            color: "primary.main",
            "&:hover": {
              textDecoration: "underline",
            },
          }}
        >
          <Home sx={{ mr: 0.5, fontSize: 20 }} />
          Home
        </Link>

        {currentPath.map((folder, index) => {
          const isLast = index === currentPath.length - 1

          if (isLast) {
            return (
              <Typography key={folder.id} color="text.primary">
                {folder.folderName}
              </Typography>
            )
          }

          return (
            <Link
              key={folder.id}
              component="button"
              variant="body1"
              onClick={() => onNavigate(folder.id)}
              sx={{
                textDecoration: "none",
                color: "primary.main",
                "&:hover": {
                  textDecoration: "underline",
                },
              }}
            >
              {folder.folderName}
            </Link>
          )
        })}
      </Breadcrumbs>
    </Box>
  )
}
