module.exports = {
  plugins: [
    require("postcss-easy-import"),
    require("tailwindcss"),
    process.env.NODE_ENV === "production"
      ? require("@fullhuman/postcss-purgecss")({
          content: ["./pages/**/*.tsx", "./components/**/*.tsx"],
          defaultExtractor: content => content.match(/[A-Za-z0-9-_:/]+/g) || []
        })
      : undefined,
    require("autoprefixer"),
    require("cssnano")
  ].filter(v => !!v)
};
