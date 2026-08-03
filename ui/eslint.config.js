import vuetify from 'eslint-config-vuetify'

export default vuetify({
  ts: true,
  rules: {
    'vue/max-attributes-per-line': [
      'error',
      {
        singleline: {
          max: 3, // Up to 3 props allowed on a single line
        },
        multiline: {
          max: 1, // If it breaks into multiple lines, strictly 1 prop per line
        },
      },
    ],
  },
})
