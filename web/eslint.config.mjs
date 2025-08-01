import antfu from "@antfu/eslint-config"
import nextPlugin from "@next/eslint-plugin-next"

export default antfu(
  {
    react: true,
    lessOpinionated: true,
    markdown: true,
    stylistic: {
      quotes: "double",
    },
    ignores: [
      "src/components/ui/**",
      ".vscode/**",
      ".cursor/**",
      "pnpm-lock.yaml",
      "**/*.md",
      ".source/",
      ".contentlayer/",
    ],
    rules: {
      "node/prefer-global/process": "off",
      "style/eol-last": "off",
      "react-hooks-extra/no-direct-set-state-in-use-effect": "off",
      "antfu/curly": "error",
      "curly": "off",
      "style/brace-style": ["error", "1tbs"],
      "no-console": "warn",
      "antfu/no-top-level-await": "off",
      "unused-imports/no-unused-vars": "warn",
    },
  },
  {
    plugins: {
      "@next/next": nextPlugin,
    },
    rules: {
      ...nextPlugin.configs.recommended.rules,
      ...nextPlugin.configs["core-web-vitals"].rules,
      "@next/next/no-img-element": "off",
    },
  },
)
