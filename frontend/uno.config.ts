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
  shortcuts: {
    standard: 'bg-neutral-200 dark:bg-neutral-800 text-neutral-700 dark:text-neutral-100',
  },
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
