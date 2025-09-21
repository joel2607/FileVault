"use client"

import { Drawer, Toolbar, Box, Button, Typography, LinearProgress, Divider } from "@mui/material"
import { CloudUpload, Storage } from "@mui/icons-material"
import { useAuth } from "@/hooks/use-auth"
import { useSubscription } from "@apollo/client"
import { STORAGE_STATISTICS_SUBSCRIPTION } from "@/lib/graphql/subscriptions"

const DRAWER_WIDTH = 280

interface SidebarProps {
  onUploadClick: () => void
}

export function Sidebar({ onUploadClick }: SidebarProps) {
  const { user } = useAuth()

  const { data: storageData } = useSubscription(STORAGE_STATISTICS_SUBSCRIPTION, {
    variables: { userID: user?.id },
    skip: !user?.id,
  })

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return "0 Bytes"
    const k = 1024
    const sizes = ["Bytes", "KB", "MB", "GB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
  }

  const usedStorage = storageData?.storageStatistics?.usedStorageKB || user?.usedStorageKb || 0
  const totalStorage = user?.storageQuotaKb || 0
  const usedPercentage = totalStorage > 0 ? (usedStorage / totalStorage) * 100 : 0

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: DRAWER_WIDTH,
          boxSizing: "border-box",
        },
      }}
    >
      <Toolbar />
      <Box sx={{ p: 2 }}>
        <Button
          variant="contained"
          fullWidth
          startIcon={<CloudUpload />}
          onClick={onUploadClick}
          size="large"
          sx={{ mb: 3 }}
        >
          Upload Files
        </Button>

        <Divider sx={{ mb: 2 }} />

        <Box sx={{ mb: 2 }}>
          <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
            <Storage sx={{ mr: 1, fontSize: 20 }} />
            <Typography variant="subtitle2">Storage</Typography>
          </Box>

          <Typography variant="body2" color="text.secondary" gutterBottom>
            {formatBytes(usedStorage * 1024)} of {formatBytes(totalStorage * 1024)} used
          </Typography>

          <LinearProgress
            variant="determinate"
            value={Math.min(usedPercentage, 100)}
            sx={{ height: 8, borderRadius: 4 }}
          />

          <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: "block" }}>
            {usedPercentage.toFixed(1)}% used
          </Typography>
        </Box>

        {storageData?.storageStatistics?.savedStorageKB > 0 && (
          <Box>
            <Typography variant="body2" color="success.main">
              {formatBytes(storageData.storageStatistics.savedStorageKB * 1024)} saved
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {storageData.storageStatistics.percentageSaved.toFixed(1)}% compression
            </Typography>
          </Box>
        )}
      </Box>
    </Drawer>
  )
}
