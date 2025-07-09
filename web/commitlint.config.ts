export default {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "type-enum": [
      2,
      "always",
      [
        "feat", // new feature
        "fix", // fix bug
        "docs", // docs update
        "style", // code style, not affect logic
        "refactor", // refactor code
        "perf", // performance optimization
        "test", // test related
        "chore", // build process or auxiliary tool changes
        "ci", // CI/CD related
        "build", // build system or external dependency changes
        "revert", // revert commit
      ],
    ],
    "subject-case": [0], // no limit on subject case
  },
}
