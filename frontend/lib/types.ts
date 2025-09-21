export interface User {
  id: string
  username: string
  email: string
  storageQuotaKb: number
  usedStorageKb: number
  savedStorageKb: number
  role: "USER" | "ADMIN"
}

export interface File {
  id: string
  fileName: string
  mimeType: string
  size: number
  isPublic: boolean
  downloadCount: number
  parentFolderId?: string
}

export interface Folder {
  id: string
  folderName: string
  parentFolderId?: string
  isPublic: boolean
  files?: File[]
  folders?: Folder[]
}

export interface StorageStatistics {
  usedStorageKB: number
  savedStorageKB: number
  percentageSaved: number
}

export interface AuthResponse {
  token: string
  user: User
}
