import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'
import express from 'express'
import 'dotenv/config'

const app = express()
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

app.use(express.static(path.join(__dirname, 'dist')))

app.get('*', (req, res) => {
  console.log(`Serving index.html for: ${req.originalUrl}`)
  res.sendFile(path.join(__dirname, 'dist/index.html'))
})

app.listen(process.env.PORT, () => {
  console.log(`Static server started and listening on port ${process.env.PORT}`)
})
