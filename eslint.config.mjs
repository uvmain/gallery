import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  {
    ignores: ['node_modules', 'dist', '**/dist/**', 'public', '**/public/**'],
    files: ['**/*.{vue,ts,js}'],
    rules: {
      'no-console': 'off',
      'brace-style': ['error', 'stroustrup'],
      'curly': ['off'],
      'semi': ['error', 'never'],
      'vue/multi-word-component-names': 'off',
      'vue/no-v-html': 'off'
    }
  },
)


