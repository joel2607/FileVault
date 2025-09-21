"use client"

import { useState, useEffect } from "react"
import { useMutation, useQuery } from "@apollo/client"
import { LOGIN_MUTATION, REGISTER_MUTATION } from "@/lib/graphql/mutations"
import { ME_QUERY } from "@/lib/graphql/queries"
import type { User } from "@/lib/types"

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const { data: meData, refetch: refetchMe } = useQuery(ME_QUERY, {
    skip: typeof window === "undefined" || !localStorage.getItem("token"),
    onCompleted: (data) => {
      setUser(data.me)
      setLoading(false)
    },
    onError: () => {
      localStorage.removeItem("token")
      setUser(null)
      setLoading(false)
    },
  })

  const [loginMutation] = useMutation(LOGIN_MUTATION)
  const [registerMutation] = useMutation(REGISTER_MUTATION)

  useEffect(() => {
    const token = localStorage.getItem("token")
    if (!token) {
      setLoading(false)
    }
  }, [])

  const login = async (email: string, password: string) => {
    try {
      const { data } = await loginMutation({
        variables: { email, password },
      })

      if (data?.login) {
        localStorage.setItem("token", data.login.token)
        setUser(data.login.user)
        return { success: true }
      }
    } catch (error: any) {
      return { success: false, error: error.message }
    }
  }

  const register = async (username: string, email: string, password: string) => {
    try {
      const { data } = await registerMutation({
        variables: { input: { username, email, password } },
      })

      if (data?.register) {
        // Auto-login after registration
        return await login(email, password)
      }
    } catch (error: any) {
      return { success: false, error: error.message }
    }
  }

  const logout = () => {
    localStorage.removeItem("token")
    setUser(null)
    window.location.href = "/login"
  }

  return {
    user,
    loading,
    login,
    register,
    logout,
    isAuthenticated: !!user,
  }
}
