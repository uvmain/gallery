export default defineNitroPlugin(async() => {
  await createDatabaseDirectory()
  await createMetadataTable()
  await createImagesDirectory()
  await createThumbnailsDirectory()
  await createThumbnailsForAllImages()
  await createMetaDataForAllImages()
})