"use client"

import { useState, useEffect } from "react"
import {
  Box,
  TextField,
  InputAdornment,
  IconButton,
  Collapse,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
  Chip,
  Button,
  Typography,
  Slider,
} from "@mui/material"
import { Search, Clear, FilterList } from "@mui/icons-material"
import { useLazyQuery } from "@apollo/client"
import { SEARCH_FILES_QUERY } from "@/lib/graphql/queries"
import { useDebounce } from "@/hooks/use-debounce"
import type { File } from "@/lib/types"

interface SearchBarProps {
  onSearchResults: (files: File[]) => void
  onClearSearch: () => void
}

const MIME_TYPES = [
  "image/jpeg",
  "image/png",
  "application/pdf",
  "text/plain",
  "video/mp4",
  "audio/mpeg",
]

export function SearchBar({ onSearchResults, onClearSearch }: SearchBarProps) {
  const [searchTerm, setSearchTerm] = useState("")
  const [showFilters, setShowFilters] = useState(false)
  const [filters, setFilters] = useState({
    mimeTypes: [] as string[],
    sizeRange: [0, 100], // MB
    dateRange: { start: "", end: "" },
    isPublic: null as boolean | null,
  })
  const debouncedSearchTerm = useDebounce(searchTerm, 300)

  const [searchFiles, { data, loading }] = useLazyQuery(SEARCH_FILES_QUERY)

  useEffect(() => {
    if (debouncedSearchTerm.length >= 3 || JSON.stringify(filters) !== JSON.stringify({ mimeTypes: [], sizeRange: [0, 100], dateRange: { start: "", end: "" }, isPublic: null })) {
      handleSearch()
    } else if (debouncedSearchTerm.length === 0) {
      onClearSearch()
    }
  }, [debouncedSearchTerm, filters])

  useEffect(() => {
    if (data) {
      onSearchResults(data.searchFiles)
    }
  }, [data, onSearchResults])

  const handleSearch = () => {
    const filterVariables: any = {
      mimeTypes: filters.mimeTypes,
      minSize: filters.sizeRange[0] * 1024 * 1024, // Convert MB to Bytes
      maxSize: filters.sizeRange[1] * 1024 * 1024, // Convert MB to Bytes
    }

    if (filters.dateRange.start) {
      filterVariables.startDate = new Date(filters.dateRange.start).toISOString()
    }
    if (filters.dateRange.end) {
      filterVariables.endDate = new Date(filters.dateRange.end).toISOString()
    }
    if (filters.isPublic !== null) {
      filterVariables.isPublic = filters.isPublic
    }

    searchFiles({
      variables: {
        query: debouncedSearchTerm,
        filter: filterVariables,
      },
    })
  }

  const handleClear = () => {
    setSearchTerm("")
    onClearSearch()
  }

  const handleFilterChange = (name: string, value: any) => {
    setFilters((prev) => ({ ...prev, [name]: value }))
  }

  const resetFilters = () => {
    setFilters({
      mimeTypes: [],
      sizeRange: [0, 100],
      dateRange: { start: "", end: "" },
      isPublic: null,
    })
  }

  return (
    <Box sx={{ mb: 3 }}>
      <TextField
        fullWidth
        variant="outlined"
        placeholder="Search for files by name, tag, or author..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Search />
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position="end">
              {searchTerm && (
                <IconButton onClick={handleClear}>
                  <Clear />
                </IconButton>
              )}
              <IconButton onClick={() => setShowFilters(!showFilters)}>
                <FilterList />
              </IconButton>
            </InputAdornment>
          ),
        }}
      />

      <Collapse in={showFilters}>
        <Box sx={{ p: 2, border: "1px solid", borderColor: "divider", borderRadius: 1, mt: 1 }}>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>File Types</InputLabel>
                <Select
                  multiple
                  value={filters.mimeTypes}
                  onChange={(e) => handleFilterChange("mimeTypes", e.target.value)}
                  input={<OutlinedInput label="File Types" />}
                  renderValue={(selected) => (
                    <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                      {(selected as string[]).map((value) => (
                        <Chip key={value} label={value} />
                      ))}
                    </Box>
                  )}
                >
                  {MIME_TYPES.map((type) => (
                    <MenuItem key={type} value={type}>
                      {type}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <Typography gutterBottom>
                File Size (MB): {filters.sizeRange[0]} - {filters.sizeRange[1]}
              </Typography>
              <Slider
                value={filters.sizeRange}
                onChange={(_, newValue) => handleFilterChange("sizeRange", newValue)}
                valueLabelDisplay="auto"
                min={0}
                max={1000}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Start Date"
                type="date"
                value={filters.dateRange.start}
                onChange={(e) => handleFilterChange("dateRange", { ...filters.dateRange, start: e.target.value })}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="End Date"
                type="date"
                value={filters.dateRange.end}
                onChange={(e) => handleFilterChange("dateRange", { ...filters.dateRange, end: e.target.value })}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Visibility</InputLabel>
                <Select
                  value={filters.isPublic === null ? "" : filters.isPublic ? "public" : "private"}
                  onChange={(e) =>
                    handleFilterChange(
                      "isPublic",
                      e.target.value === "" ? null : e.target.value === "public"
                    )
                  }
                  label="Visibility"
                >
                  <MenuItem value="">
                    <em>Any</em>
                  </MenuItem>
                  <MenuItem value="public">Public</MenuItem>
                  <MenuItem value="private">Private</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
          <Box sx={{ mt: 2, display: "flex", justifyContent: "flex-end" }}>
            <Button onClick={resetFilters} sx={{ mr: 1 }}>
              Reset Filters
            </Button>
            <Button onClick={handleSearch} variant="contained">
              Apply Filters
            </Button>
          </Box>
        </Box>
      </Collapse>
    </Box>
  )
}
