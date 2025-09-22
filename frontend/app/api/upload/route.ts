import { type NextRequest, NextResponse } from "next/server"
import { ApolloClient, InMemoryCache, createHttpLink } from "@apollo/client"
import { setContext } from "@apollo/client/link/context"
import { UPLOAD_FILES_MUTATION } from "@/lib/graphql/mutations"

const httpLink = createHttpLink({
  uri: process.env.INTERNAL_GRAPHQL_ENDPOINT || "http://backend:8080/query",
})

const authLink = setContext((_, { headers }) => {
  return {
    headers: {
      ...headers,
    },
  }
})

const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
})

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData()
    const files = formData.getAll("files") as File[]
    const parentFolderID = formData.get("parentFolderID") as string
    const token = request.headers.get("authorization")?.replace("Bearer ", "")

    if (!token) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const { data } = await client.mutate({
      mutation: UPLOAD_FILES_MUTATION,
      variables: {
        files,
        parentFolderID: parentFolderID || null,
      },
      context: {
        headers: {
          authorization: `Bearer ${token}`,
        },
      },
    })

    return NextResponse.json({
      success: true,
      files: data.uploadFiles,
    })
  } catch (error: any) {
    console.error("Upload error:", error)
    return NextResponse.json({ error: error.message || "Upload failed" }, { status: 500 })
  }
}

