let serverBaseUrl: string | undefined

export function getCachedServerUrl() {
  return serverBaseUrl
}

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

export async function backendFetchRequest(path: string, options = {}): Promise<Response> {
  serverBaseUrl = await getServerUrl()
  const url = `${serverBaseUrl}/api/${path}`
  const response = await fetch(url, options)
  return response
}
