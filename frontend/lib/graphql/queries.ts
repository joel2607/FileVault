import { gql } from "@apollo/client"

export const ME_QUERY = gql`
  query Me {
    me {
      id
      username
      email
      storageQuotaKb
      usedStorageKb
      savedStorageKb
      role
    }
  }
`

export const ROOT_QUERY = gql`
  query Root {
    root {
      files {
        id
        fileName
        mimeType
        size
        isPublic
        downloadCount
        parentFolderId
      }
      folders {
        id
        folderName
        parentFolderId
        isPublic
      }
    }
  }
`

export const FOLDER_QUERY = gql`
  query Folder($id: ID!) {
    folder(id: $id) {
      id
      folderName
      parentFolderId
      isPublic
      files {
        id
        fileName
        mimeType
        size
        isPublic
        downloadCount
        parentFolderId
      }
      folders {
        id
        folderName
        parentFolderId
        isPublic
      }
    }
  }
`

export const SEARCH_USERS_QUERY = gql`
  query SearchUsers($query: String!) {
    searchUsers(query: $query) {
      id
      username
      email
    }
  }
`

export const GET_FILE_ACCESS_QUERY = gql`
  query GetUsersWithAccess($fileID: ID!) {
    getUsersWithAccess(fileID: $fileID) {
      id
      user {
        id
        username
        email
      }
    }
  }
`

export const SEARCH_FILES_QUERY = gql`
  query SearchFiles($query: String, $filter: FileFilterInput) {
    searchFiles(query: $query, filter: $filter) {
      id
      fileName
      mimeType
      size
      isPublic
      downloadCount
      parentFolderId
    }
  }
`
