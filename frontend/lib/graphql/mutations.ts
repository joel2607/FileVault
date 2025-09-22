import { gql } from "@apollo/client"

export const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      token
      user {
        id
        username
        email
        storageQuotaKb
        usedStorageKb
        savedStorageKb
        role
      }
    }
  }
`

export const REGISTER_MUTATION = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
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

export const UPLOAD_FILES_MUTATION = gql`
  mutation UploadFiles($files: [Upload!]!, $parentFolderID: ID) {
    uploadFiles(files: $files, parentFolderID: $parentFolderID) {
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

export const CREATE_FOLDER_MUTATION = gql`
  mutation CreateFolder($input: NewFolder!) {
    createFolder(input: $input) {
      id
      folderName
      parentFolderId
      isPublic
    }
  }
`

export const UPDATE_FOLDER_MUTATION = gql`
  mutation UpdateFolder($input: UpdateFolder!) {
    updateFolder(input: $input) {
      id
      folderName
      parentFolderId
      isPublic
    }
  }
`

export const DELETE_FOLDER_MUTATION = gql`
  mutation DeleteFolder($id: ID!) {
    deleteFolder(id: $id) {
      id
    }
  }
`

export const UPDATE_FILE_MUTATION = gql`
  mutation UpdateFile($input: UpdateFile!) {
    updateFile(input: $input) {
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

export const DELETE_FILE_MUTATION = gql`
  mutation DeleteFile($id: ID!) {
    deleteFile(id: $id) {
      id
    }
  }
`

export const GENERATE_DOWNLOAD_URL_MUTATION = gql`
  mutation GenerateDownloadUrl($fileID: ID!) {
    generateDownloadUrl(fileID: $fileID)
  }
`

export const SET_FILE_PUBLIC_MUTATION = gql`
  mutation SetFilePublic($fileID: ID!) {
    setFilePublic(fileID: $fileID) {
      id
      isPublic
    }
  }
`

export const SET_FILE_PRIVATE_MUTATION = gql`
  mutation SetFilePrivate($fileID: ID!) {
    setFilePrivate(fileID: $fileID) {
      id
      isPublic
    }
  }
`

export const SHARE_FILE_WITH_USER_MUTATION = gql`
  mutation ShareFileWithUser($fileID: ID!, $userID: ID!) {
    shareFileWithUser(fileID: $fileID, userID: $userID) {
      id
      sharedWithUser {
        id
        username
        email
      }
      file {
        id
        fileName
      }
    }
  }
`

export const REMOVE_FILE_ACCESS_MUTATION = gql`
  mutation RemoveFileAccess($fileID: ID!, $userID: ID!) {
    removeFileAccess(fileID: $fileID, userID: $userID)
  }
`
