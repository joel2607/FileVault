import { gql } from "@apollo/client"

export const STORAGE_STATISTICS_SUBSCRIPTION = gql`
  subscription StorageStatistics($userID: ID) {
    storageStatistics(userID: $userID) {
      usedStorageKB
      savedStorageKB
      percentageSaved
    }
  }
`

export const FILE_DOWNLOAD_COUNT_SUBSCRIPTION = gql`
  subscription FileDownloadCount($fileID: ID!) {
    fileDownloadCount(fileID: $fileID) {
      fileID
      downloadCount
    }
  }
`
