import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetTypography,
  presetUno,
  presetWebFonts,
  presetWind,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

function getSafelist(): string[] {
  const base = 'prose prose-sm m-auto text-left'.split(' ')
  const unusedSafelist: string[] = []
  return [...unusedSafelist, ...base]
}

export default defineConfig({
  shortcuts: {},
  theme: {
    colors: {},
  },
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons(),
    presetWind(),
    presetTypography(),
    presetWebFonts({
      fonts: {
        sans: 'Open Sans',
        serif: 'DM Serif Display',
        mono: 'DM Mono',
      },
    }),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
  safelist: getSafelist(),
})
