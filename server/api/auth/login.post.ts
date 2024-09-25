export default defineEventHandler(async (event) => {
  try {
    const body = await readBody(event)
    if (!body) {
      console.error("Request body is empty or undefined")
      return createError({
        statusCode: 400,
        statusMessage: "Request body is empty or undefined",
      })
    }

    const { username, password } = body

    if (!username || !password) {
      console.error("Username or password missing")
      return createError({
        statusCode: 400,
        statusMessage: "Username and password are required",
      })
    }

    if (username != serverConfiguration.adminUser || password != serverConfiguration.adminPassword) {
      console.error(`Invalid username or password for admin user: ${username}`)
      return createError({
        statusCode: 401,
        statusMessage: "Invalid username or password",
      })
    }
    else {
      await setUserSession(event, {
        user: serverConfiguration.adminUser,
        loggedInAt: new Date(),
      },{
        maxAge: (serverConfiguration.sessionMaxAge || 0) as number
      })
    }

    return { success: true, statusMessage: "success", username }

  }
  catch (error) {
    console.error("Error handling login request:", error)
    return createError({
      statusCode: 500,
      statusMessage: "Failed to process request",
    })
  }
})
