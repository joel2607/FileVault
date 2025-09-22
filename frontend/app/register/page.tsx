"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Container, Card, CardContent, TextField, Button, Typography, Box, Alert, Link as MuiLink } from "@mui/material"
import Link from "next/link"
import { useAuth } from "@/hooks/use-auth"

export default function RegisterPage() {
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const { register } = useAuth()
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError("")

    if (password !== confirmPassword) {
      setError("Passwords do not match")
      setLoading(false)
      return
    }

    if (password.length < 6) {
      setError("Password must be at least 6 characters long")
      setLoading(false)
      return
    }

    const result = await register(username, email, password)

    if (result.success) {
      router.push("/")
    } else {
      setError(result.error || "Registration failed")
    }

    setLoading(false)
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Card elevation={3}>
        <CardContent sx={{ p: 4 }}>
          <Box textAlign="center" mb={3}>
            <Typography variant="h4" component="h1" gutterBottom>
              CloudDrive
            </Typography>
            <Typography variant="h6" color="text.secondary">
              Create your account
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              margin="normal"
              required
              autoComplete="username"
              autoFocus
            />

            <TextField
              fullWidth
              label="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              margin="normal"
              required
              autoComplete="email"
            />

            <TextField
              fullWidth
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              margin="normal"
              required
              autoComplete="new-password"
            />

            <TextField
              fullWidth
              label="Confirm Password"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              margin="normal"
              required
              autoComplete="new-password"
            />

            <Button type="submit" fullWidth variant="contained" size="large" disabled={loading} sx={{ mt: 3, mb: 2 }}>
              {loading ? "Creating account..." : "Sign Up"}
            </Button>

            <Box textAlign="center">
              <Typography variant="body2">
                Already have an account?{" "}
                <Link href="/login" passHref>
                  <MuiLink component="span" sx={{ cursor: "pointer" }}>
                    Sign in
                  </MuiLink>
                </Link>
              </Typography>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Container>
  )
}
