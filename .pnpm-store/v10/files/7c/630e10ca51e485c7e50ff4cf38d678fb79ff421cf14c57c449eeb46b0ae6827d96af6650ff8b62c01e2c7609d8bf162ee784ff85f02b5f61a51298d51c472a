import { repositoryExecSync } from "./repositoryCommand.js";
function findGitRoot(cwd) {
  return repositoryExecSync(cwd, "git", [
    "rev-parse",
    "--quiet",
    "--show-toplevel"
  ]);
}
function findSaplingRoot(cwd) {
  return repositoryExecSync(cwd, "sl", ["--quiet", "root"]);
}
function findMercurialRoot(cwd) {
  return repositoryExecSync(cwd, "hg", ["--quiet", "root"]);
}
function findSubversionRoot(cwd) {
  return repositoryExecSync(cwd, "svn", ["info", "--show-item", "wc-root"]);
}
const findMethods = [
  findGitRoot,
  // most used and fastest check
  findSubversionRoot,
  // next highest market share
  findSaplingRoot,
  // much smaller market share, but much faster
  findMercurialRoot
  // slow af
];
function findRootSync(directory) {
  for (const method of findMethods) {
    try {
      const result = method(directory);
      if (result != null && result !== "") {
        return result;
      }
    } catch (e) {
      if (e instanceof Error && "code" in e && e.code === "ENOTDIR") {
        throw new Error(`Given path isn't a directory: ${directory}`);
      }
    }
  }
  return null;
}
export {
  findRootSync
};
