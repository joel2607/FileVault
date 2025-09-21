"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { Box, TextField, InputAdornment, IconButton, Menu, MenuItem, Chip, Typography, Divider } from "@mui/material"
import { Search, FilterList, Clear } from "@mui/icons-material"
import { useLazyQuery } from "@apollo/client"
import { SEARCH_FILES_QUERY } from "@/lib/graphql/queries"
import type { File } from "@/lib/types"

interface SearchBarProps {
  onSearchResults: (files: File[]) => void
  onClearSearch: () => void
}

interface SearchFilters {
  mimeType?: string
  isPublic?: boolean
  minSize?: number
  maxSize?: number
}

export function SearchBar({ onSearchResults, onClearSearch }: SearchBarProps) {
  const [searchTerm, setSearchTerm] = useState("")
  const [filters, setFilters] = useState<SearchFilters>({})
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [activeFilters, setActiveFilters] = useState<string[]>([])

  const [searchFiles, { loading }] = useLazyQuery(SEARCH_FILES_QUERY)

  const mimeTypeOptions = [
    { value: "image/", label: "Images" },
    { value: "video/", label: "Videos" },
    { value: "audio/", label: "Audio" },
    { value: "application/pdf", label: "PDFs" },
    { value: "application/msword", label: "Documents" },
    { value: "application/vnd.ms-excel", label: "Spreadsheets" },
  ]

  const sizeOptions = [
    { value: { min: 0, max: 1024 * 1024 }, label: "< 1 MB" },
    { value: { min: 1024 * 1024, max: 10 * 1024 * 1024 }, label: "1-10 MB" },
    { value: { min: 10 * 1024 * 1024, max: 100 * 1024 * 1024 }, label: "10-100 MB" },
    { value: { min: 100 * 1024 * 1024, max: undefined }, label: "> 100 MB" },
  ]

  useEffect(() => {
    const delayedSearch = setTimeout(() => {
      if (searchTerm.trim() || Object.keys(filters).length > 0) {
        handleSearch()
      }
    }, 500)

    return () => clearTimeout(delayedSearch)
  }, [searchTerm, filters])

  const handleSearch = async () => {
    try {
      const searchFilter: any = {}

      if (searchTerm.trim()) {
        searchFilter.fileName = searchTerm.trim()
      }

      if (filters.mimeType) {
        searchFilter.mimeType = filters.mimeType
      }

      if (filters.isPublic !== undefined) {
        searchFilter.isPublic = filters.isPublic
      }

      if (filters.minSize !== undefined) {
        searchFilter.minSize = filters.minSize
      }

      if (filters.maxSize !== undefined) {
        searchFilter.maxSize = filters.maxSize
      }

      const { data } = await searchFiles({
        variables: { filter: searchFilter },
      })

      if (data?.searchFiles) {
        onSearchResults(data.searchFiles)
      }
    } catch (error) {
      console.error("Search failed:", error)
    }
  }

  const handleFilterClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleFilterClose = () => {
    setAnchorEl(null)
  }

  const handleMimeTypeFilter = (mimeType: string, label: string) => {
    const newFilters = { ...filters, mimeType }
    setFilters(newFilters)
    setActiveFilters((prev) => [...prev.filter((f) => !f.startsWith("Type:")), `Type: ${label}`])
    handleFilterClose()
  }

  const handleVisibilityFilter = (isPublic: boolean) => {
    const newFilters = { ...filters, isPublic }
    setFilters(newFilters)
    setActiveFilters((prev) => [
      ...prev.filter((f) => !f.startsWith("Visibility:")),
      `Visibility: ${isPublic ? "Public" : "Private"}`,
    ])
    handleFilterClose()
  }

  const handleSizeFilter = (sizeRange: { min: number; max?: number }, label: string) => {
    const newFilters = { ...filters, minSize: sizeRange.min, maxSize: sizeRange.max }
    setFilters(newFilters)
    setActiveFilters((prev) => [...prev.filter((f) => !f.startsWith("Size:")), `Size: ${label}`])
    handleFilterClose()
  }

  const removeFilter = (filterToRemove: string) => {
    const newActiveFilters = activeFilters.filter((f) => f !== filterToRemove)
    setActiveFilters(newActiveFilters)

    const newFilters = { ...filters }
    if (filterToRemove.startsWith("Type:")) {
      delete newFilters.mimeType
    } else if (filterToRemove.startsWith("Visibility:")) {
      delete newFilters.isPublic
    } else if (filterToRemove.startsWith("Size:")) {
      delete newFilters.minSize
      delete newFilters.maxSize
    }
    setFilters(newFilters)
  }

  const clearAllFilters = () => {
    setSearchTerm("")
    setFilters({})
    setActiveFilters([])
    onClearSearch()
  }

  const hasActiveSearch = searchTerm.trim() || activeFilters.length > 0

  return (
    <Box sx={{ mb: 3 }}>
      <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
        <TextField
          fullWidth
          placeholder="Search files..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search />
              </InputAdornment>
            ),
            endAdornment: hasActiveSearch && (
              <InputAdornment position="end">
                <IconButton onClick={clearAllFilters} size="small">
                  <Clear />
                </IconButton>
              </InputAdornment>
            ),
          }}
          disabled={loading}
        />

        <IconButton onClick={handleFilterClick} color={activeFilters.length > 0 ? "primary" : "default"}>
          <FilterList />
        </IconButton>
      </Box>

      {/* Active Filters */}
      {activeFilters.length > 0 && (
        <Box sx={{ mt: 2, display: "flex", gap: 1, flexWrap: "wrap", alignItems: "center" }}>
          <Typography variant="body2" color="text.secondary">
            Filters:
          </Typography>
          {activeFilters.map((filter) => (
            <Chip
              key={filter}
              label={filter}
              onDelete={() => removeFilter(filter)}
              size="small"
              color="primary"
              variant="outlined"
            />
          ))}
        </Box>
      )}

      {/* Filter Menu */}
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleFilterClose}>
        <MenuItem disabled>
          <Typography variant="subtitle2">File Type</Typography>
        </MenuItem>
        {mimeTypeOptions.map((option) => (
          <MenuItem key={option.value} onClick={() => handleMimeTypeFilter(option.value, option.label)}>
            {option.label}
          </MenuItem>
        ))}

        <Divider />

        <MenuItem disabled>
          <Typography variant="subtitle2">Visibility</Typography>
        </MenuItem>
        <MenuItem onClick={() => handleVisibilityFilter(true)}>Public Files</MenuItem>
        <MenuItem onClick={() => handleVisibilityFilter(false)}>Private Files</MenuItem>

        <Divider />

        <MenuItem disabled>
          <Typography variant="subtitle2">File Size</Typography>
        </MenuItem>
        {sizeOptions.map((option) => (
          <MenuItem key={option.label} onClick={() => handleSizeFilter(option.value, option.label)}>
            {option.label}
          </MenuItem>
        ))}
      </Menu>
    </Box>
  )
}
