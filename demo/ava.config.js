import path from "path";

export default {
  compileEnhancements: false,
  extensions: ["ts"],
  require: ["tsconfig-paths/register", "ts-node/register"],
  files: ["!packages/**/*"],
  environmentVariables: {
    TS_NODE_PROJECT: path.resolve(__dirname, "tsconfig.node.json")
  }
};
