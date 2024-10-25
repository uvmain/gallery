let serverBaseUrl: string | undefined

export async function getServerUrl(): Promise<string> {
  if (serverBaseUrl !== undefined) {
    return serverBaseUrl
  }
  try {
    const response = await fetch(`/api/slugs?offset=0&limit=1`)
    await response.json() as string[]
    serverBaseUrl = ''
    return serverBaseUrl
  }
  catch {
    serverBaseUrl = 'http://localhost:8080'
    return serverBaseUrl
  }
}
