import { repositoryExec } from "./repositoryCommand.js";
async function findGitRoot(directory) {
  return repositoryExec(directory, "git", [
    "rev-parse",
    "--quiet",
    "--show-toplevel"
  ]);
}
async function findSaplingRoot(directory) {
  return repositoryExec(directory, "sl", ["--quiet", "root"]);
}
async function findMercurialRoot(directory) {
  return repositoryExec(directory, "hg", ["--quiet", "root"]);
}
async function findSubversionRoot(directory) {
  return repositoryExec(directory, "svn", ["info", "--show-item", "wc-root"]);
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
async function findRoot(directory) {
  for (const method of findMethods) {
    try {
      const result = await method(directory);
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
  findRoot
};
