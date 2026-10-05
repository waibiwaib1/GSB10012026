import { VersionLabel } from './release';

enum SEMVER {
  major = 'major',
  premajor = 'premajor',
  minor = 'minor',
  preminor = 'preminor',
  patch = 'patch',
  prepatch = 'prepatch',
  noVersion = ''
}

export type IVersionLabels = Map<VersionLabel, string[]>;

export default SEMVER;

export function getHigherSemverTag(left: SEMVER, right: string): SEMVER {
  if (left === SEMVER.major || right === SEMVER.major) {
    return SEMVER.major;
  }

  if (left === SEMVER.minor || right === SEMVER.minor) {
    return SEMVER.minor;
  }

  return SEMVER.patch;
}

interface ISemVerOptions {
  onlyPublishWithReleaseLabel?: boolean;
  skipReleaseLabels?: string[];
}

export function calculateSemVerBump(
  labels: string[][],
  labelMap: IVersionLabels,
  { onlyPublishWithReleaseLabel, skipReleaseLabels = [] }: ISemVerOptions = {}
) {
  const labelSet = new Set<string>();

  const skipLabels = labelMap.get('skip-release')!;

  const allSkipReleaseLabels = [
    ...new Set([...skipReleaseLabels, ...skipLabels])
  ];

  labels.map(pr => {
    pr.forEach(label => {
      const userLabel = [...labelMap.entries()].find(pair =>
        pair[1].includes(label)
      );
      labelSet.add(userLabel ? userLabel[0] : label);
    });
  });

  let skipRelease = false;
  let isPrerelease = false;

  if (labels.length > 0 && labels[0].length > 0) {
    isPrerelease = labels[0].some(label =>
      labelMap.get('prerelease')!.includes(label)
    );
    skipRelease = onlyPublishWithReleaseLabel
      ? !labels[0].some(label => labelMap.get('release')!.includes(label))
      : !!labels[0].find(label => allSkipReleaseLabels.includes(label));
  }

  const version = [...labelSet].reduce(getHigherSemverTag, SEMVER.patch);

  if (skipRelease) {
    return SEMVER.noVersion;
  }

  if (isPrerelease) {
    if (version === SEMVER.major) {
      return SEMVER.premajor;
    }

    if (version === SEMVER.minor) {
      return SEMVER.preminor;
    }

    return SEMVER.prepatch;
  }

  return version;
}
